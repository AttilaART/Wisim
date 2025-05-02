<script>
  import { format_number } from "./helper";
  import Slider from "./slider.svelte";
  let balance_without_loans = 100;
  let slider_value = $state();
  let real_balance = $state(110);
  let calculated_balance = $derived(balance_without_loans + slider_value);
</script>

<div style="display: flex; flex-direction: column; height: calc(100% - 60px);">
  <div style="padding: 10px 20px 10px 20px;">
    <h1 style="text-align: left;">Increase / Decrease Bank Loan</h1>
    <div style="display: flex; flex-direction: row; margin-bottom: 10px;">
      <span style="padding: 5px; flex: 1 1 60%; ">
        <Slider
          min={0}
          max={100}
          options={{
            default_value: 50,
            show_min_value: true,
            show_current_value: true,
            show_max_value: true,
            unit: " CHF",
          }}
          bind:Value={slider_value}
        ></Slider></span
      >

      <span style="width: 10%;"></span>
      <button style="flex 0 0 20%">Confirm</button>
    </div>
  </div>
  <span style="height: 3px; width: 100%; background-color: black;"></span>
  <div style="padding: 10px 20px 10px 20px; text-align: left;">
    Est. Balance: {calculated_balance} CHF
    <span
      style="color: {calculated_balance > real_balance
        ? '#2f9e44'
        : calculated_balance < real_balance
          ? '#e03131'
          : 'grey'}">({format_number(calculated_balance - real_balance)})</span
    >
    <br />
    Est. Monthly Interest: {calculated_balance} CHF
    <span
      style="color: {calculated_balance > real_balance
        ? '#2f9e44'
        : calculated_balance < real_balance
          ? '#e03131'
          : 'grey'}">({format_number(calculated_balance - real_balance)})</span
    >
  </div>
</div>

<style>
</style>
