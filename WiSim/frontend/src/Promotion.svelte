<script>
  import Donut from "./Donut.svelte";
  import { format_number } from "./helper";
  import Slider from "./slider.svelte";
  import Info from "./assets/images/Info.svelte";
  import Tooltip from "./Tooltip.svelte";

  let promotion_quantity_input = $state("150'000 CHF");
  let promotion_quantity = $state(150000);

  let promotion_style = $state({
    Quality: 0.5,
    Ecology: 0.5,
    Ethics: 0.5,
    Durability: 0.5,
  });
</script>

<div style="display: flex; height: calc(100% - 58px);">
  <div
    style="flex: 1 1 50%; border-right: var(--border); padding: 10px; overflow-y: scroll;"
  >
    <h2>Promotion Quantity</h2>
    <div style="display: flex;">
      <h2 style="flex 1 0 100%;">
        <input
          type="text"
          bind:value={promotion_quantity_input}
          onfocus={() =>
            (promotion_quantity_input = String(promotion_quantity))}
          onfocusout={() => {
            if (!isNaN(parseFloat(promotion_quantity_input))) {
              promotion_quantity = parseFloat(promotion_quantity_input);
            }
            promotion_quantity_input = `${format_number(promotion_quantity, false)} CHF`;
          }}
        />
      </h2>
      <div style="flex: 1 1 100%;"></div>
      <button style="flex: 0 0 fit-content;">Confirm</button>
    </div>
    <br />
    <h2>Promotion Style</h2>
    <p>Quality</p>
    <Slider
      bind:Value={promotion_style.Quality}
      min={0}
      max={1}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ecology</p>
    <Slider
      bind:Value={promotion_style.Ecology}
      min={0}
      max={1}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ethics</p>
    <Slider
      bind:Value={promotion_style.Ethics}
      min={0}
      max={1}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Durability</p>
    <Slider
      bind:Value={promotion_style.Durability}
      min={0}
      max={1}
      options={{ step: 0.1 }}
    ></Slider>
    <div style="text-align: right; margin-top: 10px;">
      <button style="padding: 10px;">Confirm Style</button>
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
