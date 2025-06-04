<script>
  import BarChart from "./BarChart.svelte";
  import { format_currency, format_number } from "./helper";
  import Slider from "./slider.svelte";
  import Info from "./assets/images/Info.svelte";
  import { number } from "echarts";

  // TODO: Make loans affect balance

  let slider_value = $state();
  let balance_without_loans = $state(100000 - 70000);
  let bridge_loans = $state(30000);
  let existing_loans = $state(40000);
  let credit_limit = $state(150000);

  let new_loans = $derived(slider_value - existing_loans);

  let balance_before = $derived(
    balance_without_loans + bridge_loans + existing_loans,
  );
  let balance_after = $derived(balance_before + new_loans);
  slider_value = existing_loans;

  let interest_rate = 0.0015;
  let bridge_loan_interest_rate = 0.003;
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
            existing_loans = slider_value;
          }}>Apply</button
        >
      </div>
    </div>
  </div>
  <span class="sep_horisontal"></span>
  <div style="padding: 10px 20px 10px 20px; text-align: left;">
    <span class="balance"
      >Est. Balance: {format_number(balance_after, false, 0)} CHF
      <span
        style="color: {balance_after > balance_before
          ? 'var(--green)'
          : balance_after < balance_before
            ? 'var(--red)'
            : 'var(--grey)'}"
        >({format_number(balance_after - balance_before, true)})</span
      >
      <Info></Info>
      <div class="details">
        Balance without loans: {format_number(balance_without_loans, false, 0)}
        <br />
        Bridge Loans: {format_number(bridge_loans, false, 0)} <br />
        Bank Loans: {format_number(existing_loans + new_loans, false, 0)}
      </div></span
    >

    <span class="balance">
      Est. Monthly Interest: {format_number(interest_after, false, 0)} CHF
      <span
        style="color: {interest_after < interest_before
          ? 'var(--green)'
          : interest_after > interest_before
            ? 'var(--red)'
            : 'var(--grey)'}"
        >({format_number(interest_after - interest_before, true, 0)})</span
      >
      <Info></Info>
      <div class="details">
        From Bridge Loans: {format_number(
          bridge_loan_interest_rate * bridge_loans,
          false,
          0,
        )} <br />
        From Bank Loans: {format_number(
          interest_rate * (existing_loans + new_loans),
          false,
          0,
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
