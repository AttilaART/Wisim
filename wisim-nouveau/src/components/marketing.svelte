<script>
	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void}} */
	let { clientState = $bindable(), updateDecisions } = $props();
</script>

<label for="price">Price</label>
<input
	id="price"
	bind:value={
		() => {
			return clientState.company.Offer.Price;
		},
		(value /** @type {number} */) => {
			clientState.company.Offer.Price = value;
			updateDecisions(clientState.decisions);
		}
	}
	type="number"
/>
