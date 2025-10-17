<script>
	import { preventPageReload } from '$lib/helper.svelte';

	/** @typedef {Object} Props
	 * @property {import("$lib/javascript/simulation").clientState} clientState,
	 * @property {(Decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions,
	 * @property {string} productID
	 */

	/** @type {Props} */

	let { clientState = $bindable(), updateDecisions, productID } = $props();

	/**
	 * @param {{Quantity: number, Price: number, Quality: number, Ecology: number, Ethics: number, Durability: number}} promotion
	 * @param {string} changed which paremeter was changed
	 */
	function normalise(promotion, changed) {
		let lengthExceptChanged = 0;

		for (let field of Object.keys(promotion)) {
			if (field == 'Quantity') {
				continue;
			} else if (field == changed) {
				continue;
			}

			console.log(promotion[field]);
			// @ts-ignore
			lengthExceptChanged += promotion[field] * promotion[field];
		}

		lengthExceptChanged = Math.sqrt(lengthExceptChanged);

		if (lengthExceptChanged == 0) {
			return promotion;
		}

		let multiplicationFactor = Math.sqrt(1 - promotion[changed]) / lengthExceptChanged;
		console.log(multiplicationFactor);

		for (let field of Object.keys(promotion)) {
			if (field == 'Quantity') {
				continue;
			} else if (field == changed) {
				continue;
			}
			promotion[field] *= multiplicationFactor;
		}

		return promotion;
	}
</script>

<form
	use:preventPageReload
	onchange={() => {
		updateDecisions(clientState.Decisions);
	}}
>
	<label for="PromotionQuantity">
		<h2>Advertisment Budget</h2>
		<input
			id="quantity"
			bind:value={clientState.Decisions.Products[productID].Promotion.Quantity}
			type="number"
		/>
	</label>

	<h2>Advertisment Style</h2>
	<label for="PromotionQuality"
		>Quality
		<input
			oninput={() => {
				clientState.Decisions.Products[productID].Promotion = normalise(
					clientState.Decisions.Products[productID].Promotion,
					'Quality'
				);
			}}
			id="PromotionQuality"
			bind:value={clientState.Decisions.Products[productID].Promotion.Quality}
			type="range"
			max={1}
			step={0.01}
		/>
	</label>

	<label for="PromotionEcology"
		>Ecology
		<input
			oninput={() => {
				clientState.Decisions.Products[productID].Promotion = normalise(
					clientState.Decisions.Products[productID].Promotion,
					'Ecology'
				);
			}}
			id="PromotionEcology"
			bind:value={clientState.Decisions.Products[productID].Promotion.Ecology}
			type="range"
			max={1}
			step={0.01}
		/>
	</label>

	<label for="PromotionEthicals"
		>Ethics
		<input
			oninput={() => {
				clientState.Decisions.Products[productID].Promotion = normalise(
					clientState.Decisions.Products[productID].Promotion,
					'Ethics'
				);
			}}
			id="PromotionEthicals"
			bind:value={clientState.Decisions.Products[productID].Promotion.Ethics}
			type="range"
			max={1}
			step={0.01}
		/>
	</label>

	<label for="PromotionDurability"
		>Durability
		<input
			oninput={() => {
				clientState.Decisions.Products[productID].Promotion = normalise(
					clientState.Decisions.Products[productID].Promotion,
					'Durability'
				);
			}}
			id="PromotionDurability"
			bind:value={clientState.Decisions.Products[productID].Promotion.Durability}
			type="range"
			max={1}
			step={0.01}
		/>
	</label>
	<label for="PromotionDurability"
		>Price
		<input
			oninput={() => {
				clientState.Decisions.Products[productID].Promotion = normalise(
					clientState.Decisions.Products[productID].Promotion,
					'Price'
				);
			}}
			id="PromotionDurability"
			bind:value={clientState.Decisions.Products[productID].Promotion.Price}
			type="range"
			max={1}
			step={0.01}
		/>
	</label>
</form>
