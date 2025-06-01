<script lang="ts">
  import { format_currency, format_number } from "./helper";
  import NumberInput from "./number_input.svelte";
  import { decisions } from "./store.svelte";
  let global_price = $state(decisions.Marketing.Price);
  let global_price_input = $state(format_currency(global_price, 2));
  let unapplied_changes = $state(false);

  function cancel() {
    global_price = decisions.Marketing.Price;
    global_price_input = format_currency(global_price, 2);
  }

  function apply() {
    decisions.Marketing.Price = global_price;
  }

  $effect(() => {
    if (global_price == decisions.Marketing.Price) {
      unapplied_changes = false;
    } else {
      unapplied_changes = true;
    }
  });
</script>

<div style="text-align: left; padding: 10px;">
  <h2>Set Price</h2>
  <div>
    <p>Global Price</p>
    <h2>
      <NumberInput
        formatter={(value: number) => {
          return format_currency(value, 2);
        }}
        bind:value={global_price}
      ></NumberInput>
    </h2>
    <button class={unapplied_changes ? "" : "greyed_out"} onclick={cancel}
      >Cancel</button
    >
    <button class={unapplied_changes ? "" : "greyed_out"} onclick={apply}
      >Apply</button
    >
  </div>
</div>

<style>
  button {
    padding: 10px;
    margin: 10px;
  }
</style>
