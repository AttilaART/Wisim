package simulation

// Logistics
func (c *Company) calculate_logistics(decisions Decisions) {
	if len(c.Warehouses) == 0 && len(decisions.Production.Logistics) == 0 {
		return
	}

	// Purchase Machines
	println("Purchasing warehouses")
	var warehouses_to_delete_id []int
	for _, w := range decisions.Production.Logistics {
		if w.Change == Delta_New {
			c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Purchase of warehouse", production, "", true, -float64(w.Item.Value))
			c.Warehouses = append(c.Warehouses, w.Item)
		} else if w.Change == Delta_Remove {
			c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Selling of warehouse", production, "", true, float64(w.Item.Value))
			warehouses_to_delete_id = append(warehouses_to_delete_id, w.Item.ID)
		}
	}

	c.Warehouses = deleteByID(c.Warehouses, warehouses_to_delete_id...)

	for _, w := range c.Warehouses {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_equity("Warehouse", logistics, "", false, float64(w.Value))
	}
}
