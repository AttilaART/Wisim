<script lang="ts">
  import {
    Get_accounting_data,
    Get_financial_report,
  } from "../wailsjs/go/main/App";

  import Table from "./Table.svelte";
  import Report from "./Report.svelte";

  import { month_counter, financial_data_store } from "./store";

  let month = $state(0);

  month_counter.subscribe((value) => {
    month = value;
    get_income_statement_data();
    get_assets_data();
    get_liabilities_data();
  });

  type accounting_data = {
    data: object[] | undefined;
    loading: boolean;
    error: Error | null;
  };

  let income_statement: accounting_data = $state({
    data: undefined,
    loading: true,
    error: null,
  });
  let assets: accounting_data = $state({
    data: undefined,
    loading: true,
    error: null,
  });
  let liabilities: accounting_data = $state({
    data: undefined,
    loading: true,
    error: null,
  });

  async function get_financial_data(month: number) {
    let data = await Get_financial_report(0, month);
    return data;
  }

  async function get_income_statement_data() {
    income_statement.loading = true;
    income_statement.error = null;
    try {
      income_statement.data = await Get_accounting_data(
        0,
        month,
        "Income_statement",
      );
    } catch (exception) {
      income_statement.error = exception;
    }
    income_statement.loading = false;
  }

  async function get_assets_data() {
    assets.loading = true;
    assets.error = null;
    try {
      assets.data = await Get_accounting_data(0, month, "Assets");
    } catch (exception) {
      assets.error = exception;
    }
    assets.loading = false;
  }
  async function get_liabilities_data() {
    liabilities.loading = true;
    liabilities.error = null;
    try {
      liabilities.data = await Get_accounting_data(0, month, "Liabilities");
    } catch (exception) {
      liabilities.error = exception;
    }
    liabilities.loading = false;
  }

  let reports_tabs_state = $state({
    overview: true,
    finances: false,
    sales: false,
    employees: false,
    production: false,
  });
  function load_reports_overview() {
    reports_tabs_state.overview = true;
    reports_tabs_state.finances = false;
    reports_tabs_state.sales = false;
    reports_tabs_state.employees = false;
    reports_tabs_state.production = false;
  }

  function load_reports_finances() {
    reports_tabs_state.overview = false;
    reports_tabs_state.finances = true;
    reports_tabs_state.sales = false;
    reports_tabs_state.employees = false;
    reports_tabs_state.production = false;
    get_income_statement_data();
    get_liabilities_data();
    get_assets_data();
  }
  function load_reports_sales() {
    reports_tabs_state.overview = false;
    reports_tabs_state.finances = false;
    reports_tabs_state.sales = true;
    reports_tabs_state.employees = false;
    reports_tabs_state.production = false;
  }
  function load_reports_employees() {
    reports_tabs_state.overview = false;
    reports_tabs_state.finances = false;
    reports_tabs_state.sales = false;
    reports_tabs_state.employees = true;
    reports_tabs_state.production = false;
  }
  function load_reports_production() {
    reports_tabs_state.overview = false;
    reports_tabs_state.finances = false;
    reports_tabs_state.sales = false;
    reports_tabs_state.employees = false;
    reports_tabs_state.production = true;
  }
</script>

<div class="report_div">
  <div class="tab_selector">
    <button class="button tab_button" onclick={load_reports_overview}
      >Overview</button
    >
    <button class="button tab_button" onclick={load_reports_finances}
      >Finances</button
    >
    <button class="button tab_button" onclick={load_reports_sales}>Sales</button
    >
    <button class="button tab_button" onclick={load_reports_employees}
      >Employees</button
    >
    <button class="button tab_button" onclick={load_reports_production}
      >Production</button
    >
  </div>
  {#if reports_tabs_state.overview}
    <div>Overview is under construction</div>
  {/if}
  {#if reports_tabs_state.finances}
    <div
      class="grid_container"
      style="grid-template-columns: auto auto; grid-template-rows: auto auto;"
    >
      <div
        style="grid-row: 1; grid-column: 1 / span 2; display: flex; flex-direction: row"
      >
        <div class="grid_item" style="flex: 1 1 100%;;">
          A graph will go here
        </div>
        <div class="grid_item" style="flex: 1 1 100%;">
          A second graph will go here
        </div>
        <div class="grid_item" style="flex: 1 1 100%;">
          Another graph will go here
        </div>
      </div>
      <div class="grid_item" style="grid-column: 1 / span 2; grid-row: 2;">
        <h2>Financial report</h2>
        <Report get_report={get_financial_data}></Report>
      </div>
      <div class="grid_item" style="grid-column: 1 / span 2; grid-row: 3;">
        <h2>Income and loss statement</h2>
        {#if income_statement.loading}
          loading data...
        {:else if income_statement.error === null}
          <Table data={income_statement.data} table_type={"Income_statement"}
          ></Table>
        {:else}
          failed to load data: <br />
          {income_statement.error}
        {/if}
      </div>
      <div class="grid_item" style="grid-column: 1; grid-row: 4;">
        <h2>Assets</h2>
        {#if assets.loading}
          loading data...
        {:else if assets.error === null}
          <Table data={assets.data} table_type={"Assets"}></Table>
        {:else}
          failed to load data: <br />
          {assets.error}
        {/if}
      </div>
      <div class="grid_item" style="grid-column: 2; grid-row: 4;">
        <h2>Liabilities</h2>
        {#if liabilities.loading}
          loading data...
        {:else if liabilities.error === null}
          <Table data={liabilities.data} table_type={"Liabilities"}></Table>
        {:else}
          failed to load data: <br />
          {liabilities.error}
        {/if}
      </div>
    </div>
  {/if}
  {#if reports_tabs_state.sales}
    <div>Sales is under construction</div>
  {/if}
  {#if reports_tabs_state.employees}
    <div>Employees is under construction</div>
  {/if}
  {#if reports_tabs_state.production}
    <div>Production is under construction</div>
  {/if}
</div>

<style>
  .grid_item {
    background-color: rgba(0, 0, 0, 0.5);
    border-radius: 10px;
    margin: 10px;
    padding: 10px;
    display: flex;
    flex-direction: column;
  }

  .grid_container {
    display: grid;
  }

  .report_div {
    overflow-y: scroll;
    max-height: calc(100vh - 20px - 50px);
  }
</style>
