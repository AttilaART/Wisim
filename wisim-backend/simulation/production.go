package simulation

import (
	"fmt"
	"log"
	"slices"
)

func (c *Company) research(decisions Decisions) {
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Quality research", research, "", true, float64(-decisions.Research.Quality))
	c.Tech.Quality += decisions.Research.Quality / 1000
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Ecology research", research, "", true, float64(-decisions.Research.Ecology))
	c.Tech.Quality += decisions.Research.Ecology / 1000
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Durability research", research, "", true, float64(-decisions.Research.Durability))
	c.Tech.Quality += decisions.Research.Durability / 1000
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Production cost research", research, "", true, float64(-decisions.Research.Production_cost))
	c.Tech.Quality += decisions.Research.Production_cost / 1000
}

// Production functions
func (c *Company) calculateProduction(decisions Decisions, externalFactors ExternalFactors) {
	productionPersonelle := c.employeePool.Get_employees_of_company(c.ID, Employee_type_production)
	if len(productionPersonelle) == 0 {
		println("Warning: no production employees!")
	}

	productionReport := &c.Reports[len(c.Reports)-1].ProductionReport

	// Purchase Machines
	println("Purchasing machines")
	var machinesToDeleteID []int
	for _, m := range decisions.Production.Machines {
		switch m.Change {
		case Delta_New:
			productionReport.MachinesPurchased += 1
			c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Purchase of machine", production, "", true, -float64(m.Item.Value))
			c.Machines = append(c.Machines, m.Item)
		case Delta_Remove:
			productionReport.MachinesSold += 1
			c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Selling of machine", production, "", true, float64(m.Item.Value))
			machinesToDeleteID = append(machinesToDeleteID, m.Item.ID)
		case Delta_Change:
			*getMachineByID(m.Item.ID, c.Machines) = m.Item
		}
	}

	c.Machines = deleteByID(c.Machines, machinesToDeleteID...)

	calculate_machines_value(
		&c.Machines,
		&c.Reports[len(c.Reports)-1].ProductionReport,
		&c.Reports[len(c.Reports)-1].BalanceSheet,
		externalFactors,
	)
	// Calculate machine upkeep
	var machineUpkeep float32
	for _, m := range c.Machines {
		machineUpkeep += m.MaintananceCost
	}
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Machine upkeep", production, "The upkeep of out production machines", true, float64(-machineUpkeep))

	// Produce
	println("Assigning workers")
	c.Machines, productionReport.WorkerSurplus = assign_workers(
		c.employeePool,
		c.Machines,
		c.employeePool.Get_employees_of_company(
			c.ID, Employee_type_production),
	)

	println("Producing products")

	productionReport.ProductSpecificReport = make(map[string]struct {
		TotalProduction       int
		BaseProduction        int
		BonusProduction       int
		ExcessProduction      int
		TotalProductsProduced int
		BaseProductsProduced  int
		BonusProductsProduced int
	})

	produce(
		c.employeePool,
		c.Machines,
		c.Offers,
		productionReport,
		&c.Reports[len(c.Reports)-1].BalanceSheet,
		externalFactors,
	)
}

func getMachineByID(ID int, machines []Machine) *Machine {
	for i := range machines {
		if machines[i].ID == ID {
			return &machines[i]
		}
	}
	return nil
}

