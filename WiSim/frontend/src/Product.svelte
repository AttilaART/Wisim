<script lang="ts">
  import Slider from "./slider.svelte";
  import {
    company,
    company_id,
    current_decisions,
    external_factors,
  } from "./store.svelte";
  import { format_number, isEqual } from "./helper.svelte";
  import NumberInput from "./number_input.svelte";
  import { Calculate_product_stats } from "../wailsjs/go/main/App";
  import { get } from "svelte/store";
  import { simulation } from "../wailsjs/go/models";

  let unapplied_changes: boolean = $state(false);

  let product_decisions: simulation.Decisions_product = $state(
    JSON.parse(JSON.stringify(current_decisions.Marketing.Product)),
  );

  function cancel() {
    product_decisions = JSON.parse(
      JSON.stringify(current_decisions.Marketing.Product),
    );
  }

  function apply() {
    current_decisions.Marketing.Product = JSON.parse(
      JSON.stringify(product_decisions),
    );
  }

  let temp_product_promise = $state(
    Calculate_product_stats(
      get(company_id),
      product_decisions,
      current_decisions.Research,
    ),
  );

  $effect(() => {
    if (isEqual(product_decisions, current_decisions.Marketing.Product)) {
      unapplied_changes = false;
    } else {
      unapplied_changes = true;
    }
  });

  $effect(() => {
    temp_product_promise = Calculate_product_stats(
      get(company_id),
      product_decisions,
      current_decisions.Research,
    );
  });
</script>

<div
  class="grid_container"
  style="grid-template-columns: 40% auto 40%; height: calc(100% - 50px);"
>
  <div class="grid_item" style="grid-row: 1; grid-column: 1; text-align: left;">
    <h2>Materials</h2>
    <p>Quality</p>
    <Slider
      bind:Value={product_decisions.Materials.Quality}
      min={0.1}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ecology</p>
    <Slider
      bind:Value={product_decisions.Materials.Ecology}
      min={0.1}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ethical Sourcing (Not implemented)</p>
    <Slider
      bind:Value={product_decisions.Materials.Ethical_sourcing}
      min={0.1}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
  </div>
  <div class="grid_item" style="grid-row: 1; grid-column: 2;">
    Product Image
  </div>

  <div
    class="grid_item"
    style="grid-row: 2; grid-column: 1 /span 2; border-right: var(--border-thin); border-top: var(--border-thin); display: flex; flex-direction: row; border-top-right-radius: var(--window-border-radius);"
  >
    <div style="flex: 1 1 40%;">
      {#await temp_product_promise}
        <h3>Product Stats</h3>
        <h4>Calculating...</h4>
      {:then temp_product}
        <h3>Product Stats</h3>
        <table style="width: 100%;">
          <tbody>
            <tr>
              <td>Quality:</td>
              <td style="text-align: right;">
                {format_number(temp_product.Quality_factor, false, 3)}
              </td>
            </tr>
            <tr>
              <td>Ecology:</td>
              <td style="text-align: right;"
                >{format_number(temp_product.Ecology_factor, false, 3)}</td
              >
            </tr>
            <tr>
              <td>Ethics:</td>
              <td style="text-align: right;"
                >{format_number(temp_product.Ethics_factor, false, 3)}</td
              >
            </tr>
            <tr>
              <td>Durability:</td>
              <td style="text-align: right;"
                >{format_number(temp_product.Durabilty, false, 0)} Mon</td
              >
            </tr>
          </tbody>
        </table>
      {/await}
    </div>
    <div style="flex: 1 1 60%;">
      {#await temp_product_promise then temp_product}
        <h3>Production Stats</h3>
        <table style="width: 100%;">
          <tbody>
            <tr>
              <td>Production cost:</td>
              <td style="text-align: right;"
                >{format_number(temp_product.Production_cost, false, 3)}x</td
              >
            </tr>
            <tr>
              <td>Material use:</td>
              <td style="text-align: right;"
                >{format_number(temp_product.Material_use, false, 3)}</td
              >
            </tr>
            <tr>
              <td>Weight:</td>
              <td style="text-align: right;"
                >{format_number(temp_product.Weight, false, 3)}</td
              >
            </tr>
          </tbody>
        </table>
      {/await}
    </div>
  </div>
  <div
    class="grid_item"
    style="grid-row: 1 /span 2; grid-column: 3; text-align: right;"
  >
    <h2>Manufacturing</h2>
    <p>Quality</p>
    <Slider
      bind:Value={product_decisions.Manufacturing.Quality}
      min={0.1}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Durability</p>
    <Slider
      bind:Value={product_decisions.Manufacturing.Durability}
      min={0.1}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ecological Energy usage</p>
    <Slider
      bind:Value={product_decisions.Manufacturing.Ecological_energy}
      min={0.1}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Material Efficiency</p>
    <Slider
      bind:Value={product_decisions.Manufacturing.Material_efficiency}
      min={0.1}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Max durability</p>
    <h2>
      <NumberInput
        bind:value={product_decisions.Manufacturing.Max_durability}
        align_right={true}
        formatter={(value) => {
          return format_number(value, false, 0) + " Months";
        }}
      ></NumberInput>
    </h2>
    <div style="display: flex; margin-top: auto; width: 100%;">
      <button
        class={unapplied_changes ? "" : "greyed_out"}
        style="flex: 1 1 50%;"
        onclick={cancel}>Cancel</button
      >
      <button
        class={unapplied_changes ? "" : "greyed_out"}
        style="flex: 1 1 50%;"
        onclick={apply}>Apply</button
      >
    </div>
  </div>
</div>

<style>
  p {
    margin: 0;
    margin-left: 5px;
  }

  button {
    padding: 10px;
    margin: 10px;
  }

  h3 {
    text-align: left;
  }

  table {
    text-align: left;
    padding: 5px;
    padding-left: 10px;
  }
</style>
