<script lang="ts">
  type Button = {
    Text: string;
    Style: string;
    Show: number; // -1 == hide, 0==grey out, 1==show
    Onclick_function: () => void;
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
    --width: 200px;
    flex: 1 1 var(--width);
    width: var(--width);
    height: 100%;
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
    overflow: hidden;
  }

  .sidebar.horizontal {
    height: fit-content;
    width: 100%;
    flex-direction: row;
    border-bottom: var(--border);
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
    background-color: var(--accent-color2);
  }

  .sidebar_button:active {
    background-color: var(--accent-color);
    color: var(--second-color);
  }

  .selected,
  .selected:hover {
    background-color: rgba(0, 0, 0, 1);
    color: var(--second-color);
  }

  .sidebar_button.horizontal {
    margin: 0 0 0 0;
    text-align: center;
    flex: 1 1 fit-content;
    font-size: 1rem;
    padding: 5px;
  }
</style>
