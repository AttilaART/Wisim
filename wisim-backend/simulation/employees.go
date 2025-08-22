package simulation

import (
	"errors"
	"fmt"
	"math/rand"
)

// Employee functions

func (c *Company) simulate_employees(employee_deltas []Delta[Employee], external_factors External_factors, severance_pay float32, employee_type Employee_type) error {
	// Get correct employee slice
	company_employees := c.employee_pool.Get_employees_of_company(c.Id, employee_type)

	// get corresponding pointers for employees & update values
	var employee_deltas_ptr []Delta[*Employee]
	for _, e_delta := range employee_deltas {
		employee_ptr := c.employee_pool.Find_employee_by_id(e_delta.Item.Id)
		if employee_ptr == nil {
			panic("Employee Not Found")
		}

		*employee_ptr = e_delta.Item
		employee_deltas_ptr = append(employee_deltas_ptr, Delta[*Employee]{Change: e_delta.Change, Item: employee_ptr})
	}

	if company_employees == nil {
		panic("No employee list loaded")
	}

	// Fire/"layoff" employees
	var employees_layed_off int
	for i := range employee_deltas_ptr {
		if employee_deltas_ptr[i].Change == Delta_Remove {
			employees_layed_off += 1
			employee_deltas_ptr[i].Item.Employer = Employee_employer_none
		}
	}

	// Calcualte Turnover
	var num_employees_who_left int
	for _, e_ptr := range company_employees {
		if e_ptr.Employer != Employee_employer_none {
			if rand.Float32() <= external_factors.Turnover {
				e_ptr.Employer = Employee_employer_none
				num_employees_who_left += 1
			}
		}
	}

	// calculate severance pay

	for range employees_layed_off {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
			"Severance pay for employee",
			severance,
			"When you layoff an employee, you have to pay them severance. Sometimes it's more expensive to fire someone, than just letting them sit idle.",
			true,
			-float64(severance),
		)
	}

	// Refresh employees after some left
	company_employees = c.employee_pool.Get_employees_of_company(c.Id, employee_type)

	// Calcualate Payroll
	var group int
	switch employee_type {
	case Employee_type_production:
		group = production_personelle
	case Employee_type_marketing:
		group = marketing_personelle
	default:
		group = other_personelle
	}
	for i := range company_employees {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(fmt.Sprintf("Pay for %s employee %d", company_employees[i].Employee_type.to_string(), company_employees[i].Id), group, "", true, round(float64(-company_employees[i].Pay)/12, 2))
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(fmt.Sprintf("Bonus for %s employee %d", company_employees[i].Employee_type.to_string(), company_employees[i].Id), group, "", true, round(float64(-company_employees[i].Bonus)/12, 2))
	}

	// Calculate training
	err := c.train_employees(company_employees, 0.01)
	if err != nil {
		return err
	}

	// Calculate Motivation
	if employee_type == Employee_type_production {
		calculate_motivation(company_employees, external_factors.Production_minimum_wage, 1)
	} else if employee_type == Employee_type_marketing {
		calculate_motivation(company_employees, external_factors.Marketing_minimum_wage, 1)
	} else {
		panic("Unsupported Employee_type")
	}

	// finalise

	return nil
}

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

func (c Company) turnover(employees []*Employee, turnover_rate float32) (employees_who_left_count int) { // returns new list of employees & number employees that left
	for _, e_ptr := range employees {
		if e_ptr.Employer != Employee_employer_none {
			if rand.Float32() <= turnover_rate {
				e_ptr.Employer = Employee_employer_none
			}
		}
	}

	return employees_who_left_count
}

func (c *Company) train_employees(employees []*Employee, passive_training float32) error {
	for i := range employees {
		if employees[i].Extra_training > 0 {
			c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Extra training for employee: "+employees[i].Name, employee_training, "Training improves speed and quality of employeess' work", true, float64(employees[i].Extra_training))
		} else if employees[i].Extra_training < 0 {
			return errors.New(fmt.Sprintf("Training for %s employee %d (%s) is less than 0", employees[i].Employee_type.to_string(), employees[i].Id, employees[i].Name))
		}

		employees[i].Skill += float32(employees[i].Extra_training) / 1000.0
		employees[i].Skill += passive_training

	}

	return nil
}

func (c *Company) calculate_payroll(employees []*Employee, employee_type Employee_type) {
	var group int
	switch employee_type {
	case Employee_type_production:
		group = production_personelle
	case Employee_type_marketing:
		group = marketing_personelle
	default:
		group = other_personelle
	}

	for i := range employees {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
			fmt.Sprintf("Pay for %s employee %d", employees[i].Employee_type.to_string(), employees[i].Id),
			group,
			"",
			true,
			round(float64(-employees[i].Pay)/12, 2),
		)

		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
			fmt.Sprintf("Bonus for %s employee %d", employees[i].Employee_type.to_string(), employees[i].Id),
			group,
			"",
			true,
			round(float64(-employees[i].Bonus)/12, 2),
		)
	}
}

func calculate_motivation(employees []*Employee, minimum_wage float32, base_motivation float32) {
	// REDO MOTIVATION
	for _, a := range employees {
		pay_factor := float32(a.Pay / (minimum_wage * 1.2))
		raise_factor := float32(a.Pay / a.Pay)
		working_hours_factor := float32(a.Working_hours / 8)
		time_off_factor := float32(a.Working_hours / a.Working_hours)
		training_factor := float32(a.Extra_training / 1000)

		a.Motivation = ((base_motivation*2 +
			pay_factor +
			raise_factor*2 +
			working_hours_factor +
			time_off_factor*2 +
			training_factor) /
			9)
	}
}
