package simulation

import (
	"fmt"
	"slices"
)

func (c *Company) simulate_company(decisions Decisions, external_factors External_factors) error {
	c.Reports = append(c.Reports, Report{Month: external_factors.Month})
	c.Decision_history = append(c.Decision_history, decisions)

	// Personelle
	println("Simulatig personelle")

	err := c.simulate_employees(decisions.Employees.Production_deltas, external_factors, decisions.Employees.Severance_pay, Employee_type_production)
	if err != nil {
		return err
	}
	err = c.simulate_employees(decisions.Employees.Marketing_deltas, external_factors, decisions.Employees.Severance_pay, Employee_type_marketing)
	if err != nil {
		return err
	}

	// Offer
	println("Calculating product stats")
	c.Offer.Price = decisions.Marketing.Price

	c.Base_marketing_strength += decisions.Research.Promotion / 1000 * c.Base_marketing_strength
	c.Offer.Promotion_quality = promotion_quality(c.Base_marketing_strength, c.employee_pool.Get_employees_of_company(c.Id, Employee_type_marketing))

	c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
		"Advertisement costs",
		marketing,
		"Cost of your ads (equals promotion quantity)",
		true,
		float64(-decisions.Marketing.Promotion.Quantity),
	)
	// Offer: Product

	c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Quality research", research, "", true, float64(-decisions.Research.Quality))
	c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Ecology research", research, "", true, float64(-decisions.Research.Ecology))
	c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Durability research", research, "", true, float64(-decisions.Research.Durability))

	// calculate avg skill of production personelle -> influences Quality_factor
	production_personelle := c.employee_pool.Get_employees_of_company(c.Id, Employee_type_production)
	var total_production_skill float32 = 0
	for _, e := range production_personelle {
		total_production_skill += e.Skill
	}
	avg_production_skill := total_production_skill / float32(len(production_personelle))

	c.Offer.Product.calculate_quality(
		avg_production_skill,
		decisions.Marketing.Product.Materials.Quality,
		decisions.Marketing.Product.Manufacturing.Quality,
	)

	c.Offer.Product.Base_ecology *= 1 + decisions.Research.Ecology/1000

	c.Offer.Product.calculate_ecology(
		c.Offer.Product.Base_ecology,
		c.Offer.Product.Material_use,
		decisions.Marketing.Product.Manufacturing.Ecological_energy,
	)

	c.Offer.Product.Base_durability *= 1 + decisions.Research.Durability/1000
	c.Offer.Product.calculate_durability(decisions.Marketing.Product.Manufacturing.Max_durability, decisions.Marketing.Product.Manufacturing.Durability)

	// Production
	println("Calculating production")

	c.Offer.Product.Base_production_speed *= 1 + decisions.Research.Speed/1000

	if c.Offer.Product.Base_production_speed == 0 {
		panic("Base_production_speed is 0")
	}

	c.calculate_production(decisions, external_factors, calculate_production_speed(
		decisions.Marketing.Product.Manufacturing,
		float32(c.Offer.Product.Base_production_speed),
	))

	// Logistics
	c.Items_in_storage += c.Reports[len(c.Reports)-1].Production_report.Base_production +
		c.Reports[len(c.Reports)-1].Production_report.Bonus_production

	println("Calculating logistics")
	c.calculate_logistics(decisions)

	// Finances
	return nil
}

// Production functions
func (c *Company) calculate_production(decisions Decisions, external_factors External_factors, production_speed float32) {
	production_personelle := c.employee_pool.Get_employees_of_company(c.Id, Employee_type_production)
	if len(production_personelle) == 0 {
		println("Warning: no production employees!")
	}

	production_report := &c.Reports[len(c.Reports)-1].Production_report

	// Purchase Machines
	println("Purchasing machines")
	var machines_to_delete_id []int
	for _, m := range decisions.Production.Machines {
		if m.Change == Delta_New {
			production_report.Machines_purchased += 1
			c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Purchase of machine", production, "", true, -float64(m.Item.Value))
			c.Machines = append(c.Machines, m.Item)
		} else if m.Change == Delta_Remove {
			production_report.Machines_sold += 1
			c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Selling of machine", production, "", true, float64(m.Item.Value))
			machines_to_delete_id = append(machines_to_delete_id, m.Item.Id)
		}
	}

	c.Machines = delete_by_id(c.Machines, machines_to_delete_id...)

	// fix pointer stuff (json doesn't transmit pointers)
	for i := range c.Machines {
		for ii, e := range c.Machines[i].Assigned_workers_ptr {

			c.Machines[i].Assigned_workers_ptr[ii] = c.employee_pool.Find_employee_by_id(e.Id)
			if c.Machines[i].Assigned_workers_ptr[ii] == nil {
				panic("Invalid employee assigned to machine")
			}
		}
	}

	calculate_machines_value(
		&c.Machines,
		&c.Reports[len(c.Reports)-1].Production_report,
		&c.Reports[len(c.Reports)-1].Balance_sheet,
		external_factors,
	)

	// Produce
	println("Assigning workers")
	c.Machines, production_report.Worker_surplus = assign_workers(c.Machines, c.employee_pool.Get_employees_of_company(c.Id, Employee_type_production))

	println("Producing products")
	produce(
		c.Machines,
		c.Offer.Product,
		production_speed,
		production_report,
		&c.Reports[len(c.Reports)-1].Balance_sheet,
		external_factors,
		c.employee_pool.Get_employees_of_company(c.Id, Employee_type_production),
	)
}

