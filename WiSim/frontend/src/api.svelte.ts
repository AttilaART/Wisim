import { get } from "svelte/store";
import {
  Get_company,
  Get_decisions,
  Get_external_factors,
  Initial_app_load,
  New_simulation,
  Submit_decisions,
  Trigger_simulation,
} from "../wailsjs/go/main/App";
import { int, simulation } from "../wailsjs/go/models";
import { company_id, current_decisions, update_decisions, update_external_factors, update_reports, latest_reports, update_company, company, external_factors } from "./store.svelte";
import { Statement, Invoice } from "./IncomeAndLoss.svelte";

export async function initial_app_load(): Promise<void> {
  return Initial_app_load()
}

export async function start_new_game() {
  update_external_factors(await New_simulation(1));
  update_company(await Get_company(get(company_id)))

  update_decisions(await Get_decisions(get(company_id), external_factors.Month));
}

export async function trigger_simulation(force?: boolean) {
  await Submit_decisions(get(company_id), current_decisions)

  update_external_factors(await Trigger_simulation(force))
  update_company(await Get_company(get(company_id)))
  update_reports(company.Reports[company.Reports.length - 1])

  update_decisions(await Get_decisions(get(company_id), external_factors.Month));
}

export async function Get_budget(month: number, company: number): Promise<Statement> {
  let budget = [
    {
      Name: "Est. Income",
      Period: "Budget 03/0001",
      Lines: [
        { Name: "Gross Sales", Value: 123 },
        { Name: "Cost of Goods Sold", Value: -123 },
        { Name: "Gross Profit", Value: 123, line_above: true },
      ],
    },
    {
      Name: "Est. Operating Expenses",
      Lines: [
        { Name: "Advertising", Value: -123 },
        { Name: "Facilities & Logistics", Value: -123 },
        { Name: "Equipment (machines)", Value: -123 },
        { Name: "Research & Development", Value: -123 },
        { Name: "Total Operating Expenses", Value: -123, line_above: true },
      ],
    },
    {
      Name: "Est. Non-Operating Expenses",
      Lines: [
        { Name: "Write-Offs", Value: -123 },
        { Name: "Loan interest", Value: -123 },
        { Name: "Bridge Loan interest", Value: -123 },
        { Name: "Total Non-Operating Expenses", Value: -123, line_above: true },
        { Name: "    Taxes", Value: -123 },
        { Name: "Net Income", Value: 123, line_above: true },
        { Name: "Cashflow", Value: 123, line_above: true },
      ],
    },
  ];

  return budget
}

export async function Get_budget_invoices(month: number, company: number): Promise<Invoice[]> {
  let invoice: Invoice[] = [
    {
      Name: "Sale of products",
      Info: "Sold 14'000 Products",
      Category: "Sales",
      Value: 1800000,
    },
    {
      Name: "Employee pay",
      Info: "Payed 150 production personelle",
      Category: "Personelle",
      Value: -670000,
    },
  ];

  return invoice
}

export async function Get_income_statement(month: number, company: number): Promise<Statement> {
  if (Object.entries(latest_reports).length == 0) {
    throw new Error("No report availible")
  }

  let income_report: simulation.Financial_Report = latest_reports.Financial_Report

  let statement: Statement = []
  for (let sec of Object.entries(income_report)) {
    let section = {
      Name: sec[0],
      Lines: []
    }

    for (let l of Object.entries(sec[1])) {
      let line = {
        Name: l[0],
        Value: l[1]
      }
      section.Lines.push(line)
    }
    statement.push(section)
  }
  return statement
}

export async function Get_invoices(month: number, company: number): Promise<Invoice[]> {
  let entries: simulation.FinanceReportEntry[]
  try {

    entries = latest_reports.Balance_sheet.Invoice_log
  } catch (exception) {
    if (exception == TypeError) {
      throw new Error("No invoice log availible")
    } else {
      throw (exception)
    }
  }


  let invoices: Invoice[] = []
  for (let entry of entries) {
    let invoice: Invoice = {
      Name: entry.Name,
      Info: entry.Info,
      Category: int.int[entry.Group],
      Value: entry.Value
    }

    invoices.push(invoice)
  }

  return invoices
}
