<script lang="ts">
  import Report from "./Report.svelte";

  import {
    Get_product_statistics,
    Get_sales_statistics,
    Get_marketing_statistics,
  } from "../wailsjs/go/main/App";

  import { month_counter } from "./store";
  import Chart from "./Chart.svelte";
  import { fade, fly } from "svelte/transition";

  let month = $state(0);

  month_counter.subscribe((value) => {
    month = value;
  });

  async function get_marketing_data(month: number) {
    let data = await Get_marketing_statistics(0, month);
    delete data.Product;
    return data;
  }

  async function get_company_sales_statistics(month: number) {
    let data = await Get_sales_statistics(0, month);
    return data;
  }
  async function get_market_sales_statistics(month: number) {
    let data = await Get_sales_statistics(-1, month);
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
    <div class="grid_item" style="flex: 1 2 100%;;">
      <h3>Profit</h3>
      <Chart height={300}></Chart>
    </div>
    <div class="grid_item" style="flex: 1 1 100%;">
      product statistics will go here
    </div>
  </div>
  <div class="grid_item" style="grid-column: 1 / span 2; grid-row: 2;">
    <h2>Marketing statistics</h2>
    <Report
      get_report={get_marketing_data}
      show_plus={false}
      colour_values={false}
      decimal_places={0}
      custom_rows={[]}
    ></Report>
  </div>
  <div class="grid_item" style="grid-column: 1 / span 2; grid-row: 3;">
    <h2>Company_sales_statistics</h2>
    <Report
      get_report={get_company_sales_statistics}
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
  <div class="grid_item" style="grid-column: 1 / span 2; grid-row: 4;">
    <h2>Market_sales_statistics</h2>
    <Report
      get_report={get_market_sales_statistics}
      show_plus={false}
      colour_values={false}
      decimal_places={0}
      custom_rows={[]}
    ></Report>
  </div>
</div>
