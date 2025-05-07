<script lang="ts">
  import Info from "./assets/images/Info.svelte";
  import { fade } from "svelte/transition";

  let { text, position }: { text: string; position?: string } = $props();
  let show_text: boolean = $state(false);

  if (position == undefined) {
    position = "bottom";
  }

  let offset_from: string;

  switch (position) {
    case "top":
      offset_from = "bottom";
      break;
    case "bottom":
      offset_from = "top";
      break;
    case "left":
      offset_from = "right";
      break;
    case "right":
      offset_from = "left";
      break;
  }
</script>

<!-- svelte-ignore a11y_mouse_events_have_key_events-->
<span
  style="position: relative; display: inline;"
  onmouseover={() => {
    show_text = true;
  }}
  onmouseleave={() => {
    show_text = false;
  }}
  role="tooltip"
>
  <Info></Info>
  {#if show_text}
    <div
      class="tooltip-text"
      transition:fade={{ duration: 250 }}
      style="{offset_from}: var(--offset); transform:  translateX(50%);"
    >
      {text}
    </div>
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
