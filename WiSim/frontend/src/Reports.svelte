<script lang="ts">
  import Financial_report from "./Financial_report.svelte";
  import PersonelleReport from "./PersonelleReport.svelte";
  import ProductionReport from "./ProductionReport.svelte";
  import SalesReport from "./Sales_report.svelte";
  import Sidebar from "./Sidebar.svelte";
  import { fade, fly } from "svelte/transition";

  import { month_counter } from "./store";

  let reports_tabs_state = $state({
    overview: true,
    finances: false,
    sales: false,
    personelle: false,
    production: false,
  });

  let page_width: number = $state();

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

<svelte:window bind:outerWidth={page_width} />

{#if $month_counter >= 0}
  <div
    out:fade={{ duration: 300 }}
    in:fly={{ duration: 300, delay: 300, y: -40 }}
  >
    <Sidebar
      expand={true}
      keep_pressed={true}
      horisontal={true}
      buttons={[
        {
          Text: "Overview",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            load_reports_overview();
          },
        },
        {
          Text: "Finances",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            load_reports_finances();
          },
        },
        {
          Text: "Sales",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            load_reports_sales();
          },
        },
        {
          Text: "Personelle",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            load_reports_personelle();
          },
        },
        {
          Text: "Production",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            load_reports_production();
          },
        },
      ]}
    ></Sidebar>
    {#key page_width}
      {#if reports_tabs_state.overview}
        <div
          out:fade={{ duration: 300 }}
          in:fly={{ duration: 300, delay: 300, y: -40 }}
        >
          Overview is under construction
        </div>
      {:else if reports_tabs_state.finances}
        <Financial_report></Financial_report>
      {:else if reports_tabs_state.sales}
        <SalesReport></SalesReport>
      {:else if reports_tabs_state.personelle}
        <PersonelleReport></PersonelleReport>
      {:else if reports_tabs_state.production}
        <div>
          <ProductionReport></ProductionReport>
        </div>
      {/if}
    {/key}
  </div>
{:else}
  <div>
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
  </div>
{/if}

<style>
</style>
