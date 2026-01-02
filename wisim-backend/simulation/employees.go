package simulation

import (
	"fmt"
	"math/rand"
)

// Employee functions

func (c *Company) simulateEmployees(externalFactors ExternalFactors) {
	employeeIDs := c.Get_employees_ids(Employee_type_all)
	c.turnover(employeeIDs)

	// refresh employee list
	employeeIDs = c.Get_employees_ids(Employee_type_all)
	productioEemployeeIDs := c.Get_employees_ids(Employee_type_production)
	marketingEemployeeIDs := c.Get_employees_ids(Employee_type_marketing)
	c.pay(employeeIDs)
	c.motivateAndTrain(productioEemployeeIDs, externalFactors.ProductionMinimumWage, 0.01)
	c.motivateAndTrain(marketingEemployeeIDs, externalFactors.MarketingMinimumWage, 0.01)

	// Increment months at company
}

func (c *Company) turnover(employeeIDs []int) {
	for _, ID := range employeeIDs {
		if c.employeePool[ID].MonthsAtCompany == 0 {
			continue
		}

		if c.employeePool[ID].Motivation <= 0.5 {
			c.employeePool[ID].Employer = Employee_employer_none
		} else {
			c.employeePool[ID].MonthsAtCompany += 1
		}
	}
}

func (c *Company) pay(employeeIDs []int) {
	var totalPay float64
	for _, ID := range employeeIDs {
		totalPay += c.employeePool[ID].Pay

		invoiceTag := production_personelle
		if c.employeePool[ID].EmployeeType == Employee_type_marketing {
			invoiceTag = marketing_personelle
		} else if c.employeePool[ID].EmployeeType != Employee_type_production {
			invoiceTag = other_personelle
		}

		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(fmt.Sprint("Pay for employee %s (%d)", c.employeePool[ID].Name, c.employeePool[ID].ID), invoiceTag, "", true, -float64(c.employeePool[ID].Pay))
	}

	println("Total employee pay:", totalPay)
}

func (c *Company) motivateAndTrain(employeeIDs []int, minimumWage, passiveTraining float64) {
	var averagePay float64

	// Motivation
	if len(employeeIDs) != 0 {
		var totalPay float64 = 0

		for _, ID := range employeeIDs {
			totalPay += c.employeePool[ID].Pay
		}

		averagePay = totalPay / float64(len(employeeIDs))
	}

	for _, ID := range employeeIDs {
		var motivationFactor float64 = 0

		motivationFactor += (c.employeePool[ID].Pay - averagePay) / 1000
		motivationFactor += (c.employeePool[ID].Pay - minimumWage) / 100
		motivationFactor += (c.employeePool[ID].ExtraTraining) / 10000
		motivationFactor += (8 - c.employeePool[ID].WorkingHours) / 20

		const motivationNormalisationFactor = 0.01

		scaleMotivationFactor := func(factor, maxvalue float64) float64 {
			if factor < 0 {
				return -(-1/(1-factor) + 1) * maxvalue
			}
			return (-1/(1+factor) + 1) * maxvalue
		}

		c.employeePool[ID].Motivation += float64(scaleMotivationFactor(float64(motivationFactor), 1))

		// passive motivation normalisation
		if c.employeePool[ID].Motivation > 1 {
			c.employeePool[ID].Motivation -= motivationNormalisationFactor
		} else if c.employeePool[ID].Motivation < 1 {
			c.employeePool[ID].Motivation += motivationNormalisationFactor
		}

		// clamp motivation
		c.employeePool[ID].Motivation = min(c.employeePool[ID].Motivation, 2)

		c.employeePool[ID].Skill += passiveTraining
		// clamp skill
		c.employeePool[ID].Skill = min(c.employeePool[ID].Skill, 2)

		if c.employeePool[ID].ExtraTraining > 0 {
			c.employeePool[ID].Skill += c.employeePool[ID].ExtraTraining / 1000
			c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Extra training for employee: "+c.employeePool[ID].Name, employee_training, "Training improves speed and quality of employeess' work", true, float64(c.employeePool[ID].ExtraTraining))
		}

		fmt.Printf("Motivation emloyee %d: %f, %f\n", ID, c.employeePool[ID].Motivation, motivationFactor)
	}
}

/*
func layoff(employees []Employee, size_of_layoff int) ([]Employee, int) {
	// find the x worst performing employees
	if size_of_layoff >= 0 {
		return employees, 0
	} else {
		size_of_layoff = -size_of_layoff
	}

	find_worst_employee := func() int {
		worst_employee := 0

		for i, e := range employees {
			if (e.Motivation * e.Working_hours * e.Skill) < (employees[worst_employee].Motivation * employees[worst_employee].Working_hours * employees[worst_employee].Skill) {
				worst_employee = i
			}
		}
		return worst_employee
	}

	number_of_employees_who_left := 0

	for range size_of_layoff {
		to_fire := find_worst_employee()

		if len(employees) >= 1 {
			employees = delete_by_index(employees, to_fire)
			number_of_employees_who_left += 1
		} else {
			return employees, number_of_employees_who_left
		}
	}

	return employees, number_of_employees_who_left
}
*/

func (g *GameState) handleEmployeeDeltas() {
	// Handle Changes & fires
	type hiringDelta struct {
		Company int
		Delta   Delta[Employee]
	}

	hiring := make(map[int][]hiringDelta)

	applyChanges := func(target *Employee, source Employee) {
		target.Bonus = source.Bonus
		target.EmployeeType = source.EmployeeType
		target.Pay = source.Pay
		target.WorkingHours = source.WorkingHours
	}

	handleDeltas := func(deltas []Delta[Employee], company int) {
		for _, d := range deltas {
			switch d.Change {
			case Delta_Change:
				applyChanges(g.Employees[d.Item.ID], d.Item)
			case Delta_Remove:
				g.Employees[d.Item.ID].Employer = Employee_employer_none
			case Delta_New:
				hiring[d.Item.ID] = append(hiring[d.Item.ID], hiringDelta{Company: company, Delta: d})
			}
		}
	}

	for i, decisions := range g.CurrentDecisions {
		handleDeltas(decisions.Employees.MarketingDeltas, i)
		handleDeltas(decisions.Employees.ProductionDeltas, i)
	}

	for i, hd := range hiring {
		winnerIndex := 0
		if len(hd) > 0 {
			winnerIndex = rand.Intn(len(hd))
		}

		applyChanges(g.Employees[i], hd[winnerIndex].Delta.Item)
		g.Employees[i].Employer = hd[winnerIndex].Company
	}
}
