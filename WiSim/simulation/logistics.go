package simulation

// Logistics
func (c *Company) calculate_logistics(decisions Decisions) {
	if len(c.Warehouses) == 0 && len(decisions.Production.Logistics) == 0 {
		return
	}

	for i := range c.Warehouses {
		c.Machines[i].Status = Existing
	}

	// Purchase Machines
	println("Purchasing warehouses")
	var warehouses_to_delete_index []int
	for i, w := range decisions.Production.Logistics {
		if w.Status == New {
			c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Purchase of warehouse", production, "", true, -float64(w.Value))
		} else if w.Status == Sold {
			c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Selling of warehouse", production, "", true, float64(w.Value))
			warehouses_to_delete_index = append(warehouses_to_delete_index, i)
		}
	}

	c.Warehouses = delete_by_index(decisions.Production.Logistics, warehouses_to_delete_index...)

	for _, w := range c.Warehouses {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_equity("Warehouse", logistics, "", false, float64(w.Value))
	}
}
