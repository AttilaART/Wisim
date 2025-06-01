<script lang="ts">
  import Sidebar from "./Sidebar.svelte";
  import IncomeAndLoss from "./IncomeAndLoss.svelte";
  import {
    Get_budget,
    Get_budget_invoices,
    Get_income_statement,
    Get_invoices,
  } from "./getters.svelte";
  import Debt from "./Debt.svelte";
  import BalanceSheet from "./BalanceSheet.svelte";

  let page = $state("budget");
</script>

<Sidebar
  buttons={[
    {
      Text: "Budget",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "budget";
      },
      selected_by_default: true,
    },
    {
      Text: "Income & Loss",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "income and loss";
      },
    },
    {
      Text: "Debt",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "debt";
      },
    },
    {
      Text: "Balance Sheet",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "balance sheet";
      },
    },
    {
      Text: "Graphs",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "graphs";
      },
    },
  ]}
  expand={true}
  horisontal={true}
  keep_pressed={true}
></Sidebar>
{#if page == "budget"}
  <IncomeAndLoss
    get_income_statement={Get_budget}
    get_invoice_log={Get_budget_invoices}
    is_budget={true}
  ></IncomeAndLoss>
{:else if page == "income and loss"}
  <IncomeAndLoss
    get_income_statement={Get_income_statement}
    get_invoice_log={Get_invoices}
    is_budget={false}
  ></IncomeAndLoss>
{:else if page == "debt"}
  <Debt></Debt>
{:else if page == "balance sheet"}
  <BalanceSheet></BalanceSheet>
{:else if page == "graphs"}
  <h2>Coming Soon</h2>
{/if}
