```go
func calculateProductStatsWrapped() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 2 {
			return fmt.Sprintf("Invalid arguments: Expected 2, Got %d",
				len(args))
		}

		err := checkArguments(args)
		if err != nil {
			return err.Error()
		}

		productStats, productionLineCost := simulation.
			CalculateProductStats(product, productCompontents)

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
```
