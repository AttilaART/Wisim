<script>
	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void}} */
	let { clientState = $bindable(), updateDecisions } = $props();
</script>

<form
	class="grid"
	onchange={() => {
		updateDecisions(clientState.decisions);
	}}
>
	<div>
		<label for="PromotionQuantity">
			<h2>Advertisment Budget</h2>
			<input
				id="quantity"
				bind:value={clientState.decisions.Marketing.Promotion.Quantity}
				type="number"
			/>
		</label>

		<h2>Advertisment Style</h2>
		<label for="PromotionQuality"
			>Quality
			<input
				id="PromotionQuality"
				bind:value={clientState.decisions.Marketing.Promotion.Style_quality}
				type="range"
			/>
		</label>

		<label for="PromotionEcology"
			>Ecology
			<input
				id="PromotionEcology"
				bind:value={clientState.decisions.Marketing.Promotion.Style_ecology}
				type="range"
			/>
		</label>

		<label for="PromotionEthicals"
			>Ethics
			<input
				id="PromotionEthicals"
				bind:value={clientState.decisions.Marketing.Promotion.Style_ethics}
				type="range"
			/>
		</label>

		<label for="PromotionDurability"
			>Durability
			<input
				id="PromotionDurability"
				bind:value={clientState.decisions.Marketing.Promotion.Style_durability}
				type="range"
			/>
		</label>
	</div>
</form>

<style>
	#price::after {
		content: ' CHF';
	}
</style>
