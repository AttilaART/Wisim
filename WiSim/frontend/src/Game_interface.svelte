<script lang="ts">
  import Popup from "./Popup.svelte";
  import Reports from "./Reports.svelte";
  import Sidebar from "./Sidebar.svelte";
  import {
    Get_simulation_step,
    Get_bank_balance,
    Get_current_stock,
    Trigger_simulation,
    Revert_simulation,
  } from "../wailsjs/go/main/App";

  import { month_counter } from "./store";

  let { return_function } = $props();
  let is_error: boolean = $state(false);
  let error: Error | undefined = $state(undefined);

  type Menu_state = {
    Dashboard: boolean;
    Decisions: boolean;
    Reports: boolean;
  };

  let is_simulating: boolean = $state(false);
  let menu_state: Menu_state = $state({
    Dashboard: true,
    Decisions: false,
    Reports: false,
  });

  let bank_balance: number | string = $state("loading...");
  let products_in_stock: number | string = $state("loading...");

  month_counter.subscribe((value) => {
    get_balance();
    get_products_in_storage();
  });

  async function get_month() {
    $month_counter = await Get_simulation_step();
  }

  async function get_balance() {
    try {
      bank_balance = await Get_bank_balance(0);
    } catch {
      bank_balance = "An error has occured";
    }
  }

  async function get_products_in_storage() {
    try {
      products_in_stock = await Get_current_stock(0);
    } catch {
      products_in_stock = "An error has occured";
    }
  }
  get_month();
  get_balance();
  get_products_in_storage();

  async function cancel_simulation() {
    is_simulating = false;
    try {
      Revert_simulation();
    } catch (exception) {
      is_error = true;
      error = exception;
    }
  }

  async function trigger_simulation() {
    is_simulating = true;
    let month_promise: Promise<number> = Trigger_simulation();
    month_promise.then(
      (value) => {
        month_counter.set(value);
      },
      (reason) => {
        (is_error = true), (error = reason);
        console.log("error due to simulation");
      },
    );

    setTimeout(() => {
      is_simulating = false;
    }, 1000);
    //try {
    //  month_temp = await Trigger_simulation();
    //} catch (exception) {
    //  is_error = true;
    //  error = exception;
    //}
    //is_simulating = false;
    //month_counter.set(month_temp);
  }

  function load_menus(mode: string) {
    if (mode == "Dashboard") {
      menu_state.Dashboard = true;
      menu_state.Decisions = false;
      menu_state.Reports = false;
    } else if (mode == "Decisions") {
      menu_state.Dashboard = false;
      menu_state.Decisions = true;
      menu_state.Reports = false;
    } else if (mode == "Reports") {
      menu_state.Dashboard = false;
      menu_state.Decisions = false;
      menu_state.Reports = true;
    }
  }
</script>

{#if is_simulating}
  <Popup
    button_data={{
      Text: "Cancel simulation",
      Style: "",
      Class: "popup_button",
      Onclick_function: cancel_simulation,
    }}
    content={'<div class="loader"></div>'}
  ></Popup>
{/if}
{#if is_error}
  <Popup
    button_data={{
      Text: "OK",
      Style: "",
      Class: "popup_button",
      Onclick_function: () => {
        is_error = false;
        error = undefined;
      },
    }}
    content={`<div> <h2>An Error has occured</h2> <br> ${error}</div>`}
  ></Popup>
{/if}
<div class="game_interface">
  <div
    style="display: flex; flex-direction: column; grid-column: 1; grid-row: 1 / span 2;"
  >
    <Sidebar
      expand={false}
      buttons={[
        {
          Text: "Return to main menu",
          Style: "",
          Onclick_function: () => {
            return_function();
          },
        },
      ]}
    ></Sidebar>

    <Sidebar
      expand={true}
      buttons={[
        {
          Text: "Dashboard",
          Style: "",
          Onclick_function: () => {
            load_menus("Dashboard");
          },
        },
        {
          Text: "Decisions",
          Style: "",
          Onclick_function: () => {
            load_menus("Decisions");
          },
        },
        {
          Text: "Reports",
          Style: "",
          Onclick_function: () => {
            load_menus("Reports");
          },
        },
        {
          Text: "Simulate",
          Style: "margin-top: auto",
          Onclick_function: () => {
            trigger_simulation();
          },
        },
      ]}
    ></Sidebar>
  </div>
  <div style="grid-column: 2; grid-row: 1; width: 100%;">
    {#if menu_state.Dashboard}
      <div>Dashboard is under construction</div>
    {/if}
    {#if menu_state.Decisions}
      <div>Decisions is under construction</div>
    {/if}
    {#if menu_state.Reports}
      <Reports></Reports>
    {/if}
  </div>
  <div class="bottom" style="grid-column: 2; grid-row: 2;">
    <div class="bottom_data">Bank balance: {bank_balance}</div>
    <div class="bottom_data">Products in storage: {products_in_stock}</div>
    <div class="bottom_data">Month {$month_counter}</div>
  </div>
</div>

<style>
  .game_interface {
    height: calc(100% - 20px);
    width: auto;
    max-width: 100%;
    max-height: calc(100% - 20px);
    padding: 10px;
    background-image: linear-gradient(
      to right,
      rgba(0, 0, 0, 0.5),
      rgba(0, 0, 0, 0)
    );
    display: grid;
    grid-template-columns: calc(200px + 20px) auto;
    grid-template-rows: auto 50px;
    overflow: hidden;
  }

  .bottom {
    background-color: rgba(0, 0, 0, 0.5);
    border-radius: 10px;
    margin: 10px;
    padding: 5px;
    display: flex;
  }

  .bottom_data {
    margin-left: 10px;
    margin-right: 10px;
    flex-grow: 1;
    text-align: center;
  }
</style>
