<script lang="ts">
  import GameInterface from "./Game_interface.svelte";
  import Sidebar from "./Sidebar.svelte";
  import Popup from "./Popup.svelte";
  import { New_simulation, Initial_app_load } from "../wailsjs/go/main/App";
  import { fade } from "svelte/transition";

  import { month_counter } from "./store.svelte";

  let background_image_blurred = $state("");
  let is_loading = $state(false);

  let mode = $state({
    main_menu: true,
    game_interface: false,
  });

  let Initial_app_load_promise = Initial_app_load();

  function load_singleplayer() {
    is_loading = true;
    start_new_game().then(() => {
      background_image_blurred = " blurred";
      mode.main_menu = false;
      mode.game_interface = true;
      is_loading = false;
    });
  }

  function load_main_menu() {
    background_image_blurred = "";
    mode.main_menu = true;
    mode.game_interface = false;
  }

  async function start_new_game() {
    $month_counter = await New_simulation();
  }
</script>

<main>
  {#if is_loading}
    <div>
      <Popup button_data={null} content={'<div class="loader"></div>'}></Popup>
    </div>
  {/if}
  <div
    style="position: absolute; height: 100vh; width: 100%; z-index: 0; overflow: hidden;"
  >
    <div
      class="background_image{background_image_blurred}"
      style="background-color: white;"
    ></div>

    <div
      style="position:absolute; width: 110%; height: 110%; left: -10px; top: -10px;
    background-image: linear-gradient(
      to right,
      rgba(0, 0, 0, 0.0),
      rgba(0, 0, 0, 0)
    );"
    ></div>
  </div>

  {#await Initial_app_load_promise}
    <Popup content={"<div class='loader'></div>"} button_data={null}></Popup>
  {:then}
    <div class="main_div" style="">
      {#if mode.main_menu}
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
                Onclick_function: load_singleplayer,
              },
              {
                Text: "Host game",
                Style: "",
                Show: 1,
                Onclick_function: () => {},
              },
              {
                Text: "Join game",
                Style: "",
                Show: 1,
                Onclick_function: () => {},
              },
              {
                Text: "Settings",
                Style: "margin-top: auto",
                Show: 1,
                Onclick_function: () => {},
              },
            ]}
          ></Sidebar>
        </div>
      {:else if mode.game_interface}
        <GameInterface></GameInterface>
      {/if}
    </div>
  {:catch error}
    <div>
      <Popup button_data={{ text: "OK" }} content={`<div>${error}</div>`}
      ></Popup>
    </div>
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

  .main_div {
    position: absolute;
    background-size: cover;
    background-position: center center;
    background-repeat: no-repeat;
    height: 100vh;
    max-height: 100vh;
    width: 100%;
    padding: 0px;
    margin: 0px;
  }

  .title_menu {
    display: flex;
    flex-direction: column;
    height: 100%;
    width: 50%;
    text-align: left;
    padding: 10px;
  }

  .menu_button.button {
    width: 200px;
    height: fit-content;
    padding: 10px;
    margin: 0px 0px 0px 0px;
    border-radius: 3;
  }

  .background_image {
    filter: blur(0);
    position: absolute;
    width: 110%;
    height: 110%;
    left: -10px;
    top: -10px;
    transition: filter 300ms;
  }

  .background_image.blurred {
    filter: blur(5px);
  }
</style>
