<script lang="ts">
  import { format_number } from "./helper";
  import { month_counter, company, error } from "./store.svelte";

  export type Statement = Section[] | undefined | null;
  type StatementLine = {
    Name: string;
    Value: number;
    line_above?: boolean;
  };
  type Section = {
    Name: string;
    Lines: StatementLine[];
  };

  export type Invoice = {
    Name: string;
    Info: string;
    Category: string;
    Value: number;
  };

  let {
    get_income_statement,
    get_invoice_log,
  }: {
    get_income_statement: (
      month: number,
      company: number,
    ) => Promise<Statement>;
    get_invoice_log: (month: number, company: number) => Promise<Invoice[]>;
  } = $props();
  let income: Statement = $state(undefined);
  let invoice_log: undefined | null | Invoice[] = $state(undefined);

  async function s() {
    try {
      income = await get_income_statement($month_counter, $company);
    } catch (exeption) {
      invoice_log = null;
      $error = exeption;
    }
  }

  async function i() {
    try {
      invoice_log = await get_invoice_log($month_counter, $company);
    } catch (exeption) {
      invoice_log = null;
      $error = exeption;
    }
  }

  s();
  i();
</script>

<div style="display: flex; height: calc(100% - 60px);">
  <div style="padding: 10px; overflow-y: scroll; height: auto; flex: 1 1 60%;">
    {#if income === undefined}
      Loading report...
    {:else if income === null}
      There was a problem
    {:else}
      <table style="border-collapse: collapse; width: calc(100% - 10px);">
        <thead>
          <tr>
            <th></th>
            <th></th>
            <th></th>
          </tr>
        </thead>
        {#each income as s}
          {@render section(s)}
        {/each}
      </table>
    {/if}
  </div>
  <span class="sep_vertical" style="flex: 0 0 var(--border-width);"></span>
  <div
    style="padding: 0; overflow-y: scroll; flex: 1 1 calc(40% - var(border-width);"
  >
    {#if invoice_log === undefined}
      Loading...
    {:else if invoice_log === null}
      Therer was a problem
    {:else}
      <h2>Invoice Log</h2>
      <div class="sep_horisontal"></div>
      {#each invoice_log as i}
        {@render invoice(i)}
      {/each}
    {/if}
  </div>
</div>

{#snippet section(s: Section)}
  <tr>
    <td style="text-align: left;" colspan="3">
      <h3>{s.Name}</h3>
    </td>
  </tr>
  {#each s.Lines as l}
    <tr>
      <td style="width: 20px;"></td>
      <td
        style="text-align: left;{l.line_above
          ? 'border-top: var(--border-thin);'
          : ''}">{l.Name}</td
      >
      <td
        style="text-align: right; font-weight: 550; color: {l.Value > 0
          ? 'var(--green)'
          : l.Value < 0
            ? 'var(--red)'
            : 'var(--main-color)'};{l.line_above
          ? 'border-top: var(--border-thin);'
          : ''}">{format_number(l.Value, true, 0)}</td
      >
    </tr>
  {/each}
{/snippet}

{#snippet invoice(i: Invoice)}
  <div
    style="width: calc(100% - 20px); border-bottom: var(--border); padding: 0 10px 0 10px;"
  >
    <table style="width: 100%;">
      <thead>
        <tr>
          <th style="text-align: left;"><h4>{i.Name}</h4></th>
          <th
            style="text-align: right; color: {i.Value > 0
              ? 'var(--green)'
              : i.Value < 0
                ? 'var(--red)'
                : 'var(--main-color)'}"
            ><h4>{format_number(i.Value, true, 0)}</h4></th
          >
        </tr>
      </thead>
      <tbody>
        <tr>
          <td style="text-align: left;" colspan="2"
            ><span class="invoice info">{i.Info}</span></td
          >
        </tr>
        <tr>
          <td style="text-align: left;" colspan="2"
            ><span class="invoice category">{i.Category}</span></td
          >
        </tr>
      </tbody>
    </table>
  </div>
{/snippet}

<style>
  .invoice {
    font-size: 0.9em;
  }

  .category {
    opacity: 0.6;
  }
</style>
