<script lang="ts">
  import Report from "./Report.svelte";

  import { Get_production_report } from "../wailsjs/go/main/App";

  import { month_counter } from "./store";
  import Chart from "./Chart.svelte";
  import { fade, fly } from "svelte/transition";

  let month = $state(0);

  month_counter.subscribe((value) => {
    month = value;
  });

  async function get_production_data() {
    let data = await Get_production_report(0, month);
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
    <h2>Production report</h2>
    <Report
      get_report={get_production_data}
      colour_values={false}
      show_plus={false}
      decimal_places={0}
      custom_rows={[
        {
          row_name: "Machines_purchased",
          properties: {
            show_plus: false,
            colour_values: true,
            decimal_places: 0,
          },
        },
        {
          row_name: "Machines_sold",
          properties: {
            show_plus: false,
            colour_values: true,
            decimal_places: 0,
          },
        },
        {
          row_name: "Worker_surplus",
          properties: {
            show_plus: false,
            colour_values: true,
            decimal_places: 0,
            custom_name: "Worker surplus (relative to machines)",
            info: "A negative surplus means you have a deficit",
          },
        },
      ]}
    ></Report>
  </div>
</div>
