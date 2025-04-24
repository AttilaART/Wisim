<script lang="ts">
  import Report from "./Report.svelte";

  import { Get_sales_report } from "../wailsjs/go/main/App";

  import { month_counter } from "./store";
  import Chart from "./Chart.svelte";

  let month = $state(0);

  month_counter.subscribe((value) => {
    month = value;
  });

  async function get_sales_data(month: number) {
    let data = await Get_sales_report(0, month);
    return data;
  }
</script>

<div
  class="grid_container"
  style="grid-template-columns: auto auto; grid-template-rows: auto auto;"
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
    <h2>Sales report</h2>
    <Report
      get_report={get_sales_data}
      show_plus={false}
      colour_values={false}
      decimal_places={0}
      custom_rows={[
        {
          row_name: "Difference_to_previous_month",
          properties: {
            colour_values: true,
            show_plus: true,
            decimal_places: 0,
          },
        },
      ]}
    ></Report>
  </div>
</div>
