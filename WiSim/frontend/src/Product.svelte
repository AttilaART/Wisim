<script lang="ts">
  import Slider from "./slider.svelte";
  import { decisions } from "./store.svelte";
  import { isEqual } from "./helper";

  let unapplied_changes: boolean = $state(false);

  let materials: {
    Quality: number;
    Ecology: number;
    Ethical_sourcing: number;
  } = $state({ ...decisions.Marketing.Product.Materials });

  let manufacturing: {
    Quality: number;
    Durability: number;
    Ecological_energy: number;
    Material_efficiency: number;
  } = $state({ ...decisions.Marketing.Product.Manufacturing });

  function cancel() {
    materials = { ...decisions.Marketing.Product.Materials };
    manufacturing = { ...decisions.Marketing.Product.Manufacturing };
  }

  function apply() {
    decisions.Marketing.Product.Materials = { ...materials };
    decisions.Marketing.Product.Manufacturing = { ...manufacturing };
  }

  let product_stats = $state({
    Quality: 0.5,
    Ecology: 0.5,
    Ethics: 0.5,
    Durability: 0.5,
  });

  let manufacturing_stats = $state({
    Speed: 0.5,
    Cost: 0.5,
    Weight: 0.5,
  });

  $effect(() => {
    if (
      isEqual(manufacturing, decisions.Marketing.Product.Manufacturing) &&
      isEqual(materials, decisions.Marketing.Product.Materials)
    ) {
      unapplied_changes = false;
    } else {
      unapplied_changes = true;
    }
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
      bind:Value={materials.Quality}
      min={0}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ecology</p>
    <Slider
      bind:Value={materials.Ecology}
      min={0}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ethical Sourcing</p>
    <Slider
      bind:Value={materials.Ethical_sourcing}
      min={0}
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
      <h3>Product Stats</h3>
      <table style="width: 100%;">
        <tbody>
          <tr>
            <td>Quality:</td>
            <td style="text-align: right;">{product_stats.Quality}</td>
          </tr>
          <tr>
            <td>Ecology:</td>
            <td style="text-align: right;">{product_stats.Ecology}</td>
          </tr>
          <tr>
            <td>Ethics:</td>
            <td style="text-align: right;">{product_stats.Ethics}</td>
          </tr>
          <tr>
            <td>Durability:</td>
            <td style="text-align: right;">{product_stats.Durability}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div style="flex: 1 1 60%;">
      <h3>Manufacturing Stats</h3>
      <table style="width: 100%;">
        <tbody>
          <tr>
            <td>Manufacturing Speed:</td>
            <td style="text-align: right;">{manufacturing_stats.Speed}x</td>
          </tr>
          <tr>
            <td>Material Cost:</td>
            <td style="text-align: right;">{manufacturing_stats.Cost}</td>
          </tr>
          <tr>
            <td>Weight:</td>
            <td style="text-align: right;">{manufacturing_stats.Weight}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
  <div
    class="grid_item"
    style="grid-row: 1 /span 2; grid-column: 3; text-align: right;"
  >
    <h2>Manufacturing</h2>
    <p>Quality</p>
    <Slider
      bind:Value={manufacturing.Quality}
      min={0}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Durability</p>
    <Slider
      bind:Value={manufacturing.Durability}
      min={0}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Ethical Energy Sourcing</p>
    <Slider
      bind:Value={manufacturing.Ecological_energy}
      min={0}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
    <p>Material Efficiency</p>
    <Slider
      bind:Value={manufacturing.Material_efficiency}
      min={0}
      max={5}
      options={{ step: 0.1 }}
    ></Slider>
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
