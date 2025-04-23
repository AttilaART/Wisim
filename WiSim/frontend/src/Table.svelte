<script lang="ts">
  import triangle from "./assets/images/Triangle(right)(white).png";
  import { int } from "../wailsjs/go/models";
  import { format_number, red, green } from "./helper";

  let currency = "CHF";

  let { data, table_type } = $props();
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

  let make_table = parse_accounting_data(data, table_type);
  async function parse_accounting_data(
    data: FinanceReportEntry[],
    table_type: string,
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
  <table style="text-align: left; ">
    <thead>
      <tr>
        <th style="width: 15; word-wrap: break-word;"></th>
        <th>Name</th>
        <th>Info</th>
        {#if table_type == "Income_statement"}
          <th>Cash cost?</th>
        {/if}
        <th style="text-align: right;">Value ({currency})</th>
      </tr>
    </thead>
    <tbody>
      {#each result.entries as folded_entry}
        {@render folded_table_row(folded_entry)}
      {/each}
    </tbody>
    <tfoot>
      {#if table_type == "Income_statement"}
        <tr>
          <td></td>
          <td>TOTAL CASHFLOW</td>

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
      {/if}
    </tfoot>
  </table>
{:catch error}
  an error occured: <br /> {error}
{/await}

{#snippet folded_table_row(folded_entry: FoldedFinanceReportEntry)}
  {#if folded_entry.sub_entries.length > 1}
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
            class="fold_image"
            height="10px"
            width="10px"
            alt="Triangle for folding entries"
            style="transform: rotate({Number(
              folded_entries_state[folded_entry.index],
            ) * 90}deg);"
          /></button
        ></td
      >
      <td>{folded_entry.entry.Name}</td>
      <td>{folded_entry.entry.Info}</td>

      {#if table_type == "Income_statement"}
        <td>{add_the_third_bool(folded_entry.entry.Cash_cost)}</td>
      {/if}
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
        {@render table_row(sub_entry, 0.8)}
      {/each}
    {/if}
  {:else}
    {@render table_row(folded_entry.sub_entries[0], 1)}
  {/if}
{/snippet}

{#snippet table_row(entry: FinanceReportEntry, opacity: number)}
  <tr style="opacity: {opacity * 100}%;">
    <td></td>
    <td>{entry.Name}</td>
    <td>{entry.Info}</td>
    {#if table_type == "Income_statement"}
      <td>{entry.Cash_cost}</td>
    {/if}
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
    width: 15px;
    padding: 0px;
  }

  .fold_image {
    position: relative;
    top: -2px;
  }
</style>
