<script>
	import { format } from '$lib/javascript/format';

	/** @type {{report: import("$lib/javascript/simulation").Report}}*/
	let { report } = $props();

	/**
	 * @param {import("$lib/javascript/simulation").FinanceReportEntry[]} data
	 * @returns { import("$lib/javascript/simulation").FinanceReportEntry[] }
	 */
	function aggregateEntries(data) {
		/** @type {import("$lib/javascript/simulation").FinanceReportEntry[]} */
		let aggregatedData = [];
		for (let entry of data) {
			let i = aggregatedData.findIndex((e) => e.Name == entry.Name);
			if (i < 0) {
				aggregatedData.push(entry);
			} else {
				aggregatedData[i].Value += entry.Value;
			}
		}

		return aggregatedData;
	}

	/**
	 * @param {import("$lib/javascript/simulation").FinanceReportEntry[]} data
	 * @returns { number }
	 */
	function total(data) {
		let total = 0;
		data.map((e) => {
			total += e.Value;
		});
		return total;
	}
</script>

<div class="grid">
	{@render table(aggregateEntries(JSON.parse(JSON.stringify(report.BalanceSheet.Liabilities))))}
	{@render table(aggregateEntries(JSON.parse(JSON.stringify(report.BalanceSheet.Assets))))}
</div>

{#snippet table(/** @type {import("$lib/javascript/simulation").FinanceReportEntry[]}*/ entries)}
	<table>
		<tbody>
			{#each entries as entry}
				<tr>
					<td>{entry.Name}</td>
					<td>{format.currency(entry.Value, true, 0)}</td>
				</tr>
			{/each}
		</tbody>
		<tfoot>
			<tr>
				<td>Total</td>
				<td>{format.currency(total(entries), true, 0)}</td>
			</tr>
		</tfoot>
	</table>
{/snippet}
