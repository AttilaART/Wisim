package simulation

import (
	"cmp"
	"slices"
)

// Finances
func (c *Company) calculate_budget(decisions Decisions, external_factors External_factors) {
	// using a pointer to avoid refactiong :)
	financial_report := &c.Reports[len(c.Reports)-1].FinancialReport

	localStorageCapacity := 0
	local_storage_costs := 0.0
	for _, w := range c.Warehouses {
		local_storage_costs += float64(-w.OperatingCosts)
		localStorageCapacity += w.Capacity
	}
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(
		"Local storage costs",
		logistics,
		"The operating costs of our own warehouses",
		true,
		local_storage_costs)

	totalProductsInStorage := 0
	for productId := range c.ProductsInStorage {
		totalProductsInStorage += c.ProductsInStorage[productId]
	}

	itemsInExternalStorage := (totalProductsInStorage - localStorageCapacity)

	if itemsInExternalStorage > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(
			"External storage",
			logistics,
			"The cost of storing products in external warehouses, this happens when our own are full",
			true,
			-(float64(external_factors.ExternalStoragePrice) * float64(itemsInExternalStorage)))
	}

	// Loans

	// Get loans from last year
	if len(c.Reports) >= 2 {
		for _, e := range c.Reports[len(c.Reports)-2].BalanceSheet.Liabilities {
			if e.Group == loans || e.Group == bridge_loans {
				c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities = append(c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities, e)
			}
		}
	}

	// Calculate interest
	loan := 0.0
	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
		if e.Group == loans {
			loan += e.Value
		}
	}
	intrest := loan * float64(external_factors.IntrestRate)
	// update company
	c.Loans = loan

	if intrest > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Intrest payments", loan_intrest, "", true, intrest)
	}
	// Calculate bridge loan intrest
	bl := 0.0 // Bridge_loans
	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
		if e.Group == bridge_loans {
			bl += e.Value
		}
	}

	bl_intrest := bl * float64(external_factors.BridgeLoansIntrestRate)
	// update company
	c.BridgeLoans = bl

	if bl_intrest > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Bridge loan intrest payments", bridge_loans, "", true, intrest)
	}

	// Increase or decrease loans
	increase_of_loans := decisions.Finances.Set_bank_loan - c.loan_quantity()
	if increase_of_loans > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_liabilities("Bank loan", loans, "", true, -float64(increase_of_loans))
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Income from bank loan", loans, "", true, float64(increase_of_loans))
	} else if increase_of_loans < 0 {
		money_remaining := float64(-increase_of_loans)
		var loans_to_delete []int
		for i, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
			if e.Group == 1 {
				if money_remaining >= e.Value {
					loans_to_delete = append(loans_to_delete, i)
					money_remaining -= e.Value
					c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Payement of loan", loans, "", true, -e.Value)
				} else if money_remaining < e.Value {
					c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities[i].Value -= money_remaining
					c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Payement of loan", loans, "", true, -money_remaining)
					money_remaining = 0
					break
				}
			}
		}

		seq := func(yield func(int) bool) {
			for i := range loans_to_delete {
				if !yield(i) {
					return
				}
			}
		}
		loans_to_delete = slices.SortedFunc(seq, func(a, b int) int { return cmp.Compare(b, a) })

		for _, i := range loans_to_delete {
			c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities = slices.Delete(c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities, i, i)
		}
	}

	// Totals
	insert_in_finance_report := func(e FinanceReportEntry) {
		if e.CashCost {
			financial_report.Totals.Cashflow += e.Value
		}

		if e.Value > 0 {
			if e.Group == sales { // Gross sales
				financial_report.Income.GrossSales += e.Value
				financial_report.Income.GrossProfit += e.Value
			} else { // Other income
				financial_report.Income.OtherIncome += e.Value
				financial_report.Income.GrossProfit += e.Value

				println("Other income: " + e.Name + e.Info)
			}
		} else {
			switch e.Group {
			// gross profit
			// operating expenses
			case production_personelle, production, materials, energy: // Cost of sales
				financial_report.Income.CostOfSales += e.Value
				financial_report.Income.GrossProfit += e.Value

			case marketing, marketing_personelle: // advertising
				financial_report.OperatingExpenses.Advertising += e.Value
				financial_report.Totals.TotalOperatingExpenses += e.Value
			case logistics, facilities: // facilities and logistics
				financial_report.OperatingExpenses.FacilitiesAndLogistics += e.Value
				financial_report.Totals.TotalOperatingExpenses += e.Value
			case research: // research & development
				financial_report.OperatingExpenses.ResearchAndDevelopment += e.Value
				financial_report.Totals.TotalOperatingExpenses += e.Value

			// non operating expenses
			case write_off: // write offs
				financial_report.NonOperatingExpenses.WriteOffs += e.Value
				financial_report.Totals.TotalNonOperatingExpenses += e.Value
			case loan_intrest: // loan interest
				financial_report.NonOperatingExpenses.LoanInterest += e.Value
				financial_report.Totals.TotalNonOperatingExpenses += e.Value
			case loans: // loan repayment
				financial_report.NonOperatingExpenses.LoanRepayment += e.Value
				financial_report.Totals.TotalNonOperatingExpenses += e.Value
			case bridge_loan_intrest:
				financial_report.NonOperatingExpenses.BridgeLoanIntrest += e.Value
				financial_report.Totals.TotalNonOperatingExpenses += e.Value
			case bridge_loans:
				financial_report.NonOperatingExpenses.BridgeLoanRepayment += e.Value
				financial_report.Totals.TotalNonOperatingExpenses += e.Value
			case taxes:
				financial_report.NonOperatingExpenses.Taxes += e.Value
			default:
				financial_report.NonOperatingExpenses.Other += e.Value
				financial_report.Totals.TotalNonOperatingExpenses += e.Value
				// totals
				// case taxes:
				//	financial_report.Non_operating_expenses.Taxes += e.Value
			}
		}
		financial_report.Totals.IncomeBeforeTax += e.Value
	}

	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.InvoiceLog {
		insert_in_finance_report(e)
	}

	financial_report.Totals.Cashflow += financial_report.NonOperatingExpenses.Taxes
	ptr_taxes_entry := c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(
		"Taxes",
		taxes,
		"Taxes paid on our profit",
		true,
		tax(financial_report.Totals.IncomeBeforeTax, external_factors))
	insert_in_finance_report(*ptr_taxes_entry)

	financial_report.Totals.NetIncome = financial_report.Totals.IncomeBeforeTax + financial_report.NonOperatingExpenses.Taxes
	// Calculate bridge loans

	// try to pay off existing bridge loans
	if c.Balance+financial_report.Totals.Cashflow > 0 {
		var loans_to_delete []int
		for i, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
			if e.Group != bridge_loans {
				continue
			}
			if c.Balance+financial_report.Totals.Cashflow >= e.Value {
				ptr_repayment := c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Payment of bridge loan", bridge_loans, "", true, -e.Value)
				insert_in_finance_report(*ptr_repayment)

				loans_to_delete = append(loans_to_delete, i)
			} else if c.Balance+financial_report.Totals.Cashflow < e.Value {
				c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities[i].Value = e.Value - (c.Balance + financial_report.Totals.Cashflow)
				ptr_repayment := c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Payement of bridge loan", bridge_loans, "", true, -c.Balance)
				insert_in_finance_report(*ptr_repayment)

				break
			}
		}

		c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities = delete_by_index(c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities, loans_to_delete...)

	} else if c.Balance+financial_report.Totals.Cashflow < 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(
			"Bridge loan",
			bridge_loans,
			"You are automatically lent out bridge loans when your balance goes beneath 0",
			true,
			-(c.Balance + financial_report.Totals.Cashflow))
		ptr_bridge_loan := c.Reports[len(c.Reports)-1].BalanceSheet.add_to_liabilities(
			"Bridge loan",
			bridge_loans,
			"You are automatically lent out bridge loans when your balance goes beneath 0",
			true,
			-c.Balance)

		insert_in_finance_report(*ptr_bridge_loan)
	}

	c.Balance += financial_report.Totals.Cashflow

	// calculate Liabilities
	total_assets := 0.0
	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Assets {
		total_assets += e.Value
	}
	total_liabilities := 0.0
	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
		total_liabilities += e.Value
	}

	equity := total_assets - total_liabilities
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_liabilities(
		"Private equity",
		other,
		"The amount of money that is owned exclusively by the company",
		false,
		equity)
}

func tax(EBIT float64, external_factors External_factors) float64 {
	if EBIT > 0 {
		return -round(EBIT*float64(external_factors.TaxRate), 2)
	}
	return 0
}

func (c *Company) loan_quantity() (loan_value float64) {
	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
		if e.Group == loans {
			loan_value -= e.Value
		}
	}

	return loan_value
}

func clean_up_financeReportEntries(entries []FinanceReportEntry) []FinanceReportEntry {
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Value == 0 {
			entries = delete_by_index(entries, i)
		}
	}

	return entries
}
