<script lang="ts">
  import Sidebar from "./Sidebar.svelte";
  import { fade, fly, slide } from "svelte/transition";
  import Window from "./Window.svelte";
  import Finances from "./Finances.svelte";
  import Close from "./assets/images/Close.svelte";
  import Marketing from "./Marketing.svelte";
  import { company, current_decisions, month } from "./store.svelte";
  import lttwalpaper from "./assets/images/lttwalpaper.jpeg";
  import Production from "./Production.svelte";
  import Employees from "./Employees.svelte";
  import { trigger_simulation } from "./api.svelte";
  import { format_number } from "./helper.svelte";
  import { latest_reports } from "./store.svelte";

  let windows: { [key: string]: { Zindex: number; Hidden: boolean } } = $state(
    {},
  );

  let desktop_canvas_size: ResizeObserverSize[] = $state();

  $effect(() => {
    console.log("decisions have been updated");
    console.log($state.snapshot(current_decisions));
  });

  function new_window(name: string) {
    console.log(windows[name]);
    if (windows[name] == undefined) {
      let max_Zindex: number = Math.min(
        ...Object.values(windows).map((value) => {
          return value.Zindex;
        }),
      );
      if (max_Zindex == Infinity) max_Zindex = 0;
      windows[name] = { Zindex: max_Zindex + 1, Hidden: false };
    }
  }

  function move_window_to_top(name: string) {
    let current_z_index = windows[name].Zindex;
    for (let i in windows) {
      if (windows[i].Zindex >= current_z_index && i != name) {
        windows[i].Zindex -= 1;
      }
    }
    windows[name].Zindex = Object.keys(windows).length;
    // console.log($state.snapshot(windows))
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
          onClick: () => {
            new_window("employees");
          },
          dont_keep_pressed: true,
        },
        {
          Text: "Production",
          Style: "",
          Show: 1,
          onClick: () => {
            new_window("production");
          },
          dont_keep_pressed: true,
        },
        {
          Text: "Marketing",
          Style: "",
          Show: 1,
          onClick: () => {
            new_window("marketing");
          },
          dont_keep_pressed: true,
        },
        {
          Text: "Finances",
          Style: "",
          Show: 1,
          onClick: () => {
            new_window("finances");
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
      <div style="flex: 1 1 ">
        Balance: {company.Balance
          ? format_number(company.Balance, true, 2)
          : "Loading..."}
      </div>
      <div style="flex: 1 1 ">Month {$month}</div>
      <div style="flex: 1 0 ">
        Time until next step: <span style="color: red;">5 min</span>
        <button
          onclick={() => {
            trigger_simulation(true);
          }}>Ready</button
        >
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
      {#if windows["finances"] != undefined}
        <Window
          title="Finances"
          canvas_size={{
            x: desktop_canvas_size[0].inlineSize,
            y: desktop_canvas_size[0].blockSize,
          }}
          onClose={() => delete windows["finances"]}
          onDrag={() => move_window_to_top("finances")}
          onHide={() => (windows["finances"].Hidden = true)}
          bind:Zindex={windows["finances"].Zindex}
          bind:Hidden={windows["finances"].Hidden}
        >
          <Finances></Finances>
        </Window>
      {/if}
      {#if windows["marketing"] != undefined}
        <Window
          title="Marketing"
          canvas_size={{
            x: desktop_canvas_size[0].inlineSize,
            y: desktop_canvas_size[0].blockSize,
          }}
          onClose={() => delete windows["marketing"]}
          onDrag={() => move_window_to_top("marketing")}
          onHide={() => (windows["marketing"].Hidden = true)}
          bind:Zindex={windows["marketing"].Zindex}
          bind:Hidden={windows["marketing"].Hidden}
        >
          <Marketing></Marketing></Window
        >
      {/if}
      {#if windows["production"] != undefined}
        <Window
          title="Production"
          canvas_size={{
            x: desktop_canvas_size[0].inlineSize,
            y: desktop_canvas_size[0].blockSize,
          }}
          onClose={() => delete windows["production"]}
          onDrag={() => move_window_to_top("production")}
          onHide={() => (windows["production"].Hidden = true)}
          bind:Zindex={windows["production"].Zindex}
          bind:Hidden={windows["production"].Hidden}
        >
          <Production></Production></Window
        >
      {/if}
      {#if windows["employees"]}
        <Window
          title="Employees"
          canvas_size={{
            x: desktop_canvas_size[0].inlineSize,
            y: desktop_canvas_size[0].blockSize,
          }}
          onClose={() => delete windows["employees"]}
          onDrag={() => move_window_to_top("employees")}
          onHide={() => (windows["employees"].Hidden = true)}
          bind:Zindex={windows["employees"].Zindex}
          bind:Hidden={windows["employees"].Hidden}
          ><Employees></Employees></Window
        >
      {/if}
    </div>
    <span class="bottom-bar">
      {#each Object.keys(windows) as w}
        <div
          style="height: calc(var(--height) - var(--border-width); width: 200px; position: relative;"
          transition:slide={{ axis: "x" }}
          class="app_button {windows[w].Hidden ? '' : 'shown'}"
        >
          <button
            style="width: fit-content; border: none; mix-blend-mode: difference; background-color: transparent; position: absolute; right: 0; margin-top: 2px;"
            onclick={() => {
              delete windows[w];
            }}><Close></Close></button
          >
          <button
            onclick={() => {
              windows[w].Hidden = false;
            }}
            style="padding: auto 10px auto 10px; border: none; background-color: inherit; color: inherit; margin-right: 10px; margin-top: 2px;"
            >{w}
          </button>
        </div>
      {/each}
    </span>
  </div>
</div>

<style>
  .game_interface {
    height: 100vh;
    width: calc(100%);
    display: grid;
    background-color: var(--second-color);
  }
  .top-bar {
    padding: 5px 10px 5px 10px;
    width: calc(100% - 20px);
    display: flex;
    flex-direction: row;
    align-items: center;
    border-bottom: var(--border-thin);
  }

  .bottom-bar {
    display: flex;
    width: 100%;
    --height: 35px;
    height: var(--height);
    flex-direction: row;
    overflow-x: scroll;
    background-color: var(--second-color);
    border-top: var(--border-thin);
    z-index: 9999;
  }

  .app_button {
    position: relative;
    border: none;
    padding: 0 20px;
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
    display: flex;
    flex-direction: column;
    position: absolute;
    height: 100%;
    border-right: var(--border-thin);
    transform: translateX(-190px);
    background-color: var(--second-color);
    z-index: 10000;
    transition: 1s;
  }

  .sidebar:hover {
    transform: translateX(0px);
  }
</style>
