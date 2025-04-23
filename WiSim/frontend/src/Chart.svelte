<script>
  import * as echarts from "../node_modules/echarts/";
  import { onMount } from "svelte";
  import { writable } from "svelte/store";

  let { height } = $props();

  let myChart = $state();
  let width_store = writable(10);

  onMount(() => {
    myChart = echarts.init(document.getElementById("chart"), null, {
      renderer: "canvas",
    });

    // Specify the configuration items and data for the chart
    let option = {
      //title: {
      //  text: "",
      //},
      // tooltip: {},
      //legend: {
      //  data: ["sales"],
      //},
      xAxis: {
        data: ["Shirts", "Cardigans", "Chiffons", "Pants", "Heels", "Socks"],
      },
      yAxis: {},
      series: [
        {
          name: "sales",
          type: "line",
          data: [5, 20, 36, 10, 10, 20],
        },
        {
          name: "sales",
          type: "line",
          data: [10, 20, 46, 20, 10, 10],
        },
      ],
    };

    // Display the chart using the configuration items and data just specified.
    myChart.setOption(option);

    width_store.subscribe((width) => {
      myChart.resize({ width: width - 10, height: height });
    });
  });
</script>

<div bind:clientWidth={$width_store} style="height: {height}px;">
  <div class="container" id="chart" style="width: 10px; height: 10px"></div>
</div>
