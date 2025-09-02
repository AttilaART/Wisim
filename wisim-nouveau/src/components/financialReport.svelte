<script>
	import { format } from '$lib/javascript/format';
	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void}} */
	let { clientState = $bindable(), updateDecisions } = $props();

	let financial_Report = $derived(
		/** @returns {import("$lib/javascript/simulation").Financial_report?}*/ (() => {
			if (!(clientState.company.Reports == null || clientState.company.Reports.length == 0)) {
				return clientState.company.Reports[clientState.company.Reports.length - 1].Financial_report;
			}
			return null;
		})()
	);
</script>

<div style="min-width: 600px;">
	{#if financial_Report}
		{#each Object.entries(financial_Report ? financial_Report : {}) as reportSection}
			<h3>{reportSection[0].replaceAll('_', ' ')}</h3>
			{@render renderSection(reportSection[1])}
		{/each}
	{:else}
		<h3>No Data</h3>
		<p>Data Will appear here when it's ready</p>
	{/if}
</div>

{#snippet renderSection(/** @type {Object.<String, Number>}}*/ section)}
	{#each Object.entries(section) as entry}
		<div class="grid">
			<p>{entry[0].replaceAll('_', ' ')}</p>
			<p style="text-align: right;">{format.currency(entry[1], true, 2)}</p>
		</div>
	{/each}
{/snippet}
