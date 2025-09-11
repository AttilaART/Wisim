<script>
	import { format } from '$lib/javascript/format';
	import { financeReportCategories } from '$lib/javascript/simulation';

	/**
	 * @typedef {Object} props
	 * @property {import("$lib/javascript/simulation").clientState} clientState
	 * @property {(decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions
	 */
	/** @type {props} */
	let { clientState = $bindable(), updateDecisions } = $props();

	/** @typedef {Object.<string, import("$lib/javascript/simulation").FinanceReportEntry[]>} groupedInvoices
  /** @type {groupedInvoices} */
	let invoices = $derived(
		/** @returns {groupedInvoices}*/ (() => {
			if (!(clientState.company.Reports == null || clientState.company.Reports.length == 0)) {
				let rawInvoices =
					clientState.company.Reports[clientState.company.Reports.length - 1].Balance_sheet
						.Invoice_log;

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
		})()
	);

	/**
	 * @param {import("$lib/javascript/simulation").FinanceReportEntry[]} group
	 * @returns {number}
	 */
	function getGroupTotal(group) {
		let totalValue = 0;
		for (let i of group) {
			totalValue += i.Value;
		}
		return totalValue;
	}
</script>

{#each Object.entries(invoices) as group (group[0])}
	<details>
		<summary style="min-width: 500px;">
			<div class="grid">
				<div>
					{format.titleCase(financeReportCategories[group[0]].replaceAll('_', ' '))}
				</div>
				<div style="text-align: right;">
					{format.currency(getGroupTotal(group[1]), true, 2)}
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
