<script lang="ts">
  type Button = {
    Text: string;
    Style: string;
    Show: number; // -1 == hide, 0==grey out, 1==show
    Onclick_function: () => void;
    dont_keep_pressed?: boolean;
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

  let button_selection = $state(
    (() => {
      let button_selection = [];
      for (let _ in buttons) {
        button_selection.push(false);
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
    if (horisontal) {
      if (
        button_selection[index] &&
        keep_pressed &&
        !buttons[index].dont_keep_pressed
      ) {
        return "sidebar_button horizontal selected";
      } else {
        return "sidebar_button horizontal";
      }
    } else {
      if (
        button_selection[index] &&
        keep_pressed &&
        !buttons[index].dont_keep_pressed
      ) {
        return "sidebar_button selected";
      } else {
        return "sidebar_button";
      }
    }
  }

  let style = $state("");
  if (!expand) {
    style = "flex: 0 0 auto;";
  }
</script>

{#if !horisontal}
  <div class="sidebar" {style}>
    {#each buttons as button, index}
      {@render sidebar_button(button, index)}
    {/each}
  </div>
{:else}
  <div class="sidebar horizontal" {style}>
    {#each buttons as button, index}
      {@render sidebar_button(button, index)}
    {/each}
  </div>
{/if}

{#snippet sidebar_button(button_data: Button, index: number)}
  {#if button_data.Show == 1}
    <button
      class={set_class(index)}
      style={button_data.Style}
      onclick={() => {
        select_button(index);
        button_data.Onclick_function();
      }}>{@html button_data.Text}</button
    >
  {:else}
    <button
      class={set_class(index)}
      style="opacity: 60%; {button_data.Style}"
      onclick={() => {
        select_button(index);
        button_data.Onclick_function();
      }}>{@html button_data.Text}</button
    >
  {/if}
{/snippet}

<style>
  .sidebar {
    width: fit-content;
    display: flex;
    flex-direction: column;
    flex: 1 1 100%;
    flex-wrap: wrap;
    background-color: rgba(0, 0, 0, 0.5);
    border-radius: 10px;
    margin: 10px 10px 10px 10px;
    overflow: hidden;
  }

  .sidebar.horizontal {
    height: fit-content;
    width: calc(100% - 20px);
    flex-direction: row;
  }

  .sidebar_button {
    width: 200px;
    height: fit-content;
    padding: 10px;
    margin: 0px 0px 0px 0px;
    border-radius: 0;
  }

  .sidebar_button:hover {
    background-color: rgba(255, 255, 255, 0.3);
  }

  .sidebar_button:active {
    background-color: rgba(255, 255, 255, 1);
    color: #000;
  }

  .sidebar_button.selected {
    background-color: rgba(255, 255, 255, 0.5);
    font-weight: bold;
  }

  .sidebar_button.horizontal {
    text-align: center;
    flex: 1 1 fit-content;
  }

  .sidebar_button.horizontal.selected {
    background-color: rgba(255, 255, 255, 0.5);
    font-weight: bold;
  }
</style>
