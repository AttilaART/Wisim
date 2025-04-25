package simulation

import (
	"cmp"
	"slices"
)

// Finances
func (company *Company) calculate_budget(decisions Decisions, external_factors External_factors) Financial_Report {
	var financial_report Financial_Report

	local_storage_capacity := 0
	local_storage_costs := 0.0
	for _, w := range company.Warehouses {
		local_storage_costs += float64(-w.Operating_costs)
		local_storage_capacity += w.Capacity
	}
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
		"Internal storage costs",
		logistics,
		"The operating costs of our own warehouses",
		true,
		local_storage_costs)
	items_in_external_storage := (company.Items_in_storage - local_storage_capacity)
	if items_in_external_storage > 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
			"External storage",
			logistics,
			"The cost of storing products in external warehouses, this happens when our own are full",
			true,
			-(float64(external_factors.External_storage_price) * float64(company.Items_in_storage)))
	}

	// Loans

	// Get loans from last year
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
		if e.Group == loans || e.Group == bridge_loans {
			company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities = append(company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities, e)
		}
	}

	// Calculate interest
	intrest := 0.0
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
		if e.Group == loans {
			intrest += e.Value * float64(external_factors.Intrest_rate)
		}
	}

	if intrest > 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Intrest payments", loan_intrest, "", true, -intrest)
	}
	// Calculate bride loan intrest
	bl_intrest := 0.0
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
		if e.Group == bridge_loans {
			bl_intrest += e.Value * float64(external_factors.Bridge_loans_intrest_rate)
		}
	}

	if bl_intrest > 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Bridge loan intrest payments", bridge_loans, "", true, -intrest)
	}

	// Increase or decrease loans
	if decisions.Increase_of_loans > 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_liabilities("Bank loan", loans, "", true, float64(decisions.Increase_of_loans))
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Income from bank loan", loans, "", true, float64(decisions.Increase_of_loans))
	} else if decisions.Increase_of_loans < 0 {
		money_remaining := float64(-decisions.Increase_of_loans)
		var loans_to_delete []int
		for i, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
			if e.Group == 1 {
				if money_remaining >= e.Value {
					loans_to_delete = append(loans_to_delete, i)
					money_remaining -= e.Value
					company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Payement of loan", loans, "", true, -e.Value)
				} else if money_remaining < e.Value {
					company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities[i].Value -= money_remaining
					company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Payement of loan", loans, "", true, -money_remaining)
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
			company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities = slices.Delete(company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities, i, i)
		}
	}

	// Totals
	financial_report.Total_income = 0
	financial_report.Cash_costs = 0
	financial_report.Non_cash_costs = 0
	financial_report.Total_expenses_before_tax = 0
	financial_report.Operating_profit = 0
	financial_report.Taxes = 0
	financial_report.Bridge_loan_repayments = 0
	financial_report.Loan_repayments = 0
	financial_report.Total_expenses_after_tax = 0
	financial_report.Net_Profit = financial_report.Total_expenses_after_tax + financial_report.Total_income

	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement {
		financial_report.Net_Profit += e.Value
		if e.Value > 0 && (e.Group != loans || e.Group != bridge_loans) {
			financial_report.Total_income += e.Value
		} else {
			financial_report.Total_expenses_after_tax += e.Value
			financial_report.Total_expenses_before_tax += e.Value
			if e.Cash_cost {
				financial_report.Cash_costs += e.Value
			} else {
				financial_report.Non_cash_costs += e.Value
			}
		}
		if e.Group != loan_intrest && e.Group != bridge_loan_intrest {
			financial_report.Operating_profit += e.Value
		}
		if e.Group == loans {
			financial_report.Loan_repayments += e.Value
		}
		if e.Group == bridge_loans {
			financial_report.Bridge_loan_repayments += e.Value
		}
		if e.Cash_cost {
			company.Balance += e.Value
		}
	}

	financial_report.Taxes = profit_taxes(financial_report.Operating_profit, external_factors)
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
		"Taxes",
		taxes,
		"Taxes paid on our profit",
		true,
		financial_report.Taxes)

	financial_report.Net_Profit += financial_report.Taxes
	financial_report.Total_expenses_after_tax += financial_report.Taxes
	financial_report.Cash_costs += financial_report.Taxes
	company.Balance += financial_report.Taxes

	// Calculate bridge loans

	// try to pay off existing bridge loans
	if company.Balance > 0 {
		var loans_to_delete []int
		for i, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
			if e.Group != bridge_loans {
				continue
			}
			if company.Balance >= e.Value {
				company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Payment of bridge loan", bridge_loans, "", true, -e.Value)
				loans_to_delete = append(loans_to_delete, i)
				company.Balance -= e.Value
				financial_report.Cash_costs -= e.Value
				financial_report.Total_expenses_before_tax -= e.Value
				financial_report.Total_expenses_after_tax -= e.Value
				financial_report.Net_Profit -= e.Value
				financial_report.Bridge_loan_repayments -= e.Value
			} else if company.Balance < e.Value {
				company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities[i].Value = e.Value - company.Balance
				company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Payement of bridge loan", bridge_loans, "", true, -company.Balance)
				financial_report.Cash_costs -= company.Balance
				financial_report.Total_expenses_before_tax -= company.Balance
				financial_report.Total_expenses_after_tax -= company.Balance
				financial_report.Net_Profit -= company.Balance
				financial_report.Bridge_loan_repayments -= company.Balance
				company.Balance = 0
				break
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
			company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities = slices.Delete(company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities, i, i)
		}
	} else if company.Balance < 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
			"Bridge loan",
			bridge_loans,
			"You are automatically lent out bridge loans when your balance goes beneath 0",
			true,
			-company.Balance)
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_liabilities(
			"Bridge loan",
			bridge_loans,
			"You are automatically lent out bridge loans when your balance goes beneath 0",
			true,
			-company.Balance)
		// financial_report.Total_income += -company.Balance
		// financial_report.EBIT += -company.Balance
		// financial_report.Profit += -company.Balance
		financial_report.New_bridge_loan += -company.Balance
		// company.Balance = 0
		company.Balance = -company.Balance
	}

	// calculate Liabilities
	total_assets := 0.0
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Assets {
		total_assets += e.Value
	}
	total_liabilities := 0.0
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
		total_liabilities += e.Value
	}

	equity := total_assets - total_liabilities
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_liabilities(
		"Private equity",
		other,
		"The amount of money that is owned exclusively by the company",
		false,
		equity)

	return financial_report
}

func profit_taxes(EBIT float64, external_factors External_factors) float64 {
	if EBIT > 0 {
		return -round(EBIT*float64(external_factors.Tax_rate), 2)
	}
	return 0
}
