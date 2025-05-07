package simulation

import "errors"

// Logistics
func (company *Company) calculate_logistics(decisions Decisionsold) ([]FinanceReportEntry, []FinanceReportEntry, error) {
	var Income_entries []FinanceReportEntry
	var Assets_entries []FinanceReportEntry

	if decisions.Purchase_of_warehouses > 0 {
		existing_warehouses := len(company.Warehouses)

		for i := range decisions.Purchase_of_warehouses {
			company.Warehouses = append(company.Warehouses, Warehouse{
				existing_warehouses + i,
				20000,
				20000,
				100000,
			})
			Income_entries = append(Income_entries, FinanceReportEntry{"Purchase of warehouse", logistics, "", false, -100000})
		}
	} else if (-decisions.Purchase_of_warehouses) == len(company.Warehouses) {
		company.Reports[len(company.Reports)].Production_report.Warehouses_bought = len(company.Warehouses)
		company.Warehouses = []Warehouse{}
	} else if decisions.Purchase_of_warehouses < 0 {
		if -decisions.Purchase_of_warehouses > len(company.Warehouses) {
			return Income_entries, Assets_entries, errors.New("can't sell more warehouses than you have")
		}
		for range -decisions.Purchase_of_warehouses {
			Income_entries = append(Income_entries, FinanceReportEntry{
				"Selling of warehouse",
				logistics,
				"",
				true,
				float64(company.Warehouses[len(company.Warehouses)-1].Value),
			})
			company.Warehouses = company.Warehouses[1:len(company.Warehouses)]
		}

		for _, w := range company.Warehouses {
			Assets_entries = append(Assets_entries, FinanceReportEntry{"Warehouse as asset", logistics, "", false, float64(w.Value)})
		}
	}

	company.Reports[len(company.Reports)-1].Production_report.Warehouses_bought = decisions.Purchase_of_warehouses
	return Income_entries, Assets_entries, nil
}
