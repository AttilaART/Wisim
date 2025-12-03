<script>
	import { format } from '$lib/javascript/format';

	/** @type {{MarketStatistics: Object.<string, import("$lib/javascript/simulation").CompanyMarketStatistics>}} */
	let { MarketStatistics } = $props();
</script>

<table>
	<thead>
		<tr>
			<th> Rank </th>
			<th> Name </th>
			<th> CEO </th>
			<th> Value </th>
		</tr>
	</thead>
	<tbody>
		{#each Object.values(MarketStatistics).sort((a, b) => b.Value - a.Value) as c, i}
			<tr>
				<td>{i + 1}</td>
				<td>{c.Name}</td>
				<td>{c.CEO}</td>
				<td><strong>{format.currency(c.Value, false, 0)}</strong></td>
			</tr>
		{/each}
	</tbody>
</table>
