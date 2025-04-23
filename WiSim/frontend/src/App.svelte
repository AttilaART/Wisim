<script lang="ts">
  import background_image from "./assets/images/Temp_background_image.png";
  import GameInterface from "./Game_interface.svelte";
  import Sidebar from "./Sidebar.svelte";
  import Popup from "./Popup.svelte";
  import { New_simulation } from "../wailsjs/go/main/App";

  import { month_counter, loading } from "./store";

  let background_image_blur = $state(0);
  let is_loading = $state($loading);

  let mode = $state({
    main_menu: true,
    game_interface: false,
  });

  function load_singleplayer() {
    $loading = true;
    start_new_game().then(() => {
      background_image_blur = 4;
      mode.main_menu = false;
      mode.game_interface = true;
      $loading = false;
    });
  }

  function load_main_menu() {
    background_image_blur = 0;
    mode.main_menu = true;
    mode.game_interface = false;
  }

  async function start_new_game() {
    $month_counter = await New_simulation();
  }
</script>

<main>
  <!--{#if is_loading}
    <Popup button_data={null} content={'<div class="loader"></div>'}></Popup>
  {/if}-->
  <div
    style="position: absolute; height: 100vh; width: 100%; z-index: 0; overflow: hidden;"
  >
    <div
      style="background-image: url({background_image}); filter: blur({background_image_blur}px); position: absolute; width: 110%; height: 110%; left: -10px; top: -10px"
    ></div>
  </div>
  <div class="main_div" style="">
    {#if mode.main_menu}
      <div class="main_menu">
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
    {/if}
    {#if mode.game_interface}
      <GameInterface return_function={load_main_menu}></GameInterface>
    {/if}
  </div>
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

  .main_menu {
    display: flex;
    flex-direction: column;
    background-image: linear-gradient(
      to right,
      rgba(0, 0, 0, 0.5),
      rgba(0, 0, 0, 0)
    );
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
</style>
