<script lang="ts">
  // import { Component } from "svelte";
  import { draggable } from "@neodrag/svelte";
  import close from "./assets/images/close.svg";
  import min from "./assets/images/min.svg";
  import fullscreen_icon from "./assets/images/fullscreen.svg";
  import { tick } from "svelte";

  let {
    title,
    onClose,
    onDrag,
    onHide,
    Zindex = $bindable(),
    Hidden = $bindable(),
    canvas_size,
  }: {
    title: string;
    onClose: () => void;
    onDrag: () => void;
    onHide: () => void;
    Zindex: number;
    Hidden: boolean;
    canvas_size?: { x: number; y: number };
  } = $props();

  let fullscreen: boolean = $state(false);
  // let position = $state({ x: -canvas_size.x / 4, y: 0 });
  let position = $state({ x: 0, y: 0 });
  let last_non_fullscreen_position = { ...position };

  let window_div: HTMLDivElement = $state();
  let titlebar: HTMLDivElement = $state();
  function update_position(x: number, y: number) {
    if (y < 0) {
      y = 0;
    }
    position.x = x;
    position.y = y;
  }

  (async () => {
    await tick();
    let spawn_pos: { x: number; y: number } = {
      x: canvas_size.x / 4,
      y: canvas_size.y / 4,
    };
    update_position(spawn_pos.x, spawn_pos.y);
  })();
</script>

<div
  role="none"
  class="window {Hidden ? 'hidden' : ''}"
  style="z-index: {Zindex};
    transform: translate({position.x}px, {position.y}px); 
  {fullscreen
    ? `--width: ${canvas_size.x - 5}px; --height: ${canvas_size.y - 5}px;`
    : `--width: ${750}px; --height: ${650}px`}"
  use:draggable={{
    //bounds: "parent",
    handle: titlebar,
    position,
    onDragEnd: ({ offsetX, offsetY }) => update_position(offsetX, offsetY),
    onDragStart: () => {
      onDrag();
      // console.log($state.snapshot(position));
    },
  }}
  bind:this={window_div}
  onmousedown={() => {
    onDrag();
  }}
>
  <div class="window-titlebar" bind:this={titlebar}>
    {title}
    <div style="position: absolute; right: 0; top: 0;">
      <button
        class="window_button"
        onclick={() => {
          onHide();
        }}><img class="window-icon" src={min} alt="minimise" /></button
      >
      <button
        class="window_button"
        onclick={() => {
          if (fullscreen) {
            position = { ...last_non_fullscreen_position };
          } else {
            last_non_fullscreen_position = { ...position };
            update_position(0, 0);
          }
          fullscreen = !fullscreen;
        }}
        ><img class="window-icon" src={fullscreen_icon} alt="maximise" />
      </button>
      <button
        class="window_button"
        onclick={() => {
          onClose();
        }}
      >
        <img class="window-icon" src={close} alt="close" />
      </button>
    </div>
  </div>
  <slot />
</div>

<style>
  .window {
    position: absolute;
    --height: 450px;
    --width: 650px;
    width: var(--width);
    height: var(--height);
    border: var(--window-border);
    overflow: hidden;
    padding: 0px;
    background-color: var(--second-color);
    pointer-events: all;
    transition: opacity 1s;
    border-radius: var(--window-border-radius);
  }

  .window.hidden {
    opacity: 0;
  }

  .window-titlebar {
    width: 100%;
    height: 25px;
    font-size: 1rem;
    text-align: center;
    /*border-bottom: var(--border);*/
    position: relative;
    background-color: var(--accent-color);
  }

  .window_button {
    padding: 0;
    height: 100%;
    width: 25px;
    height: 25px;
    border: none;
    background-color: transparent;

    img {
      width: 15px;
      height: 15px;
      mix-blend-mode: lighten;
      margin: 5px;
      line-height: 0;
    }

    &:hover {
      backdrop-filter: brightness(1.1);
    }
  }
</style>
