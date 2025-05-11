<script>
  import { format_currency, format_number } from "./helper";
  import { decisions } from "./store.svelte";
  let global_price = $state(decisions.Marketing.Price);
  let global_price_input = $state(format_currency(global_price, 2));
</script>

<div style="text-align: left; padding: 10px;">
  <h2>Set Price</h2>
  <div>
    <p>Global Price</p>
    <h2>
      <input
        type="text"
        bind:value={global_price_input}
        onfocus={() => (global_price_input = String(global_price))}
        onfocusout={() => {
          if (!isNaN(parseFloat(global_price_input))) {
            global_price = parseFloat(global_price_input);
          }
          global_price_input = format_currency(global_price, 2);
        }}
      />
    </h2>
    <button>Confirm</button>
  </div>
</div>

<style>
  p {
    padding-left: 5px;
    margin: 0px;
  }
</style>
