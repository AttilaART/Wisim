<script lang="ts">
  import { simulation } from "../wailsjs/go/models";

  interface Assets {
    Current_assets: simulation.FinanceReportEntry[];
    NonCurrent_assets: simulation.FinanceReportEntry[];
  }

  interface Liablilities {
    Liabilities: simulation.FinanceReportEntry[];
    Private_equity: simulation.FinanceReportEntry[];
  }

  type row = {
    data: (string | number | undefined)[];
    isAssetsHeader?: boolean;
    isLiabilitesHeader?: boolean;
  };

  let assets_length: number = 1;
  let liabilities_length: number = 1;

  function add_row(
    table: row[],
    key: string,
    value: number | undefined,
    side: boolean,
    isHeader: boolean,
  ) {
    if (side == false) {
      if (assets_length >= table.length) {
        table.push({
          data: [key, value, undefined, undefined],
          isAssetsHeader: isHeader,
        });
      } else {
        table[assets_length - 1].data[0] = key;
        table[assets_length - 1].data[1] = value;
        table[liabilities_length - 1].isAssetsHeader = isHeader;
      }

      assets_length += 1;
    } else {
      if (liabilities_length >= table.length) {
        table.push({
          data: [undefined, undefined, key, value],
          isLiabilitesHeader: isHeader,
        });
      } else {
        table[liabilities_length - 1].data[2] = key;
        table[liabilities_length - 1].data[3] = value;
        table[liabilities_length - 1].isLiabilitesHeader = isHeader;
      }
      liabilities_length += 1;
    }
  }

  function make_table(assets: Assets, equity: Liablilities): row[] {
    let table: row[] = [];

    let total_assets: number = 0;
    let total_liabilites: number = 0;
    add_row(table, "<h4>Current Assets</h4>", undefined, false, true);

    for (let row of assets.Current_assets) {
      add_row(table, row.Name, row.Value, false, false);
      total_assets += row.Value;
    }

    add_row(table, "<h4>Non-Current Assets</h4>", undefined, false, true);

    for (let row of assets.NonCurrent_assets) {
      add_row(table, row.Name, row.Value, false, false);
      total_assets += row.Value;
    }

    add_row(table, "<h4>Liabilities</h4>", undefined, true, true);

    for (let row of equity.Liabilities) {
      add_row(table, row.Name, row.Value, true, false);
      total_liabilites += row.Value;
    }

    add_row(table, "<h4>Private Equity</h4>", undefined, true, true);

    for (let row of equity.Private_equity) {
      add_row(table, row.Name, row.Value, true, false);
      total_liabilites += row.Value;
    }

    table.push({
      data: [
        "<h4>Total Assets</h4>",
        total_assets,
        "<h4>Total Liabilities and Private Equity</h4>",
        total_liabilites,
      ],
      isAssetsHeader: true,
      isLiabilitesHeader: true,
    });

    return table;
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
  let table: row[] = $derived(make_table(assets, liabilities));
</script>

<table style="width: 100%;">
  <thead>
    <tr style="border-bottom: var(--border-thin);">
      <th style="width: 35%;"><h2>Assets</h2></th>
      <th style="border-right: var(--border-thin); width: 15%;"></th>
      <th></th>
      <th style="text-align: right;"><h2>Liabilities</h2></th>
    </tr>
  </thead>
  <tbody>
    {#each table as row, i}
      <tr
        style={i + 1 == table.length
          ? "border-top: var(--border-thin); border-bottom: var(--border-thin);"
          : ""}
      >
        <td style="padding-left: {row.isAssetsHeader ? 10 : 15}px;"
          >{@html row.data[0]}</td
        >
        <td style="text-align: right; border-right: var(--border-thin);"
          >{@html row.data[1]}</td
        >
        <td style="padding-left: {row.isLiabilitesHeader ? 10 : 15}px;"
          >{@html row.data[2]}</td
        >
        <td style="text-align: right;">{@html row.data[3]}</td>
      </tr>
    {/each}
  </tbody>
</table>

<style>
  table {
    text-align: left;
    border-collapse: collapse;
  }

  h2 {
    padding-left: 5px;
    padding-right: 5px;
  }

  td {
    padding-right: 10px;
  }
</style>
