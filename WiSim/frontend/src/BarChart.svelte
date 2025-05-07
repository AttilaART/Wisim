<script lang="ts">
  import { generateGradient } from "./helper";
  import type { Series } from "./helper";

  type Bar = Series[];

  type MaxLine = {
    Show: boolean;
    Value: number;
    Label: string;
    OffsetPx?: number;
  };
  type Options = {
    MaxLine?: MaxLine;
  };

  let { Data, opts }: { Data: Bar[]; opts?: Options } = $props();
  let chart_element: HTMLDivElement = $state();
  let loaded = $state(false);

  let show_max_line = $state(false);
  let is_max_line_offset = $state(false);
  if (opts) {
    if (opts.MaxLine) {
      if (opts.MaxLine.Show) {
        show_max_line = true;
      }
      if (opts.MaxLine.OffsetPx) {
        is_max_line_offset = true;
      }
    }
  }

  let maximum: number = opts.MaxLine ? opts.MaxLine.Value : getMaxValue(Data);

  function getMaxValue(Data: Bar[]): number {
    let max_value: number = 0;
    for (let i in Data) {
      let bar_total = getBarLength(Data[i]);
      if (bar_total > max_value) {
        max_value = bar_total;
      }
    }

    return max_value;
  }

  function getBarLength(bar: Bar): number {
    let bar_total: number = 0;
    for (let ii in bar) {
      bar_total += bar[ii].Value;
    }
    return bar_total;
  }

  setTimeout(() => (loaded = true));
</script>

<div
  class="chart"
  style="--number-of-bars: {Data.length}; {is_max_line_offset
    ? `--maxline-offset: ${opts.MaxLine.OffsetPx}px`
    : ''}"
>
  <div class="data_canvas">
    <div class="baseline"></div>
    {#if show_max_line}
      <div class="maxline"></div>
      <div class="maxline_label">{opts.MaxLine.Label}</div>
    {/if}
    {#each Data as b}
      {@render bar(b)}
    {/each}
  </div>
</div>

{#snippet bar(b: Bar)}
  <div
    class="bar{loaded ? ' after' : ''}"
    style="background: {generateGradient(
      b,
      'linear-gradient',
      'to right',
      maximum,
    )}"
  ></div>
{/snippet}

<style>
  .chart {
    --bar-spacing: 10px;
    --bar-width: 20px;
    --number-of-bars: 1;
    --maxline-offset: 100px;
  }

  .data_canvas {
    position: relative;
    padding: 10px;
    padding-left: 0px;
    width: 100%;
  }

  .bar {
    position: relative;
    height: var(--bar-width);
    width: 0;
    margin-top: var(--bar-spacing);
    margin-bottom: var(--bar-spacing);
    transition: width 1s;
  }

  .bar.after {
    width: calc(100% - var(--maxline-offset));
  }

  .baseline,
  .maxline {
    z-index: 10;
    position: absolute;
    border-right: var(--border);
    height: calc(
      (var(--bar-spacing) + var(--bar-width)) * var(--number-of-bars) +
        var(--bar-spacing)
    );
  }

  .maxline {
    border-right-style: dotted;
    right: calc(var(--maxline-offset) + 10px);
  }

  .maxline_label {
    position: absolute;
    padding: 10px;
    right: var(--maxline-offset);
    transform: translateX(100%);
  }
</style>
