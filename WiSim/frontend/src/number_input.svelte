<script lang="ts">
  let {
    formatter,
    value = $bindable(),
  }: {
    formatter: (value: number) => string;
    value: number;
  } = $props();
  let input_text: string = $state(formatter(value));
</script>

<input
  type="text"
  bind:value={input_text}
  onfocus={() => (input_text = String(value))}
  onfocusout={() => {
    if (!isNaN(parseFloat(input_text))) {
      value = parseFloat(input_text);
    }
    input_text = formatter(value);
  }}
/>
