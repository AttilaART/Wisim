<script>
	import { format } from '$lib/javascript/format';

	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void, product: import("$lib/javascript/simulation").Product, externalFactors: import("$lib/javascript/simulation").External_factors}} */
	let {
		clientState = $bindable(),
		updateDecisions,
		product = $bindable(),
		externalFactors
	} = $props();
</script>

<div>
	<table>
		<thead>
			<tr>
				<th colspan="2"> <h3>Product Stats</h3> </th>
				<th colspan="2"> <h3>Production Stats</h3> </th>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>Quality: </td>
				<td>{format.number(product.Quality_factor, false, 1)}</td>
				<td>Production Cost:</td>
				<td>{format.number(product.Production_cost, false, 1)}</td>
			</tr>
			<tr>
				<td>Ecology: </td>
				<td>{format.number(product.Ecology_factor, false, 1)}</td>
				<td>Material use:</td>
				<td>{format.number(product.Material_use, false, 1)}</td>
			</tr>
			<tr>
				<td>Ethics: </td>
				<td>{format.number(product.Ethics_factor, false, 1)}</td>
				<td>Material cost:</td>
				<td>{format.number(product.Material_use * externalFactors.Material_price, false, 1)}</td>
			</tr>
			<tr>
				<td>Durability: </td>
				<td>{product.Durabilty}</td>
				<td>Weight:</td>
				<td>{format.number(product.Weight, false, 1)}</td>
			</tr>
		</tbody>
	</table>
</div>
<form
	onchange={() => {
		updateDecisions(clientState.decisions);
	}}
>
	<div>
		<h2>General</h2>
		<label for="price"
			>Price
			<input id="price" bind:value={clientState.decisions.Marketing.Price} type="number" />
		</label>

		<h2>Materials</h2>
		<label for="quality"
			>Quality
			<input
				id="Quality"
				bind:value={clientState.decisions.Marketing.Product.Materials.Quality}
				type="range"
			/>
		</label>

		<label for="ProductEcology"
			>Ecology
			<input
				id="ProductEcology"
				bind:value={clientState.decisions.Marketing.Product.Materials.Ecology}
				type="range"
			/>
		</label>

		<label for="ProductEthicalSourcing"
			>Ethical Sourcing
			<input
				id="ProductEthicalSourcing"
				bind:value={clientState.decisions.Marketing.Product.Materials.Ethical_sourcing}
				type="range"
			/>
		</label>
	</div>

	<div>
		<h2>Manufacturing</h2>
		<label for="ProductQuality"
			>Quality
			<input
				id="ProductQuality"
				bind:value={clientState.decisions.Marketing.Product.Manufacturing.Quality}
				type="range"
			/>
		</label>

		<label for="ProductMaterialEfficiency"
			>Material Efficiency
			<input
				id="ProductMaterialEfficiency"
				bind:value={clientState.decisions.Marketing.Product.Manufacturing.Material_efficiency}
				type="range"
			/>
		</label>

		<label for="ProductEcoEnergy"
			>Eco-Energy
			<input
				id="ProductEcoEnergy"
				bind:value={clientState.decisions.Marketing.Product.Manufacturing.Ecological_energy}
				type="range"
			/>
		</label>

		<label for="ProductDurability"
			>Durability
			<input
				id="ProductDurability"
				bind:value={clientState.decisions.Marketing.Product.Manufacturing.Durability}
				type="range"
			/>
		</label>

		<label for="ProductMaxDurability"
			>Max Durability
			<input
				id="ProductMaxDurability"
				bind:value={clientState.decisions.Marketing.Product.Manufacturing.Max_durability}
				type="number"
			/>
		</label>
	</div>
</form>

<style>
	form {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--pico-spacing);
	}

	#price::after {
		content: ' CHF';
	}
</style>