func calculate_machines_value(
	machines *[]Machine,
	production_report *Production_report,
	balance_sheet *Balance_sheet,
	external_factors ExternalFactors,
) {
	for i, m := range *machines {
		(*machines)[i].Value = float32(round(float64(m.Value*(1-external_factors.MachineDepreciationRate)), 2))
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
func assign_workers(employee_pool Employee_pool, machines []Machine, workers_ids []int) ([]Machine, int) {
	println("Sorting Employees")

	workers_ids = slices.SortedFunc(func(yield func(int) bool) {
		for _, id := range workers_ids {
			if !yield(id) {
				return
			}
		}
	}, func(a, b int) int {
		e_a := employee_pool[a]
		e_b := employee_pool[b]
		vala := (e_a.Motivation * e_a.Skill * e_a.WorkingHours)
		valb := (e_b.Motivation * e_b.Skill * e_b.WorkingHours)
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
		if a.ProductionCapacity < b.ProductionCapacity {
			return 1
		} else if a.ProductionCapacity == b.ProductionCapacity {
			return 0
		}
		return -1
	})

	var Worker_surplus int
	ii := 0
	for i := range machines {
		machines[i].AssignedWorkersIDs = make([]int, 0)

		for range machines[i].RequiredWorkers {
			if ii >= (len(workers_ids) - 1) {
				break
			}
			machines[i].AssignedWorkersIDs = append(machines[i].AssignedWorkersIDs, workers_ids[ii])
			ii++
		}
		Worker_surplus += len(machines[i].AssignedWorkersIDs) - machines[i].RequiredWorkers
	}

	if Worker_surplus >= 0 {
		Worker_surplus = len(workers_ids) - ii - 1
	}

	return machines, Worker_surplus
}

func produce(
	employeePool Employee_pool,
	machines []Machine,
	offers map[string]Offer,
	productionReport *Production_report,
	balanceSheet *Balance_sheet,
	externalFactors ExternalFactors,
) {
	energyUse := 0.0

	totalProduction := 0
	for _, m := range machines {
		baseProductionOfMachine, bonusProductionOfMachine := calculateMachineProduction(employeePool, m)
		productSpecificReport := productionReport.ProductSpecificReport[m.AssignedProductID]
		productSpecificReport.BaseProduction += baseProductionOfMachine
		productSpecificReport.BonusProduction += bonusProductionOfMachine
		productSpecificReport.TotalProduction += bonusProductionOfMachine + baseProductionOfMachine

		productionReport.ProductSpecificReport[m.AssignedProductID] = productSpecificReport

		totalProduction += baseProductionOfMachine + bonusProductionOfMachine
		energyUse += float64(m.EnergyUse)
	}

	if len(machines) <= 0 {
		println("Company owns no machines!")
	}

	for productID := range offers {
		productSpecificReport := productionReport.ProductSpecificReport[productID]

		productSpecificReport.BaseProductsProduced = int(float32(productSpecificReport.BaseProduction) / float32(offers[productID].ProductStats.ProductionCost))
		productSpecificReport.BonusProductsProduced = int(float32(productSpecificReport.BonusProduction) / float32(offers[productID].ProductStats.ProductionCost))
		productSpecificReport.TotalProductsProduced = productSpecificReport.BaseProductsProduced + productSpecificReport.BonusProductsProduced

		productSpecificReport.ExcessProduction = productSpecificReport.TotalProduction - (int(offers[productID].ProductStats.ProductionCost * float32(productSpecificReport.TotalProductsProduced)))

		productionReport.MaterialUsed += offers[productID].ProductStats.MaterialUse * float32(productSpecificReport.TotalProductsProduced)

		productionReport.ProductSpecificReport[productID] = productSpecificReport
		fmt.Printf("Product %s (%s) produced: \n    Base: %d\n    Bonus: %d\n    Total: %d\n",
			offers[productID].Product.Name,
			offers[productID].Product.ID,
			productSpecificReport.BaseProductsProduced,
			productSpecificReport.BonusProductsProduced,
			productSpecificReport.TotalProductsProduced)
	}

	productionReport.EnergyUsed = float32(energyUse)

	materialCosts := -round(float64(externalFactors.MaterialPrice)*float64(productionReport.MaterialUsed), 2)
	energyCosts := -round(float64(externalFactors.EnergyPrice)*float64(productionReport.EnergyUsed), 2)

	balanceSheet.add_to_income_statement("Material costs", production, "The cost of materials used in your products", true, materialCosts)
	balanceSheet.add_to_income_statement("Energy costs", production, "The cost of energy used by machines in production", true, energyCosts)

	productionReport.AvgMachineProductivity = float32(totalProduction) / float32(len(machines))
}

// return (base production, bonus production)
func calculateMachineProduction(employeePool Employee_pool, machine Machine) (int, int) {
	// calculate averages
	var skill float32 = 0
	var motivation float32 = 0
	var workingHours float32 = 0

	if machine.MinimumWorkers <= 0 {
		log.Println("machine.Minimum_workers <= 0")
		machine.MinimumWorkers = 1
	}

	if len(machine.AssignedWorkersIDs) < machine.MinimumWorkers {
		fmt.Printf("Machine has too few workers: %d instead of %d+\n", len(machine.AssignedWorkersIDs), machine.MinimumWorkers)
		return 0, 0
	} else if len(machine.AssignedWorkersIDs) > machine.RequiredWorkers {
		fmt.Printf("Machine has too many workers: %d instead of %d+\n", len(machine.AssignedWorkersIDs), machine.RequiredWorkers)
	}

	if len(machine.AssignedWorkersIDs) == 0 {
		return 0, 0
	}

	for _, ID := range machine.AssignedWorkersIDs {
		skill += employeePool[ID].Skill
		motivation += employeePool[ID].Motivation
		workingHours += employeePool[ID].WorkingHours
	}

	skill = skill / float32(len(machine.AssignedWorkersIDs))
	motivation = motivation / float32(len(machine.AssignedWorkersIDs))
	workingHours = workingHours / float32(len(machine.AssignedWorkersIDs))

	if skill <= 0 {
		panic("skill is 0 or less")
	}
	if motivation <= 0 {
		panic("motivation is 0 or less")
	}

	baseProduction := int(float32(machine.ProductionCapacity) * (workingHours / 8))
	bonusProduction := max(int(float32(baseProduction)*skill*motivation-float32(baseProduction)), 0)

	if baseProduction < 0 {
		panic(fmt.Sprintf("base_production is 0 or less (%d)", baseProduction))
	}

	fmt.Printf("------------\n")
	fmt.Printf("machine %d base_production: %d\n", machine.ID, baseProduction)
	fmt.Printf("machine %d bonus_production: %d\n", machine.ID, bonusProduction)
	fmt.Printf("machine %d skill: %f\n", machine.ID, skill)
	fmt.Printf("machine %d motivation: %f\n", machine.ID, motivation)
	fmt.Printf("machine %d working_hours: %f\n", machine.ID, workingHours)
	fmt.Printf("Workers %d assigned: %d\n", machine.ID, len(machine.AssignedWorkersIDs))

	return baseProduction, bonusProduction
}
