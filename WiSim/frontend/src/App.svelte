<script lang="ts">
  import background_image from "./assets/images/Temp_background_image.png";
  import {
    Get_accounting_data,
    Get_financial_report,
  } from "../wailsjs/go/main/App";
  import Table from "./Table.svelte";
  import Report from "./Report.svelte";

  let background_image_blur = 0;

  let mode = {
    main_menu: true,
    game_interface: false,
  };

  async function load_singleplayer() {
    background_image_blur = 4;
    mode.main_menu = false;
    mode.game_interface = true;
  }

  async function load_main_menu() {
    background_image_blur = 0;
    mode.main_menu = true;
    mode.game_interface = false;
  }

  let sidebar_state = {
    dashboard: false,
    decisons: false,
    reports: false,
  };

  let reports_tabs_state = {
    overview: true,
    finances: false,
    sales: false,
    employees: false,
    production: false,
  };

  let is_simulating = false;

  function load_dashboard() {
    sidebar_state.dashboard = true;
    sidebar_state.decisons = false;
    sidebar_state.reports = false;
  }
  function load_decisons() {
    sidebar_state.dashboard = false;
    sidebar_state.decisons = true;
    sidebar_state.reports = false;
  }

  function load_reports() {
    sidebar_state.dashboard = false;
    sidebar_state.decisons = false;
    sidebar_state.reports = true;
  }

  function trigger_simulation() {
    is_simulating = true;
    setTimeout(function () {
      is_simulating = false;
    }, 100000);
  }

  function cancel_simulation() {
    is_simulating = false;
  }

  async function get_income_statement_data() {
    let data = await Get_accounting_data(0, 0, "Income_statement");
    return data;
  }

  async function get_financial_data() {
    let data = await Get_financial_report(0, 0);
    return data;
  }

  let income_statement_promise = get_income_statement_data();

  let financial_report_promise = get_financial_data();

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

