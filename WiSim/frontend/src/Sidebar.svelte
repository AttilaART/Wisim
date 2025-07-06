<script lang="ts">
  import { preferences } from "./store.svelte";

  type Button = {
    Text: string;
    Style: string;
    Show: number; // -1 == hide, 0==grey out, 1==show
    onClick: () => void;
    dont_keep_pressed?: boolean;
    selected_by_default?: boolean;
  };
  const {
    buttons,
    expand,
    keep_pressed,
    horisontal,
  }: {
    buttons: Button[];
    expand: boolean;
    keep_pressed?: boolean;
    horisontal?: boolean;
  } = $props();

  let force_show_expanded_tabs: boolean = $state(
    (() => {
      if (preferences.hide_tabs != "always") return true;
      else return false;
    })(),
  );

  let button_selection: boolean[] = $state(
    (() => {
      let button_selected_by_default: boolean = false;
      let button_selection = [];
      for (let i in buttons) {
        if (!buttons[i].selected_by_default) {
          if (button_selected_by_default) {
            console.warn("multiple buttons selected by default");
          }
          button_selected_by_default = true;
          button_selection.push(false);
        } else {
          button_selection.push(true);
        }
      }
      return button_selection;
    })(),
  );

  function select_button(index: number) {
    if (typeof index === "undefined") {
      console.warn("index has type undefined");
    }
    if (!buttons[index].dont_keep_pressed) {
      for (let i in button_selection) {
        button_selection[i] = false;
      }
      button_selection[index] = true;
    }
  }

  function set_class(index: number) {
    let styling: string;
    if (horisontal) {
      if (
        button_selection[index] &&
        keep_pressed &&
        !buttons[index].dont_keep_pressed
      ) {
        styling = "sidebar_button horizontal selected";
      } else {
        styling = "sidebar_button horizontal";
      }
    } else {
      if (
        button_selection[index] &&
        keep_pressed &&
        !buttons[index].dont_keep_pressed
      ) {
        styling = "sidebar_button selected";
      } else {
        styling = "sidebar_button";
      }
    }
    return styling;
  }

  let style = $state("");
  if (!expand) {
    style = "flex: 0 0 auto;";
  }
</script>

{#if !horisontal}
  <div class="sidebar" {style}>
    <hr />
    {#each buttons as button, index}
      {@render sidebar_button(button, index)}
    {/each}
  </div>
{:else}
  <div
    class="sidebar_horisontal_container {force_show_expanded_tabs
      ? 'initial'
      : ''}"
    onmouseleave={() => {
      if (
        preferences.hide_tabs == "after_hover" ||
        preferences.hide_tabs == "always"
      )
        force_show_expanded_tabs = false;
    }}
  >
    <div class="sidebar horizontal" {style}>
      {#each buttons as button, index}
        {@render sidebar_button(button, index)}
      {/each}
    </div>
    <div id="hover_thing">⌄</div>
  </div>
{/if}

{#snippet sidebar_button(button_data: Button, index: number)}
  {#if button_data.Show == 1}
    <button
      class={set_class(index)}
      style={button_data.Style}
      onclick={() => {
        select_button(index);
        button_data.onClick();
      }}>{@html button_data.Text}</button
    >
  {:else}
    <button
      class={set_class(index)}
      style="opacity: 60%; {button_data.Style}"
      onclick={() => {
        select_button(index);
        button_data.onClick();
      }}>{@html button_data.Text}</button
    >
  {/if}
{/snippet}

<style>
  /* Generig styling */

  .sidebar {
    --width: 200px;
    flex: 1 1 var(--width);
    width: var(--width);
    height: 100%;
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
    overflow: hidden;
  }

  .sidebar_button {
    width: 100%;
    height: fit-content;
    padding: 10px;
    padding-left: 20px;
    margin: 10px 0px 10px 0px;
    text-align: left;
    border: none;
  }

  .sidebar_button:hover {
    background-color: var(--accent-color);
  }

  .sidebar_button:active {
    background-color: var(--accent-color2);
    color: var(--second-color);
  }

  .selected,
  .selected:hover {
    background-color: var(--main-color);
    color: var(--second-color);
  }

  /* horizontal sidebar */

  .sidebar_horisontal_container.initial #hover_thing,
  .sidebar_horisontal_container:hover #hover_thing {
    opacity: 0%;
    height: 0px;
  }

  .sidebar_horisontal_container {
    height: 25px;
    width: 100%;
    transition: 0.5s;
    position: relative;
    border-bottom: var(--window-border);
    overflow: hidden;
  }

  .sidebar_horisontal_container.initial,
  .sidebar_horisontal_container:hover {
    width: 100%;
    transition: 0.5s;
    height: 60px;
  }

  #hover_thing {
    position: absolute;
    transition: 0.5s;
    height: 10px;
    bottom: 0px;
    right: 50%;
    transform: translate(50%, -20px);
    font-size: 20px;
    text-align: center;
  }

  .sidebar.horizontal {
    width: calc(100% - 20px);
    flex-direction: row;
    /*border-bottom: var(--window-border);*/
    padding: 10px;
    background-color: none;
    opacity: 0%;
    transition: 1s;
  }

  .sidebar_horisontal_container.initial .sidebar.horizontal,
  .sidebar_horisontal_container:hover .sidebar.horizontal {
    opacity: 100%;
    transition: 0.5s;
  }

  .sidebar_button.horizontal {
    margin: 5px;
    text-align: center;
    flex: 1 1 fit-content;
    font-size: 1rem;
    padding: 5px;
    transition:
      background 0.25s,
      opacity 0.8s;
  }
</style>
