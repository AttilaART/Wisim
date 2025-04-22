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
    <button
      class="popup_button"
      style="{button_data.Style} "
      onclick={button_data.Onclick_function}>{button_data.Text}</button
    >
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
</style>
