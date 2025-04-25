<script lang="ts">
  import Report from "./Report.svelte";

  import { Get_personelle_report } from "../wailsjs/go/main/App";

  import { month_counter } from "./store";
  import Chart from "./Chart.svelte";
  import { fade, fly } from "svelte/transition";

  let month = $state(0);
  let selected = $state("General");
  let select_width = $state(0);

  month_counter.subscribe((value) => {
    month = value;
  });

  async function get_general_personelle_data() {
    let data = await Get_personelle_report(0, month, "General");
    return data;
  }

  async function get_production_personelle_data() {
    let data = await Get_personelle_report(0, month, "Production");
    return data;
  }

  async function get_marketing_personelle_data() {
    let data = await Get_personelle_report(0, month, "Marketing");
    return data;
  }
</script>

<div
  class="grid_container"
  style="grid-template-columns: auto auto; grid-template-rows: auto auto;"
  out:fade={{ duration: 300 }}
  in:fly={{ duration: 300, delay: 300, y: -40 }}
>
  <div
    style="grid-row: 1; grid-column: 1 / span 2; display: flex; flex-direction: row"
  >
    <div class="grid_item" style="flex: 1 1 100%;;">
      <h3>Profit</h3>
      <Chart height={300}></Chart>
    </div>
    <div class="grid_item" style="flex: 1 1 100%;">
      A second graph will go here
    </div>
    <div class="grid_item" style="flex: 1 1 100%;">
      Another graph will go here
    </div>
  </div>
  <div class="grid_item" style="grid-column: 1 / span 2; grid-row: 2;">
    <h2>
      <div class="custom-select">
        <span class="custom-arrow" style="transform: translateX({select_width})"
        ></span>
        <select id="Select" bind:value={selected}>
          <option value="General">General</option>
          <option value="Marketing">Marketing</option>
          <option value="Production">Production</option>
        </select>
      </div>
      <!--<CustomSelect
        options={[
          {
            id: "General",
            text: "General",
            onselect: () => {
              selected = "General";
            },
          },
          {
            id: "Marketing",
            text: "Marketing",
            onselect: () => {
              selected = "Marketing";
            },
          },
          {
            id: "Production",
            text: "Production",
            onselect: () => {
              selected = "Production";
            },
          },
        ]}
      ></CustomSelect>-->
      Personelle report
    </h2>
    {#if selected == "General"}
      <Report
        get_report={get_general_personelle_data}
        colour_values={false}
        show_plus={false}
        decimal_places={0}
        custom_rows={[]}
      ></Report>
    {:else if selected == "Production"}
      <Report
        get_report={get_production_personelle_data}
        colour_values={false}
        show_plus={false}
        decimal_places={0}
        custom_rows={[]}
      ></Report>
    {:else}
      <Report
        get_report={get_marketing_personelle_data}
        colour_values={false}
        show_plus={false}
        decimal_places={0}
        custom_rows={[]}
      ></Report>
    {/if}
  </div>
</div>

<style>
  .custom-select {
    display: inline;
    position: relative;
  }

  select {
    -webkit-appearance: none;

    color: white;
    font-size: 25px;
    font-weight: bold;
    background: transparent;
    border-width: 0;

    font-family: "Nunito";
    font-style: normal;
    src:
      local(""),
      url("assets/fonts/nunito-v16-latin-regular.woff2") format("woff2");

    text-align-last: right;
    padding-left: 15px;
  }

  select:hover {
    text-decoration: underline;
  }

  .dummy {
    position: absolute;
    left: -10000px;
  }

  option {
    direction: rtl;
    text-align: right;
  }

  .custom-arrow {
    -webkit-appearance: unset;
    position: absolute;
    top: 0;
    left: 0;
    display: block;
    height: 100%;
    pointer-events: none;
  }

  .custom-arrow::before,
  .custom-arrow::after {
    content: "";
    --size: 7px;
    position: absolute;
    width: 0;
    height: 0;
    top: calc(50% - 5px);
  }

  .custom-arrow::before {
    border-left: var(--size) solid transparent;
    border-right: var(--size) solid transparent;
    border-top: var(--size) solid rgba(255, 255, 2255, 0.5);
    transform: translateY(calc(50% + 2px));
  }

  .custom-arrow::after {
    border-left: var(--size) solid transparent;
    border-right: var(--size) solid transparent;
    border-bottom: var(--size) solid rgba(255, 255, 2255, 0.5);
    transform: translateY(calc(-50% - 2px));
  }
</style>
