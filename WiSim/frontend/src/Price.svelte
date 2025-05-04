<script>
  import { format_number } from "./helper";
  let global_price_input = $state("100 CHF");
  let global_price = $state(100);
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
          global_price_input = `${format_number(global_price, false, 2)} CHF`;
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
