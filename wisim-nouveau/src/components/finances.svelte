<script>
	import { format } from '$lib/javascript/format';
	/**
	 * @typedef {Object} props
	 * @property {import("$lib/javascript/simulation").clientState} clientState
	 * @property {(decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions
	 * @property {()=>void} openDebtWindow
	 * @property {()=>void} openFinancialReportWindow
	 * @property {()=>void} openInvoicesWindow
	 */

	/** @type {props} */
	let {
		clientState = $bindable(),
		updateDecisions,
		openDebtWindow,
		openFinancialReportWindow,
		openInvoicesWindow
	} = $props();

	let financialReport = $derived(
		/** @returns {import("$lib/javascript/simulation").FinancialReport?}*/ (() => {
			if (!(clientState.Company.Reports == null || clientState.Company.Reports.length == 0)) {
				return clientState.Company.Reports[clientState.Company.Reports.length - 1].FinancialReport;
			}
			return null;
		})()
	);
</script>

<div class="dashboard">
	<article class="dashboard-item">
		{#if financialReport}
			<div>
				<h3><span data-tooltip="Income from all sources">Gross Profit</span></h3>
				<h2>{format.currency(financialReport?.Income.GrossProfit, true, 2)}</h2>
			</div>

			<div>
				<h3>Cost of Sales</h3>
				<h2>
					{format.currency(financialReport?.Income.CostOfSales, true, 2)}
				</h2>
			</div>

			<div>
				<!--
				<h3>Profit per product sold</h3>
				<h2>
					{format.currency(
						financialReport
							? financialReport.Totals.NetIncome /
									clientState.Company.Reports[clientState.Company.Reports.length - 1].SalesReport
										.CompanySalesStatistics.ProductsSold
							: 0,
						true,
						2
					)}
				</h2>
        -->
			</div>

			<div>
				<h3>Net Income</h3>
				<h2>{format.currency(financialReport?.Totals.NetIncome, true, 2)}</h2>
			</div>

			<button onclick={openFinancialReportWindow}> See full report </button>
			<button onclick={openInvoicesWindow}> See invoice log</button>
		{:else}
			<h3>No Data</h3>
			<small>Data will show up here after the first step</small>
		{/if}
	</article>
	<article class="dashboard-item">
		<div style="grid-column: span 2;">
			<div class="grid">
				<h3>Bank Loan</h3>
				<h4>Monthly Intrest</h4>
			</div>
			<div class="grid">
				<h2>{format.currency(clientState.Decisions.Finances.SetBankLoan, false, 0)}</h2>
				<strong>
					{format.currency(
						-clientState.Decisions.Finances.SetBankLoan * clientState.ExternalFactors.IntrestRate,
						false,
						2
					)}
				</strong>
			</div>
		</div>

		<div style="grid-column: span 2;">
			<div class="grid">
				<h3>Bridge Loan</h3>
				<h4>Monthly Intrest</h4>
			</div>
			<div class="grid">
				<h2>{format.currency(clientState.Company.BridgeLoans, false, 0)}</h2>
				<strong>
					{format.currency(
						-clientState.Company.BridgeLoans * clientState.ExternalFactors.BridgeLoansIntrestRate,
						false,
						2
					)}
				</strong>
			</div>
		</div>

		<button style="grid-column: span 2;" onclick={openDebtWindow}>Change Loans</button>
	</article>
</div>

<style>
	.dashboard-item {
		display: grid;
		grid-template-columns: 1fr 1fr;
		grid-template-rows: 1fr 1fr 3rem;
		gap: var(--pico-spacing);
		padding: var(--pico-spacing);
		article {
			margin-bottom: 0;
		}
	}

	.dashboard {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--pico-spacing);
		padding: var(--pico-spacing);
		article {
			margin-bottom: 0;
		}
	}
</style>
