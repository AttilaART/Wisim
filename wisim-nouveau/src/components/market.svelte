<script>
	import { format } from '$lib/javascript/format';
	import { chart } from '$lib/helper.svelte';
	import { financeReportCategories } from '$lib/javascript/simulation';
	import Window from './window.svelte';
	import FinancialReport from './financialReport.svelte';
	import MonthlyOverview from './monthlyOverview.svelte';
	import SalesReport from './salesReport.svelte';
	import BalanceSheet from './balanceSheet.svelte';
	import { position, transform } from '@neodrag/svelte';
	import Leaderboard from './leaderboard.svelte';

	/**
	 * @typedef {Object} props
	 * @property {import("$lib/javascript/simulation").clientState} clientState
	 * @property {(decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions
	 * @property {(contents: import("svelte").Snippet<[string]>)=>string} newWindow,
	 * @property {(windowId: string)=>void} deleteWindow,
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

	let tab = $state('overview');
	let selectedCompanyID = $state(0);

	async function getLatestMarketOverview() {
		/** @type {Object.<string, import("$lib/javascript/simulation").CompanyMarketStatistics>}*/
		let v = JSON.parse(
			await (
				await fetch('http://' + serverAdress.replace('localhost', '127.0.0.1') + '/' + 'market/')
			).text()
		);
		return v;
	}

	/**
	 * @param {Object.<string, import("$lib/javascript/simulation").CompanyMarketStatistics>} MarketStatistics
	 * @param {(p: import("$lib/javascript/simulation").CompanyMarketStatisticsProduct)=>(any)} getData
	 * @returns {( number|string )[]}
	 */
	function getProductData(MarketStatistics, getData) {
		let returnList = [];
		for (let c of Object.values(MarketStatistics)) {
			for (let p of Object.values(c.Products)) {
				returnList.push(getData(p));
			}
		}

		return returnList;
	}
</script>

<div style="min-width: 50rem;">
	<div class="grid" style="margin-bottom: 1rem;">
		<button
			class={tab == 'all' ? '' : 'outline'}
			onclick={() => {
				tab = 'all';
			}}>All</button
		>
		<button
			class={tab == 'overview' ? '' : 'outline'}
			onclick={() => {
				tab = 'overview';
			}}>Overview</button
		>
		<button
			disabled
			class={tab == 'finances' ? '' : 'outline'}
			onclick={() => {
				tab = 'finances';
			}}>Finances</button
		>
		<button
			disabled
			class={tab == 'sales' ? '' : 'outline'}
			onclick={() => {
				tab = 'sales';
			}}>Sales</button
		>
	</div>
	{#key clientState.ExternalFactors}
		{#await getLatestMarketOverview() then marketOverview}
			{#if tab != 'all'}
				<select bind:value={selectedCompanyID} name="company" id="">
					{#each Object.values(marketOverview) as c}
						<option value={c.CompanyID}>{c.Name} ({c.CEO})</option>
					{/each}
				</select>
			{/if}

			{#if tab == 'all'}
				{@render all(marketOverview)}
			{:else if tab == 'overview'}
				{@render overview(marketOverview[selectedCompanyID])}
			{:else if tab == 'finances'}
				{@render finances(marketOverview[selectedCompanyID])}
			{:else}
				{@render sales(marketOverview[selectedCompanyID])}
			{/if}
		{:catch error}
			<center style="margin: 2rem">
				<h1>No Data</h1>
				Wait until the next turn to see data
				{error}
			</center>
		{/await}
	{/key}
</div>

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
					<td>{CompanyMarketStatistics.CEO}</td>
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

{#snippet all(
	/** @type {Object.<string, import("$lib/javascript/simulation").CompanyMarketStatistics>}*/ MarketStatistics
)}
	<div class="grid">
		<Leaderboard {MarketStatistics}></Leaderboard>
	</div>
	<div class="grid">
		{#if Object.values(MarketStatistics)
			.map((c) => {
				return { value: c.MonthlySales, name: `${c.Name}` + '\u200B'.repeat(c.CompanyID) };
			})
			.filter((v) => v.value > 0).length >= 1}
			<div
				use:chart={{
					title: {
						text: 'Market Share (Sales)'
					},
					tooltip: {},
					series: [
						{
							name: 'Monthly Sales',
							type: 'pie',
							data: Object.values(MarketStatistics)
								.map((c) => {
									return {
										value: c.MonthlySales,
										name: `${c.Name}` + '\u200B'.repeat(c.CompanyID)
									};
								})
								.filter((v) => v.value > 0)
						}
					]
				}}
				style="height: 20rem"
			></div>

			<div
				use:chart={{
					title: {
						text: 'Market Share (Value)'
					},
					tooltip: {},
					series: [
						{
							name: 'Value of Monthly Sales (CHF)',
							type: 'pie',
							data: Object.values(MarketStatistics)
								.map((c) => {
									return {
										value: c.ValueOfMonthlySales,
										name: `${c.Name}` + '\u200B'.repeat(c.CompanyID)
									};
								})
								.filter((v) => v.value > 0)
						}
					]
				}}
				style="height: 20rem"
			></div>
		{:else}
			<div class="chart-placeholder">
				<span> No Sales </span>
			</div>
			<div class="chart-placeholder">
				<span> No Sales </span>
			</div>
		{/if}
	</div>

	<div class="grid">
		{#if getProductData(MarketStatistics, (p) => p).length >= 1}
			<div
				use:chart={{
					title: {
						text: 'Product Prices'
					},
					tooltip: {},
					series: [
						{
							name: 'Price',
							type: 'bar',
							data: getProductData(MarketStatistics, (e) => {
								return { value: e.Marketing_statistics.Price, name: e.Marketing_statistics.Name };
							})
						}
					]
				}}
				style="height: 20rem"
			></div>
		{:else}
			<div class="chart-placeholder">
				<span> No Products </span>
			</div>
		{/if}

		<div
			use:chart={{
				title: {
					text: 'Market Share (Value)'
				},
				tooltip: {},
				series: [
					{
						name: 'Value of Monthly Sales (CHF)',
						type: 'pie',
						data: Object.values(MarketStatistics).map((c) => {
							return {
								value: c.ValueOfMonthlySales,
								name: `${c.Name}` + '\u200B'.repeat(c.CompanyID)
							};
						})
					}
				]
			}}
			style="height: 20rem"
		></div>
	</div>
{/snippet}

{#snippet finances(
	/** @type {import("$lib/javascript/simulation").CompanyMarketStatistics}*/ CompanyMarketStatistics
)}{/snippet}

{#snippet sales(
	/** @type {import("$lib/javascript/simulation").CompanyMarketStatistics}*/ CompanyMarketStatistics
)}{/snippet}

<style>
	.grid {
		margin-bottom: var(--spacing);
	}
	.chart-placeholder {
		position: relative;
		height: 20rem;
		background: rgba(0, 0, 0, 0.5);
		span {
			position: absolute;
			top: 50%;
			left: 50%;
			transform: translate(-50%, -50%);
		}
	}
</style>
