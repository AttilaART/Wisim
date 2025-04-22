<script lang="ts">
  import background_image from "./assets/images/Temp_background_image.png";
  import GameInterface from "./Game_interface.svelte";
  import Sidebar from "./Sidebar.svelte";

  let background_image_blur = $state(0);

  let mode = $state({
    main_menu: true,
    game_interface: false,
  });

  function load_singleplayer() {
    background_image_blur = 4;
    mode.main_menu = false;
    mode.game_interface = true;
  }

  function load_main_menu() {
    background_image_blur = 0;
    mode.main_menu = true;
    mode.game_interface = false;
  }
</script>

<main>
  <div
    style="position: absolute; height: 100vh; width: 100%; z-index: 0; overflow: hidden;"
  >
    <div
      style="background-image: url({background_image}); filter: blur({background_image_blur}px); position: absolute; width: 110%; height: 110%; left: -10px; top: -10px"
    ></div>
  </div>
  <div class="main_div" style="z-index: 20">
    {#if mode.main_menu}
      <div class="main_menu">
        <h1 style="padding: 0 8px;">WiSim</h1>

        <Sidebar
          expand={true}
          buttons={[
            {
              Text: "Singleplayer",
              Style: "",
              Onclick_function: load_singleplayer,
            },
            {
              Text: "Host game",
              Style: "",
              Onclick_function: () => {},
            },
            {
              Text: "Join game",
              Style: "",
              Onclick_function: () => {},
            },
            {
              Text: "Settings",
              Style: "margin-top: auto",
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
