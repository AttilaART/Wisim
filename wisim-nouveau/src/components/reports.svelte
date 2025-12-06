<script module>
	/**
	 * @param {import("$lib/javascript/simulation").Report} report
	 */
	export function getMonths(report) {
		return `${report.Month}`;
	}

	/**
	 * @param {import("$lib/javascript/simulation").Report} report
	 */
	function getNetProfit(report) {
		return `${report.FinancialReport.Totals.NetIncome}`;
	}

	/**
	 * @param {import("$lib/javascript/simulation").Report} report
	 */
	function getBankBalance(report) {
		return `${report.BalanceSheet.Bank_balance}`;
	}
</script>

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
	 */

	/** @type {props} */
	let { clientState = $bindable(), updateDecisions, newWindow, deleteWindow } = $props();

	let onCurrentReport = $state(true);

	let selectedReportIndex = $state(clientState.Company.Reports.length);
	let currentReportIndex = $derived.by(() => {
		if (selectedReportIndex < 1) {
			return 0;
		} else if (selectedReportIndex > clientState.Company.Reports.length) {
			return clientState.Company.Reports.length - 1;
		}
		return selectedReportIndex - 1;
	});

	$effect(() => {
		if (onCurrentReport) {
			selectedReportIndex = clientState.Company.Reports.length;
		}
	});

	let currentReport = $derived(clientState.Company.Reports[currentReportIndex]);

	let selectedProduct = $state(currentReport ? Object.keys(currentReport.SalesReport)[0] : '');

	let timespan = $state(12);

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
				cashflow += currentReport.FinancialReport.Totals.Cashflow;
			}
		}

		if (months <= 0) return 0;

		cashflow = (cashflow / months) * 4;

		return cashflow;
	}

	function totalAssets() {
		let assets = 0;

		if (clientState.Company.Reports.length <= 0) return 0;
		for (let e of currentReport.BalanceSheet.Assets) {
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

	/**
	 * @param {(report: import("$lib/javascript/simulation").Report)=>number | string} getData
	 */
	function aggragateData(getData) {
		/** @type {Array.<number | string>} */
		let data = [];

		if (clientState.Company.Reports.length == 0) return data;

		for (let i = currentReportIndex - timespan; i <= currentReportIndex; i++) {
			if (i < 0) continue;
			data.push(getData(clientState.Company.Reports[i]));
		}

		return data;
	}
	/** @typedef {Object.<string, import("$lib/javascript/simulation").FinanceReportEntry[]>} groupedInvoices
  /** 
   * @param {import("$lib/javascript/simulation").Report} report
   * @returns {groupedInvoices} 
   */
	function groupInvoices(report) {
		if (!(report == null)) {
			let rawInvoices = report.BalanceSheet.InvoiceLog;

			/** @type {groupedInvoices} */
			let groupedInvoices = {};
			for (let i of rawInvoices) {
				if (groupedInvoices[i.Group] == undefined) {
					groupedInvoices[i.Group] = [i];
				} else {
					groupedInvoices[i.Group].push(i);
				}
			}

			return groupedInvoices;
		}
		return {};
	}

	/**
	 * @param {import("$lib/javascript/simulation").FinanceReportEntry[]} group
	 * @returns {number}
	 */
	function getInvoiceGroupTotal(group) {
		let totalValue = 0;
		for (let i of group) {
			totalValue += i.Value;
		}
		return totalValue;
	}
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
			class={tab == 'monthly' ? '' : 'outline'}
			onclick={() => {
				tab = 'monthly';
			}}>Monthly</button
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

	<fieldset role="group">
		<input type="text" disabled value="Report Month" style="width: 10rem;" />
		<input
			type="number"
			bind:value={selectedReportIndex}
			min={1}
			max={clientState.Company.Reports.length}
			onchange={() => {
				if (selectedReportIndex > clientState.Company.Reports.length) {
					clientState.Company.Reports.length;
				} else if (selectedReportIndex < 1) {
					selectedReportIndex = 1;
				}
				onCurrentReport = selectedReportIndex == clientState.Company.Reports.length;
			}}
		/>

		<input type="text" disabled value="Timespan (months)" style="width: 15rem;" />
		<input
			type="number"
			bind:value={timespan}
			min={1}
			max={clientState.Company.Reports.length}
			onchange={() => {
				if (timespan < 1) {
					timespan = 1;
				}
			}}
		/>
	</fieldset>

	{#if tab == 'overview'}
		{@render overview()}
	{:else if tab == 'monthly'}
		{@render monthly()}
	{:else if tab == 'finances'}
		{@render finances()}
	{:else if tab == 'assets'}
		{@render assets()}
	{:else if tab == 'sales'}
		{@render sales()}
	{/if}
</div>

{#snippet overview()}
	<h1>
		<input
			bind:value={clientState.Decisions.General.CompanyName}
			type="text"
			autocomplete="off"
			onchange={() => {
				updateDecisions(clientState.Decisions);
			}}
		/>
	</h1>

	<div class="grid">
		<img src="" alt="logo" />
		<table>
			<tbody>
				<tr>
					<td>CEO: </td>
					<td>{clientState.Company.CEO}</td>
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

{#snippet monthly()}
	{#if selectedReportIndex != -1}
		<MonthlyOverview report={currentReport}></MonthlyOverview>
	{:else}
		No Data
	{/if}
{/snippet}

{#snippet finances()}
	{#key currentReportIndex - timespan}
		<div class="grid" style="margin-bottom: 1rem;">
			<div
				use:chart={{
					title: {
						text: 'Net Income'
					},
					tooltip: {},
					yAxis: {},
					xAxis: {
						data: aggragateData(getMonths)
					},
					series: [
						{
							name: 'Net Income',
							type: 'line',
							data: aggragateData(getNetProfit)
						}
					]
				}}
				style="height: 20rem"
			></div>

			<div
				use:chart={{
					title: {
						text: 'Bank Balance'
					},
					tooltip: {},
					xAxis: {
						data: aggragateData(getMonths)
					},
					yAxis: {},
					series: [
						{
							name: 'Balance',
							type: 'line',
							data: aggragateData(getBankBalance)
						}
					]
				}}
				style="height: 20rem"
			></div>
		</div>
	{/key}

	<h1>Key Metrics</h1>
	{#if clientState.Company.Reports.length >= 1}
		<article id="key-metrics">
			<label for="">
				Net Income
				<h2>{format.currency(currentReport.FinancialReport.Totals.NetIncome, true, 0)}</h2>
			</label>

			<label for="">
				Total Cashflow
				<h2>{format.currency(currentReport.FinancialReport.Totals.Cashflow, true, 0)}</h2>
			</label>

			<label for="">
				Total assets
				<h2>{format.currency(totalAssets(), false, 0)}</h2>
			</label>

			<label for="">
				Operating Expenses
				<h2>
					{format.currency(currentReport.FinancialReport.Totals.TotalOperatingExpenses, false, 0)}
				</h2>
			</label>

			<label for="">
				Non-operating Expenses
				<h2>
					{format.currency(
						currentReport.FinancialReport.Totals.TotalNonOperatingExpenses,
						false,
						0
					)}
				</h2>
			</label>

			<label for="">
				Non-operating and Operating Expenses
				<h2>
					{format.currency(
						currentReport.FinancialReport.Totals.TotalNonOperatingExpenses +
							currentReport.FinancialReport.Totals.TotalOperatingExpenses,
						false,
						0
					)}
				</h2>
			</label>
		</article>

		<div style="text-align: right;">
			<button
				class="outline"
				onclick={() => {
					newWindow(invoiceLog);
				}}>See Invoice Log</button
			>
			<button
				onclick={() => {
					newWindow(financialReoport);
				}}>See Full Financial Report</button
			>
		</div>
	{:else}
		<article>No Data</article>
	{/if}
{/snippet}

{#snippet invoiceLog(/** @type {number} */ windowID)}
	<Window
		title="Invoice Log {selectedReportIndex}"
		closeWindow={() => {
			deleteWindow(windowID);
		}}
	>
		{#each Object.entries(groupInvoices(currentReport)) as group (group[0])}
			<details>
				<summary style="min-width: 500px;">
					<div class="grid">
						<div>
							{format.titleCase(financeReportCategories[group[0]].replaceAll('_', ' '))}
						</div>
						<div style="text-align: right;">
							{format.currency(getInvoiceGroupTotal(group[1]), true, 2)}
						</div>
					</div>
				</summary>
				{#each group[1] as i}
					<article>
						<h6>{i.Name}</h6>
						<p>{i.Info}</p>
						<h5>{format.currency(i.Value, true, 2)}</h5>
					</article>
				{/each}
			</details>
		{/each}
	</Window>
{/snippet}

{#snippet financialReoport(/** @type {number} */ windowID)}
	<Window
		title="Financial Report {selectedReportIndex}"
		closeWindow={() => {
			deleteWindow(windowID);
		}}
	>
		<FinancialReport financial_Report={currentReport.FinancialReport}></FinancialReport>
	</Window>
{/snippet}

{#snippet sales()}
	{#if currentReport}
		<select bind:value={selectedProduct}>
			{#each Object.keys(clientState.Company.Offers) as product}
				<option value={product}>{clientState.Company.Offers[product].Product.Name}</option>
			{/each}
		</select>
		{#key [selectedReportIndex, currentReportIndex, timespan, selectedProduct]}
			<SalesReport
				report={currentReport}
				aggrageteData={aggragateData}
				{selectedProduct}
				{clientState}
			></SalesReport>
		{/key}
	{:else}
		No Data
	{/if}
{/snippet}

{#snippet assets()}
	{#if currentReport}
		<BalanceSheet report={currentReport}></BalanceSheet>
	{:else}
		No Data
	{/if}
{/snippet}

<style>
	#key-metrics {
		display: grid;
		grid-template-columns: 33% 33% 33%;
		gap: 1rem;
	}
</style>
