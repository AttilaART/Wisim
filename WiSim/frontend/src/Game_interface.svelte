<script lang="ts">
  import Sidebar from "./Sidebar.svelte";
  import { fade, fly, slide } from "svelte/transition";
  import Window from "./window.svelte";
  import Finances from "./Finances.svelte";
  import { windows } from "./store.svelte";

  type window = {
    Id: number;
    Loaded: boolean;
  };

  let finance_window: window = $state({
    Id: undefined,
    Loaded: false,
  });
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
          Onclick_function: () => {},
          dont_keep_pressed: true,
        },
        {
          Text: "Finances",
          Style: "",
          Show: 1,
          Onclick_function: () => {
            finance_window.Loaded = true;
            windows[finance_window.Id].hidden = false;
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
    <div class="desktop">
      {#if finance_window.Loaded}
        <Window
          title="Finances"
          content={Finances}
          content_args={[]}
          bind:window_id={finance_window.Id}
          bind:loaded={finance_window.Loaded}
        ></Window>
      {/if}
    </div>
    <span class="bottom-bar">
      <span class="app_button">Windows:</span>
      {#each windows as w}
        <button
          transition:slide={{ axis: "x" }}
          onclick={() => {
            w.hidden = false;
          }}
          class="app_button {w.hidden ? '' : 'shown'}">{w.name}</button
        >
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
    width: fit-content;
    height: 36px;
    min-width: 10px;
    flex-direction: row;
    border-right: var(--border);
    border-top: var(--border);
  }

  .app_button {
    border: none;
    padding: 5px 10px 6px 10px;
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
