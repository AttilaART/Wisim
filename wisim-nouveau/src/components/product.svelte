<script>
	import { format } from '$lib/javascript/format';
	import { Calculate_product } from '../calculateProduct';
	import Window from './window.svelte';
	/** @typedef {Object} Props
	 * @property {import("$lib/javascript/simulation").clientState} clientState,
	 * @property {(Decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions,
	 * @property {Object.<string, import("$lib/javascript/simulation").Offer>} offers,
	 * @property {import("$lib/javascript/simulation").ExternalFactors} externalFactors,
	 * @property {(contents: import("svelte").Snippet<[number]>)=>number} newWindow,
	 * @property {(windowId: number)=>void} deleteWindow,
	 */

	/**
	 * @type {Props}
	 */
	let {
		clientState = $bindable(),
		updateDecisions,
		offers = $bindable(),
		externalFactors,
		newWindow,
		deleteWindow
	} = $props();

	let newProductID = $state(Math.trunc(Math.random() * 100000000));
	let newProductWindowID = $state(0);
	/** @type {Object.<string, string>} Key: windowID, Value: ProductID*/
	let productConfugureWindows = $state({});

	/** @type {import("$lib/javascript/simulation").Decisions_product} */
	const baseProductDecisions = {
		Price: 150,
		Name: 'Unnambed Product',

		Materials: {
			Quality: 1,
			Ecology: 1,
			EthicalSourcing: 1
		},

		Manufacturing: {
			Quality: 1,
			EcologicalEnergy: 1,
			MaterialEfficiency: 1,
			Durability: 1,
			MaxDurability: 1
		},

		Promotion: {
			Quantity: 10000,
			StyleQuality: 0.25,
			StyleEcology: 0.25,
			StyleEthics: 0.25,
			StyleDurability: 0.25
		}
	};

	const baseProduct = Calculate_product(
		clientState.Company,
		{
			ID: -1,
			CompanyID: clientState.Company.ID,
			Name: 'Unnambed Product',
			Weight: 1,
			MaterialUse: 1,
			ProductionCost: 1,

			Ethics: 1,
			Quality: 1,
			Ecology: 1,
			Durabilty: 1
		},
		baseProductDecisions,
		clientState.Employees.production
	);

	function generateProductID() {
		let id = 0;
		console.log(clientState.Company.Offers);
		console.log(id);
		while (clientState.Company.Offers[id] ? true : false) {
			id = Math.trunc(Math.random() * 10000000000);
		}
		return id;
	}
</script>