func calculate_machines_value(
	machines *[]Machine,
	production_report *Production_report,
	balance_sheet *Balance_sheet,
	external_factors External_factors,
) {
	for i, m := range *machines {
		(*machines)[i].Value = float32(round(float64(m.Value*(1-external_factors.Machine_depreciation_rate)), 2))
		balance_sheet.add_to_income_statement(
			"Machine depreciation",
			write_off,
			"",
			false,
			round(-float64(m.Value-(*machines)[i].Value), 2),
		)
	}

	for _, m := range *machines {
		balance_sheet.add_to_equity(
			"Machine",
			production,
			"Machines, with help of employees, produce your product",
			false,
			float64(m.Value),
		)
	}
}

// returns machines & number of unassigned wokers (if not enough workers for the machines, it return a negative int)
func assign_workers(machines []Machine, workers []*Employee) ([]Machine, int) {
	println("Sorting Employees")

	workers = slices.SortedFunc(func(yield func(*Employee) bool) {
		for _, w := range workers {
			if !yield(w) {
				return
			}
		}
	}, func(a, b *Employee) int {
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
		machines[i].Assigned_workers_ptr = make([]*Employee, 0)

		for range machines[i].Required_workers {
			if ii >= (len(workers) - 1) {
				break
			}
			machines[i].Assigned_workers_ptr = append(machines[i].Assigned_workers_ptr, workers[ii])
			ii++
		}
		Worker_surplus += len(machines[i].Assigned_workers_ptr) - machines[i].Required_workers
	}

	if Worker_surplus >= 0 {
		Worker_surplus = len(workers) - ii - 1
	}

	return machines, Worker_surplus
}

func produce(
	machines []Machine,
	product Product,
	production_speed float32,
	production_report *Production_report,
	balance_sheet *Balance_sheet,
	external_factors External_factors,
	employees []*Employee,
) {
	base_production := 0
	bonus_production := 0

	energy_use := 0.0
	for _, m := range machines {
		base_prod_of_machine, bonus_prod_of_machine := calculate_machine_production(m, production_speed, &employees)
		base_production += base_prod_of_machine
		bonus_production += bonus_prod_of_machine

		energy_use += float64(m.Energy_use)
	}

	production_report.Products_produced = base_production + bonus_production
	production_report.Base_production = base_production
	production_report.Bonus_production = bonus_production

	production_report.Material_used = product.Material_use * float32(production_report.Products_produced)
	production_report.Energy_used = float32(energy_use)

	material_costs := -round(float64(external_factors.Material_price)*float64(production_report.Material_used), 2)
	energy_costs := -round(float64(external_factors.Energy_price)*float64(production_report.Energy_used), 2)

	balance_sheet.add_to_income_statement("Material costs", production, "The cost of materials used in your products", true, material_costs)
	balance_sheet.add_to_income_statement("Energy costs", production, "The cost of energy used by machines in production", true, energy_costs)

	production_report.Avg_machine_productivity = float32(production_report.Products_produced) / float32(len(machines))
}

// return (base production, bonus production)
func calculate_machine_production(machine Machine, production_speed float32, employees *[]*Employee) (int, int) {
	// calculate averages
	var skill float32 = 0
	var motivation float32 = 0
	var working_hours float32 = 0

	if machine.Minimum_workers <= 0 {
		panic("machine.Minimum_workers <= 0")
	}

	if len(machine.Assigned_workers_ptr) < machine.Minimum_workers {
		return 0, 0
	}

	for _, employee_ptr := range machine.Assigned_workers_ptr {
		skill += employee_ptr.Skill
		motivation += employee_ptr.Motivation
		working_hours += employee_ptr.Working_hours

	}

	num_assigned_workers := len(machine.Assigned_workers_ptr)

	if num_assigned_workers > 0 {
		return 0, 0
	}

	skill = skill / float32(num_assigned_workers)
	motivation = motivation / float32(num_assigned_workers)
	working_hours = working_hours / float32(num_assigned_workers)

	if skill <= 0 {
		panic("skill is 0 or less")
	}
	if motivation <= 0 {
		panic("motivation is 0 or less")
	}

	base_production := int(float32(machine.Production_capacity) * production_speed * (working_hours / 8))
	bonus_production := int(float32(base_production)*skill*motivation - float32(base_production))
	if bonus_production < 0 {
		bonus_production = 0
	}

	if base_production <= 0 {
		panic(fmt.Sprintf("base_production is 0 or less (%d)", base_production))
	}

	fmt.Printf("------------\n")
	fmt.Printf("machine base_production: %d\n", base_production)
	fmt.Printf("machine bonus_production: %d\n", bonus_production)
	fmt.Printf("machine skill: %f\n", skill)
	fmt.Printf("machine motivation: %f\n", skill)
	fmt.Printf("machine working_hours: %f\n", skill)
	fmt.Printf("Workers assigned: %d\n", len(machine.Assigned_workers_ptr))

	return base_production, bonus_production
}
