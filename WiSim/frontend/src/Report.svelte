<script lang="ts">
  let report = $props();

  type key_value_pair = {
    key: String;
    value: String | number;
  };

  function sep_key_values(report) {
    let keys = report.keys();
    let values = report.values();

    let key_values: key_value_pair[] = [];
    for (let i in keys) {
      key_values.push({ key: keys[i], value: values[i] });
    }

    return key_values;
  }
</script>

<table>
  <thead>
    <tr>
      <th></th>
      <th></th>
    </tr>
  </thead>
  <tbody>
    {#each sep_key_values(report) as pair}
      {@render row(pair.key, pair.value)}
    {/each}
  </tbody>
</table>

{#snippet row(key, value)}
  <tr>
    <td>{key}</td>
    <td style="text-align: right;">{value}</td>
  </tr>
{/snippet}
