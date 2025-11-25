<script>
	import { format } from '$lib/javascript/format';
	import { chart } from '$lib/helper.svelte';
	import { financeReportCategories } from '$lib/javascript/simulation';
	import Window from './window.svelte';
	import FinancialReport from './financialReport.svelte';
	import MonthlyOverview from './monthlyOverview.svelte';
	import SalesReport from './salesReport.svelte';
	import BalanceSheet from './balanceSheet.svelte';

	/**
	 * @typedef {Object} props
	 * @property {import("$lib/javascript/simulation").clientState} clientState
	 * @property {(decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions
	 * @property {(contents: import("svelte").Snippet<[number]>)=>number} newWindow,
	 * @property {(windowId: number)=>void} deleteWindow,
	 * @property {string} serverAdress
	 */

	/** @type {props} */
	let {
		clientState = $bindable(),
		updateDecisions,
		newWindow,
		deleteWindow,
		serverAdress
	} = $props();

	async function getLatestMarketOverview() {
		/** @type {Object.<string, import("$lib/javascript/simulation").CompanyMarketStatistics>}*/
		let v = JSON.parse(
			await (
				await fetch('http://' + serverAdress.replace('localhost', '127.0.0.1') + '/market/')
			).text()
		);
		return v;
	}
</script>

{#snippet overview(
	/** @type {import("$lib/javascript/simulation").CompanyMarketStatistics}*/ CompanyMarketStatistics
)}
	<h1>
		{CompanyMarketStatistics.Name}
	</h1>

	<div class="grid">
		<img src="" alt="logo" />
		<table>
			<tbody>
				<tr>
					<td>CEO: </td>
					<td>PLAYER NAME</td>
				</tr>

				<tr>
					<td>Employee Count: </td>
					<td>{format.number(CompanyMarketStatistics.EmployeeCount, false, 0)}</td>
				</tr>

				<tr>
					<td>Quartarly Net Income</td>
					<td>{format.currency(CompanyMarketStatistics.QuartalyNetIncome, true, 0)}</td>
				</tr>

				<tr>
					<td>Company Assets</td>
					<td>{format.currency(CompanyMarketStatistics.Assets, true, 0)}</td>
				</tr>

				<tr>
					<td>Company value</td>
					<td>{format.currency(CompanyMarketStatistics.Value, true, 0)}</td>
				</tr>
			</tbody>
		</table>
	</div>
{/snippet}

{#key clientState.ExternalFactors.Month}
	{#await getLatestMarketOverview() then marketOverview}
		{#each Object.entries(marketOverview) as c}
			{@render overview(c[1])}
		{/each}
	{/await}
{/key}
