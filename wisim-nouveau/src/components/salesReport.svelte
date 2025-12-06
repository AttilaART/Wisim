<script>
	import { chart } from '$lib/helper.svelte';
	import { format } from '$lib/javascript/format';
	import { getMonths } from '../components/reports.svelte';
	/**
	 * @typedef {Object} Props
	 * @property {import("$lib/javascript/simulation").clientState} clientState
	 * @property {import("$lib/javascript/simulation").Report} report
	 * @property {string} selectedProduct
	 * @property {(getData: (report: import("$lib/javascript/simulation").Report)=>number|string)=>(string | number)[]} aggrageteData
	 */
	/** @type {Props}*/
	let { clientState, report, selectedProduct = $bindable(), aggrageteData } = $props();

	let productExists = $derived(undefined == report.SalesReport[selectedProduct] ? false : true);

	/**
	 * @param {import("$lib/javascript/simulation").Report} report
	 * @param {string} productID
	 */
	function getProductSales(report, productID) {
		if (report.SalesReport[productID] == undefined) {
			return 0;
		}

		return report.SalesReport[productID].ProductSalesStatistics.ProductsSold;
	}

	/**
	 * @param {import("$lib/javascript/simulation").Report} report
	 * @param {string} productID
	 */
	function getProductDemand(report, productID) {
		if (report.SalesReport[productID] == undefined) {
			return 0;
		}

		return report.SalesReport[productID].ProductSalesStatistics.ProductDemand;
	}

	/**
	 * @param {import("$lib/javascript/simulation").Report} report
	 * @param {string} productID
	 */
	function getProductsProduced(report, productID) {
		if (report.ProductionReport.ProductSpecificReport[productID] == undefined) {
			return 0;
		}

		return report.ProductionReport.ProductSpecificReport[productID].TotalProductsProduced;
	}

	/**
	 * @param {import("$lib/javascript/simulation").Report} report
	 */
	function getProductRatios(report) {
		return Object.entries(report.SalesReport).map((e) => {
			return {
				value: e[1].ProductSalesStatistics.ProductsSold,
				name: clientState.Company.Offers[e[0]].Product.Name
			};
		});
	}
</script>

{#if selectedProduct == ''}
	No product selected
{:else if !productExists}
	<div>
		{selectedProduct}
		This Product didn't exist yet.
	</div>
{:else}
	<div class="grid" style="margin-bottom: 1rem;">
		<div
			use:chart={{
				title: {
					text: 'Supply vs Demand'
				},
				tooltip: {},
				yAxis: {},
				xAxis: {
					data: aggrageteData(getMonths)
				},
				series: [
					{
						name: 'Demand',
						type: 'line',
						data: aggrageteData(
							(/** @type {import("$lib/javascript/simulation").Report}*/ report) => {
								return getProductDemand(report, selectedProduct);
							}
						)
					},
					{
						name: 'Products Produced',
						type: 'line',
						data: aggrageteData(
							(/** @type {import("$lib/javascript/simulation").Report}*/ report) => {
								return getProductsProduced(report, selectedProduct);
							}
						)
					}
				]
			}}
			style="height: 20rem"
		></div>
		<div
			use:chart={{
				title: {
					text: 'Products Sold'
				},
				tooltip: {},
				series: [
					{
						name: 'Products Sold',
						type: 'pie',
						data: getProductRatios(report)
					}
				]
			}}
			style="height: 20rem"
		></div>
	</div>

	<h1>Sales Report</h1>

	<article id="key-metrics">
		<label for="">
			Products Sold
			<h2>
				{format.number(
					report.SalesReport[selectedProduct].ProductSalesStatistics.ProductsSold,
					false,
					0
				)}
			</h2>
		</label>

		<label for="">
			Product Demand
			<h2>
				{format.number(
					report.SalesReport[selectedProduct].ProductSalesStatistics.ProductDemand,
					false,
					0
				)}
			</h2>
		</label>

		<label for="">
			Products Produced
			<h2>
				{format.number(
					report.ProductionReport.ProductSpecificReport[selectedProduct].TotalProductsProduced,
					false,
					0
				)}
			</h2>
		</label>

		<label for="">
			Advertisment Budget
			<h2>
				{format.currency(
					report.SalesReport[selectedProduct].MarketingStatistics.PromotionQuantity,
					false,
					0
				)}
			</h2>
		</label>

		<label for="">
			Advertisement Impressions
			<h2>
				{format.number(
					report.SalesReport[selectedProduct].MarketingStatistics.ImpressionCount,
					false,
					0
				)}
			</h2>
		</label>

		<label for="">
			Advertisement Quality
			<h2>
				{format.number(
					report.SalesReport[selectedProduct].MarketingStatistics.PromotionQuality,
					false,
					1
				)}
			</h2>
		</label>
	</article>
{/if}

<style>
	#key-metrics {
		display: grid;
		grid-template-columns: 33% 33% 33%;
		gap: 1rem;
	}
</style>