<div class="productsGrid">
	<button
		onclick={() => {
			deleteWindow(newProductWindowID);
			newProductID = generateProductID();
			clientState.Company.Offers[newProductID] = {
				Product: JSON.parse(JSON.stringify(baseProduct)),
				Price: 150,
				Promotion: {
					Quantity: 0,
					StyleQuality: 0,
					StyleDurability: 0,
					StyleEcology: 0,
					StyleEthics: 0
				},
				PromotionQuality: 0
			};
			clientState.Company.Offers[newProductID].Product.ID = newProductID;

			clientState.Decisions.Products[newProductID] = JSON.parse(
				JSON.stringify(baseProductDecisions)
			);
			console.log(clientState.Decisions);

			newProductWindowID = newWindow(newProduct);
		}}
	>
		<center>
			<h1 style="margin: 0px;">+</h1>
		</center>
		New Product
	</button>
	{#if offers}
		{#each Object.entries(offers) as offer}
			<article>
				<h3><input type="text" bind:value={offer[1].Product.Name} /></h3>
				<table>
					<tbody>
						<tr>
							<td>Quality: {format.number(offer[1].Product.Quality, false, 2)}</td>
							<td>Production Cost: {format.number(offer[1].Product.ProductionCost, false, 2)}</td>
						</tr>
						<tr>
							<td>Durability: {format.number(offer[1].Product.Durabilty, false, 2)}</td>
							<td>Materials Use: {format.number(offer[1].Product.MaterialUse, false, 2)}</td>
						</tr>
						<tr>
							<td>Ecology: {format.number(offer[1].Product.Ecology, false, 2)}</td>
							<td>Ethics: {format.number(offer[1].Product.Ethics, false, 2)}</td>
						</tr>
					</tbody>
				</table>
				<small
					>In storage: {clientState.Company.ProductsInStorage[offer[0]]
						? format.number(clientState.Company.ProductsInStorage[offer[0]], false, 0)
						: '0'}</small
				> <br />
				<small>Sold for</small>
				<h4><input type="number" bind:value={offer[1].Price} /></h4>
				<button
					onclick={() => {
						productConfugureWindows[newWindow(configurePromotion)] = offer[0];
					}}>Configure Promotion</button
				>
			</article>
		{/each}
	{/if}
</div>

{#snippet newProduct(/** @type {number} */ id)}
	<Window
		title="New Product"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
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
						<td
							>{format.number(
								clientState.Company.Offers[newProductID].Product.Quality,
								false,
								1
							)}</td
						>
						<td>Production Cost:</td>
						<td
							>{format.number(
								clientState.Company.Offers[newProductID].Product.ProductionCost,
								false,
								1
							)}</td
						>
					</tr>
					<tr>
						<td>Ecology: </td>
						<td
							>{format.number(
								clientState.Company.Offers[newProductID].Product.Ecology,
								false,
								1
							)}</td
						>
						<td>Material use:</td>
						<td
							>{format.number(
								clientState.Company.Offers[newProductID].Product.MaterialUse,
								false,
								1
							)}</td
						>
					</tr>
					<tr>
						<td>Ethics: </td>
						<td
							>{format.number(
								clientState.Company.Offers[newProductID].Product.Ethics,
								false,
								1
							)}</td
						>
						<td>Material cost:</td>
						<td
							>{format.number(
								clientState.Company.Offers[newProductID].Product.MaterialUse *
									externalFactors.MaterialPrice,
								false,
								1
							)}</td
						>
					</tr>
					<tr>
						<td>Durability: </td>
						<td>{clientState.Company.Offers[newProductID].Product.Durabilty}</td>
						<td>Weight:</td>
						<td
							>{format.number(
								clientState.Company.Offers[newProductID].Product.Weight,
								false,
								1
							)}</td
						>
					</tr>
				</tbody>
			</table>
		</div>
		<form
			onchange={() => {
				updateDecisions(clientState.Decisions);
				clientState.Company.Offers[newProductID].Product = Calculate_product(
					clientState.Company,
					clientState.Company.Offers[newProductID].Product,
					clientState.Decisions.Products[newProductID],
					clientState.Employees.production
				);
			}}
		>
			<div>
				<h2>General</h2>
				<label for="name"
					>Name
					<input
						id="name"
						bind:value={clientState.Decisions.Products[newProductID].Name}
						type="text"
					/>
				</label>
				<label for="price"
					>Price
					<input
						id="price"
						bind:value={clientState.Decisions.Products[newProductID].Price}
						type="number"
					/>
				</label>
				<!--
				<label for="maxProduction"
					>Product goal
					<input
						id="maxProduction"
						bind:value={clientState.Decisions.Products[newProductID].ProductionGoal}
						type="number"
					/>
				</label>-->

				<h2>Materials</h2>
				<label for="quality"
					>Quality
					<input
						id="Quality"
						bind:value={clientState.Decisions.Products[newProductID].Materials.Quality}
						type="range"
					/>
				</label>

				<label for="ProductEcology"
					>Ecology
					<input
						id="ProductEcology"
						bind:value={clientState.Decisions.Products[newProductID].Materials.Ecology}
						type="range"
					/>
				</label>

				<label for="ProductEthicalSourcing"
					>Ethical Sourcing
					<input
						id="ProductEthicalSourcing"
						bind:value={clientState.Decisions.Products[newProductID].Materials.EthicalSourcing}
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
						bind:value={clientState.Decisions.Products[newProductID].Manufacturing.Quality}
						type="range"
					/>
				</label>

				<label for="ProductMaterialEfficiency"
					>Material Efficiency
					<input
						id="ProductMaterialEfficiency"
						bind:value={
							clientState.Decisions.Products[newProductID].Manufacturing.MaterialEfficiency
						}
						type="range"
					/>
				</label>

				<label for="ProductEcoEnergy"
					>Eco-Energy
					<input
						id="ProductEcoEnergy"
						bind:value={clientState.Decisions.Products[newProductID].Manufacturing.EcologicalEnergy}
						type="range"
					/>
				</label>

				<label for="ProductDurability"
					>Durability
					<input
						id="ProductDurability"
						bind:value={clientState.Decisions.Products[newProductID].Manufacturing.Durability}
						type="range"
					/>
				</label>

				<label for="ProductMaxDurability"
					>Max Durability
					<input
						id="ProductMaxDurability"
						bind:value={clientState.Decisions.Products[newProductID].Manufacturing.MaxDurability}
						type="number"
					/>
				</label>
			</div>
			<input
				type="submit"
				value="Cancel"
				onclick={() => {
					delete clientState.Decisions.Products[newProductID];
					deleteWindow(newProductWindowID);
				}}
			/>
			<input
				type="submit"
				onclick={() => {
					updateDecisions(clientState.Decisions);
					deleteWindow(newProductWindowID);
				}}
				value="Confirm"
			/>
		</form>
	</Window>
{/snippet}

{#snippet configurePromotion(/** @type {number}*/ windowID)}
	<Window
		title={`Configure ${clientState.Company.Offers[productConfugureWindows[windowID]].Product.Name}`}
		closeWindow={() => {
			deleteWindow(windowID);
			delete productConfugureWindows[windowID];
		}}
	>
		<form
			onchange={() => {
				updateDecisions(clientState.Decisions);
			}}
		>
			<label for="PromotionQuantity">
				<h2>Advertisment Budget</h2>
				<input
					id="quantity"
					bind:value={
						clientState.Decisions.Products[productConfugureWindows[windowID]].Promotion.Quantity
					}
					type="number"
				/>
			</label>

			<h2>Advertisment Style</h2>
			<label for="PromotionQuality"
				>Quality
				<input
					id="PromotionQuality"
					bind:value={
						clientState.Decisions.Products[productConfugureWindows[windowID]].Promotion.StyleQuality
					}
					type="range"
				/>
			</label>

			<label for="PromotionEcology"
				>Ecology
				<input
					id="PromotionEcology"
					bind:value={
						clientState.Decisions.Products[productConfugureWindows[windowID]].Promotion.StyleEcology
					}
					type="range"
				/>
			</label>

			<label for="PromotionEthicals"
				>Ethics
				<input
					id="PromotionEthicals"
					bind:value={
						clientState.Decisions.Products[productConfugureWindows[windowID]].Promotion.StyleEthics
					}
					type="range"
				/>
			</label>

			<label for="PromotionDurability"
				>Durability
				<input
					id="PromotionDurability"
					bind:value={
						clientState.Decisions.Products[productConfugureWindows[windowID]].Promotion
							.StyleDurability
					}
					type="range"
				/>
			</label>
		</form>
	</Window>
{/snippet}

<style>
	form {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--pico-spacing);
	}

	#price::after {
		content: ' CHF';
	}

	.productsGrid {
		display: flex;
		flex-direction: column;
		gap: var(--pico-spacing);
	}
</style>
