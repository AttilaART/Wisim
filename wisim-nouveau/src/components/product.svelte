<script>
	import { format } from '$lib/javascript/format';
	import { calculateProductStats } from '../calculateProduct';
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

	const images = import.meta.glob(['$lib/images/Products/Base_blueprint/*.svg'], {
		eager: true,
		as: 'url'
	});

	let newProductID = $state(`${Math.trunc(Math.random() * 100000000)}`);
	let newProductWindowID = $state(0);
	/** @type {Object.<string, string>} Key: wi:windowndowID, Value: ProductID*/
	let productConfugureWindows = $state({});

	/** @type {import("$lib/javascript/simulation").Product} */
	const baseProduct = {
		ID: '-1',
		CompanyID: clientState.Company.ID,
		Name: 'Unnamed Product',

		Components: {
			FormFactor: null,
			Frame: null,
			Body: null,
			Mechanism: null,
			Misc: []
		},

		MaterialQuality: 0,
		ExtraDurability: 0,
		ExtraQuality: 0
	};

	/** @type {import("$lib/javascript/simulation").Decisions_product} */
	const baseProductDecisions = {
		Price: 150,
		Name: baseProduct.Name,
		Promotion: {
			Quantity: 10000,
			Quality: 0.2,
			Durability: 0.2,
			Price: 0.2,
			Ecology: 0.2,
			Ethics: 0.2
		},
		Product: baseProduct
	};

	/** @type {Object.<string, number>} */
	let productionLineCosts = $state({});
	/** @type {import("svelte").Snippet? } */
	let currentComponentSnippet = $state(null);

	/** @param {HTMLElement} el */
	function focusElement(el) {
		el.focus();
	}

	function generateProductID() {
		let id = 0;
		console.log(clientState.Company.Offers);
		console.log(id);
		while (clientState.Company.Offers[id] ? true : false) {
			id = Math.trunc(Math.random() * 10000000000);
		}
		return String(id);
	}
</script>

