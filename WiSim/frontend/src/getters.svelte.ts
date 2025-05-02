import { Statement, Invoice } from "./IncomeAndLoss.svelte";
export async function Get_budget(month: number, company: number): Promise<Statement> {
  let budget = [
    {
      Name: "Est. Income",
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
  let income = [
    {
      Name: "Income",
      Lines: [
        { Name: "Gross Sales", Value: 123 },
        { Name: "Cost of Goods Sold", Value: -123 },
        { Name: "Gross Profit", Value: 123, line_above: true },
      ],
    },
    {
      Name: "Operating Expenses",
      Lines: [
        { Name: "Advertising", Value: -123 },
        { Name: "Facilities & Logistics", Value: -123 },
        { Name: "Equipment (machines)", Value: -123 },
        { Name: "Research & Development", Value: -123 },
        { Name: "Total Operating Expenses", Value: -123, line_above: true },
      ],
    },
    {
      Name: "Non-Operating Expenses",
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

  return income
}

export async function Get_invoices(month: number, company: number): Promise<Invoice[]> {
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
