package simulation

import (
	"cmp"
	"slices"
)

// Finances
func (c *Company) calculate_budget(decisions Decisions, external_factors External_factors) {
	// using a pointer to avoid refactiong :)
	financial_report := &c.Reports[len(c.Reports)-1].Financial_Report

	local_storage_capacity := 0
	local_storage_costs := 0.0
	for _, w := range c.Warehouses {
		local_storage_costs += float64(-w.Operating_costs)
		local_storage_capacity += w.Capacity
	}
	c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
		"Local storage costs",
		logistics,
		"The operating costs of our own warehouses",
		true,
		local_storage_costs)
	items_in_external_storage := (c.Items_in_storage - local_storage_capacity)
	if items_in_external_storage > 0 {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
			"External storage",
			logistics,
			"The cost of storing products in external warehouses, this happens when our own are full",
			true,
			-(float64(external_factors.External_storage_price) * float64(c.Items_in_storage)))
	}

	// Loans

	// Get loans from last year
	if len(c.Reports) >= 2 {
		for _, e := range c.Reports[len(c.Reports)-2].Balance_sheet.Liabilities {
			if e.Group == loans || e.Group == bridge_loans {
				c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities = append(c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities, e)
			}
		}
	}

	// Calculate interest
	loan := 0.0
	for _, e := range c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities {
		if e.Group == loans {
			loan += e.Value
		}
	}
	intrest := loan * float64(external_factors.Intrest_rate)
	// update company
	c.Loans = loan

	if intrest > 0 {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Intrest payments", loan_intrest, "", true, intrest)
	}
	// Calculate bridge loan intrest
	bl := 0.0 // Bridge_loans
	for _, e := range c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities {
		if e.Group == bridge_loans {
			bl += e.Value
		}
	}

	bl_intrest := bl * float64(external_factors.Bridge_loans_intrest_rate)
	// update company
	c.Bridge_loans = bl

	if bl_intrest > 0 {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Bridge loan intrest payments", bridge_loans, "", true, intrest)
	}

	// Increase or decrease loans
	increase_of_loans := decisions.Finances.Set_bank_loan - c.loan_quantity()
	if increase_of_loans > 0 {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_liabilities("Bank loan", loans, "", true, -float64(increase_of_loans))
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Income from bank loan", loans, "", true, float64(increase_of_loans))
	} else if increase_of_loans < 0 {
		money_remaining := float64(-increase_of_loans)
		var loans_to_delete []int
		for i, e := range c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities {
			if e.Group == 1 {
				if money_remaining >= e.Value {
					loans_to_delete = append(loans_to_delete, i)
					money_remaining -= e.Value
					c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Payement of loan", loans, "", true, -e.Value)
				} else if money_remaining < e.Value {
					c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities[i].Value -= money_remaining
					c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Payement of loan", loans, "", true, -money_remaining)
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
			c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities = slices.Delete(c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities, i, i)
		}
	}

	// Totals
	insert_in_finance_report := func(e FinanceReportEntry) {
		if e.Cash_cost {
			financial_report.Non_operating_expenses.Cashflow += e.Value
		}

		if e.Value > 0 {
			if e.Group == sales { // Gross sales
				financial_report.Income.Gross_sales += e.Value
				financial_report.Income.Gross_profit += e.Value
			} else { // Other income
				financial_report.Income.Other_income += e.Value
				financial_report.Income.Gross_profit += e.Value

				println("Other income: " + e.Name + e.Info)
			}
		} else {
			switch e.Group {
			// gross profit
			// operating expenses
			case production_personelle, production, materials, energy: // Cost of sales
				financial_report.Income.Cost_of_sales += e.Value
				financial_report.Income.Cost_of_sales += e.Value

			case marketing, marketing_personelle: // advertising
				financial_report.Operating_expenses.Advertising += e.Value
				financial_report.Operating_expenses.Total_operating_expenses += e.Value
			case logistics, facilities: // facilities and logistics
				financial_report.Operating_expenses.Facilities_and_logistics += e.Value
				financial_report.Operating_expenses.Total_operating_expenses += e.Value
			case research: // research & development
				financial_report.Operating_expenses.Research_and_development += e.Value
				financial_report.Operating_expenses.Total_operating_expenses += e.Value

			// non operating expenses
			case write_off: // write offs
				financial_report.Non_operating_expenses.Write_offs += e.Value
				financial_report.Non_operating_expenses.Total_non_operating_expenses += e.Value
			case loan_intrest: // loan interest
				financial_report.Non_operating_expenses.Loan_interest += e.Value
				financial_report.Non_operating_expenses.Total_non_operating_expenses += e.Value
			case loans: // loan repayment
				financial_report.Non_operating_expenses.Loan_repayment += e.Value
				financial_report.Non_operating_expenses.Total_non_operating_expenses += e.Value
			case bridge_loan_intrest:
				financial_report.Non_operating_expenses.Bridge_loan_intrest += e.Value
				financial_report.Non_operating_expenses.Total_non_operating_expenses += e.Value
			case bridge_loans:
				financial_report.Non_operating_expenses.Bridge_loan_repayment += e.Value
				financial_report.Non_operating_expenses.Total_non_operating_expenses += e.Value
			case taxes:
				financial_report.Non_operating_expenses.Taxes += e.Value
			default:
				financial_report.Non_operating_expenses.Other += e.Value
				financial_report.Non_operating_expenses.Total_non_operating_expenses += e.Value
				// totals
				// case taxes:
				//	financial_report.Non_operating_expenses.Taxes += e.Value
			}
		}
		financial_report.Non_operating_expenses.Income_before_tax += e.Value
	}

	for _, e := range c.Reports[len(c.Reports)-1].Balance_sheet.Invoice_log {
		insert_in_finance_report(e)
	}

	financial_report.Non_operating_expenses.Cashflow += financial_report.Non_operating_expenses.Taxes
	ptr_taxes_entry := c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
		"Taxes",
		taxes,
		"Taxes paid on our profit",
		true,
		tax(financial_report.Non_operating_expenses.Income_before_tax, external_factors))
	insert_in_finance_report(*ptr_taxes_entry)

	financial_report.Non_operating_expenses.Net_income = financial_report.Non_operating_expenses.Income_before_tax + financial_report.Non_operating_expenses.Taxes
	// Calculate bridge loans

	// try to pay off existing bridge loans
	if c.Balance+financial_report.Non_operating_expenses.Cashflow > 0 {
		var loans_to_delete []int
		for i, e := range c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities {
			if e.Group != bridge_loans {
				continue
			}
			if c.Balance+financial_report.Non_operating_expenses.Cashflow >= e.Value {
				ptr_repayment := c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Payment of bridge loan", bridge_loans, "", true, -e.Value)
				insert_in_finance_report(*ptr_repayment)

				loans_to_delete = append(loans_to_delete, i)
			} else if c.Balance+financial_report.Non_operating_expenses.Cashflow < e.Value {
				c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities[i].Value = e.Value - (c.Balance + financial_report.Non_operating_expenses.Cashflow)
				ptr_repayment := c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement("Payement of bridge loan", bridge_loans, "", true, -c.Balance)
				insert_in_finance_report(*ptr_repayment)

				break
			}
		}

		c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities = delete_by_index(c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities, loans_to_delete...)

	} else if c.Balance+financial_report.Non_operating_expenses.Cashflow < 0 {
		c.Reports[len(c.Reports)-1].Balance_sheet.add_to_income_statement(
			"Bridge loan",
			bridge_loans,
			"You are automatically lent out bridge loans when your balance goes beneath 0",
			true,
			-(c.Balance + financial_report.Non_operating_expenses.Cashflow))
		ptr_bridge_loan := c.Reports[len(c.Reports)-1].Balance_sheet.add_to_liabilities(
			"Bridge loan",
			bridge_loans,
			"You are automatically lent out bridge loans when your balance goes beneath 0",
			true,
			-c.Balance)

		insert_in_finance_report(*ptr_bridge_loan)
	}

	c.Balance += financial_report.Non_operating_expenses.Cashflow

	// calculate Liabilities
	total_assets := 0.0
	for _, e := range c.Reports[len(c.Reports)-1].Balance_sheet.Assets {
		total_assets += e.Value
	}
	total_liabilities := 0.0
	for _, e := range c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities {
		total_liabilities += e.Value
	}

	equity := total_assets - total_liabilities
	c.Reports[len(c.Reports)-1].Balance_sheet.add_to_liabilities(
		"Private equity",
		other,
		"The amount of money that is owned exclusively by the company",
		false,
		equity)
}

func tax(EBIT float64, external_factors External_factors) float64 {
	if EBIT > 0 {
		return -round(EBIT*float64(external_factors.Tax_rate), 2)
	}
	return 0
}

func (c *Company) loan_quantity() (loan_value float64) {
	for _, e := range c.Reports[len(c.Reports)-1].Balance_sheet.Liabilities {
		if e.Group == loans {
			loan_value -= e.Value
		}
	}

	return loan_value
}
