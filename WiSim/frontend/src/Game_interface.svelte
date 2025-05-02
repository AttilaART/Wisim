<script lang="ts">
  import Sidebar from "./Sidebar.svelte";
  import { fade, fly, slide } from "svelte/transition";
  import Window from "./window.svelte";
  import Finances from "./Finances.svelte";
</script>

<div
  class="game_interface"
  style="display: flex; flex-direction: row;"
  in:fade={{ duration: 300, delay: 300 }}
  out:fade={{ duration: 300 }}
>
  <div
    style="width: fit-content; height: 100%; border-right: 3px solid black; display: flex; flex-direction: column;"
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
          Onclick_function: () => {},
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
        <span
          style="height: 100%; border-left: 3px solid black; margin-right: 10px;"
        ></span>
        <button style=" height: 100%; border: none;">Messages</button>
      </div>
    </div>
    <div class="desktop">
      <Window title="Finances" content={Finances} content_args={[]}></Window>
    </div>
  </div>
</div>

<style>
  .game_interface {
    height: 100%;
    width: calc(100%);
    display: grid;
    overflow: unset;
  }
  .top-bar {
    padding: 5px 10px 5px 10px;
    width: 100%;
    display: flex;
    flex-direction: row;
    align-items: center;
    border-bottom: 3px solid black;
  }

  .desktop {
    height: 100%;
    width: 100%;
    padding: 0;
  }
</style>
