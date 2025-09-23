<script>
	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void}} */
	let { clientState = $bindable(), updateDecisions } = $props();
</script>

<form
	class="grid"
	onchange={() => {
		updateDecisions(clientState.Decisions);
	}}
>
	<div>
		<label for="PromotionQuantity">
			<h2>Advertisment Budget</h2>
			<input
				id="quantity"
				bind:value={clientState.Decisions.Marketing.Promotion.Quantity}
				type="number"
			/>
		</label>

		<h2>Advertisment Style</h2>
		<label for="PromotionQuality"
			>Quality
			<input
				id="PromotionQuality"
				bind:value={clientState.Decisions.Marketing.Promotion.StyleQuality}
				type="range"
			/>
		</label>

		<label for="PromotionEcology"
			>Ecology
			<input
				id="PromotionEcology"
				bind:value={clientState.Decisions.Marketing.Promotion.StyleEcology}
				type="range"
			/>
		</label>

		<label for="PromotionEthicals"
			>Ethics
			<input
				id="PromotionEthicals"
				bind:value={clientState.Decisions.Marketing.Promotion.StyleEthics}
				type="range"
			/>
		</label>

		<label for="PromotionDurability"
			>Durability
			<input
				id="PromotionDurability"
				bind:value={clientState.Decisions.Marketing.Promotion.StyleDurability}
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
