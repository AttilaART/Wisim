<script lang="ts">
  // import { Component } from "svelte";
  import { draggable } from "@neodrag/svelte";
  import {
    windows,
    new_window,
    move_window_to_top,
    delete_window,
  } from "./window_manager.svelte";
  import Close from "./assets/images/Close.svelte";
  import Min from "./assets/images/Min.svelte";
  import Fullscreen from "./assets/images/Fullscreen.svelte";
  import { tick } from "svelte";

  let {
    title,
    content,
    close_window,
    window_id = $bindable(),
    canvas_size,
  }: {
    title: string;
    content: any;
    content_args: any[];
    close_window: () => void;
    window_id?: number;
    canvas_size?: { x: number; y: number };
  } = $props();
  window_id = new_window(title, close_window);

  let fullscreen: boolean = $state(false);
  // let position = $state({ x: -canvas_size.x / 4, y: 0 });
  let position = $state({ x: 0, y: 0 });

  let window_div: HTMLDivElement = $state();
  let titlebar: HTMLDivElement = $state();
  function update_position(x: number, y: number) {
    position.x = x;
    position.y = y;
    console.log($state.snapshot(position));
  }

  (async () => {
    await tick();
    update_position(-canvas_size.x / 4, 0);
  })();

  $effect(() => {
    console.log(windows[window_id].z_index);
  });
</script>

{#if windows[window_id]}
  <div
    role="none"
    class="window"
    style="z-index: {windows[window_id].z_index};
    transform: translate({position.x}px, {position.y}px); display: {windows[
      window_id
    ].hidden
      ? 'none'
      : 'unset'};
  {fullscreen
      ? `--width: ${canvas_size.x - 5}px; --height: ${canvas_size.y - 5}px;`
      : `--width: ${canvas_size.x / 2 - 5}px; --height: ${canvas_size.y / 2 - 5}px`}"
    use:draggable={{
      bounds: "parent",
      handle: titlebar,
      position,
      onDragEnd: ({ offsetX, offsetY }) => update_position(offsetX, offsetY),
      onDragStart: () => {
        move_window_to_top(window_id);
        // console.log($state.snapshot(position));
      },
    }}
    bind:this={window_div}
    onmousedown={() => {
      move_window_to_top(window_id);
    }}
  >
    <div class="window-titlebar" bind:this={titlebar}>
      {title}
      <div style="position: absolute; right: 0; top: 0;">
        <button
          class="window_button"
          onclick={() => (windows[window_id].hidden = true)}><Min></Min></button
        >
        <button
          class="window_button"
          onclick={() => {
            fullscreen = !fullscreen;
            position.x = -canvas_size.x / 2;
            position.y = 0;
          }}
          ><Fullscreen></Fullscreen>
        </button>
        <button
          class="window_button"
          onclick={() => {
            delete_window(window_id);
          }}
        >
          <Close></Close>
        </button>
      </div>
    </div>
    {@render content?.()}
  </div>
{/if}

<style>
  .window {
    position: absolute;
    --height: 450px;
    --width: 650px;
    width: var(--width);
    height: var(--height);
    border: var(--border);
    overflow: hidden;
    padding: 0px;
    background-color: var(--second-color);
    pointer-events: all;
  }

  .window-titlebar {
    width: 100%;
    height: 25px;
    font-size: 1rem;
    text-align: center;
    border-bottom: var(--border);
    position: relative;
  }

  .window_button {
    padding: 0;
    height: 100%;
    width: 25px;
    height: 25px;
    border: none;
  }
</style>
