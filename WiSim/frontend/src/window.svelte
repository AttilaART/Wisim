<script lang="ts">
  // import { Component } from "svelte";
  import { draggable } from "@neodrag/svelte";
  import { windows, new_window, move_window_to_top } from "./store.svelte";
  import min from "./assets/images/min.svg";
  import fullscreen from "./assets/images/fullscreen.svg";
  import close from "./assets/images/close.svg";

  let {
    title,
    content,
    window_id = $bindable(),
    loaded = $bindable(),
  }: {
    title: string;
    content: any;
    content_args: any[];
    window_id?: number;
    loaded?: boolean;
  } = $props();
  window_id = new_window(title);

  let position = $state({ x: 0, y: 0 });

  let window_div: HTMLDivElement = $state();
  let titlebar: HTMLDivElement = $state();
  function update_position(x: number, y: number) {
    position.x = x;
    position.y = y;
    console.log($state.snapshot(position));
  }
</script>

<div
  role="none"
  class="window"
  style="z-index: {windows[window_id]
    .z_index}; transform: translate({position.x}px, {position.y}px); display: {windows[
    window_id
  ].hidden
    ? 'none'
    : 'unset'}"
  use:draggable={{
    bounds: "parent",
    handle: titlebar,
    position,
    onDragEnd: ({ offsetX, offsetY }) => update_position(offsetX, offsetY),
    onDragStart: () => {
      move_window_to_top(window_id);
      console.log($state.snapshot(position));
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
        onclick={() => (windows[window_id].hidden = true)}
        ><img
          aria-label="minimise"
          alt="minimise"
          style="display: block;margin: auto; width: 50%;"
          src={min}
        /></button
      >
      <button class="window_button"
        ><img
          aria-label="fullscreen"
          alt="fullscreen"
          style="display: block;margin: auto; width: 50%;"
          src={fullscreen}
        /></button
      >
      <button
        class="window_button"
        onclick={() => {
          windows.splice(window_id, 1);
          loaded = false;
        }}
        ><img
          aria-label="close window"
          alt="close window"
          style="display: block;margin: auto; width: 50%;"
          src={close}
        /></button
      >
    </div>
  </div>
  {@render content?.()}
</div>

<style>
  .window {
    position: absolute;
    width: 650px;
    height: 450px;
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