<div class="productsGrid">
	<button
		onclick={() => {
			deleteWindow(newProductWindowID);
			newProductID = generateProductID();
			clientState.Company.Offers[newProductID] = {
				Product: JSON.parse(JSON.stringify(baseProduct)),
				productStats: (() => {
					/** @type {{productStats: import("$lib/javascript/simulation").ProductStats, productionLineCost: number}}*/ let {
						productStats,
						productionLineCost
					} = calculateProductStats(baseProduct, clientState.productComponents);
					return productStats;
				})(),

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
							<td>Quality: {format.number(offer[1].productStats.Quality, false, 2)}</td>
							<td
								>Production Cost: {format.number(
									offer[1].productStats.ProductionCost,
									false,
									2
								)}</td
							>
						</tr>
						<tr>
							<td>Durability: {format.number(offer[1].productStats.Durability, false, 2)}</td>
							<td>Materials Use: {format.number(offer[1].productStats.MaterialUse, false, 2)}</td>
						</tr>
						<tr>
							<td>Ecology: {format.number(offer[1].productStats.Ecology, false, 2)}</td>
							<td>Ethics: {format.number(offer[1].productStats.Ethics, false, 2)}</td>
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
		<form
			class="main-product-div"
			onchange={() => {
				updateDecisions(clientState.Decisions);
				let { productStats, productionLineCost } = calculateProductStats(
					clientState.Company.Offers[newProductID].Product,
					clientState.productComponents
				);

				clientState.Company.Offers[newProductID].productStats = productStats;
				productionLineCosts[newProductID] = productionLineCost;
			}}
		>
			<div>
				<div>
					<h2>
						<input
							id="name"
							bind:value={clientState.Decisions.Products[newProductID].Name}
							type="text"
						/>
					</h2>
				</div>
				<div class="productDesignerGrid">
					<button
						onclick={() => {
							currentComponentSnippet = productFormFactor;
						}}
					>
						{#if !clientState.Decisions.Products[newProductID].Product.Components.FormFactor}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.FormFactor[
											`${clientState.Decisions.Products[newProductID].Product.Components.FormFactor}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
							{console.log(
								clientState.Decisions.Products[newProductID].Product.Components.FormFactor
							)}
						{/if}
					</button>
					<button
						onclick={() => {
							currentComponentSnippet = productFormFactor;
						}}
					>
						{#if !clientState.Decisions.Products[newProductID].Product.Components.FormFactor}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.FormFactor[
											`${clientState.Decisions.Products[newProductID].Product.Components.FormFactor}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
						{/if}
					</button>
					<div>+</div>
					<div>+</div>

					{#if currentComponentSnippet == null}
						<div
							style="grid-column: 1 /span 4; grid-row: 2  /span 5; height: 100%; background: #4A6DE5;"
						>
							<center>
								<img
									src={images[
										'/src/lib/images/' +
											clientState.productComponents.FormFactor[
												`${clientState.Decisions.Products[newProductID].Product.Components.FormFactor}`
											]?.Image
									]}
									style="mix-blend-mode: lighten; max-height: 10rem;"
									alt=""
								/>
							</center>
						</div>
					{:else}
						<div
							tabindex="-1"
							use:focusElement
							style="grid-column: 1 /span 4; grid-row: 2  /span 5; height: 100%;"
							onfocus={() => {
								console.log('aasd');
							}}
						>
							{@render currentComponentSnippet()}
						</div>
					{/if}

					<div>+</div>
					<div>+</div>
					<div>+</div>
					<div>+</div>
				</div>
				<!--
				<label for="maxProduction"
					>Product goal
					<input
						id="maxProduction"
						bind:value={clientState.Decisions.Products[newProductID].ProductionGoal}
						type="number"
					/>
				</label>-->
				<fieldset role="group">
					<label for="ProductMaterialQuality"
						>Material Quality
						<input
							id="ProductMaterialQuality"
							bind:value={clientState.Decisions.Products[newProductID].Product.MaterialQuality}
							type="number"
						/>
					</label>

					<label for="ProductDurability"
						>Durability
						<input
							id="ProductDurability"
							bind:value={clientState.Decisions.Products[newProductID].Product.ExtraDurability}
							type="number"
						/>
					</label>

					<label for="quality"
						>Quality
						<input
							id="Quality"
							bind:value={clientState.Decisions.Products[newProductID].Product.ExtraQuality}
							type="number"
						/>
					</label>
				</fieldset>
			</div>
			<div>
				<table>
					<thead>
						<tr>
							<th colspan="2"> <h3>Manufacturing Stats</h3> </th>
							<th colspan="2"> <h3>Product Stats</h3> </th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>Production Cost:</td>
							<td
								>{format.number(
									clientState.Company.Offers[newProductID].productStats.ProductionCost,
									false,
									1
								)}</td
							>
							<td>Quality: </td>
							<td
								>{format.number(
									clientState.Company.Offers[newProductID].productStats.Quality,
									false,
									1
								)}</td
							>
						</tr>
						<tr>
							<td>Material use:</td>
							<td
								>{format.number(
									clientState.Company.Offers[newProductID].productStats.MaterialUse,
									false,
									1
								)}</td
							>
							<td>Ecology: </td>
							<td
								>{format.number(
									clientState.Company.Offers[newProductID].productStats.Ecology,
									false,
									1
								)}</td
							>
						</tr>
						<tr>
							<td>Material cost:</td>
							<td
								>{format.number(
									clientState.Company.Offers[newProductID].productStats.MaterialUse *
										externalFactors.MaterialPrice,
									false,
									1
								)}</td
							>
							<td>Ethics: </td>
							<td
								>{format.number(
									clientState.Company.Offers[newProductID].productStats.Ethics,
									false,
									1
								)}</td
							>
						</tr>
						<tr>
							<td></td>
							<td></td>
							<td>Durability: </td>
							<td>{clientState.Company.Offers[newProductID].productStats.Durability}</td>
						</tr>
						<tr>
							<td></td>
							<td></td>
							<td>Weight:</td>
							<td
								><!--{format.number(
								clientState.Company.Offers[newProductID].productStats.Weight,
								false,
								1
							)}--></td
							>
						</tr>
					</tbody>
				</table>

				<label for="price"
					>Price
					<input
						id="price"
						bind:value={clientState.Decisions.Products[newProductID].Price}
						type="number"
					/>
				</label>
			</div>
			<footer>
				<input
					type="submit"
					value="Cancel"
					onclick={() => {
						delete clientState.Decisions.Products[newProductID];
						delete clientState.Company.Offers[newProductID];
						deleteWindow(newProductWindowID);
					}}
				/>
				<div id="spacer" style="flex: 1 1 0;"></div>
				<div>
					Production cost 🔧 {clientState.Company.Offers[newProductID].productStats.ProductionCost}
				</div>
				<input
					type="submit"
					onclick={() => {
						updateDecisions(clientState.Decisions);
						deleteWindow(newProductWindowID);
					}}
					value="Confirm"
				/>
			</footer>
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
						clientState.Decisions.Products[productConfugureWindows[windowID]].Promotion.Quality
					}
					type="range"
				/>
			</label>

			<label for="PromotionEcology"
				>Ecology
				<input
					id="PromotionEcology"
					bind:value={
						clientState.Decisions.Products[productConfugureWindows[windowID]].Promotion.Ecology
					}
					type="range"
				/>
			</label>

			<label for="PromotionEthicals"
				>Ethics
				<input
					id="PromotionEthicals"
					bind:value={
						clientState.Decisions.Products[productConfugureWindows[windowID]].Promotion.Ethics
					}
					type="range"
				/>
			</label>

			<label for="PromotionDurability"
				>Durability
				<input
					id="PromotionDurability"
					bind:value={
						clientState.Decisions.Products[productConfugureWindows[windowID]].Promotion.Durability
					}
					type="range"
				/>
			</label>
		</form>
	</Window>
{/snippet}

{#snippet productFormFactor()}
	<div class="productComponentSelector">
		{#each Object.entries(clientState.productComponents.FormFactor) as c}
			<button
				onclick={() => {
					clientState.Company.Offers[newProductID].Product.Components.FormFactor = c[0];
					clientState.Decisions.Products[newProductID].Product.Components.FormFactor = c[0];
					currentComponentSnippet = null;
				}}
			>
				<img
					style="mix-blend-mode: lighten;"
					src={images['/src/lib/images/' + c[1].Image]}
					alt=""
				/>
				{c[1].Name}
			</button>
		{/each}
	</div>
{/snippet}

<style>
	.main-product-div {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--pico-spacing);
		footer {
			grid-column: 1 / span 2;
			display: flex;
			flex-direction: row;
			justify-content: flex-end;
			gap: var(--pico-spacing);
			* {
				flex: 0 0;
			}
		}
	}

	.productDesignerGrid {
		width: 100%;
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		grid-template-rows: repeat(6, 1fr);
	}

	.productDesignerGrid > button {
		text-align: center;
		text-justify: center;
		padding: 3px;
		height: 2.5rem;
	}

	.productComponentSelector {
		display: grid;
		padding: var(--pico-spacing);
		gap: var(--pico-spacing);
		grid-template-columns: 1fr 1fr 1fr;
		overflow: scroll;
		height: 100% !important;
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
