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

	let financial_Report = $derived(
		/** @returns {import("$lib/javascript/simulation").Financial_report?}*/ (() => {
			if (!(clientState.company.Reports == null || clientState.company.Reports.length == 0)) {
				return clientState.company.Reports[clientState.company.Reports.length - 1].Financial_report;
			}
			return null;
		})()
	);
</script>

<div class="dashboard">
	<article class="dashboard-item">
		{#if financial_Report}
			<div>
				<h3><span data-tooltip="Income from all sources">Gross Profit</span></h3>
				<h2>{format.currency(financial_Report?.Income.Gross_profit, true, 2)}</h2>
			</div>

			<div>
				<h3>Cost of Sales</h3>
				<h2>
					{format.currency(financial_Report?.Income.Cost_of_sales, true, 2)}
				</h2>
			</div>

			<div>
				<h3>Profit per product sold</h3>
				<h2>
					{format.currency(
						financial_Report
							? financial_Report.Totals.Net_income /
									clientState.company.Reports[clientState.company.Reports.length - 1].Sales_report
										.Company_sales_statistics.Products_sold
							: 0,
						true,
						2
					)}
				</h2>
			</div>

			<div>
				<h3>Net Income</h3>
				<h2>{format.currency(financial_Report?.Totals.Net_income, true, 2)}</h2>
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
				<h2>{format.currency(clientState.decisions.Finances.Set_bank_loan, false, 0)}</h2>
				<strong>
					{format.currency(
						-clientState.decisions.Finances.Set_bank_loan *
							clientState.external_factors.Intrest_rate,
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
				<h2>{format.currency(clientState.company.Bridge_loans, false, 0)}</h2>
				<strong>
					{format.currency(
						-clientState.company.Bridge_loans *
							clientState.external_factors.Bridge_loans_intrest_rate,
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
