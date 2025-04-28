<script lang="ts">
  import Popup from "./Popup.svelte";
  import Reports from "./Reports.svelte";
  import Sidebar from "./Sidebar.svelte";
  import { format_number, capitalise_first_letter } from "./helper";
  import {
    Get_simulation_step,
    Get_bank_balance,
    Get_current_stock,
    Trigger_simulation,
    Revert_simulation,
    Submit_decisions,
  } from "../wailsjs/go/main/App";
  import { fade, fly, slide } from "svelte/transition";

  import { month_counter, company, decisions } from "./store";
  import Decisions from "./Decisions.svelte";
  import { simulation } from "../wailsjs/go/models";

  let { return_function } = $props();
  let is_error: boolean = $state(false);
  let error: Error | undefined = $state(undefined);
  let try_cancel_sim: boolean = false;
  let reports_tab__greyed_out: number = $state(0);
  let bottom_data_height: number = $state();
  $decisions = new simulation.Decisions();
  $company = 0;

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

    if (value < 0) {
      reports_tab__greyed_out = 0;
    } else {
      reports_tab__greyed_out = 1;
    }
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
    try_cancel_sim = true;
  }

  async function trigger_simulation() {
    await Submit_decisions(0, $decisions);
    is_simulating = true;
    let month_promise: Promise<number> = Trigger_simulation(false);
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

      if (try_cancel_sim) {
        let month_promise = Revert_simulation();
        month_promise.then(
          (value) => {
            $month_counter = value;
          },
          (reason) => {
            is_error = true;
            error = reason;
          },
        );
      }
      try_cancel_sim = false;
    }, 1000);

    try_cancel_sim = false;

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
        error = null;
      },
    }}
    content={`<div> <h2>An Error has occured</h2> <br> ${capitalise_first_letter(error.toString())}</div>`}
  ></Popup>
{/if}
<div
  class="game_interface"
  style="grid-template-rows: auto {bottom_data_height + 22}px;"
  in:fade={{ duration: 300, delay: 300 }}
  out:fade={{ duration: 300 }}
>
  <div
    style="display: flex; flex-direction: column; grid-column: 1; grid-row: 1 / span 2;"
  >
    <Sidebar
      expand={false}
      buttons={[
        {
          Text: "Return to main menu",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            return_function();
          },
        },
      ]}
    ></Sidebar>

    <Sidebar
      expand={true}
      keep_pressed={true}
      buttons={[
        {
          Text: "Dashboard",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            load_menus("Dashboard");
          },
        },
        {
          Text: "Decisions",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            load_menus("Decisions");
          },
        },
        {
          Text: "Reports",
          Style: "",
          Show: reports_tab__greyed_out,
          Onclick_function: () => {
            //if (reports_tab__greyed_out == 1) {
            load_menus("Reports");
            //}
          },
        },
        {
          Text: "Simulate",
          Style: "margin-top: auto",
          Show: 1,
          Onclick_function: () => {
            trigger_simulation();
          },
          dont_keep_pressed: true,
        },
      ]}
    ></Sidebar>
  </div>
  <div style="grid-column: 2; grid-row: 1; width: 100%;">
    {#if menu_state.Dashboard}
      <div
        class="report_div"
        out:fade={{ duration: 300 }}
        in:fly={{ duration: 300, delay: 300, y: -40 }}
      >
        <div>Dashboard is under construction</div>
      </div>
    {/if}
    {#if menu_state.Decisions}
      <div
        class="report_div"
        out:fade={{ duration: 300 }}
        in:fly={{ duration: 300, delay: 300, y: -40 }}
      >
        <Decisions></Decisions>
      </div>
    {/if}
    {#if menu_state.Reports}
      <div
        class="report_div"
        out:fade={{ duration: 300 }}
        in:fly={{ duration: 300, delay: 300, y: -40 }}
      >
        <Reports></Reports>
      </div>
    {/if}
  </div>
  <div
    class="bottom"
    style="grid-column: 2; grid-row: 2; height: fit-content;"
    bind:clientHeight={bottom_data_height}
  >
    <div class="bottom_data">
      Bank balance: {format_number(bank_balance, true)}
    </div>
    <div class="bottom_data">
      Products in storage: {format_number(products_in_stock, false, 0)}
    </div>
    <div class="bottom_data">
      Month / Year: {($month_counter + 1) % 13} / {Math.floor(
        $month_counter / 12,
      )}
    </div>
  </div>
</div>

<style>
  .game_interface {
    height: calc(100% - 20px);
    width: calc(100% - 20px);
    max-height: calc(100% - 20px);
    padding: 10px;
    display: grid;
    grid-template-columns: calc(200px + 20px) auto;
    grid-template-rows: auto 55px;
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
