package simulation

import (
	"cmp"
	"slices"
)

// Finances
func (c *Company) calculateBudget(decisions Decisions, externalFactors ExternalFactors) {
	// using a pointer to avoid refactiong :)
	financialReport := &c.Reports[len(c.Reports)-1].FinancialReport

	localStorageCapacity := 0
	localStorageCosts := 0.0
	for _, w := range c.Warehouses {
		localStorageCosts += float64(-w.OperatingCosts)
		localStorageCapacity += w.Capacity
	}
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(
		"Local storage costs",
		logistics,
		"The operating costs of our own warehouses",
		true,
		localStorageCosts)

	totalProductsInStorage := 0
	for productID := range c.ProductsInStorage {
		totalProductsInStorage += c.ProductsInStorage[productID]
	}

	itemsInExternalStorage := (totalProductsInStorage - localStorageCapacity)

	if itemsInExternalStorage > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(
			"External storage",
			logistics,
			"The cost of storing products in external warehouses, this happens when our own are full",
			true,
			-(float64(externalFactors.ExternalStoragePrice) * float64(itemsInExternalStorage)))
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
	intrest := loan * float64(externalFactors.IntrestRate)
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

	blIntrest := bl * float64(externalFactors.BridgeLoansIntrestRate)
	// update company
	c.BridgeLoans = bl

	if blIntrest > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Bridge loan intrest payments", bridge_loans, "", true, intrest)
	}

	// Increase or decrease loans
	increaseOfLoans := decisions.Finances.SetBankLoan - c.loanQuantity()
	if increaseOfLoans > 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_liabilities("Bank loan", loans, "", true, -float64(increaseOfLoans))
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Income from bank loan", loans, "", true, float64(increaseOfLoans))
	} else if increaseOfLoans < 0 {
		moneyRemaining := float64(-increaseOfLoans)
		var loansToFelete []int
		for i, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
			if e.Group == 1 {
				if moneyRemaining >= e.Value {
					loansToFelete = append(loansToFelete, i)
					moneyRemaining -= e.Value
					c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Payement of loan", loans, "", true, -e.Value)
				} else if moneyRemaining < e.Value {
					c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities[i].Value -= moneyRemaining
					c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Payement of loan", loans, "", true, -moneyRemaining)
					moneyRemaining = 0
					break
				}
			}
		}

		seq := func(yield func(int) bool) {
			for i := range loansToFelete {
				if !yield(i) {
					return
				}
			}
		}
		loansToFelete = slices.SortedFunc(seq, func(a, b int) int { return cmp.Compare(b, a) })

		for _, i := range loansToFelete {
			c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities = slices.Delete(c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities, i, i)
		}
	}

	// Totals
	insertInFinanceReport := func(e FinanceReportEntry) {
		if e.CashCost {
			financialReport.Totals.Cashflow += e.Value
		}

		if e.Value > 0 {
			if e.Group == sales { // Gross sales
				financialReport.Income.GrossSales += e.Value
				financialReport.Income.GrossProfit += e.Value
			} else { // Other income
				financialReport.Income.OtherIncome += e.Value
				financialReport.Income.GrossProfit += e.Value

				println("Other income: " + e.Name + e.Info)
			}
		} else {
			switch e.Group {
			// gross profit
			// operating expenses
			case production_personelle, production, materials, energy: // Cost of sales
				financialReport.Income.CostOfSales += e.Value
				financialReport.Income.GrossProfit += e.Value

			case marketing, marketing_personelle: // advertising
				financialReport.OperatingExpenses.Advertising += e.Value
				financialReport.Totals.TotalOperatingExpenses += e.Value
			case logistics, facilities: // facilities and logistics
				financialReport.OperatingExpenses.FacilitiesAndLogistics += e.Value
				financialReport.Totals.TotalOperatingExpenses += e.Value
			case research: // research & development
				financialReport.OperatingExpenses.ResearchAndDevelopment += e.Value
				financialReport.Totals.TotalOperatingExpenses += e.Value

			// non operating expenses
			case write_off: // write offs
				financialReport.NonOperatingExpenses.WriteOffs += e.Value
				financialReport.Totals.TotalNonOperatingExpenses += e.Value
			case loan_intrest: // loan interest
				financialReport.NonOperatingExpenses.LoanInterest += e.Value
				financialReport.Totals.TotalNonOperatingExpenses += e.Value
			case loans: // loan repayment
				financialReport.NonOperatingExpenses.LoanRepayment += e.Value
				financialReport.Totals.TotalNonOperatingExpenses += e.Value
			case bridge_loan_intrest:
				financialReport.NonOperatingExpenses.BridgeLoanIntrest += e.Value
				financialReport.Totals.TotalNonOperatingExpenses += e.Value
			case bridge_loans:
				financialReport.NonOperatingExpenses.BridgeLoanRepayment += e.Value
				financialReport.Totals.TotalNonOperatingExpenses += e.Value
			case taxes:
				financialReport.NonOperatingExpenses.Taxes += e.Value
			default:
				financialReport.NonOperatingExpenses.Other += e.Value
				financialReport.Totals.TotalNonOperatingExpenses += e.Value
				// totals
				// case taxes:
				//	financial_report.Non_operating_expenses.Taxes += e.Value
			}
		}
		financialReport.Totals.IncomeBeforeTax += e.Value
	}

	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.InvoiceLog {
		insertInFinanceReport(e)
	}

	financialReport.Totals.Cashflow += financialReport.NonOperatingExpenses.Taxes
	ptrTaxesEntry := c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(
		"Taxes",
		taxes,
		"Taxes paid on our profit",
		true,
		tax(financialReport.Totals.IncomeBeforeTax, externalFactors))
	insertInFinanceReport(*ptrTaxesEntry)

	financialReport.Totals.NetIncome = financialReport.Totals.IncomeBeforeTax + financialReport.NonOperatingExpenses.Taxes
	// Calculate bridge loans

	// try to pay off existing bridge loans
	if c.Balance+financialReport.Totals.Cashflow > 0 {
		var loansToDelete []int
		for i, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
			if e.Group != bridge_loans {
				continue
			}
			if c.Balance+financialReport.Totals.Cashflow >= e.Value {
				ptrRepayment := c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Payment of bridge loan", bridge_loans, "", true, -e.Value)
				insertInFinanceReport(*ptrRepayment)

				loansToDelete = append(loansToDelete, i)
			} else if c.Balance+financialReport.Totals.Cashflow < e.Value {
				c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities[i].Value = e.Value - (c.Balance + financialReport.Totals.Cashflow)
				ptrRepayment := c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement("Payement of bridge loan", bridge_loans, "", true, -c.Balance)
				insertInFinanceReport(*ptrRepayment)

				break
			}
		}

		c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities = delete_by_index(c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities, loansToDelete...)

	} else if c.Balance+financialReport.Totals.Cashflow < 0 {
		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(
			"Bridge loan",
			bridge_loans,
			"You are automatically lent out bridge loans when your balance goes beneath 0",
			true,
			-(c.Balance + financialReport.Totals.Cashflow))
		ptrBridgeLoan := c.Reports[len(c.Reports)-1].BalanceSheet.add_to_liabilities(
			"Bridge loan",
			bridge_loans,
			"You are automatically lent out bridge loans when your balance goes beneath 0",
			true,
			-c.Balance)

		insertInFinanceReport(*ptrBridgeLoan)
	}

	c.Balance += financialReport.Totals.Cashflow

	// calculate Liabilities
	totalAssets := 0.0
	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Assets {
		totalAssets += e.Value
	}
	totalLiabilities := 0.0
	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
		totalLiabilities += e.Value
	}

	equity := totalAssets - totalLiabilities
	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_liabilities(
		"Private equity",
		other,
		"The amount of money that is owned exclusively by the company",
		false,
		equity)
}

func tax(EBIT float64, externalFactors ExternalFactors) float64 {
	if EBIT > 0 {
		return -round(EBIT*float64(externalFactors.TaxRate), 2)
	}
	return 0
}

func (c *Company) loanQuantity() (loanValue float64) {
	for _, e := range c.Reports[len(c.Reports)-1].BalanceSheet.Liabilities {
		if e.Group == loans {
			loanValue -= e.Value
		}
	}

	return loanValue
}

func cleanUpFinanceReportEntries(entries []FinanceReportEntry) []FinanceReportEntry {
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Value == 0 {
			entries = delete_by_index(entries, i)
		}
	}

	return entries
}
