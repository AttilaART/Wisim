<script lang="ts">
  import triangle from "./assets/images/Triangle(right)(white).png";
  import { int } from "../wailsjs/go/models";

  let red = "rgb(255, 128, 128)";
  let green = "rgb(128, 255, 128)";
  let currency = "CHF";

  let { data } = $props();
  type FinanceReportEntry = {
    Name: string;
    Group: number;
    Info: string;
    Cash_cost: number;
    Value: number;
  };

  type FoldedFinanceReportEntry = {
    index: number;
    entry: FinanceReportEntry;
    sub_entries: FinanceReportEntry[];
  };

  let make_table = parse_accounting_data(data);
  async function parse_accounting_data(
    data: FinanceReportEntry[],
  ): Promise<{ entries: FoldedFinanceReportEntry[]; total_cashflow: number }> {
    let groups: number[] = [];
    let has_folded_entry: boolean[] = [];
    for (let r of data) {
      if (!groups.includes(r.Group)) {
        groups.push(r.Group);
      }
    }
    let folded_entries: FoldedFinanceReportEntry[] = [];

    for (let g of groups) {
      folded_entries.push({
        index: 0,
        entry: { Name: "", Group: g, Info: "", Cash_cost: -1, Value: 0 },
        sub_entries: [],
      });
    }

    function capitalise_first_letter(val: string): string {
      return String(val).charAt(0).toUpperCase() + String(val).slice(1);
    }

    for (let r of data) {
      for (let i in groups) {
        if (r.Group == groups[i]) {
          if (has_folded_entry[i]) {
            folded_entries[i].sub_entries.push(r);
          } else {
            folded_entries[i] = {
              index: Number(i),
              entry: {
                Name: capitalise_first_letter(int.int[groups[i]]).replace(
                  "_",
                  " ",
                ),
                Group: groups[i],
                Info: r.Info,
                Cash_cost: r.Cash_cost,
                Value: 0,
              },
              sub_entries: [],
            };

            has_folded_entry[i] = true;
            folded_entries[i].sub_entries.push(r);
          }
        }
      }
    }

    var total_cashflow = 0;
    for (let i in folded_entries) {
      let value: number = 0;
      folded_entries[i].sub_entries = folded_entries[i].sub_entries.toSorted(
        (a, b) => {
          if (a.Value > b.Value) {
            return -1;
          } else if (a.Value == b.Value) {
            return 0;
          }
          return 1;
        },
      );
      for (let e of folded_entries[i].sub_entries) {
        value += e.Value;
        if (e.Cash_cost == 1) {
          total_cashflow += e.Value;
        }

        if (
          (e.Cash_cost == 1 && folded_entries[i].entry.Cash_cost == 0) ||
          (e.Cash_cost == 0 && folded_entries[i].entry.Cash_cost == 1)
        ) {
          folded_entries[i].entry.Cash_cost = 0;
        }
      }
      folded_entries[i].entry.Value = value;
    }

    folded_entries = folded_entries.toSorted((a, b) => {
      if (a.entry.Value > b.entry.Value) {
        return -1;
      } else if (a.entry.Value == b.entry.Value) {
        return 0;
      } else {
        return 1;
      }
    });

    return { entries: folded_entries, total_cashflow: total_cashflow };
  }

  function format_number(num: number): string {
    return num.toLocaleString("de-CH", {
      maximumFractionDigits: 2,
      minimumFractionDigits: 2,
    });
  }

  function add_the_third_bool(i: number | boolean): string {
    if (i === true || i === false) {
      return i.toString();
    } else {
      return "---";
    }
  }

  let folded_entries_state: boolean[] = $state([]);
</script>

{#await make_table}
  loading table...
{:then result}
  <table style="text-align: left; overflow-y: scroll;">
    <thead>
      <tr>
        <th></th>
        <th>Name</th>
        <th>Info</th>
        <th>Group</th>
        <th>Cash cost?</th>
        <th style="text-align: right;">Value ({currency})</th>
      </tr>
    </thead>
    <tbody>
      {#each result.entries as folded_entry}
        {@render folded_table_row(folded_entry)}
      {/each}
    </tbody>
    <tfoot>
      <tr>
        <td></td>
        <td>TOTAL CASHFLOW</td>
        <td></td>
        <td></td>
        <td></td>
        {#if result.total_cashflow > 0}
          <td class="value_field" style="color: {green};"
            >+{format_number(result.total_cashflow)}</td
          >
        {:else if result.total_cashflow == 0}
          <td class="value_field">{format_number(result.total_cashflow)}</td>
        {:else}
          <td class="value_field" style="color: {red};"
            >{format_number(result.total_cashflow)}</td
          >
        {/if}
      </tr>
    </tfoot>
  </table>
{:catch error}
  an error occured: <br /> {error}
{/await}

{#snippet folded_table_row(folded_entry: FoldedFinanceReportEntry)}
  <tr>
    <td>
      <div hidden>
        {() => {
          folded_entries_state[folded_entry.index] = false;
        }}
      </div>
      <button
        class="fold_button"
        onclick={() => {
          folded_entries_state[folded_entry.index] =
            !folded_entries_state[folded_entry.index];

          folded_entries_state = folded_entries_state;
        }}
        ><img
          src={triangle}
          height="10px"
          width="10px"
          alt="Triangle for folding entries"
          style="transform: rotate({Number(
            folded_entries_state[folded_entry.index],
          ) * 90}deg)"
        /></button
      ></td
    >
    <td>{folded_entry.entry.Name}</td>
    <td>{folded_entry.entry.Info}</td>
    <td>{folded_entry.entry.Group}</td>
    <td>{add_the_third_bool(folded_entry.entry.Cash_cost)}</td>
    {#if folded_entry.entry.Value > 0}
      <td class="value_field" style="color: {green};"
        >+{format_number(folded_entry.entry.Value)}</td
      >
    {:else if folded_entry.entry.Value == 0}
      <td class="value_field">{format_number(folded_entry.entry.Value)}</td>
    {:else}
      <td class="value_field" style="color: {red};"
        >{format_number(folded_entry.entry.Value)}</td
      >
    {/if}
  </tr>
  {#if folded_entries_state[folded_entry.index]}
    {#each folded_entry.sub_entries as sub_entry}
      {@render table_row(sub_entry)}
    {/each}
  {/if}
{/snippet}

{#snippet table_row(entry: FinanceReportEntry)}
  <tr style="opacity: 80%;">
    <td></td>
    <td>{entry.Name}</td>
    <td>{entry.Info}</td>
    <td>{entry.Group}</td>
    <td>{entry.Cash_cost}</td>
    {#if entry.Value > 0}
      <td class="value_field" style="color: {green};"
        >+{format_number(entry.Value)}</td
      >
    {:else if entry.Value == 0}
      <td class="value_field">{format_number(entry.Value)}</td>
    {:else}
      <td class="value_field" style="color: {red};"
        >{format_number(entry.Value)}</td
      >
    {/if}
  </tr>
{/snippet}

<style>
  .fold_button {
    background-color: transparent;
    border-color: transparent;
    --folded: false;
  }

  .value_field {
    font-family: "JetBrains mono", "Courier New", Courier, monospace;
    text-align: right;
  }
</style>
