<script lang="ts">
  type Button = {
    Text: string;
    Style: string;
    Onclick_function: () => void;
  };
  const { buttons, expand }: { buttons: Button[]; expand: boolean } = $props();

  let style = $state("");
  if (!expand) {
    style = "flex: 0 0 auto;";
  }
</script>

<div class="sidebar" {style}>
  {#each buttons as button}
    {@render sidebar_button(button)}
  {/each}
</div>

{#snippet sidebar_button(button_data: Button)}
  <button
    class="sidebar_button"
    style={button_data.Style}
    onclick={button_data.Onclick_function}>{@html button_data.Text}</button
  >
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

  .sidebar_button {
    width: 200px;
    height: fit-content;
    padding: 10px;
    margin: 0px 0px 0px 0px;
    border-radius: 0;
  }
</style>
