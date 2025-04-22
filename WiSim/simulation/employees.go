package simulation

import (
	"errors"
	"fmt"
	"math/rand"
)

// Employee functions

func (company *Company) simulate_employees(decisions Decisions, external_factors External_factors) error {
	// Calcualte Turnover
	company.Production_personelle, _ = turnover(company.Production_personelle, external_factors.Turnover)
	company.Marketing_personelle, _ = turnover(company.Marketing_personelle, external_factors.Turnover)

	// Fire/"layoff" employees
	var base_motivation float32 = 1

	var production_employees_layed_off int
	company.Production_personelle, production_employees_layed_off = layoff(company.Production_personelle, decisions.New_hires_in_production)
	company.Production_personelle = hire(company.Production_personelle, decisions.New_hires_in_production, decisions.Base_pay_for_production, decisions.Working_hours_for_production, production_employee, base_motivation)

	var marketing_employees_layed_off int
	company.Marketing_personelle, marketing_employees_layed_off = layoff(company.Marketing_personelle, decisions.New_hires_in_marketing)
	company.Marketing_personelle = hire(company.Marketing_personelle, decisions.New_hires_in_marketing, decisions.Base_pay_for_marketing, decisions.Working_hours_for_marketing, marketing_employee, base_motivation)
	// calculate severance pay
	for range production_employees_layed_off {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Severance pay for a production employee", severance, "When you layoff an employee, you have to pay them severance. Sometimes it's more expensive to fire someone, than just letting them be idle.", true, -float64(decisions.Severance_for_production_personelle))
	}
	for range marketing_employees_layed_off {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Severance pay for a marketing employee", severance, "When you layoff an employee, you have to pay them severance. Sometimes it's more expensive to fire someone, than just letting them be idle.", true, -float64(decisions.Severance_for_marketing_personelle))
	}

	// Calcualate Payroll
	var production_payroll []FinanceReportEntry
	var marketing_payroll []FinanceReportEntry
	company.Production_personelle, production_payroll = calculate_payroll(company.Production_personelle, decisions.Raise_for_production_personelle)
	company.Marketing_personelle, marketing_payroll = calculate_payroll(company.Marketing_personelle, decisions.Raise_for_marketing_personelle)

	company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement = append(company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement, production_payroll...)
	company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement = append(company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement, marketing_payroll...)

	// copy employees for calculate_motivation function
	old_production_personelle := make([]Employee, len(company.Production_personelle))
	old_marketing_personelle := make([]Employee, len(company.Marketing_personelle))
	copy(old_production_personelle, company.Production_personelle)
	copy(old_marketing_personelle, company.Marketing_personelle)

	var err error
	company.Production_personelle, err = set_working_hours(company.Production_personelle, decisions.Working_hours_for_production)
	if err != nil {
		return err
	}

	company.Marketing_personelle, err = set_working_hours(company.Marketing_personelle, decisions.Working_hours_for_marketing)
	if err != nil {
		return err
	}

	// Calculate training
	company.Production_personelle = train_employees(company.Production_personelle, decisions.Investment_in_production_training, 0.1)
	company.Marketing_personelle = train_employees(company.Marketing_personelle, decisions.Investment_in_marketing_training, 0.1)

	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Production personelle training", employee_training, "Training improves the skill of your employees", true, -float64(decisions.Investment_in_production_training))
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Marketing personelle training", employee_training, "Training improves the skill of your employees", true, -float64(decisions.Investment_in_marketing_training))

	// Calculate Motivation
	company.Production_personelle, err = calculate_motivation(company.Production_personelle, old_production_personelle, external_factors.Production_minimum_wage, 1.0)
	if err != nil {
		return err
	}
	company.Marketing_personelle, err = calculate_motivation(company.Marketing_personelle, old_marketing_personelle, external_factors.Marketing_minimum_wage, 1.0)
	if err != nil {
		return err
	}

	return nil
}

func hire(employees []Employee, new_hires int, base_pay float32, working_hours float32, employee_type int, base_motivation float32) []Employee {
	if new_hires <= 0 {
		return employees
	}

	for range new_hires {
		employees = append(employees, generate_employee(generate_new_employee_id(&employees), base_pay, working_hours, employee_type, base_motivation))
	}

	return employees
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

func turnover(employees []Employee, turnover_rate float32) ([]Employee, int) { // returns new list of employees & number employees that left
	num_of_employees_leaving := int(round(float64(turnover_rate)*float64(len(employees)), 0))
	var employees_leaving_index []int

	for range num_of_employees_leaving {
		employees_leaving_index = append(employees_leaving_index, rand.Intn(len(employees)-1))
	}

	if len(employees_leaving_index) > len(employees) {
		employees_leaving_index = employees_leaving_index[0 : len(employees)-1]
	}

	employees = delete_by_index(employees, employees_leaving_index...)

	return employees, len(employees_leaving_index)
}

func train_employees(employees []Employee, investment float32, passive_training float32) []Employee {
	investment_per_employee := investment / float32(len(employees))
	for i, e := range employees {
		employees[i].Skill += e.Skill * investment_per_employee / 1000.0
		employees[i].Skill += passive_training
	}
	return employees
}

func calculate_payroll(employees []Employee, pay_raise float32) ([]Employee, []FinanceReportEntry) {
	employees_pay := make([]FinanceReportEntry, len(employees))
	for i := range employees {
		employees[i].Salary += employees[i].Salary * pay_raise
		employees_pay[i] = FinanceReportEntry{
			fmt.Sprintf("Pay for %s employee %d", employees[i].get_type(), employees[i].Id),
			personelle,
			"",
			true,
			round(float64(-employees[i].Salary)/12, 2),
		}
	}
	return employees, employees_pay
}

func set_working_hours(employees []Employee, working_hours float32) ([]Employee, error) {
	if working_hours <= 0 {
		return employees, errors.New("working hours must be > 0h")
	}

	for i := range employees {
		employees[i].Working_hours = working_hours
	}

	return employees, nil
}

func calculate_motivation(employees_after_calculations []Employee, employees_before_calculations []Employee, minimum_wage float32, base_motivation float32) ([]Employee, error) {
	if len(employees_after_calculations) != len(employees_before_calculations) {
		return nil, errors.New("employees_after_calculations not the same length as employees_before_calculations")
	}

	// REDO MOTIVATION
	for i, e := range employees_after_calculations {
		pay_factor := e.Salary / (minimum_wage * 1.2)
		raise_factor := e.Salary / employees_before_calculations[i].Salary
		working_hours_factor := e.Working_hours / 8
		time_off_factor := e.Working_hours / employees_before_calculations[i].Working_hours
		training_factor := e.Skill / employees_before_calculations[i].Skill

		employees_after_calculations[i].Motivation = ((base_motivation*2 +
			pay_factor +
			raise_factor*2 +
			working_hours_factor +
			time_off_factor*2 +
			training_factor) /
			9)
	}
	return employees_after_calculations, nil
}
