<script lang="ts" module>
  let windows: { [key: string]: { Zindex: number; Hidden: boolean } } = $state(
    {},
  );

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

  export function new_window(name: string) {
    console.log(`Opening ${name} window`);
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
</script>

<script lang="ts">
  import { fade, fly, slide } from "svelte/transition";
  import Window from "./Window.svelte";
  import Finances from "./Finances.svelte";
  import Marketing from "./Marketing.svelte";
  import { company, current_decisions, loading, month } from "./store.svelte";
  import wallpaper from "./assets/images/gnomewallpaper.jpeg";
  import Production from "./Production.svelte";
  import Employees from "./Employees.svelte";
  import MenuBar from "./MenuBar.svelte";
  import { trigger_simulation } from "./api.svelte";

  let loading_promise = $state();
  let popup: HTMLDialogElement = $state();

  let desktop_canvas_size: ResizeObserverSize[] = $state();

  $effect(() => {
    console.log("decisions have been updated");
    console.log($state.snapshot(current_decisions));
    console.log("company data edited");
    console.log($state.snapshot(company));
  });
</script>

{#await loading_promise}
  <dialog open><div class="loader"></div></dialog>
{:catch error}
  <dialog open bind:this={popup}>
    <article>
      {error}
      <footer>
        <button
          onclick={() => {
            popup.close();
          }}>OK</button
        >
      </footer>
    </article>
  </dialog>
{/await}

<div
  id="game_interface"
  style="display: flex; flex-direction: row;"
  in:fade={{ duration: 300, delay: 300 }}
  out:fade={{ duration: 300 }}
>
  <button id="ready" onclick={() => (loading_promise = trigger_simulation())}
    >Ready?</button
  >
  <div
    id="desktop"
    style="background: no-repeat url({wallpaper}); background-size: cover;"
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
  <MenuBar></MenuBar>
</div>

<style>
  #game_interface {
    height: 100vh;
    width: calc(100%);
    display: grid;
    background-color: var(--second-color);
  }

  #desktop {
    height: 100%;
    width: 100%;
    padding: 0;
  }

  #ready {
    position: fixed;
    right: 1rem;
    top: 1rem;
    padding: 10px;
    opacity: 0.9;
  }

  #ready:hover,
  #ready:active {
    opacity: 1;
  }
</style>
