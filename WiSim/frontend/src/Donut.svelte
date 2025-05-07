<script lang="ts">
  import { generateGradient } from "./helper";
  import type { Series } from "./helper";

  let { data }: { data: Series[] } = $props();
</script>

<div style="height: auto;">
  <table style="width: 100%;">
    <tbody>
      {#each data as series}
        <tr>
          <td style="text-align: left;">
            <div
              class="legend_icon"
              style="background-color: {series.Color};"
            ></div>
            {series.Name}
          </td>
          <td style="text-align: right;">{series.Value}</td>
        </tr>
      {/each}
    </tbody>
  </table>
  <center style="position: relative; height: 100%;">
    <div
      class="circle"
      style="position: relative; top: 50%; background-image: {generateGradient(
        data,
        'conic-gradient',
        'from 0deg',
      )};"
    >
      <div class="inner_circle" style=""></div>
    </div>
  </center>
</div>

<style>
  .circle {
    border-radius: 50%;
    width: 50%;
    aspect-ratio: 1 / 1;
    margin: 10px;
    border: var(--border);
  }

  .inner_circle {
    border-radius: 50%;
    width: 80%;
    aspect-ratio: 1 / 1;
    border: var(--border);
    background-color: var(--second-color);
    transform: translateY(calc(13.5% - var(--border-width) * 2));
    margin: 0;
  }

  .legend_icon {
    display: inline-block;
    width: 10px;
    height: 10px;
    aspect-ratio: 1 /1;
  }
</style>
