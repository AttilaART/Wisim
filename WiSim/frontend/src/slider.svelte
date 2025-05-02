<script lang="ts">
  import { tick } from "svelte";
  import { format_number } from "./helper";
  interface Options {
    default_value?: number;
    show_min_value?: boolean;
    show_max_value?: boolean;
    show_current_value?: boolean;
    unit?: string;
  }

  let {
    min,
    max,
    options,
    Value = $bindable(),
  }: {
    min: number;
    max: number;
    default_value?: number;
    options: Options;
    Value: number;
  } = $props();
  let slider_element: HTMLInputElement = $state();
  let min_value_element: HTMLDivElement = $state();
  let current_value_element: HTMLDivElement = $state();
  let max_value_element: HTMLDivElement = $state();
  let collision: boolean = $state(false);
  let slider_width: number = $state();

  async function adjust_background_when_loaded() {
    await tick();
    adjust_background();
    slider_width = get_slider_width();
  }
  adjust_background_when_loaded();

  function get_slider_width(): number {
    if (slider_element) {
      return slider_element.getBoundingClientRect().width - 25;
    }
  }

  function adjust_background() {
    slider_element.style.setProperty(
      "--ratio",
      `${((Value - min) / (max - min)) * 100}%`,
    );
  }

  function check_collision(): boolean {
    if (!current_value_element) {
      return false;
    }

    // Get bounding boxes
    let current_value_box = current_value_element.getBoundingClientRect();

    let min_value_box: undefined | DOMRect;

    if (min_value_element) {
      min_value_box = min_value_element.getBoundingClientRect();
    }

    let max_value_box: undefined | DOMRect;

    if (max_value_element) {
      max_value_box = max_value_element.getBoundingClientRect();
    }

    // Check for collision
    if (min_value_box) {
      // Check only horisontal
      if (current_value_box.left <= min_value_box.right) {
        return true;
      }
    }

    if (min_value_box) {
      // Check only horisontal
      if (max_value_box.left <= current_value_box.right) {
        return true;
      }
    }

    return false;
  }
</script>

