<script>
	import { format } from '$lib/javascript/format';
	import { chart } from '$lib/helper.svelte';
	/**
	 * @typedef {Object} props
	 * @property {import("$lib/javascript/simulation").clientState} clientState
	 * @property {(decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions
	 */

	/** @type {props} */
	let { clientState = $bindable(), updateDecisions } = $props();

	let latestReport = $derived(clientState.Company.Reports[clientState.Company.Reports.length - 1]);

	let tab = $state('overview');

	function quartarlyNetIncome() {
		let income = 0;

		let months = 0;
		for (let i = 0; i < 4; i++) {
			if (clientState.Company.Reports.length - 1 - i >= 0) {
				months += 1;
				income +=
					clientState.Company.Reports[clientState.Company.Reports.length - 1 - i].FinancialReport
						.Totals.NetIncome;
			}
		}

		if (months <= 0) return 0;

		income = (income / months) * 4;

		return income;
	}

	function quartarlyCashflow() {
		let cashflow = 0;

		let months = 0;
		for (let i = 0; i < 4; i++) {
			if (clientState.Company.Reports.length - 1 - i >= 0) {
				months += 1;
				cashflow +=
					clientState.Company.Reports[clientState.Company.Reports.length - 1 - i].FinancialReport
						.Totals.Cashflow;
			}
		}

		if (months <= 0) return 0;

		cashflow = (cashflow / months) * 4;

		return cashflow;
	}

	function totalAssets() {
		let assets = 0;

		if (clientState.Company.Reports.length <= 0) return 0;
		for (let e of clientState.Company.Reports[clientState.Company.Reports.length - 1].BalanceSheet
			.Assets) {
			assets += e.Value;
		}

		return assets;
	}

	/**
	 * Calcuate valueation based on "Discounted Cash Flows" + assets
	 * https://online.hbs.edu/blog/post/how-to-value-a-company
	 */
	function valuation() {
		if (clientState.Company.Reports.length <= 0) return 0;

		let cashflowPerYear = quartarlyCashflow() * 4;
		let years = 5;

		return (
			(cashflowPerYear / (1 + clientState.ExternalFactors.IntrestRate)) * years + totalAssets()
		);
	}

	$effect(() => {
		clientState.Company.Name = clientState.Decisions.General.CompanyName;
	});
</script>

<div style="min-width: 50rem;">
	<div class="grid" style="margin-bottom: 1rem;">
		<button
			class={tab == 'overview' ? '' : 'outline'}
			onclick={() => {
				tab = 'overview';
			}}>Overview</button
		>
		<button
			class={tab == 'finances' ? '' : 'outline'}
			onclick={() => {
				tab = 'finances';
			}}>Finances</button
		>
		<button
			class={tab == 'assets' ? '' : 'outline'}
			onclick={() => {
				tab = 'assets';
			}}>Assets</button
		>
		<button
			class={tab == 'sales' ? '' : 'outline'}
			onclick={() => {
				tab = 'sales';
			}}>Sales</button
		>
	</div>

	{#if tab == 'overview'}
		{@render overview()}
	{:else if tab == 'finances'}
		{@render finances()}
	{:else if tab == 'assets'}
		Assets
	{:else if tab == 'sales'}
		Sales
	{/if}
</div>

{#snippet overview()}
	<h1>
		<input bind:value={clientState.Decisions.General.CompanyName} type="text" />
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
					<td
						>{format.number(
							clientState.Employees.marketing.length + clientState.Employees.production.length,
							false,
							0
						)}</td
					>
				</tr>

				<tr>
					<td>Quartarly Net Income</td>
					<td>{format.currency(quartarlyNetIncome(), true, 0)}</td>
				</tr>

				<tr>
					<td>Company Assets</td>
					<td>{format.currency(totalAssets(), true, 0)}</td>
				</tr>

				<tr>
					<td>Company value</td>
					<td>{format.currency(valuation(), true, 0)}</td>
				</tr>
			</tbody>
		</table>
	</div>
{/snippet}

{#snippet finances()}
	<div class="grid">
		<div>
			<div
				use:chart={{
					title: {
						text: 'ECharts Getting Started Example'
					},
					tooltip: {},
					xAxis: {
						data: ['shirt', 'cardigan', 'chiffon', 'pants', 'heels', 'socks']
					},
					yAxis: {},
					series: [
						{
							name: 'sales',
							type: 'line',
							data: [5, 20, 36, 10, 10, 20]
						}
					]
				}}
				style="height: 20rem"
			></div>
		</div>

		<div>
			<div
				use:chart={{
					title: {
						text: 'ECharts Getting Started Example'
					},
					tooltip: {},
					xAxis: {
						data: ['shirt', 'cardigan', 'chiffon', 'pants', 'heels', 'socks']
					},
					yAxis: {},
					series: [
						{
							name: 'sales',
							type: 'line',
							data: [5, 20, 36, 10, 10, 20]
						}
					]
				}}
				style="height: 20rem"
			></div>
		</div>
	</div>

	<h1>Key Metrics</h1>
	{#if clientState.Company.Reports.length >= 1}
		<article class="grid">
			<label for="">
				Net Income
				<h2>{format.currency(latestReport.FinancialReport.Totals.NetIncome, true, 0)}</h2>
			</label>

			<label for="">
				Total Cashflow
				<h2>{format.currency(latestReport.FinancialReport.Totals.Cashflow, true, 0)}</h2>
			</label>

			<label for="">
				Total assets
				<h2>{format.currency(totalAssets(), false, 0)}</h2>
			</label>
		</article>
	{:else}
		<article>No Data</article>
	{/if}
{/snippet}
