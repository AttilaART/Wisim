package main

import (
	"errors"
	"log"
	"slices"
)

func (company *Company) simulate_company(decisions Decisions, external_factors External_factors) error {
	company.Reports = append(company.Reports, Report{Month: external_factors.Month})

	// Personelle
	println("Simulatig personelle")

	err := company.simulate_employees(decisions, external_factors)
	if err != nil {
		return err
	}

	// Offer
	println("Calculating product stats")
	company.Offer.Marketing_strength = marketing_strength(decisions.Marketing, company.Marketing_personelle)

	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
		"Marketing investment",
		marketing,
		"Basic investment in marketing, the quantity of marketing",
		true,
		float64(-decisions.Marketing),
	)
	// Offer: Product

	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Quality development investment", product_development, "", true, float64(-decisions.Investment_in_quality_development))
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Ecology development investment", product_development, "", true, float64(-decisions.Investment_in_ecological_production))
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Durability development investment", product_development, "", true, float64(-decisions.Investment_in_durability_development))

	// calculate avg skill of production personelle -> influences Quality_factor
	var total_production_skill float32 = 0
	for _, e := range company.Production_personelle {
		total_production_skill += e.Skill
	}
	avg_production_skill := total_production_skill / float32(len(company.Production_personelle))

	company.Offer.Product.calculate_quality(
		decisions.Investment_in_quality_development,
		avg_production_skill,
		decisions.Material_quality,
	)

	company.Offer.Product.calculate_ecology(
		1,
		company.Offer.Product.Material_use,
		decisions.Percentage_of_ecological_energy,
	)

	// Production
	println("Calculating production")
	company.calculate_production(decisions, external_factors, &company.Reports[len(company.Reports)-1].Balance_sheet)

	// Logistics
	company.Items_in_storage += company.Reports[len(company.Reports)-1].Production_report.Base_production + company.Reports[len(company.Reports)-1].Production_report.Bonus_production

	println("Calculating logistics")
	Balance_entries, Assets_entries, err := company.calculate_logistics(decisions)
	if err != nil {
		return err
	}

	company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement = append(company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement, Balance_entries...)
	company.Reports[len(company.Reports)-1].Balance_sheet.Assets = append(company.Reports[len(company.Reports)-1].Balance_sheet.Assets, Assets_entries...)
	// Finances

	company.Decision_history = append(company.Decision_history, decisions)

	return nil
}

// Production functions
func (company *Company) calculate_production(decisions Decisions, external_factors External_factors, finance_report *Balance_sheet) {
	var production_report Production_report
	// Purchase Machines
	println("Purchasing machines")
	Machines, production_report, Balance_entries, Assets_entries, err := calculate_machines(
		company.Machines,
		decisions.Purchase_of_machines,
		decisions.Purchase_of_machines,
		external_factors.Machine_on_offer,
		production_report,
		external_factors,
	)
	if err != nil {
		println(err.Error())
		log.Fatal("purchase_machines has errored, fix this properly, Attila!")
	}

	finance_report.Income_statement = append(finance_report.Income_statement, Balance_entries...)
	finance_report.Assets = append(finance_report.Assets, Assets_entries...)
	company.Machines = Machines // doing this to avoid compiler error

	// Produce
	println("Assigning workers")
	company.Machines, production_report.Worker_surplus = assign_workers(company.Machines, company.Production_personelle)

	println("Producing products")
	production_report, Income_entries := produce(company.Machines, company.Offer.Product, production_report, external_factors, company.Production_personelle)
	company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement = append(company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement, Income_entries...)

	company.Reports[len(company.Reports)-1].Production_report = production_report
}

func calculate_machines(
	machines []Machine,
	machines_to_purchase int,
	machines_to_sell int,
	machines_for_sale Machine,
	production_report Production_report,
	external_factors External_factors,
) ([]Machine, Production_report, []FinanceReportEntry, []FinanceReportEntry, error) {
	production_report.Machines_purchased = []Machine{}
	production_report.Machines_sold = []Machine{}
	var Income_entries []FinanceReportEntry
	var Assets_entries []FinanceReportEntry

	for i, m := range machines {
		Income_entries = append(Income_entries, FinanceReportEntry{"Machine depreciation", production, "", false, round(-float64(m.Value)*float64(external_factors.Machine_depreciation_rate), 2)})
		machines[i].Value = float32(round(float64(m.Value-m.Value*external_factors.Machine_depreciation_rate), 2))
	}

	if machines_to_sell > 0 {
		for range machines_to_sell {
			production_report.Machines_purchased = append(production_report.Machines_purchased, machines[0])
			Income_entries = append(Income_entries, FinanceReportEntry{"Selling of machine", production, "", true, float64(machines[0].Value)})
			machines = machines[1:]
		}
	} else if machines_to_sell < 0 {
		err := errors.New("machines_to_sell is less than 0")
		return nil, production_report, Income_entries, Assets_entries, err
	}

	if machines_to_purchase > 0 {
		for range machines_to_purchase {
			production_report.Machines_purchased = append(production_report.Machines_purchased, machines_for_sale)
			Income_entries = append(Income_entries, FinanceReportEntry{"Purchase of machine", production, "", true, float64(-machines_for_sale.Value)})
			machines = append(machines, machines_for_sale)
		}
	} else if machines_to_purchase < 0 {
		err := errors.New("machines_to_purchase is less than 0")
		return nil, production_report, Income_entries, Assets_entries, err
	}

	for _, m := range machines {
		Assets_entries = append(Assets_entries, FinanceReportEntry{"Machine", production, "Machines, with help of employees, produce your product", false, float64(m.Value)})
	}

	return machines, production_report, Income_entries, Assets_entries, nil
}

