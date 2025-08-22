<script>
  let { value } = $props();
</script>

{#if value}
  {#if value > 0}
    ▲
  {:else if value == 0}
    -
  {:else}
    ▼
  {/if}
{/if}