<main>
  <div
    style="position: absolute; height: 100vh; width: 100%; z-index: 0; overflow: hidden;"
  >
    <div
      style="background-image: url({background_image}); filter: blur({background_image_blur}px); position: absolute; width: 110%; height: 110%; left: -10px; top: -10px"
    ></div>
  </div>
  <div class="main_div" style="z-index: 20">
    {#if mode.main_menu}
      <div class="main_menu">
        <h1 style="padding: 0 8px;">WiSim</h1>
        <div
          class="sidebar"
          style="background-color: transparent; overflow: visible;"
        >
          <button class="button menu_button" on:click={load_singleplayer}>
            Singleplayer
          </button>
          <button class="button menu_button">Host game</button>
          <button class="button menu_button">Join game</button>
          <button class="button menu_button" style="margin-top: auto;"
            >Settings</button
          >
        </div>
      </div>
    {/if}
    {#if mode.game_interface}
      {#if is_simulating}
        <div class="popup">
          <div class="popup_content">
            <div class="loader"></div>
            <button
              class="button menu_button"
              style="text-align: center; position: absolute; left: 50%; bottom: 5px; transform: translate(-50%, -50%); margin: 0 auto;"
              on:click={cancel_simulation}>Cancel simulation</button
            >
          </div>
        </div>
      {/if}
      <div class="game_interface">
        <div style="display: flex; flex-direction: column; flex 0 0">
          <div class="sidebar" style="flex: 0 0 auto;">
            <button class="button sidebar_button" on:click={load_main_menu}
              >Return to main menu</button
            >
          </div>
          <div class="sidebar">
            <button class="button sidebar_button" on:click={load_dashboard}
              >Dashboard</button
            >
            <button class="button sidebar_button" on:click={load_decisons}
              >Decisions</button
            >
            <button class="button sidebar_button" on:click={load_reports}
              >Reports</button
            >
            <button
              class="button sidebar_button"
              style="margin-top: auto;"
              on:click={trigger_simulation}>Simulate</button
            >
          </div>
        </div>
        <div style="flex 1 1 100px; width: 100%;">
          {#if sidebar_state.dashboard}
            <div>Dashboard is under construction</div>
          {/if}
          {#if sidebar_state.decisons}
            <div>Decisions is under construction</div>
          {/if}
          {#if sidebar_state.reports}
            <div class="tab_selector">
              <button class="button tab_button" on:click={load_reports_overview}
                >Overview</button
              >
              <button class="button tab_button" on:click={load_reports_finances}
                >Finances</button
              >
              <button class="button tab_button" on:click={load_reports_sales}
                >Sales</button
              >
              <button
                class="button tab_button"
                on:click={load_reports_employees}>Employees</button
              >
              <button
                class="button tab_button"
                on:click={load_reports_production}>Production</button
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
                <div class="grid_item" style="grid-column: 1;">
                  A graph will go here
                </div>
                <div class="grid_item" style="grid-column: 2;">
                  <h2 style="text-align: left; padding-left: 10px;">
                    Profit and loss statement
                  </h2>
                  {#await income_statement_promise}
                    loading data...
                  {:then income_statement_data}
                    {console.log(income_statement_data)}
                    <Table data={income_statement_data}></Table>
                  {:catch error}
                    failed to load data: <br /> {error}
                  {/await}
                </div>
                <div
                  class="grid_item"
                  style="grid-column: 1 / span 2; grid-row: 2;"
                >
                  {#await financial_report_promise}
                    loading data...
                  {:then report_data}
                    <h2>Financial report</h2>
                    <Report Report={report_data}></Report>
                  {:catch}
                    failed to load report
                  {/await}
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
          {/if}
        </div>
      </div>
    {/if}
  </div>
</main>

<style>
  main {
    margin: 0;
    padding: 0;
  }
  * {
    box-sizing: border-box;
  }

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

  .loader {
    margin: auto auto;
    width: fit-content;
    font-weight: bold;
    font-family: monospace;
    font-size: 30px;
    clip-path: inset(0 3ch 0 0);
    animation: l4 1s steps(4) infinite;
  }
  .loader:before {
    content: "Loading...";
  }
  @keyframes l4 {
    to {
      clip-path: inset(0 -1ch 0 0);
    }
  }
  .popup {
    position: absolute;
    width: 100%;
    height: 100%;
    align-content: center;
    background-color: rgba(0, 0, 0, 0.2);
  }

  .popup_content {
    position: relative;
    height: 45%;
    width: 45%;
    margin: auto auto;
    background-color: rgba(20, 20, 20, 1);
    border-radius: 10px;
    align-content: center;
    box-shadow: 5px 5px 15px rgba(0, 0, 0, 0.5);
  }

  .tab_selector {
    display: flex;
    margin: 10px;
    background-color: rgba(0, 0, 0, 0.5);
    overflow: hidden;
    border-radius: 10px;
    height: fit-content;
  }
  .game_interface {
    height: 100%;
    max-height: 100%;
    width: 100%;
    padding: 10px;
    display: flex;
    background-image: linear-gradient(
      to right,
      rgba(0, 0, 0, 0.5),
      rgba(0, 0, 0, 0)
    );
    overflow-y: scroll;
  }

  .sidebar {
    width: fit-content;
    display: flex;
    flex-direction: column;
    flex: 1 1 100%;
    flex-wrap: wrap;
    background-color: rgba(0, 0, 0, 0.5);
    border-radius: 10px;
    margin: 10px 10px 10px 10px;
    overflow: hidden;
  }

  .main_div {
    position: absolute;
    background-size: cover;
    background-position: center center;
    background-repeat: no-repeat;
    height: 100vh;
    width: 100%;
    padding: 0px;
    margin: 0px;
  }

  .main_menu {
    display: flex;
    flex-direction: column;
    background-image: linear-gradient(
      to right,
      rgba(0, 0, 0, 0.5),
      rgba(0, 0, 0, 0)
    );
    height: 100%;
    width: 50%;
    text-align: left;
    padding: 10px;
  }

  .button {
    background-color: rgba(0, 0, 0, 0);
    text-align: left;
    color: rgba(255, 255, 255, 1);
    width: 120px;
    height: 30px;
    border-radius: 3px;
    border: none;
    padding: 0 10px;
    cursor: pointer;
  }

  .button:hover {
    background-color: rgba(255, 255, 255, 0.5);
  }

  .button:active {
    background-color: rgba(255, 255, 255, 1);
    color: #000;
  }

  .menu_button.button {
    width: 200px;
    height: fit-content;
    padding: 10px;
    margin: 0px 0px 0px 0px;
    border-radius: 3;
  }

  .sidebar_button.button {
    width: 200px;
    height: fit-content;
    padding: 10px;
    margin: 0px 0px 0px 0px;
    border-radius: 0;
  }

  .tab_button.button {
    text-align: center;
    height: fit-content;
    padding: 10px;
    margin: 0px 0px 0px 0px;
    border-radius: 0;
    flex: 1 1 fit-content;
  }
</style>