<span style="position: relative; text-align: left;">
  <input
    style="width: 100%;"
    class="styled-slider slider-progress"
    type="range"
    defaultValue={options.default_value ? options.default_value : 0}
    {min}
    {max}
    bind:value={Value}
    bind:this={slider_element}
    oninput={() => {
      adjust_background();
      collision = check_collision();
      slider_width = get_slider_width();
    }}
  />
  {#if options.show_min_value}
    <div style="position: absolute; top: 10px;" bind:this={min_value_element}>
      {min}{options.unit}
    </div>
  {/if}
  {#if options.show_max_value}
    <div
      style="position: absolute; right: 0; top: 10px;"
      bind:this={max_value_element}
    >
      {max}{options.unit}
    </div>
  {/if}
  {#if options.show_current_value}
    {#key Value}
      <div
        bind:this={current_value_element}
        style="font-size: 1.2em; text-align: center; position: relative; transform: translate(calc({slider_width
          ? slider_width
          : 0}px * {(Value - min) / (max - min)} - 25%), {collision
          ? 40
          : -20}%);"
      >
        {format_number(Value ? Value : 0)}{options.unit}
      </div>
    {/key}
  {/if}
</span>

<style>
  div {
    width: fit-content;
    font-size: 0.9em;
  }

  /*generated with Input range slider CSS style generator (version 20211225)
https://toughengineer.github.io/demo/slider-styler*/
  input[type="range"].styled-slider {
    height: 2.2em;
    -webkit-appearance: none;
  }

  /*progress support*/
  input[type="range"].styled-slider.slider-progress {
  }

  input[type="range"].styled-slider:focus {
    outline: none;
  }

  /*webkit*/
  input[type="range"].styled-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    width: 25px;
    height: 25px;
    border-radius: 25px;
    background: #ffffff;
    border: 3px solid #000000;
    box-shadow: none;
    margin-top: calc(max((13px - 3px - 3px) * 0.5, 0px) - max(25px * 0.5, 3px));
  }

  input[type="range"].styled-slider::-webkit-slider-runnable-track {
    height: 13px;
    border: 3px solid #000000;
    border-radius: 10px;
    background: #ffffff;
    box-shadow: none;
  }

  input[type="range"].styled-slider:hover::-webkit-slider-runnable-track {
    background: #e5e5e5;
    border-color: #4e4e4e;
  }

  input[type="range"].styled-slider:active::-webkit-slider-runnable-track {
    background: #f5f5f5;
    border-color: #4e4e4e;
  }

  input[type="range"].styled-slider.slider-progress::-webkit-slider-runnable-track {
    background-image: linear-gradient(
      to right,
      #ff8787 0%,
      #ff8787 var(--ratio),
      white var(--ratio),
      white 100%
    );
  }

  input[type="range"].styled-slider.slider-progress:hover::-webkit-slider-runnable-track {
    background-image: linear-gradient(
      to right,
      #ff7070 0%,
      #ff7070 var(--ratio),
      white var(--ratio),
      white 100%
    );
  }

  input[type="range"].styled-slider.slider-progress:active::-webkit-slider-runnable-track {
    background-image: linear-gradient(
      to right,
      #ff5454 0%,
      #ff5454 var(--ratio),
      white var(--ratio),
      white 100%
    );
  }

  /*mozilla*/
  input[type="range"].styled-slider::-moz-range-thumb {
    width: max(calc(25px - 3px - 3px), 0px);
    height: max(calc(25px - 3px - 3px), 0px);
    border-radius: 25px;
    background: #ffffff;
    border: 3px solid #000000;
    box-shadow: none;
  }

  input[type="range"].styled-slider::-moz-range-track {
    height: max(calc(13px - 3px - 3px), 0px);
    border: 3px solid #000000;
    border-radius: 10px;
    background: #ffffff;
    box-shadow: none;
  }

  input[type="range"].styled-slider:hover::-moz-range-track {
    background: #e5e5e5;
    border-color: #4e4e4e;
  }

  input[type="range"].styled-slider:active::-moz-range-track {
    background: #f5f5f5;
    border-color: #4e4e4e;
  }

  input[type="range"].styled-slider.slider-progress::-moz-range-track {
    background-image: linear-gradient(
      to right,
      #ff8787 0%,
      #ff8787 var(--ratio),
      white var(--ratio),
      white 100%
    );
  }

  input[type="range"].styled-slider.slider-progress:hover::-moz-range-track {
    background-image: linear-gradient(
      to right,
      #ff7070 0%,
      #ff7070 var(--ratio),
      white var(--ratio),
      white 100%
    );
  }

  input[type="range"].styled-slider.slider-progress:active::-moz-range-track {
    background-image: linear-gradient(
      to right,
      #ff5454 0%,
      #ff5454 var(--ratio),
      white var(--ratio),
      white 100%
    );
  }

  /*ms*/
  input[type="range"].styled-slider::-ms-fill-upper {
    background: transparent;
    border-color: transparent;
  }

  input[type="range"].styled-slider::-ms-fill-lower {
    background: transparent;
    border-color: transparent;
  }

  input[type="range"].styled-slider::-ms-thumb {
    width: 25px;
    height: 25px;
    border-radius: 25px;
    background: #ffffff;
    border: 3px solid #000000;
    box-shadow: none;
    margin-top: 0;
    box-sizing: border-box;
  }

  input[type="range"].styled-slider::-ms-track {
    height: 13px;
    border-radius: 10px;
    background: #ffffff;
    border: 3px solid #000000;
    box-shadow: none;
    box-sizing: border-box;
  }

  input[type="range"].styled-slider:hover::-ms-track {
    background: #e5e5e5;
    border-color: #4e4e4e;
  }

  input[type="range"].styled-slider:active::-ms-track {
    background: #f5f5f5;
    border-color: #4e4e4e;
  }

  input[type="range"].styled-slider.slider-progress::-ms-fill-lower {
    height: max(calc(13px - 3px - 3px), 0px);
    border-radius: 10px 0 0 10px;
    margin: -3px 0 -3px -3px;
    background: #ff8787;
    border: 3px solid #000000;
    border-right-width: 0;
  }

  input[type="range"].styled-slider.slider-progress:hover::-ms-fill-lower {
    background: #ff7070;
    border-color: #4e4e4e;
  }

  input[type="range"].styled-slider.slider-progress:active::-ms-fill-lower {
    background: #ff5454;
    border-color: #4e4e4e;
  }
</style>
