<script lang="ts">
  let {
    formatter,
    value = $bindable(),
    align_right,
  }: {
    formatter: (value: number) => string;
    value: number;
    align_right?: boolean;
  } = $props();
  let input_text: string = $state(formatter(value));
</script>

<input
  type="text"
  style={align_right ? "text-align: right;" : ""}
  bind:value={input_text}
  onfocus={() => (input_text = String(value))}
  onfocusout={() => {
    if (!isNaN(parseFloat(input_text))) {
      value = parseFloat(input_text);
    }
    input_text = formatter(value);
  }}
/>
