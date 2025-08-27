<script>
	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void}} */
	let { clientState = $bindable(), updateDecisions } = $props();
</script>

<h2>Research</h2>
<form
	onchange={() => {
		updateDecisions(clientState.decisions);
	}}
>
	<label for="">
		Quality Research
		<input
			type="number"
			step="10000"
			bind:value={clientState.decisions.Research.Quality}
			max="1000000000"
		/>
	</label>

	<label for="">
		Ecology Research
		<input type="number" step="10000" bind:value={clientState.decisions.Research.Ecology} />
	</label>

	<label for="">
		Production Speed Research
		<input type="number" step="10000" bind:value={clientState.decisions.Research.Production_cost} />
	</label>

	<label for="">
		Durability Research
		<input type="number" step="10000" bind:value={clientState.decisions.Research.Durability} />
	</label>

	<label for="">
		Marketing Research
		<input type="number" step="10000" bind:value={clientState.decisions.Research.Promotion} />
	</label>
</form>
