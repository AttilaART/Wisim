package simulation

import (
	"errors"
	"fmt"
	"math/rand"
)

// Employee functions

func (c *Company) simulate_employees(employee_actions []Employee_action, external_factors External_factors, severance_pay float32) error {
	if len(employee_actions) <= 0 {
		fmt.Printf("Warning: Company %d has 0 employee actions\n", c.Id)
		return nil
	}

	// fix pointer stuff (cus json)
	for i := range employee_actions {
		var err error
		employee_actions[i], err = c.Link_employees_to_action(employee_actions[i])
		if err != nil {
			panic(err)
		}
	}

	// Calcualte Turnover
	_ = turnover(&employee_actions, external_factors.Turnover)

	// Fire/"layoff" employees
	var employees_layed_off int
	for i := range employee_actions {
		if employee_actions[i].Status == Layed_off {
			employees_layed_off += 1
		}
	}

	// calculate severance pay

	for range employees_layed_off {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Severance pay for employee", severance, "When you layoff an employee, you have to pay them severance. Sometimes it's more expensive to fire someone, than just letting them be idle.", true, -float64(severance))
	}

	// remove leaving employees
	var index_to_delete []int
	for i := range employee_actions {
		s := employee_actions[i].Status
		switch s {
		case Quit, Layed_off, Fired:
			index_to_delete = append(index_to_delete, i)
		}
	}

	employee_actions = delete_by_index(employee_actions, index_to_delete...)

	// Calcualate Payroll
	c.calculate_payroll(employee_actions)

	// Calculate training
	err := c.train_employees(employee_actions, 0.01)
	if err != nil {
		return err
	}

	// Calculate Motivation
	if employee_actions[0].employee.Employee_type == Production_employee {
		calculate_motivation(employee_actions, external_factors.Production_minimum_wage, 1)
	} else if employee_actions[0].employee.Employee_type == Marketing_employee {
		calculate_motivation(employee_actions, external_factors.Production_minimum_wage, 1)
	} else {
		calculate_motivation(employee_actions, external_factors.Marketing_minimum_wage, 1)
	}

	// finalise

	// Update existing employee arrays
	// (Add employee pointer to array)
	var p_employees []*Employee
	for i := range employee_actions {
		p_employees = append(p_employees, employee_actions[i].employee)
	}

	if employee_actions[0].employee.Employee_type == Production_employee {
		c.Production_personelle = p_employees
	} else if employee_actions[0].employee.Employee_type == Marketing_employee {
		c.Marketing_personelle = p_employees
	} else {
		panic(fmt.Sprintf("Unknown employee type %d (%s)", employee_actions[0].employee.Employee_type, employee_actions[0].employee.get_type()))
	}

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

func turnover(actions *[]Employee_action, turnover_rate float32) (employees_who_left_count int) { // returns new list of employees & number employees that left
	num_of_employees_leaving := int(round(float64(turnover_rate)*float64(len(*actions)), 0))

	for range num_of_employees_leaving {
		if employees_who_left_count >= len(*actions) {
			break
		}

		employee_leaving_index := rand.Intn(len(*actions) - 1)

		for (*actions)[employee_leaving_index].Status != Quit {
			(*actions)[employee_leaving_index].Status = Quit
			employee_leaving_index = rand.Intn(len(*actions) - 1)
		}

		employees_who_left_count += 1
	}

	return employees_who_left_count
}

func (c *Company) train_employees(employee_actions []Employee_action, passive_training float32) error {
	for i := range employee_actions {
		if employee_actions[i].Extra_training > 0 {
			c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Extra training for employee: "+employee_actions[i].employee.Name, employee_training, "Training improves speed and quality of employeess' work", true, float64(employee_actions[i].Extra_training))
		} else if employee_actions[i].Extra_training < 0 {
			return errors.New(fmt.Sprintf("Training for %s employee %d (%s) is less than 0", employee_actions[i].employee.get_type(), employee_actions[i].employee.Id, employee_actions[i].employee.Name))
		}

		employee_actions[i].employee.Skill += float32(employee_actions[i].Extra_training) / 1000.0
		employee_actions[i].employee.Skill += passive_training

	}

	return nil
}

func (c *Company) calculate_payroll(employees_actions []Employee_action) {
	for i := range employees_actions {
		var group int
		switch employees_actions[i].employee.Employee_type {
		case Production_employee:
			group = production_personelle
		case Marketing_employee:
			group = marketing_personelle
		default:
			group = other_personelle
		}

		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
			fmt.Sprintf("Pay for %s employee %d", employees_actions[i].employee.get_type(), employees_actions[i].employee.Id),
			group,
			"",
			true,
			round(float64(-employees_actions[i].Pay)/12, 2),
		)

		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
			fmt.Sprintf("Bonus for %s employee %d", employees_actions[i].employee.get_type(), employees_actions[i].employee.Id),
			group,
			"",
			true,
			round(float64(-employees_actions[i].Bonus)/12, 2),
		)

		employees_actions[i].employee.Pay = employees_actions[i].Pay
		employees_actions[i].employee.Bonus = employees_actions[i].Bonus
	}
}

func calculate_motivation(employee_actions []Employee_action, minimum_wage float32, base_motivation float32) {
	// REDO MOTIVATION
	for _, a := range employee_actions {
		pay_factor := float32(a.Pay / (minimum_wage * 1.2))
		raise_factor := float32(a.Pay / a.employee.Pay)
		working_hours_factor := float32(a.Working_hours / 8)
		time_off_factor := float32(a.Working_hours / a.employee.Working_hours)
		training_factor := float32(a.Extra_training / 1000)

		a.employee.Motivation = ((base_motivation*2 +
			pay_factor +
			raise_factor*2 +
			working_hours_factor +
			time_off_factor*2 +
			training_factor) /
			9)
	}
}
