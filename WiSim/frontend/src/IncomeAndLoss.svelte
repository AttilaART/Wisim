<script lang="ts">
  import close from "./assets/images/close.svg";
  import { format_currency } from "./helper.svelte";
  import { company_id, error, external_factors } from "./store.svelte";

  export type Statement = Section[] | undefined | null;
  type StatementLine = {
    Name: string;
    Value: number;
    line_above?: boolean;
  };
  type Section = {
    Name: string;
    Lines: StatementLine[];
    Period?: string;
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
    is_budget,
  }: {
    get_income_statement: (
      month: number,
      company: number,
    ) => Promise<Statement>;
    get_invoice_log: (month: number, company: number) => Promise<Invoice[]>;
    is_budget: boolean;
  } = $props();
  let income_promise: Promise<Statement> = $state(
    get_income_statement(external_factors.Month, $company_id),
  );
  let invoice_promise: Promise<Invoice[]> = $state(
    get_invoice_log(external_factors.Month, $company_id),
  );

  let show_invoices = $state(false);
</script>

<div class="content" style="display: flex; height: calc(100% - 60px);">
  <div style="flex: 1 1 20%;"></div>
  <div style="padding: 10px; overflow-y: scroll; height: auto; flex: 1 1 60%;">
    {#await income_promise}
      Loading...
    {:then income: Statement}
      <table style="border-collapse: collapse; width: calc(100% - 10px);">
        <thead>
          <tr>
            <th></th>
            <th></th>
            <th></th>
          </tr>
        </thead>
        {#each income as s, i}
          {@render section(s, i)}
        {/each}
      </table>
    {:catch exeption}
      {exeption}
    {/await}
  </div>
  <div style="flex: 1 1 20%;">
    <button
      style="margin: 10px; padding: 10px;"
      onclick={() => (show_invoices = true)}>See invoice log</button
    >
  </div>
  {#if show_invoices}
    <dialog open style="height: 80%; overflow-y: scroll;">
      <button
        class="borderless"
        style="aspect-ratio: 1/1; width: 2rem; position: sticky; top: 0rem;"
        onclick={() => (show_invoices = false)}
      >
        <img class="window_button" src={close} alt="close" />
      </button>
      {#await invoice_promise}
        Loading...
      {:then invoice_log}
        <h2>Invoice Log</h2>
        <div class="sep_horisontal"></div>
        {#each invoice_log as i}
          {@render invoice(i)}
        {/each}
      {:catch exeption}
        {exeption}
      {/await}
    </dialog>
  {/if}
  <!--
  -->
</div>

{#snippet section(s: Section, index: number)}
  <tr>
    <td style="text-align: left;" colspan="3">
      <div style="position: relative;">
        <h3 style="display: inline;">{s.Name}</h3>
        {#if index == 0}
          <span style="position: absolute; right: 0; color: var(--grey);">
            {s.Period}</span
          >
        {/if}
      </div>
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
          : ''}">{format_currency(l.Value, 0, true)}</td
      >
    </tr>
  {/each}
{/snippet}

{#snippet invoice(i: Invoice)}
  <div
    style="width: calc(100% - 20px); border-bottom: var(--window-border); padding: 0 10px 0 10px;"
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
            ><h4>{format_currency(i.Value, 0, true)}</h4></th
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

{#if is_budget}
  <style>
    .content {
      --main-color: #ffecc7;
      color: var(--main-color);
    }
  </style>
{/if}

<style>
  .invoice {
    font-size: 0.9em;
  }

  .category {
    opacity: 0.6;
  }
</style>
