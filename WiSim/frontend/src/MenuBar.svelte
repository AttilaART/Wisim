<script lang="ts">
  import { simulation } from "../wailsjs/go/models";
  import { new_window } from "./Game_interface.svelte";
  import { format_currency, format_number } from "./helper.svelte";
  import { company, external_factors, latest_reports } from "./store.svelte";

  function get_total_number_of_products_sold(
    reports: simulation.Report[],
  ): number {
    let total: number = 0;
    for (let rep in reports) {
      console.log(company.Reports[rep]);
      total +=
        company.Reports[rep].Sales_report.Company_sales_statistics
          .Products_sold;
    }
    return total;
  }
</script>

<div style="width: 100%; position: fixed; bottom: 0; right: 0; z-index: 99999;">
  <div id="top-section">
    <button
      id="balance"
      onclick={() => {
        new_window("finances");
      }}
    >
      <h1
        style="position: absolute; font-size: 3.5rem; bottom: 0; transform: translateY(40%);"
      >
        $
      </h1>
      <div style="padding-left: 3rem;">
        <small style="padding-left: 0.3rem;">Balance</small> <br />
        <h3 style={company.Balance < 0 ? "color: var(--red);" : ""}>
          {format_number(company.Balance, false, 0)}
        </h3>
      </div>
    </button>
    <button
      class="item"
      onclick={() => {
        new_window("marketing");
      }}>🛒</button
    >
    <button
      class="item"
      onclick={() => {
        new_window("employees");
      }}>🧑‍🏭</button
    >
    <button
      class="item"
      onclick={() => {
        new_window("production");
      }}>🏭</button
    >
    <div class="spacer"></div>
    <button class="item" onclick={() => {}}>👨‍🔬</button>
    <div class="spacer"></div>
    <button class="item" onclick={() => {}}>🏢</button>
    <button class="item" onclick={() => {}}>✉️</button>
  </div>
  <div id="bottom-section">
    <small style="padding-top: 0.1rem;">Month {external_factors.Month}</small>
    <small style="padding-top: 0.1rem;"
      >Total Products Sold: {get_total_number_of_products_sold(
        company.Reports,
      )}</small
    >
    <small style="padding-top: 0.1rem;"
      >Sales last month: {latest_reports.Sales_report
        ? format_number(
            latest_reports.Sales_report.Company_sales_statistics.Products_sold,
            false,
            0,
          )
        : "NaN"}
      {#if latest_reports.Sales_report}
        {@render difference_triangle(
          latest_reports.Sales_report.Company_sales_statistics
            .Difference_to_previous_month,
        )}
      {/if}
    </small>
    <small style="padding-top: 0.1rem;"
      >Market share: {latest_reports.Sales_report
        ? format_number(
            latest_reports.Sales_report.Company_sales_statistics.Market_share *
              100,
            false,
            0,
          ) + "%"
        : "NaN"}
      {#if company.Reports && company.Reports.length >= 2}
        {@render difference_triangle(
          latest_reports.Sales_report.Company_sales_statistics.Market_share -
            company.Reports[company.Reports.length - 2].Sales_report
              .Company_sales_statistics.Market_share,
        )}
      {/if}
    </small>
    <small style="padding-top: 0.1rem;"
      >Cashflow: {latest_reports.Financial_Report
        ? format_currency(
            latest_reports.Financial_Report.Non_operating_expenses.Cashflow,
            0,
            true,
          )
        : "NaN"}
      {#if company.Reports && company.Reports.length >= 2}
        {@render difference_triangle(
          latest_reports.Financial_Report.Non_operating_expenses.Cashflow -
            company.Reports[company.Reports.length - 2].Financial_Report
              .Non_operating_expenses.Cashflow,
        )}
      {/if}
    </small>
  </div>
</div>

{#snippet difference_triangle(value: number, invert_color?: boolean)}
  {#if value !== undefined}
    {#if value > 0}
      <span style="color: {invert_color ? 'var(--red)' : 'var(--green)'};"
        >▲</span
      >
    {:else if value == 0}
      <span>—</span>
    {:else}
      <span style="color: {invert_color ? 'var(--green)' : 'var(--red)'};"
        >▼</span
      >
    {/if}
  {/if}
{/snippet}

<style>
  * {
    --side-padding: 16.5rem;
    --balance-width: 16.5rem;
  }

  #top-section,
  #bottom-section {
    background-color: var(--second-color);
    border-top: var(--border-thin);
    display: flex;
    gap: 1.5rem;
    width: calc(100% - var(--side-padding) - var(--balance-width));
    padding: 0.1rem var(--side-padding) 0.1rem var(--balance-width);
  }

  #top-section {
    height: 2rem;
    overflow: show;
    justify-content: center;
  }

  #bottom-section {
    overflow: scroll;
    flex-direction: row-reverse;
    vertical-align: middle;
    height: 1.5rem;
  }

  #balance {
    position: fixed;
    text-align: left;
    left: 0;
    bottom: 0;
    height: 5rem;
    width: var(--balance-width);
    border-radius: 0 3rem 0 0;
    border: none;
    border-top: var(--border-thin);
    border-right: var(--border-thin);
  }

  .item {
    position: relative;
    border-radius: 50%;
    overflow: hidden;
    aspect-ratio: 1/1;
    width: 3rem;
    height: 3rem;
    transform: translateY(-50%);
    border: var(--border-thin);
    font-size: 1.5rem;
    text-align: center;
  }

  .spacer {
    width: 3rem;
  }
</style>
