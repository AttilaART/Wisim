<script lang="ts">
  // import { Component } from "svelte";
  import { draggable } from "@neodrag/svelte";

  let { title, content }: { title: string; content: any; content_args: any[] } =
    $props();
  let height = $state(400);
  let width = $state(400);

  let window_div: HTMLElement = $state();
  let resizer: HTMLElement = $state();

  function set_mouse_position(event) {
    height = event.clientY;
    width = event.clientX;
  }

  function resize() {
    window_div.setAttribute("style", `height: ${height}px; width: ${width}px`);
  }

  function reset_resizer() {
    resizer.setAttribute("style", "bottom: 0%; right: 0%;");
  }
</script>

<div
  class="window"
  use:draggable={{ bounds: "parent", handle: ".handle" }}
  bind:this={window_div}
  on:neodrag:end={reset_resizer}
>
  <div class="handle window-titlebar">
    {title}
  </div>
  {@render content?.()}
  <!--<div
    class="resize"
    use:draggable={{}}
    on:neodrag={resize}
    on:neodrag:end={() => {
      resize();
      reset_resizer();
    }}
    bind:this={resizer}
  >
    x
  </div>-->
</div>

<svelte:window on:mousemove={set_mouse_position} />

<style>
  .window {
    position: absolute;
    width: 500px;
    height: 350px;
    background-color: rgba(0, 0, 0, 0.5);
    border: solid white;
    overflow: hidden;
    padding: 0px;
  }

  .window-titlebar {
    width: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    text-align: center;
  }

  .resize {
    position: absolute;
    width: 20px;
    height: 20px;
    right: 0%;
    bottom: 0%;
  }
</style>
