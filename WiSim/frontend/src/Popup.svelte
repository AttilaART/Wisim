<script lang="ts">
  type Button_data = {
    Text: String;
    Custom_style: String;
    Onclick_function: () => void;
  };

  let { button_data, content } = $props();

  //
</script>

<div class="popup">
  <div class="popup_content">
    {@html content}
    {#if button_data != null}
      <button
        class="popup_button"
        style="{button_data.Style} "
        onclick={button_data.Onclick_function}>{button_data.Text}</button
      >
    {/if}
  </div>
</div>

<style>
  .popup_button {
    text-align: center;
    position: absolute;
    left: 50%;
    bottom: 5px;
    transform: translate(-50%, -50%);
    margin: 0 auto;
    width: fit-content;
  }

  .popup {
    position: absolute;
    width: 100%;
    height: 100vh;
    align-content: center;
    background-color: rgba(0, 0, 0, 0.2);
    z-index: 999999;
  }

  .popup_content {
    position: relative;
    height: 45%;
    width: 45%;
    margin: auto auto;
    background-color: rgba(20, 20, 20, 1);
    border-radius: 10px;
    align-content: center;
    box-shadow: 5px 5px 15px rgba(0, 0, 0, 0.5);
  }
</style>
