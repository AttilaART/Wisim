package simulation

import (
	"errors"
	"fmt"
	"math/rand"
)

// Employee functions

func (c *Company) simulate_employees(external_factors External_factors, severance_pay float32, employee_type Employee_type) error {
	// IMPORTANT: Calculation of employee deltas now happens before individual companies

	// Get correct employee slice
	company_employees_ids := c.employee_pool.Get_employees_of_company(c.Id, employee_type)

	/*
		// Fire/"layoff" employees
		var employees_layed_off int
		for i := range employee_deltas_ptr {
			if employee_deltas_ptr[i].Change == Delta_Remove {
				employees_layed_off += 1
				employee_deltas_ptr[i].Item.Employer = Employee_employer_none
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
	*/

	// Calcualte Turnover
	var num_employees_who_left int
	for _, id := range company_employees_ids {
		if c.employee_pool[id].Employer != Employee_employer_none {
			if rand.Float32() <= external_factors.Turnover {
				e := c.employee_pool[id]
				e.Employer = Employee_employer_none
				c.employee_pool[id] = e
				num_employees_who_left += 1
			}
		}
	}

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
	for _, id := range company_employees_ids {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(fmt.Sprintf("Pay for %s employee %d", c.employee_pool[id].Employee_type.to_string(), c.employee_pool[id].Id), group, "", true, round(float64(-c.employee_pool[id].Pay)/12, 2))
		if c.employee_pool[id].Bonus > 0 {
			c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(fmt.Sprintf("Bonus for %s employee %d", c.employee_pool[id].Employee_type.to_string(), c.employee_pool[id].Id), group, "", true, round(float64(-c.employee_pool[id].Bonus)/12, 2))
		}
	}

	// Calculate training
	err := c.train_employees(company_employees_ids, 0.01)
	if err != nil {
		return err
	}

	// Calculate Motivation
	if employee_type == Employee_type_production {
		calculate_motivation(c.employee_pool, company_employees_ids, external_factors.Production_minimum_wage, 1)
	} else if employee_type == Employee_type_marketing {
		calculate_motivation(c.employee_pool, company_employees_ids, external_factors.Marketing_minimum_wage, 1)
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

func (g *Game_state) handleEmployeeDeltas() {
	// Handle Changes & fires
	for _, decisions := range g.Current_decisions {
		for _, e := range decisions.Employees.Marketing_deltas {
			if e.Change == Delta_Change {
				g.Employees[e.Item.Id] = &e.Item
			} else if e.Change == Delta_Remove {
				g.Employees[e.Item.Id].Employer = Employee_employer_none
			}
		}
	}

	// deal with new hires
	newHires := make(map[int][]int)

	for i, decisions := range g.Current_decisions {
		for _, delta := range decisions.Employees.Marketing_deltas {
			if delta.Change == Delta_New {
				if _, exists := newHires[delta.Item.Id]; exists {
					hire := newHires[delta.Item.Id]
					hire = append(hire, i)
					newHires[delta.Item.Id] = hire
				} else {
					newHires[delta.Item.Id] = []int{i}
				}
			}
		}
		for _, delta := range decisions.Employees.Production_deltas {
			if delta.Change == Delta_New {
				if _, exists := newHires[delta.Item.Id]; exists {
					hire := newHires[delta.Item.Id]
					hire = append(hire, i)
					newHires[delta.Item.Id] = hire
				} else {
					newHires[delta.Item.Id] = []int{i}
				}
			}
		}
	}

	for hire := range newHires {
		potentialEmployers := newHires[hire]
		if len(potentialEmployers) == 1 {
			g.Employees[hire].Employer = potentialEmployers[0]
		} else {
			g.Employees[hire].Employer = potentialEmployers[rand.Int()%len(potentialEmployers)]
		}
	}

	/*
		for _, e := range g.Employees {
			println("Employee: ", e.Id, "Employer: ", e.Employer)
		}*/
}

func (c *Company) turnover(employees []*Employee, turnover_rate float32) (employees_who_left_count int) { // returns new list of employees & number employees that left
	for _, e_ptr := range employees {
		if e_ptr.Employer != Employee_employer_none {
			if rand.Float32() <= turnover_rate {
				e_ptr.Employer = Employee_employer_none
			}
		}
	}

	return employees_who_left_count
}

func (c *Company) train_employees(employees_ids []int, passive_training float32) error {
	for _, id := range employees_ids {
		if c.employee_pool[id].Extra_training > 0 {
			c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Extra training for employee: "+c.employee_pool[id].Name, employee_training, "Training improves speed and quality of employeess' work", true, float64(c.employee_pool[id].Extra_training))
		} else if c.employee_pool[id].Extra_training < 0 {
			return errors.New(fmt.Sprintf("Training for %s employee %d (%s) is less than 0", c.employee_pool[id].Employee_type.to_string(), c.employee_pool[id].Id, c.employee_pool[id].Name))
		}

		e := c.employee_pool[id]
		e.Skill += float32(c.employee_pool[id].Extra_training) / 1000.0
		e.Skill += passive_training
		c.employee_pool[id] = e

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

func calculate_motivation(employee_pool Employee_pool, employees_ids []int, minimum_wage float32, base_motivation float32) {
	// REDO MOTIVATION
	for _, id := range employees_ids {
		pay_factor := float32(employee_pool[id].Pay / (minimum_wage * 1.2))
		raise_factor := float32(employee_pool[id].Pay / employee_pool[id].Pay) // TODO: Fix!
		working_hours_factor := float32(employee_pool[id].Working_hours / 8)
		time_off_factor := float32(employee_pool[id].Working_hours / employee_pool[id].Working_hours)
		training_factor := float32(employee_pool[id].Extra_training / 1000)

		employee_pool[id].Motivation = ((base_motivation*2 +
			pay_factor +
			raise_factor*2 +
			working_hours_factor +
			time_off_factor*2 +
			training_factor) /
			9)
	}
}
