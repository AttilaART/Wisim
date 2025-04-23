<script lang="ts">
  import Financial_report from "./Financial_report.svelte";

  import { month_counter } from "./store";

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

{#if $month_counter >= 0}
  <div class="report_div">
    <div class="tab_selector">
      <button class="button tab_button" onclick={load_reports_overview}
        >Overview</button
      >
      <button class="button tab_button" onclick={load_reports_finances}
        >Finances</button
      >
      <button class="button tab_button" onclick={load_reports_sales}
        >Sales</button
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
      <Financial_report></Financial_report>
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
{:else}
  <div
    class="grid_container"
    style="grid-template-columns: auto; grid-template-rows: auto;"
  >
    <div class="grid_item">
      <h2>Reports</h2>
      <p>This page will have contents after the first step is simlated.</p>
    </div>
  </div>
{/if}

<style>
</style>
