<script lang="ts">
  import Donut from "./Donut.svelte";
  import { format_currency, format_number } from "./helper";
  import Slider from "./slider.svelte";
  import Tooltip from "./Tooltip.svelte";
  import { decisions } from "./store.svelte";
  import { isEqual } from "./helper";
  import NumberInput from "./number_input.svelte";

  let unsapplied_changes: boolean = $state(false);

  let promotion: {
    Quantity: number;
    Style_quality: number;
    Style_ecology: number;
    Style_ethics: number;
    Style_durability: number;
  } = $state({ ...decisions.Marketing.Promotion });

  let promotion_quantity_input = $state(format_currency(promotion.Quantity));

  function apply() {
    decisions.Marketing.Promotion = { ...promotion };
  }
  function cancel() {
    promotion = { ...decisions.Marketing.Promotion };
    promotion_quantity_input = format_currency(promotion.Quantity);
  }

  $effect(() => {
    if (isEqual(promotion, decisions.Marketing.Promotion)) {
      unsapplied_changes = false;
    } else {
      unsapplied_changes = true;
    }
  });
</script>

<div style="display: flex; height: calc(100% - 58px);">
  <div
    style="flex: 1 1 50%; border-right: var(--border); padding: 10px; overflow-y: scroll;"
  >
    <h2>Promotion Quantity</h2>
    <div style="display: flex;">
      <h2 style="flex 1 0 100%;">
        <NumberInput
          formatter={(value) => {
            return format_currency(value);
          }}
          bind:value={promotion.Quantity}
        ></NumberInput>
      </h2>
    </div>
    <br />
    <h2>Promotion Style</h2>
    <p>Quality</p>
    <Slider
      bind:Value={promotion.Style_quality}
      min={0}
      max={1}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ecology</p>
    <Slider
      bind:Value={promotion.Style_ecology}
      min={0}
      max={1}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ethics</p>
    <Slider
      bind:Value={promotion.Style_ethics}
      min={0}
      max={1}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Durability</p>
    <Slider
      bind:Value={promotion.Style_durability}
      min={0}
      max={1}
      options={{ step: 0.1 }}
    ></Slider>
    <div style="display: flex; flex-direction: row;">
      <button
        class={unsapplied_changes ? "" : "greyed_out"}
        style="padding: 10px; margin: 10px; flex: 1 1"
        onclick={cancel}>Cancel</button
      >
      <button
        class={unsapplied_changes ? "" : "greyed_out"}
        style="padding: 10px; margin: 10px;  flex: 1 1;"
        onclick={apply}>Apply</button
      >
    </div>
  </div>
  <div style="flex: 1 1 50%; padding: 10px;">
    <h2>
      Your Typical Customers
      <Tooltip text="The averages of the preferences of your customers"
      ></Tooltip>
    </h2>

    <Donut
      data={[
        {
          Name: "Avg. Quality Preference",
          Color: "var(--chart-color1)",
          Value: 1.34,
        },
        {
          Name: "Avg. Durability Preference",
          Color: "var(--chart-color2)",
          Value: 0.34,
        },
        {
          Name: "Avg. Ecology Preference",
          Color: "var(--chart-color3)",
          Value: 4.62,
        },
        {
          Name: "Avg. Ethical Preference",
          Color: "var(--chart-color4)",
          Value: 2.34,
        },
      ]}
    ></Donut>
  </div>
</div>

<style>
  input[type="text"] {
    width: 150px;
  }
  p {
    margin: 0;
    margin-left: 5px;
    text-align: left;
  }

  h2 {
    text-align: left;
  }
</style>
