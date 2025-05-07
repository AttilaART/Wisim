<script lang="ts">
  import Info from "./assets/images/Info.svelte";
  import { fade } from "svelte/transition";

  let { text, position }: { text: string; position?: string } = $props();
  let show_text: boolean = $state(false);

  if (position == undefined) {
    position = "bottom";
  }
</script>

<span
  style="position: relative; display: inline;"
  onmouseover={() => {
    show_text = true;
  }}
  onmouseleave={() => {
    show_text = false;
  }}
>
  <Info></Info>
  {#if show_text}
    {#if position == "top"}
      <div
        class="tooltip-text"
        transition:fade={{ duration: 250 }}
        style="bottom: var(--offset); transform:  translateX(50%);"
      >
        {text}
      </div>
    {:else if position == "left"}
      <div
        class="tooltip-text"
        transition:fade={{ duration: 250 }}
        style="right: var(--offset); transform:  translateX(50%);"
      >
        {text}
      </div>
    {:else if position == "right"}
      <div
        class="tooltip-text"
        transition:fade={{ duration: 250 }}
        style="left: var(--offset); transform:  translateX(50%);"
      >
        {text}
      </div>
    {:else if position == "bottom"}
      <div
        class="tooltip-text"
        transition:fade={{ duration: 250 }}
        style="top: var(--offset); transform:  translateX(50%);"
      >
        {text}
      </div>
    {/if}
  {/if}
</span>

<style>
  .tooltip-text {
    font: initial;
    position: absolute;
    --offset: 30px;
    width: 150px;
    text-align: left;
    background: var(--accent-color);
    opacity: 0.95;
    border-radius: 5px;
    padding: 5px;
  }
</style>
