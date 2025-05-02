<script lang="ts">
  interface params {
    options: Option[];
    default_option_index?: number;
    style?: string;
  }

  interface Option {
    id: string;
    text: string;
    style?: string;
    index?: number;
    height?: number;
    onselect: () => void;
  }

  let { options, default_option_index, style }: params = $props();

  let selected_index: number = $state(
    default_option_index ? default_option_index : 0,
  );
  let menu_open: boolean = $state(false);

  let menu_height: number = $state(0);

  let heights: number[] = $state();

  for (let i in options) {
    Object.defineProperty(options[i], "index", { value: i });
  }

  function open_menu() {
    menu_open = true;
  }

  function get_heights() {
    let heights: number[] = [];
    heights.push(options[selected_index].height);

    for (let o of options) {
      if (o.index != selected_index) {
        heights.push(o.height);
      }
    }
  }

  function select(option_index: number) {
    selected_index = option_index;
    menu_open = false;
  }
</script>

<div
  style="display: inline; {style ? style : 'text-decoration: underline;'}"
  aria-roledescription="custom selector"
>
  {#if menu_open}
    <div style="position: relative; right: 200px">
      {@render option_snippet(options[selected_index])}
      {#each options as option}
        {#if option.index != selected_index}
          {@render option_snippet(option)}
        {/if}
      {/each}
    </div>
  {:else}
    <div style="position: relative; right: 200px;">
      {@render option_snippet(options[selected_index])}
    </div>
  {/if}
</div>

{#snippet option_snippet(option: Option)}
  <button
    class="option"
    style="position: absolute;{option.style ? option.style : ''}"
    onclick={() => {
      if (menu_open) {
        select(option.index ? option.index : options.indexOf(option));
        option.onselect();
      } else {
        menu_open = true;
      }
    }}
  >
    {option.text}
  </button>
{/snippet}

<style>
  button {
    all: unset;
  }
  button:hover {
    all: unset;
  }

  .option {
    position: relative;
  }
</style>
