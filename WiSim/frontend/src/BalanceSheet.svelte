<script lang="ts">
  import { number } from "echarts";
  import { simulation } from "../wailsjs/go/models";
  import { format_currency } from "./helper";

  interface Assets {
    Current_assets: simulation.FinanceReportEntry[];
    NonCurrent_assets: simulation.FinanceReportEntry[];
  }

  interface Liablilities {
    Liabilities: simulation.FinanceReportEntry[];
    Private_equity: simulation.FinanceReportEntry[];
  }

  function total(data: Assets | Liablilities): number {
    let total: number = 0;
    for (let section of Object.keys(data)) {
      for (let row of data[section]) {
        total += row.Value;
      }
    }

    return total;
  }

  let assets: Assets = $state();
  let liabilities: Liablilities = $state();

  assets = {
    Current_assets: [
      { Name: "Test1", Group: 1, Info: "", Value: 100, Cash_cost: true },
    ],
    NonCurrent_assets: [
      { Name: "Test1", Group: 1, Info: "", Value: 100, Cash_cost: true },
      { Name: "Test1", Group: 1, Info: "", Value: 100, Cash_cost: true },
      { Name: "Test1", Group: 1, Info: "", Value: 100, Cash_cost: true },
    ],
  };
  liabilities = {
    Liabilities: [
      { Name: "Test1", Group: 1, Info: "", Value: 100, Cash_cost: true },
      { Name: "Test1", Group: 1, Info: "", Value: 100, Cash_cost: true },
    ],
    Private_equity: [
      { Name: "Test1", Group: 1, Info: "", Value: 100, Cash_cost: true },
    ],
  };
</script>

<div style="display: flex; width: 100%;">
  {@render semi_table(assets, "Assets")}
  {@render semi_table(liabilities, "Liabilities", true)}
</div>
<div style="display: flex;">
  {@render semi_footer(total(assets), "Total: ")}
  {@render semi_footer(total(liabilities), "", true)}
</div>

{#snippet semi_table(
  data: Assets | Liablilities,
  title: string,
  on_right?: boolean,
)}
  <table>
    <thead>
      <tr class="bottom {on_right ? '' : 'right'}">
        <th><h2>{title}</h2></th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      {#each Object.keys(liabilities) as section}
        <tr class="bottom {on_right ? '' : 'right'}">
          <td><h3>{section}</h3></td>
          <td></td>
        </tr>
        {#each liabilities[section] as entry}
          <tr class={on_right ? "" : "right"}>
            <td> <p>{entry.Name}</p></td>
            <td style="text-align: right;"
              ><p>{format_currency(entry.Value)}</p></td
            >
          </tr>
        {/each}
      {/each}
    </tbody>
  </table>
{/snippet}

{#snippet semi_footer(total: number, title: string, on_right?: boolean)}
  <table class="bottom top {on_right ? '' : 'right'}">
    <tfoot>
      <tr>
        <td><h3>{title == "" ? "\u200B" : title}</h3></td>
        <td style="text-align: right;"> <p>{format_currency(total)}</p> </td>
      </tr>
    </tfoot>
  </table>
{/snippet}

<style>
  table {
    text-align: left;
    border-collapse: collapse;
    flex: 1 1 50%;
    width: 100%;
  }

  h2,
  h3 {
    padding-left: 5px;
    padding-right: 5px;
  }

  p {
    padding-right: 10px;
    padding-left: 10px;
  }

  .bottom {
    border-bottom: var(--border-thin);
  }
  .right {
    border-right: var(--border-thin);
  }

  .top {
    border-top: var(--border-thin);
  }
</style>
