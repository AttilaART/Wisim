<script lang="ts">
  import { format_number, red, green } from "./helper";
  import { get } from "svelte/store";
  import { month_counter } from "./store";

  interface properties {
    get_report: (month: number) => object;
    colour_values: boolean;
    show_plus: boolean;
    decimal_places: number;
    custom_rows: { row_name: string; properties: custom_row_properties }[];
  }

  interface custom_row_properties {
    colour_values: boolean;
    show_plus: boolean;
    decimal_places: number;
  }

  let {
    get_report,
    colour_values = false,
    show_plus = false,
    decimal_places = 0,
    custom_rows = [],
  }: properties = $props();

  console.log(custom_rows);

  month_counter.subscribe((month) => {
    sep_key_values(get_report, month);
  });

  type key_value_pair = {
    key: string;
    two_months_ago_value: any;
    last_month_value: any;
    current_value: any;
  };

  type Key_value_pairs_data = {
    data: key_value_pair[];
    loading: boolean;
    error: Error | undefined;
  };

  let data: Key_value_pairs_data = $state({
    data: [],
    loading: true,
    error: null,
  });

  let custom_rows_names: string[] = [];
  for (let i in custom_rows) {
    custom_rows_names.push(custom_rows[i].row_name);
  }

  async function sep_key_values(get_values, month: number) {
    let keys: string[] = [];
    let current: any[] = [];
    let last_month: any[] = [];
    let two_months_ago: any[] = [];

    data.error = null;
    try {
      keys = Object.keys(await get_values(month));
      current = Object.values(await get_values(month));
    } catch (error) {
      data.error = error;
    }

    if (month >= 2) {
      last_month = Object.values(await get_values(month - 1));
    } else {
      for (let _ in keys) {
        last_month.push(0);
      }
    }
    if (month >= 3) {
      two_months_ago = Object.values(await get_values(month - 2));
    } else {
      for (let _ in keys) {
        two_months_ago.push(0);
      }
    }

    let key_values: key_value_pair[] = [];
    for (let i in keys) {
      key_values.push({
        key: keys[i],
        two_months_ago_value: two_months_ago[i],
        last_month_value: last_month[i],
        current_value: current[i],
      });
    }
    data.data = key_values;
    data.loading = false;
  }

  sep_key_values(get_report, get(month_counter));
</script>

{#if data.loading}
  loading...
{:else if data.error === null}
  <table>
    <thead>
      <tr>
        <th style="width: auto; text-align: left;">Month / Year</th>
        <th style="width: 20%; text-align: right;"
          >{($month_counter - 1) % 13} / {Math.floor(
            ($month_counter - 2) / 12,
          )}</th
        >
        <th style="width: 20%; text-align: right;"
          >{$month_counter % 13} / {Math.floor(($month_counter - 1) / 12)}</th
        >
        <th style="width: 20%; text-align: right;"
          >{($month_counter + 1) % 13} / {Math.floor($month_counter / 12)}
        </th>
      </tr>
    </thead>
    <tbody>
      {#each data.data as pair}
        {#if custom_rows_names.includes(pair.key)}
          {@render row(
            pair.key,
            pair.two_months_ago_value,
            pair.last_month_value,
            pair.current_value,
            custom_rows[custom_rows_names.indexOf(pair.key)].properties,
          )}
        {:else}
          {@render row(
            pair.key,
            pair.two_months_ago_value,
            pair.last_month_value,
            pair.current_value,
            null,
          )}
        {/if}
      {/each}
    </tbody>
  </table>
{:else}
  failed to load data: <br />
  {data.error}
{/if}

{#snippet row(
  key: string,
  two_years_ago,
  previous_year,
  current_value,
  custom_properties,
)}
  <tr>
    <td style="text-align: left;">{key.replaceAll("_", " ")}</td>
    {#if custom_properties != null}
      {@render value_entry(
        two_years_ago,
        custom_properties.show_plus,
        custom_properties.colour_values,
        custom_properties.decimal_places,
      )}
      {@render value_entry(
        previous_year,
        custom_properties.show_plus,
        custom_properties.colour_values,
        custom_properties.decimal_places,
      )}
      {@render value_entry(
        current_value,
        custom_properties.show_plus,
        custom_properties.colour_values,
        custom_properties.decimal_places,
      )}
    {:else}
      {@render value_entry(
        two_years_ago,
        show_plus,
        colour_values,
        decimal_places,
      )}
      {@render value_entry(
        previous_year,
        show_plus,
        colour_values,
        decimal_places,
      )}
      {@render value_entry(
        current_value,
        show_plus,
        colour_values,
        decimal_places,
      )}
    {/if}
  </tr>
{/snippet}

{#snippet value_entry(value, show_plus, colour_values, decimal_places)}
  {#if typeof value == typeof 1}
    {#if value > 0 && colour_values}
      <td class="value_field" style="color: {green}"
        >{format_number(value, show_plus, decimal_places)}</td
      >
    {:else if value < 0 && colour_values}
      <td class="value_field" style="color: {red}"
        >{format_number(value, show_plus, decimal_places)}</td
      >
    {:else}
      <td class="value_field">{format_number(value)}</td>
    {/if}
  {:else if typeof value == typeof "string"}
    <td class="value_field">{value.replaceAll("_", " ")}</td>
  {/if}
{/snippet}
