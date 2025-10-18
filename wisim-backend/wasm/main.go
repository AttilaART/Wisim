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

func main() {
	println("Wisim WASM loaded Successfully")
	js.Global().Set("CalculateProductStatsGo", calculateProductStatsWrapped())

	<-make(chan struct{})
}
