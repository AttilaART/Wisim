<script lang="ts">
  import GameInterface from "./Game_interface.svelte";
  import Sidebar from "./Sidebar.svelte";
  import { start_new_game, initial_app_load } from "./api.svelte";
  import { fade } from "svelte/transition";

  let popup: HTMLDialogElement = $state();
  let mode = $state("main_menu");

  let loading_promise = $state(initial_app_load());

  async function load_singleplayer(): Promise<void> {
    await start_new_game();
    mode = "game";
  }
</script>

<main>
  {#await loading_promise}
    <dialog open><div class="loader"></div></dialog>
  {:then}
    {#if mode == "main_menu"}
      <div
        class="title_menu"
        out:fade={{ duration: 300 }}
        in:fade={{ duration: 300, delay: 300 }}
      >
        <h1 style="padding: 0 8px;">WiSim</h1>

        <Sidebar
          expand={true}
          buttons={[
            {
              Text: "Singleplayer",
              Style: "",
              Show: 1,
              onClick: () => {
                loading_promise = load_singleplayer();
              },
            },
            {
              Text: "Host game",
              Style: "grayed-out",
              Show: 0,
              onClick: () => {},
            },
            {
              Text: "Join game",
              Style: "",
              Show: 0,
              onClick: () => {},
            },
            {
              Text: "Settings",
              Style: "margin-top: auto",
              Show: 0,
              onClick: () => {},
            },
          ]}
        ></Sidebar>
      </div>
    {:else if mode == "game"}
      <GameInterface></GameInterface>
    {/if}
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
</main>

<style>
  main {
    margin: 0;
    padding: 0;
  }
  * {
    box-sizing: border-box;
  }

  .title_menu {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 50%;
    text-align: left;
    padding: 10px;
  }
</style>
