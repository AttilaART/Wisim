<script lang="ts">
  import Sidebar from "./Sidebar.svelte";
  import { fade, fly, slide } from "svelte/transition";
  import Window from "./window.svelte";
  import Finances from "./Finances.svelte";
  import { delete_window, windows } from "./window_manager.svelte";
  import Close from "./assets/images/Close.svelte";
  import Marketing from "./Marketing.svelte";

  type window = {
    Id: number;
    Loaded: boolean;
  };

  let finance_window: window = $state({
    Id: undefined,
    Loaded: false,
  });
  let marketing_window: window = $state({
    Id: undefined,
    Loaded: false,
  });

  let desktop_canvas_size: ResizeObserverSize[] = $state();
</script>

<div
  class="game_interface"
  style="display: flex; flex-direction: row;"
  in:fade={{ duration: 300, delay: 300 }}
  out:fade={{ duration: 300 }}
>
  <div
    style="width: fit-content; height: 100%; border-right: var(--border-width) solid var(--border-color); display: flex; flex-direction: column;"
  >
    <div style="flex: 0 0 100px;"></div>
    <Sidebar
      expand={true}
      keep_pressed={true}
      buttons={[
        {
          Text: "Employees",
          Style: "",
          Show: 1,
          Onclick_function: () => {},
          dont_keep_pressed: true,
        },
        {
          Text: "Production",
          Style: "",
          Show: 1,
          Onclick_function: () => {},
          dont_keep_pressed: true,
        },
        {
          Text: "Marketing",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            marketing_window.Loaded = true;
            try {
              windows[marketing_window.Id].hidden = false;
            } catch (exception) {
              console.warn(exception);
            }
          },
          dont_keep_pressed: true,
        },
        {
          Text: "Finances",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            finance_window.Loaded = true;
            try {
              windows[finance_window.Id].hidden = false;
            } catch (exception) {
              console.warn(exception);
            }
          },
          dont_keep_pressed: true,
        },
        {
          Text: "Research",
          Style: "",
          Show: 1,
          Onclick_function: () => {},
          dont_keep_pressed: true,
        },
        {
          Text: "Companies",
          Style: "",
          Show: 1,
          Onclick_function: () => {},
          dont_keep_pressed: true,
        },
        {
          Text: "Main Menu",
          Style: "margin-top: auto",
          Show: 1,
          Onclick_function: () => {},
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
    <div class="desktop" bind:contentBoxSize={desktop_canvas_size}>
      {#if finance_window.Loaded}
        <Window
          title="Finances"
          content={Finances}
          content_args={[]}
          canvas_size={{
            x: desktop_canvas_size[0].inlineSize,
            y: desktop_canvas_size[0].blockSize,
          }}
          close_window={() => (finance_window.Loaded = false)}
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
          close_window={() => (marketing_window.Loaded = false)}
          bind:window_id={marketing_window.Id}
        ></Window>
      {/if}
    </div>
    <span class="bottom-bar">
      {#each windows as w, w_id}
        <div
          style="height: calc(var(--height) - var(--border-width); width: 200px; position: relative;"
          transition:slide={{ axis: "x" }}
          class="app_button {w.hidden ? '' : 'shown'}"
        >
          <button
            style="width: fit-content; border: none; mix-blend-mode: difference; background-color: transparent; position: absolute; right: 0; margin-top: 2px;"
            onclick={() => {
              delete_window(w_id);
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
</style>
