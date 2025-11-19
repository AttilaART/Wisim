//go:build js && wasm

package main

import (
	"WiSim/simulation"
	"encoding/json"
	"fmt"
	"syscall/js"
)

func calculateProductStatsWrapped() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 2 {
			return fmt.Sprintf("Invalid arguments: Expected 2, Got %d", len(args))
		}

		var product simulation.Product
		err := json.Unmarshal([]byte(args[0].String()), &product)
		if err != nil {
			return err.Error()
		}

		var productCompontents simulation.ProductComponents
		err = json.Unmarshal([]byte(args[1].String()), &productCompontents)
		if err != nil {
			return err.Error()
		}

		productStats, productionLineCost := simulation.CalculateProductStats(product, productCompontents)

		json, err := json.Marshal(struct {
			ProductStats       simulation.ProductStats
			ProductionLineCost float64
		}{productStats, productionLineCost})
		if err != nil {
			return err.Error()
		}

		return string(json)
	})

	return f
}

func calculateProductionWrapped() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 5 {
			return fmt.Sprintf("Invalid arguments: Expected 4, Got %d", len(args))
		}

		var err error

		var args1 []simulation.Machine
		err = json.Unmarshal(([]byte)(args[1].String()), &args1)
		if err != nil {
			return err.Error()
		}

		var args2 map[string]simulation.Offer
		err = json.Unmarshal(([]byte)(args[2].String()), &args2)
		if err != nil {
			return err.Error()
		}

		var employees []simulation.Employee
		err = json.Unmarshal(([]byte)(args[3].String()), &employees)
		if err != nil {
			return err.Error()
		}

		args3 := make(simulation.Employee_pool)

		for i, e := range employees {
			employees[i].Employer = args[0].Int()
			args3[e.ID] = &employees[i]
		}

		productSpecificReportTemp := make(map[string]struct {
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

		machineProduction := make([]int, len(args1))

		_, _, _, workerSurplus := simulation.ProduceProducts(
			args[0].Int(),
			args1,
			args2,
			args3,
			productSpecificReportTemp,
			machineProduction,
			simulation.AssignmentPattern(args[4].Int()),
		)

		returnValue := struct {
			WorkerSurplus     int
			MachineProduction []int
		}{workerSurplus, machineProduction}

		json, err := json.Marshal(returnValue)
		if err != nil {
			return err.Error()
		}

		return string(json)
	})
	return f
}

func simulateMockStepWrapped() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 5 {
			return fmt.Sprint("expected 4 arguments, got ", len(args))
		}

		var company simulation.Company
		err := json.Unmarshal([]byte(args[0].String()), &company)
		if err != nil {
			return err.Error()
		}

		var decisions simulation.Decisions
		err = json.Unmarshal([]byte(args[1].String()), &decisions)
		if err != nil {
			return err.Error()
		}

		var externalFactors simulation.ExternalFactors
		err = json.Unmarshal([]byte(args[2].String()), &externalFactors)
		if err != nil {
			return err.Error()
		}

		var employeArray []simulation.Employee
		err = json.Unmarshal([]byte(args[3].String()), &employeArray)
		if err != nil {
			return err.Error()
		}

		employePool := make(simulation.Employee_pool)

		for _, e := range employeArray {
			employePool[e.ID] = &e
		}

		company.SetEmployeePool(employePool)

		json, err := json.Marshal(simulation.SimulateMockStep(
			company,
			decisions,
			externalFactors,
			employePool,
			args[4].Int(),
		))
		if err != nil {
			return err.Error()
		}

		return string(json)
	})

	return f
}

func main() {
	println("Wisim WASM loaded Successfully")
	js.Global().Set("CalculateProductStatsGo", calculateProductStatsWrapped())
	js.Global().Set("CalculateProductionGo", calculateProductionWrapped())
	js.Global().Set("SimulateMockStep", simulateMockStepWrapped())

	<-make(chan struct{})
}
