package simulation

import (
	"fmt"
	"log"
	"slices"
)

func (c *Company) research(decisions Decisions) {
	if decisions.Research.Quality > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Quality research", research, "", true, float64(-decisions.Research.Quality))
		c.Tech.Quality += decisions.Research.Quality / 1000000
	}
	if decisions.Research.Ecology > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Ecology research", research, "", true, float64(-decisions.Research.Ecology))
		c.Tech.Quality += decisions.Research.Ecology / 1000000
	}
	if decisions.Research.Durability > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Durability research", research, "", true, float64(-decisions.Research.Durability))
		c.Tech.Quality += decisions.Research.Durability / 1000000
	}
	if decisions.Research.Production_cost > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Production cost research", research, "", true, float64(-decisions.Research.Production_cost))
		c.Tech.Quality += decisions.Research.Production_cost / 1000000
	}
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

	println("Producing products")

	productionReport.ProductSpecificReport = make(map[string]struct {
		TotalProduction       int
		BaseProduction        int
		BonusProduction       int
		ExcessProduction      int
		TotalProductsProduced int
		BaseProductsProduced  int
		BonusProductsProduced int

		MaterialUsed float32
		EnergyUsed   float32
	})

	materialUsed, energyUsed, totalProduction, workerSurplus := ProduceProducts(
		c.ID,
		c.Machines,
		c.Offers,
		c.employeePool,
		productionReport.ProductSpecificReport,
		make([]int, len(c.Machines)),
		decisions.Production.MachineAssignmentPattern,
	)

	productionReport.MaterialUsed = materialUsed
	productionReport.EnergyUsed = energyUsed
	productionReport.WorkerSurplus = workerSurplus
	if len(c.Machines) != 0 {
		productionReport.AvgMachineProductivity = float32(totalProduction) / float32(len(c.Machines))
	}

	materialCosts := -round(float64(externalFactors.MaterialPrice)*float64(productionReport.MaterialUsed), 2)
	energyCosts := -round(float64(externalFactors.EnergyPrice)*float64(productionReport.EnergyUsed), 2)

	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Material costs", production, "The cost of materials used in your products", true, materialCosts)
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Energy costs", production, "The cost of energy used by machines in production", true, energyCosts)
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

type AssignmentPattern int

const (
	FillMachines = iota
	distributeWorkers
)

// returns machines & number of unassigned wokers (if not enough workers for the machines, it return a negative int)
func assignWorkers(employeePool Employee_pool, machines []Machine, workersIds []int, assignmentPattern AssignmentPattern) ([]Machine, int) {
	println("Sorting Employees")

	workersIds = slices.SortedFunc(func(yield func(int) bool) {
		for _, id := range workersIds {
			if !yield(id) {
				return
			}
		}
	}, func(a, b int) int {
		employeeA := employeePool[a]
		employeeB := employeePool[b]
		vala := (employeeA.Motivation * employeeA.Skill * employeeA.WorkingHours)
		valb := (employeeB.Motivation * employeeB.Skill * employeeB.WorkingHours)
		if vala < valb {
			return 1
		} else if vala == valb {
			return 0
		}
		return -1
	})

	slices.Reverse(workersIds)

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

	if assignmentPattern == FillMachines {
		for i := range machines {
			machines[i].AssignedWorkersIDs = make([]int, 0)

			for range machines[i].RequiredWorkers {
				if len(workersIds) >= 1 {
					machines[i].AssignedWorkersIDs = append(machines[i].AssignedWorkersIDs, workersIds[len(workersIds)-1])
					workersIds = workersIds[0 : len(workersIds)-1]
				}
			}
		}
	} else {

		alllMachinesAreFull := func() bool {
			for _, m := range machines {
				if len(m.AssignedWorkersIDs) < m.RequiredWorkers {
					return false
				}
			}

			return true
		}

	AssignmentLoop:
		for len(workersIds) > 0 {
			if alllMachinesAreFull() {
				break AssignmentLoop
			}
			for i := range machines {
				if len(workersIds) == 0 {
					break AssignmentLoop
				}

				if len(machines[i].AssignedWorkersIDs) < machines[i].RequiredWorkers {
					machines[i].AssignedWorkersIDs = append(machines[i].AssignedWorkersIDs, workersIds[len(workersIds)-1])
					workersIds = workersIds[0 : len(workersIds)-1]
				}
			}
		}
	}

	var workerSurplus int

	if len(workersIds) > 0 {
		workerSurplus = len(workersIds)
	} else {
		for _, m := range machines {
			workerSurplus -= m.RequiredWorkers - len(m.AssignedWorkersIDs)
		}
	}

	return machines, workerSurplus
}