// returns machines & number of unassigned wokers (if not enough workers for the machines, it return a negative int)
func assign_workers(machines []Machine, workers []Employee) ([]Machine, int) {
	println("Sorting Employees")

	workers = slices.SortedFunc(func(yield func(Employee) bool) {
		for _, w := range workers {
			if !yield(w) {
				return
			}
		}
	}, func(a, b Employee) int {
		vala := (a.Motivation * a.Skill * a.Working_hours)
		valb := (b.Motivation * b.Skill * a.Working_hours)
		if vala < valb {
			return 1
		} else if vala == valb {
			return 0
		}
		return -1
	})

	machines = slices.SortedFunc(func(yield func(Machine) bool) {
		for _, m := range machines {
			if !yield(m) {
				return
			}
		}
	}, func(a, b Machine) int {
		if a.Production_capacity < b.Production_capacity {
			return 1
		} else if a.Production_capacity == b.Production_capacity {
			return 0
		}
		return -1
	})

	var Worker_surplus int
	ii := 0
	for i := range machines {
		machines[i].Assigned_workers_ids = make([]int, 0)

		for range machines[i].Required_workers {
			if ii >= (len(workers) - 1) {
				break
			}
			machines[i].Assigned_workers_ids = append(machines[i].Assigned_workers_ids, workers[ii].Id)
			ii++
		}
		Worker_surplus += len(machines[i].Assigned_workers_ids) - machines[i].Required_workers
	}

	if Worker_surplus >= 0 {
		Worker_surplus = len(workers) - ii - 1
	}

	return machines, Worker_surplus
}

func produce(machines []Machine, product Product, production_report Production_report, external_factors External_factors, employees []Employee) (Production_report, []FinanceReportEntry) {
	Income_entries := make([]FinanceReportEntry, 2)
	base_production := 0
	bonus_production := 0

	energy_use := 0.0
	for _, m := range machines {
		base_prod_of_machine, bonus_prod_of_machine := production_of_machine(m, employees)
		base_production += base_prod_of_machine
		bonus_production += bonus_prod_of_machine

		energy_use += float64(m.Energy_use)
	}

	production_report.Products_produced = base_production + bonus_production
	production_report.Base_production = base_production
	production_report.Bonus_production = bonus_production

	production_report.Material_used = product.Material_use * float32(production_report.Products_produced)
	production_report.Energy_used = float32(energy_use)
	Income_entries[0] = FinanceReportEntry{"Material costs", 2, "The cost of materials used in your products", true, -round(float64(external_factors.Material_price)*float64(production_report.Material_used), 2)}
	Income_entries[1] = FinanceReportEntry{"Energy costs", 2, "The cost of energy used by machines in production", true, -round(float64(external_factors.Energy_price)*float64(production_report.Energy_used), 2)}

	production_report.Avg_machine_productivity = float32(production_report.Products_produced) / float32(len(machines))

	return production_report, Income_entries
}

// return (base production, bonus production)
func production_of_machine(machine Machine, employees []Employee) (int, int) {
	// calculate averages
	var skill float32 = 0
	var motivation float32 = 0
	var working_hours float32 = 0

	for _, id := range machine.Assigned_workers_ids {
		e, err := find_employee_by_id(id, employees)
		if err != nil {
			log.Fatal(err.Error())
		}

		skill += e.Skill
		motivation += e.Motivation
		working_hours += e.Working_hours
	}

	num_assigned_workers := len(machine.Assigned_workers_ids)
	skill = skill / float32(num_assigned_workers)
	motivation = motivation / float32(num_assigned_workers)
	working_hours = working_hours / float32(num_assigned_workers)

	base_production := int(float32(machine.Production_capacity) * (working_hours / 8))
	bonus_production := int(float32(base_production)*skill*motivation) - base_production

	return base_production, bonus_production
}
