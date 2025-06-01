<script lang="ts">
  import Sidebar from "./Sidebar.svelte";
  import { fade, fly, slide } from "svelte/transition";
  import Window from "./window.svelte";
  import Finances from "./Finances.svelte";
  import {
    delete_window,
    get_window_by_id,
    new_window,
    windows,
  } from "./window_manager.svelte";
  import Close from "./assets/images/Close.svelte";
  import Marketing from "./Marketing.svelte";
  import { decisions } from "./store.svelte";
  import lttwalpaper from "./assets/images/lttwalpaper.jpeg";
  import Production from "./Production.svelte";

  type window = {
    Id: number;
    Loaded: boolean;
    open_window: () => void;
  };

  let finance_window: window = $state({
    Id: undefined,
    Loaded: false,
    open_window: open_window,
  });
  let marketing_window: window = $state({
    Id: undefined,
    Loaded: false,
    open_window: open_window,
  });
  let production_window: window = $state({
    Id: undefined,
    Loaded: false,
    open_window: open_window,
  });

  let desktop_canvas_size: ResizeObserverSize[] = $state();

  $effect(() => {
    console.log("decisions have been updated");
    console.log($state.snapshot(decisions));
  });

  function open_window() {
    this.Loaded = true;
    try {
      windows[get_window_by_id(this.Id).index].hidden = false;
    } catch (exception) {
      console.warn(exception);
    }
  }
</script>

<div
  class="game_interface"
  style="display: flex; flex-direction: row;"
  in:fade={{ duration: 300, delay: 300 }}
  out:fade={{ duration: 300 }}
>
  <div class="sidebar">
    <div style="flex: 0 0 100px;"></div>
    <Sidebar
      expand={true}
      keep_pressed={true}
      buttons={[
        {
          Text: "Employees",
          Style: "",
          Show: 1,
          onClick: () => {},
          dont_keep_pressed: true,
        },
        {
          Text: "Production",
          Style: "",
          Show: 1,
          onClick: () => {
            production_window.open_window();
          },
          dont_keep_pressed: true,
        },
        {
          Text: "Marketing",
          Style: "",
          Show: 1,
          onClick: () => {
            marketing_window.open_window();
          },
          dont_keep_pressed: true,
        },
        {
          Text: "Finances",
          Style: "",
          Show: 1,
          onClick: () => {
            finance_window.open_window();
          },
          dont_keep_pressed: true,
        },
        {
          Text: "Research",
          Style: "",
          Show: 1,
          onClick: () => {},
          dont_keep_pressed: true,
        },
        {
          Text: "Companies",
          Style: "",
          Show: 1,
          onClick: () => {},
          dont_keep_pressed: true,
        },
        {
          Text: "Main Menu",
          Style: "margin-top: auto",
          Show: 1,
          onClick: () => {},
          dont_keep_pressed: true,
        },
      ]}
    ></Sidebar>
  </div>
  <div style="display: flex; flex-direction: column; 100%; width: 100%;">
    <div class="top-bar">
      <div style="flex: 1 1 ">Balance: 100'000 CHF</div>
      <div style="flex: 1 1 ">01/02/0001</div>
      <div style="flex: 1 0 ">
        Time until next step: <span style="color: red;">5 min</span>
        <button>Ready</button>
      </div>
      <div style="flex: 0 0 fit-content; height: 100%;">
        <button style=" height: 100%; border: none;">Messages</button>
      </div>
    </div>
    <div
      class="desktop"
      style="background: no-repeat url({lttwalpaper}); background-size: cover;"
      bind:contentBoxSize={desktop_canvas_size}
    >
      {#if finance_window.Loaded}
        <Window
          title="Finances"
          content={Finances}
          content_args={[]}
          canvas_size={{
            x: desktop_canvas_size[0].inlineSize,
            y: desktop_canvas_size[0].blockSize,
          }}
          onClose={() => (finance_window.Loaded = false)}
          bind:window_id={finance_window.Id}
        ></Window>
      {/if}
      {#if marketing_window.Loaded}
        <Window
          title="Marketing"
          content={Marketing}
          content_args={[]}
          canvas_size={{
            x: desktop_canvas_size[0].inlineSize,
            y: desktop_canvas_size[0].blockSize,
          }}
          onClose={() => (marketing_window.Loaded = false)}
          bind:window_id={marketing_window.Id}
        ></Window>
      {/if}
      {#if production_window.Loaded}
        <Window
          title="Production"
          content={Production}
          content_args={[]}
          canvas_size={{
            x: desktop_canvas_size[0].inlineSize,
            y: desktop_canvas_size[0].blockSize,
          }}
          onClose={() => (production_window.Loaded = false)}
          bind:window_id={production_window.Id}
        ></Window>
      {/if}
    </div>
    <span class="bottom-bar">
      {#each windows as w}
        <div
          style="height: calc(var(--height) - var(--border-width); width: 200px; position: relative;"
          transition:slide={{ axis: "x" }}
          class="app_button {w.hidden ? '' : 'shown'}"
        >
          <button
            style="width: fit-content; border: none; mix-blend-mode: difference; background-color: transparent; position: absolute; right: 0; margin-top: 2px;"
            onclick={() => {
              delete_window(w.id);
            }}><Close></Close></button
          >
          <button
            onclick={() => {
              w.hidden = false;
            }}
            style="padding: auto 10px auto 10px; border: none; background-color: inherit; color: inherit; margin-right: 10px; margin-top: 2px;"
            >{w.name}
          </button>
        </div>
      {/each}
    </span>
  </div>
</div>

<style>
  .game_interface {
    height: 100%;
    width: calc(100%);
    display: grid;
    overflow: unset;
    background-color: var(--second-color);
  }
  .top-bar {
    padding: 5px 10px 5px 10px;
    width: calc(100% - 20px);
    display: flex;
    flex-direction: row;
    align-items: center;
    border-bottom: var(--border);
  }

  .bottom-bar {
    display: flex;
    width: 100%;
    --height: 35px;
    height: var(--height);
    flex-direction: row;
    border-top: var(--border);
    overflow-x: scroll;
  }

  .app_button {
    position: relative;
    border: none;
    transition: all 0.25s;
  }

  .app_button.shown {
    color: var(--second-color);
    background-color: var(--main-color);
  }

  .desktop {
    height: 100%;
    width: 100%;
    padding: 0;
  }

  .sidebar {
    position: absolute;
    width: fit-content;
    height: 100%;
    border-right: var(--border-width) solid var(--border-color);
    display: flex;
    flex-direction: column;
    --button-width: 20px;
  }

  .sidebar:hover {
    --button-width: 200px;
  }
</style>