// return (base production, bonus production)
func calculateMachineProduction(employeePool Employee_pool, machine Machine) (int, int) {
	// calculate averages
	var avrgSkill float32 = 0
	var avrgMotivation float32 = 0
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
		println("Machine has no workers")
		return 0, 0
	}

	for _, ID := range machine.AssignedWorkersIDs {
		avrgSkill += employeePool[ID].Skill
		avrgMotivation += employeePool[ID].Motivation
		workingHours += employeePool[ID].WorkingHours
	}

	avrgSkill = avrgSkill / float32(len(machine.AssignedWorkersIDs))
	avrgMotivation = avrgMotivation / float32(len(machine.AssignedWorkersIDs))

	if avrgSkill <= 0 {
		log.Println("skill is 0 or less")
		return 0, 0
	}
	if avrgMotivation <= 0 {
		log.Println("motivation is 0 or less")
		return 0, 0
	}

	baseProduction := int(float32(machine.ProductionCapacity) * (workingHours / float32(8*machine.RequiredWorkers)))
	bonusProduction := max(int(float32(baseProduction)*avrgSkill*avrgMotivation-float32(baseProduction)), 0)

	if baseProduction < 0 {
		panic(fmt.Sprintf("base_production is 0 or less (%d)", baseProduction))
	}

	fmt.Printf("------------\n")
	fmt.Printf("machine %d base_production: %d\n", machine.ID, baseProduction)
	fmt.Printf("machine %d bonus_production: %d\n", machine.ID, bonusProduction)
	fmt.Printf("machine %d skill: %f\n", machine.ID, avrgSkill)
	fmt.Printf("machine %d motivation: %f\n", machine.ID, avrgMotivation)
	fmt.Printf("machine %d working_hours: %f\n", machine.ID, workingHours)
	fmt.Printf("Workers %d assigned: %d\n", machine.ID, len(machine.AssignedWorkersIDs))

	return baseProduction, bonusProduction
}

func ProduceProducts(companyID int, machines []Machine, offers map[string]Offer, employeePool Employee_pool, productSpecificReports map[string]struct {
	TotalProduction       int
	BaseProduction        int
	BonusProduction       int
	ExcessProduction      int
	TotalProductsProduced int
	BaseProductsProduced  int
	BonusProductsProduced int

	MaterialUsed float32
	EnergyUsed   float32
}, machineProduction []int,
	assignmentPattern AssignmentPattern,
) (materialUsed float32, energyUsed float32, totalProduction int, workerSurplus int) {
	machines, workerSurplus = assignWorkers(
		employeePool,
		machines,
		employeePool.Get_employees_of_company(
			companyID, Employee_type_production),
		assignmentPattern,
	)

	for i, m := range machines {
		baseProductionOfMachine, bonusProductionOfMachine := calculateMachineProduction(employeePool, m)
		productSpecificReport := productSpecificReports[m.AssignedProductID]
		productSpecificReport.BaseProduction += baseProductionOfMachine
		productSpecificReport.BonusProduction += bonusProductionOfMachine
		productSpecificReport.TotalProduction += bonusProductionOfMachine + baseProductionOfMachine

		productSpecificReport.EnergyUsed += float32(m.EnergyUse)

		productSpecificReports[m.AssignedProductID] = productSpecificReport

		machineProduction[i] = baseProductionOfMachine + bonusProductionOfMachine
		totalProduction += baseProductionOfMachine + bonusProductionOfMachine
	}

	if len(machines) <= 0 {
		println("Company owns no machines!")
	}

	for productID := range offers {
		productSpecificReport := productSpecificReports[productID]

		productSpecificReport.BaseProductsProduced = int(float32(productSpecificReport.BaseProduction) / float32(offers[productID].ProductStats.ProductionCost))
		productSpecificReport.BonusProductsProduced = int(float32(productSpecificReport.BonusProduction) / float32(offers[productID].ProductStats.ProductionCost))
		productSpecificReport.TotalProductsProduced = productSpecificReport.BaseProductsProduced + productSpecificReport.BonusProductsProduced

		productSpecificReport.ExcessProduction = productSpecificReport.TotalProduction - (int(offers[productID].ProductStats.ProductionCost * float32(productSpecificReport.TotalProductsProduced)))

		productSpecificReport.MaterialUsed += offers[productID].ProductStats.MaterialUse * float32(productSpecificReport.TotalProductsProduced)

		productSpecificReports[productID] = productSpecificReport
		fmt.Printf("Product %s (%s) produced: \n    Base: %d\n    Bonus: %d\n    Total: %d\n",
			offers[productID].Product.Name,
			offers[productID].Product.ID,
			productSpecificReport.BaseProductsProduced,
			productSpecificReport.BonusProductsProduced,
			productSpecificReport.TotalProductsProduced)
	}

	for _, r := range productSpecificReports {
		energyUsed += r.EnergyUsed
		materialUsed += r.MaterialUsed
		totalProduction += r.TotalProduction
	}

	return materialUsed, energyUsed, totalProduction, workerSurplus
}
