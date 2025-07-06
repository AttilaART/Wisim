<script lang="ts">
  import BarChart from "./BarChart.svelte";
  import { format_currency, format_number } from "./helper.svelte";
  import Slider from "./slider.svelte";
  import {
    company,
    current_decisions,
    external_factors,
    latest_reports,
    month,
  } from "./store.svelte";

  import { int, simulation } from "../wailsjs/go/models";
  import { get } from "svelte/store";

  let bridge_loans = $state(company.Bridge_loans);
  let existing_loans = $state(current_decisions.Finances.Set_bank_loan);
  let slider_value = $state(existing_loans);
  let balance_without_loans = $derived(
    company.Balance - bridge_loans - existing_loans,
  );
  let credit_limit = $state(1000000);

  let new_loans = $derived(slider_value - existing_loans);

  let balance_before = $derived(
    balance_without_loans + bridge_loans + existing_loans,
  );
  let balance_after = $derived(balance_before + new_loans);

  let interest_rate = external_factors.Intrest_rate;
  let bridge_loan_interest_rate = external_factors.Bridge_loans_intrest_rate;
  let interest_before = $derived(
    interest_rate * existing_loans + bridge_loan_interest_rate * bridge_loans,
  );
  let interest_after = $derived(
    interest_rate * (existing_loans + new_loans) +
      bridge_loan_interest_rate * bridge_loans,
  );

  let hover_over_balance = $state(false);

  let unaplied_changes = $state(false);

  $effect(() => {
    if (slider_value != existing_loans) unaplied_changes = true;
    else unaplied_changes = false;
  });
</script>

<div style="display: flex; flex-direction: column; height: calc(100% - 60px);">
  <div style="padding: 10px 20px 10px 20px;">
    {#key get(month)}
      <h1 style="text-align: left;">Increase / Decrease Bank Loan</h1>
      <div style="display: flex; flex-direction: row; margin-bottom: 10px;">
        <Slider
          min={0}
          max={credit_limit - bridge_loans}
          options={{
            default_value: existing_loans,
            show_min_value: true,
            show_current_value: true,
            show_max_value: true,
            snap: 10000,
            step: 1000,
            format: (val) => {
              return format_currency(val);
            },
          }}
          bind:Value={slider_value}
        ></Slider>
        <div
          style="display: flex; gap: 10px; padding: 0 10px; padding-left: 2rem;"
        >
          <button
            class={unaplied_changes ? "" : "greyed_out"}
            style="flex 0 0 20%"
            onclick={() => {
              slider_value = existing_loans;
            }}>Cancel</button
          >
          <button
            class={unaplied_changes ? "" : "greyed_out"}
            style="flex 0 0 20%"
            onclick={() => {
              company.Balance = balance_after;
              current_decisions.Finances.Set_bank_loan = slider_value;
              existing_loans = current_decisions.Finances.Set_bank_loan;
            }}>Apply</button
          >
        </div>
      </div>
    {/key}
  </div>
  <span class="sep_horisontal"></span>
  <div style="padding: 10px 20px 10px 20px; text-align: left;">
    <span class="balance"
      >Estimated Balance: <br />
      <h2>{format_currency(balance_after, 0)}</h2>
      {#if balance_after > balance_before}
        <span style="color: var(--green);">
          {format_currency(balance_after - balance_before, 0, true)}</span
        >
      {:else if balance_after < balance_before}
        <span style="color: var(--red);"
          >{format_currency(balance_after - balance_before, 0, true)}</span
        >
      {/if}
      <div class="details">
        Balance without loans: {format_currency(
          balance_without_loans,
          0,
          false,
        )}
        <br />
        Bridge Loan: {format_currency(bridge_loans, 0, false)} <br />
        Bank Loan: {format_currency(existing_loans + new_loans, 0, false)}
      </div></span
    >

    <span class="balance">
      Estimated Interest: <br />
      <h2>{format_currency(interest_after, 0, false)}</h2>
      {#if interest_after < interest_before}
        <span style="color: var(--green);">
          {format_currency(balance_after - balance_before, 0, true)}</span
        >
      {:else if interest_after > interest_before}
        <span style="color: var(--red);"
          >{format_currency(interest_after - interest_before, 0, true)}</span
        >
      {/if}
      <div class="details">
        From Bridge Loans: {format_currency(
          bridge_loan_interest_rate * bridge_loans,
          0,
          false,
        )} <br />
        From Bank Loans: {format_currency(
          interest_rate * (existing_loans + new_loans),
          0,
          false,
        )}
      </div>
    </span>
    <BarChart
      Data={[
        [
          { Name: "Bridge Loan", Value: bridge_loans, Color: "var(--red)" },
          {
            Name: "Current Bank Loan",
            Value: new_loans >= 0 ? existing_loans : existing_loans + new_loans,
            Color: "var(--red2)",
          },
          { Name: "New Bank Loan", Value: new_loans, Color: "var(--red3)" },
        ],
      ]}
      opts={{
        MaxLine: { Show: true, Value: credit_limit, Label: "Credit Limit" },
      }}
    ></BarChart>
  </div>
</div>

<style>
  .details {
    margin-left: 2em;
    color: var(--grey);
    overflow: hidden;
    height: 0;
    padding: 0;
    transition: height 1s;
  }

  .balance:hover .details {
    height: 4em;
  }
</style>
