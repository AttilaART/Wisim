<script>
	import { format } from '$lib/javascript/format';
	/** @type {{clientState: import("$lib/javascript/simulation").clientState}} */
	let { clientState = $bindable() } = $props();

	let latestReport = $derived(clientState.Company.Reports[clientState.Company.Reports.length - 1]);

	let { totalProductsSold, totalMarketShare } = $derived.by(() => {
		let totalProductsSold = 0;
		let totalMarketShare = 0;
		for (let r of Object.entries(latestReport.SalesReport)) {
			totalProductsSold += r[1].ProductSalesStatistics.ProductsSold;
			totalMarketShare += r[1].ProductSalesStatistics.MarketShare;

			console.log(r[1].ProductSalesStatistics);
		}

		return { totalProductsSold, totalMarketShare };
	});

	let totalAssets = $derived.by(() => {
		let totalAssets = 0;
		for (let e of Object.entries(latestReport.BalanceSheet.Assets)) {
			totalAssets += e[1].Value;
		}
		return totalAssets;
	});

	let { avgPromotionQuality, totalPromotionQuantity, totalPromotionImpressions } = $derived.by(
		() => {
			let avgPromotionQuality = 0;
			let totalPromotionQuantity = 0;
			let totalPromotionImpressions = 0;

			let entries = Object.entries(latestReport.SalesReport);
			for (let e of entries) {
				avgPromotionQuality += e[1].MarketingStatistics.PromotionQuality;
				totalPromotionQuantity += e[1].MarketingStatistics.PromotionQuantity;
				totalPromotionImpressions += e[1].MarketingStatistics.ImpressionCount;
			}

			if (entries.length > 0) {
				avgPromotionQuality = avgPromotionQuality / entries.length;
			}

			return { avgPromotionQuality, totalPromotionQuantity, totalPromotionImpressions };
		}
	);
</script>

{#if !latestReport}
	<h1>Loading</h1>
{:else}
	<div style="min-width:30rem;">
		<h2>Monthly Report</h2>
		<h3>Finances</h3>
		<div class="grid">
			<label for="">
				Net Income:
				<h4>{format.currency(latestReport.FinancialReport.Totals.NetIncome, true, 0)}</h4>
			</label>

			<label for="">
				Total Cashflow:
				<h4>{format.currency(latestReport.FinancialReport.Totals.Cashflow, true, 0)}</h4>
			</label>

			<label for="">
				Total Assets:
				<h4>{format.currency(totalAssets, false, 0)}</h4>
			</label>
		</div>

		<h3>Marketing</h3>
		<div class="grid">
			<label for="">
				Advertising Quality
				<h4>{format.number(avgPromotionQuality, false, 1)}</h4>
			</label>

			<label for="">
				Advertising Costs
				<h4>{format.currency(totalPromotionQuantity, false, 0)}</h4>
			</label>

			<label for="">
				Total Impressions
				<h4>{format.number(totalPromotionImpressions, false, 0)}</h4>
			</label>
		</div>

		<h3>Sales</h3>
		<div class="grid">
			<label for="">
				<h4>{format.number(totalProductsSold, false, 0)}</h4>
				Total Products Sold
			</label>

			<label for="">
				<h4>{format.number(totalMarketShare * 100, false, 1)}%</h4>
				Total Market Share
			</label>
		</div>
	</div>
{/if}
