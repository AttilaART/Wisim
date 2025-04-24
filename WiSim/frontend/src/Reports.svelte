<script lang="ts">
  import Financial_report from "./Financial_report.svelte";
  import PersonelleReport from "./PersonelleReport.svelte";
  import SalesReport from "./Sales_report.svelte";

  import { month_counter } from "./store";

  let reports_tabs_state = $state({
    overview: true,
    finances: false,
    sales: false,
    personelle: false,
    production: false,
  });
  function load_reports_overview() {
    reports_tabs_state.overview = true;
    reports_tabs_state.finances = false;
    reports_tabs_state.sales = false;
    reports_tabs_state.personelle = false;
    reports_tabs_state.production = false;
  }

  function load_reports_finances() {
    reports_tabs_state.overview = false;
    reports_tabs_state.finances = true;
    reports_tabs_state.sales = false;
    reports_tabs_state.personelle = false;
    reports_tabs_state.production = false;
  }
  function load_reports_sales() {
    reports_tabs_state.overview = false;
    reports_tabs_state.finances = false;
    reports_tabs_state.sales = true;
    reports_tabs_state.personelle = false;
    reports_tabs_state.production = false;
  }
  function load_reports_personelle() {
    reports_tabs_state.overview = false;
    reports_tabs_state.finances = false;
    reports_tabs_state.sales = false;
    reports_tabs_state.personelle = true;
    reports_tabs_state.production = false;
  }
  function load_reports_production() {
    reports_tabs_state.overview = false;
    reports_tabs_state.finances = false;
    reports_tabs_state.sales = false;
    reports_tabs_state.personelle = false;
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
      <button class="button tab_button" onclick={load_reports_personelle}
        >personelle</button
      >
      <button class="button tab_button" onclick={load_reports_production}
        >Production</button
      >
    </div>
    {#if reports_tabs_state.overview}
      <div>Overview is under construction</div>
    {:else if reports_tabs_state.finances}
      <Financial_report></Financial_report>
    {:else if reports_tabs_state.sales}
      <SalesReport></SalesReport>
    {:else if reports_tabs_state.personelle}
      <PersonelleReport></PersonelleReport>
    {:else if reports_tabs_state.production}
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
      <p>
        Here you will be able to see the reports generated for the last few
        months
      </p>
    </div>
  </div>
{/if}

<style>
</style>
