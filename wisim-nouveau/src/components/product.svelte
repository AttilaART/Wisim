<script>
	import { format } from '$lib/javascript/format';
	import { calculateProductStats } from '../calculateProduct';
	import storageIcon from '$lib/images/warehouse.svg';
	import Increment from './increment.svelte';
	import Window from './window.svelte';
	import ProductionIcon from '$lib/images/production.svg';
	import MarketingIcon from '$lib/images/marketing.svg';
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
			FormFactor: 'FormFactorMedium',
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

	/** @type {import("$lib/javascript/simulation").Product} */
	let currentProductHover = $state(JSON.parse(JSON.stringify(baseProduct)));
	/** @type {import("$lib/javascript/simulation").ProductStats} */
	let currentProductStatsHoverStats = $derived(
		(() => {
			return calculateProductStats(currentProductHover, clientState.productComponents).productStats;
		})()
	);
	/** @type {number} */
	let currentProductProductionLineCostsHover = $derived(
		(() => {
			return calculateProductStats(currentProductHover, clientState.productComponents)
				.productionLineCost;
		})()
	);

	/** @type {import("svelte").Snippet? } */
	let currentComponentSnippet = $state(null);

	/** @param {HTMLElement} el */
	function focusElement(el) {
		el.focus();
	}

	function generateProductID() {
		let id = 0;
		// console.log(clientState.Company.Offers);
		// console.log(id);
		while (clientState.Company.Offers[id] ? true : false) {
			id = Math.trunc(Math.random() * 10000000000);
		}
		return String(id);
	}

	function onProductChange() {
		updateDecisions(clientState.Decisions);

		currentProductHover = JSON.parse(
			JSON.stringify(clientState.Decisions.Products[newProductID].Product)
		);
		clientState.Company.Offers[newProductID].Product = JSON.parse(
			JSON.stringify(clientState.Decisions.Products[newProductID].Product)
		);

		let { productStats, productionLineCost } = calculateProductStats(
			clientState.Decisions.Products[newProductID].Product,
			$state.snapshot(clientState.productComponents)
		);

		clientState.Company.Offers[newProductID].productStats = productStats;
		productionLineCosts[newProductID] = productionLineCost;
	}

	function resetHoverProduct() {
		currentProductHover = JSON.parse(
			JSON.stringify(clientState.Decisions.Products[newProductID].Product)
		);
	}

	$effect(() => {
		console.log();
		if (clientState.Company.Offers[newProductID]) {
			clientState.Decisions.Products[newProductID].Product =
				clientState.Company.Offers[newProductID].Product;
		}
		updateDecisions(clientState.Decisions);
	});
</script>

<div class="productsGrid">
	<button
		onclick={() => {
			deleteWindow(newProductWindowID);
			newProductID = generateProductID();
			clientState.Company.Offers[newProductID] = {
				Product: JSON.parse(JSON.stringify(baseProduct)),
				productStats: calculateProductStats(baseProduct, clientState.productComponents)
					.productStats,

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
			//console.log(clientState.Decisions);

			newProductWindowID = newWindow(newProduct);
		}}
	>
		<div style="line-height: 3.5rem;">
			<span style="font-size: 3rem; font-weight: bolder; vertical-align: middle;">+</span>
			<span style="margin: 0px; vertical-align: middle;"> New Product </span>
		</div>
	</button>
	{#if offers}
		{#each Object.entries(offers) as offer}
			<article
				style="position: relative; display: grid; grid-template-columns: 70% 30%; grid-template-rows: auto auto auto; gap: 0.5rem;"
			>
				<img
					src={images[
						'/src/lib/images/' +
							clientState.productComponents.FormFactor[`${offer[1].Product.Components.FormFactor}`]
								.Image
					]}
					alt=""
					style="position: absolute; pointer-events: none; left: 50%; top: 50%; transform: translate(-50%, -50%); height: 6rem;"
				/>
				<div>
					<button
						class="inlineIcon marketingIcon"
						onclick={() => {
							productConfugureWindows[newWindow(configurePromotion)] = offer[0];
						}}
						aria-label="Marketing"
					></button>
					<h4 style="display: inline;">
						{offer[1].Product.Name}
					</h4>
				</div>
				<div>
					<h4>{format.currency(offer[1].Price, false, 0)}</h4>
					<div>
						{clientState.Company.ProductsInStorage[offer[0]]
							? format.number(clientState.Company.ProductsInStorage[offer[0]], false, 0)
							: '0'}
						<img
							class="inlineIcon"
							style="height: 1.2rem; translate: 0 -0.1rem ;"
							src={storageIcon}
							alt=""
						/>
						<small>+50</small>
					</div>
				</div>
				<label for="outdated{offer[1].Product.ID}">
					<input id="outdated{offer[1].Product.ID}" type="checkbox" />
					Mark as outdated
				</label>
				<button style="padding: 0.5rem;">Make a copy</button>
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
		<form class="main-product-div" onchange={onProductChange}>
			<div>
				<div>
					<h2>
						<input
							id="name"
							bind:value={clientState.Company.Offers[newProductID].Product.Name}
							type="text"
							autocomplete="off"
						/>
					</h2>
				</div>
				<div class="productDesignerGrid">
					<button
						class="outline contrast"
						onclick={() => {
							currentComponentSnippet = productFormFactor;
						}}
					>
						{#if !clientState.Company.Offers[newProductID].Product.Components.FormFactor}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.FormFactor[
											`${clientState.Company.Offers[newProductID].Product.Components.FormFactor}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
						{/if}
					</button>
					<button
						class="outline contrast"
						onclick={() => {
							currentComponentSnippet = productFrame;
						}}
					>
						{#if !clientState.Company.Offers[newProductID].Product.Components.Frame}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.Frame[
											`${clientState.Decisions.Products[newProductID].Product.Components.Frame}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
						{/if}
					</button>
					<button
						class="outline contrast"
						onclick={() => {
							currentComponentSnippet = productBody;
						}}
					>
						{#if !clientState.Company.Offers[newProductID].Product.Components.Body}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.Body[
											`${clientState.Company.Offers[newProductID].Product.Components.Body}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
						{/if}
					</button>
					<button
						class="outline contrast"
						onclick={() => {
							currentComponentSnippet = productCoffeeMechanism;
						}}
					>
						{#if !clientState.Company.Offers[newProductID].Product.Components.Mechanism}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.Mechanism[
											`${clientState.Company.Offers[newProductID].Product.Components.Mechanism}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
						{/if}
					</button>

					<div
						style="grid-column: 1 /span 4; grid-row: 2  /span 5; position: relative;"
						use:focusElement
						tabindex="-1"
					>
						{#if currentComponentSnippet == null}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.FormFactor[
											`${clientState.Company.Offers[newProductID].Product.Components.FormFactor}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; width: 100%; height: 10em; object-fit: contain; position: relative; left: 50%; top: 50%; translate: -50% -50%;"
								alt=""
							/>
						{:else}
							{@render currentComponentSnippet()}
						{/if}
					</div>

					<button
						class="outline contrast"
						onclick={() => {
							currentComponentSnippet = productMisc0;
						}}
					>
						{#if !clientState.Company.Offers[newProductID].Product.Components.Misc[0]}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.Misc[
											`${clientState.Company.Offers[newProductID].Product.Components.Misc[0]}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
						{/if}
					</button>
					<button
						class="outline contrast"
						onclick={() => {
							currentComponentSnippet = productMisc1;
						}}
					>
						{#if !clientState.Company.Offers[newProductID].Product.Components.Misc[1]}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.Misc[
											`${clientState.Company.Offers[newProductID].Product.Components.Misc[1]}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
						{/if}
					</button>
					<button
						class="outline contrast"
						onclick={() => {
							currentComponentSnippet = productMisc2;
						}}
					>
						{#if !clientState.Company.Offers[newProductID].Product.Components.Misc[2]}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.Misc[
											`${clientState.Company.Offers[newProductID].Product.Components.Misc[2]}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
						{/if}
					</button>
					<button
						class="outline contrast"
						onclick={() => {
							currentComponentSnippet = productMisc3;
						}}
					>
						{#if !clientState.Company.Offers[newProductID].Product.Components.Misc[3]}
							<span style="font-weight: bold; font-size: 1.5rem;">+</span>
						{:else}
							<img
								src={images[
									'/src/lib/images/' +
										clientState.productComponents.Misc[
											`${clientState.Company.Offers[newProductID].Product.Components.Misc[3]}`
										]?.Image
								]}
								style="mix-blend-mode: lighten; max-height: 1.5rem;"
								alt=""
							/>
						{/if}
					</button>
				</div>
				<!--
				<label for="maxProduction"
					>Product goal
					<input
						id="maxProduction"
						bind:value={clientState.Company.Offers[newProductID].ProductionGoal}
						type="number"
					/>
				</label>-->
				<!--<Increment
					bind:value={clientState.Company.Offers[newProductID].Product.MaterialQuality}
					label="Material Quality"
					min={0}
					max={99}
					onclick={onProductChange}
				/>-->

				<Increment
					bind:value={clientState.Company.Offers[newProductID].Product.ExtraDurability}
					label="Durability"
					min={0}
					max={99}
					onclick={onProductChange}
				/>

				<Increment
					bind:value={clientState.Company.Offers[newProductID].Product.ExtraQuality}
					label="Quality"
					min={0}
					max={99}
					onclick={onProductChange}
				/>
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
								>{format.number(currentProductStatsHoverStats.ProductionCost, false, 1)}
								{@render showHoverDifference(
									clientState.Company.Offers[newProductID].productStats.ProductionCost,
									currentProductStatsHoverStats.ProductionCost,
									(value) => {
										return format.number(value, true, 1);
									},
									true
								)}
							</td>
							<td>Quality: </td>
							<td
								>{format.number(currentProductStatsHoverStats.Quality, false, 1)}
								{@render showHoverDifference(
									clientState.Company.Offers[newProductID].productStats.Quality,
									currentProductStatsHoverStats.Quality,
									(value) => {
										return format.number(value, true, 1);
									},
									false
								)}
							</td>
						</tr>
						<tr>
							<td>Material use:</td>
							<td
								>{format.number(currentProductStatsHoverStats.MaterialUse, false, 1)}
								{@render showHoverDifference(
									clientState.Company.Offers[newProductID].productStats.MaterialUse,
									currentProductStatsHoverStats.MaterialUse,
									(value) => {
										return format.number(value, true, 1);
									},
									false
								)}
							</td>
							<td>Ecology: </td>
							<td
								>{format.number(currentProductStatsHoverStats.Ecology, false, 1)}
								{@render showHoverDifference(
									clientState.Company.Offers[newProductID].productStats.Ecology,
									currentProductStatsHoverStats.Ecology,
									(value) => {
										return format.number(value, true, 1);
									},
									false
								)}
							</td>
						</tr>
						<tr>
							<td>Material cost:</td>
							<td
								>{format.number(
									currentProductStatsHoverStats.MaterialUse * externalFactors.MaterialPrice,
									false,
									1
								)}
								{@render showHoverDifference(
									clientState.Company.Offers[newProductID].productStats.MaterialUse *
										externalFactors.MaterialPrice,
									currentProductStatsHoverStats.MaterialUse * externalFactors.MaterialPrice,
									(value) => {
										return format.number(value, true, 1);
									},
									false
								)}
							</td>
							<td>Ethics: </td>
							<td
								>{format.number(currentProductStatsHoverStats.Ethics, false, 1)}
								{@render showHoverDifference(
									clientState.Company.Offers[newProductID].productStats.Ethics,
									currentProductStatsHoverStats.Ethics,
									(value) => {
										return format.number(value, true, 1);
									},
									false
								)}
							</td>
						</tr>
						<tr>
							<td></td>
							<td></td>
							<td>Durability: </td>
							<td
								>{currentProductStatsHoverStats.Durability}
								{@render showHoverDifference(
									clientState.Company.Offers[newProductID].productStats.Durability,
									currentProductStatsHoverStats.Durability,
									(value) => {
										return format.number(value, true, 1);
									},
									false
								)}
							</td>
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
					Production cost <img class="inlineIcon" src={ProductionIcon} alt="production" />
					{clientState.Company.Offers[newProductID].productStats.ProductionCost}
				</div>
				<input
					type="submit"
					onclick={() => {
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
				onmouseover={() => {
					currentProductHover.Components.FormFactor = c[0];
				}}
				onmouseout={resetHoverProduct}
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

{#snippet productFrame()}
	<div class="productComponentSelector">
		{#each Object.entries(clientState.productComponents.Frame) as c}
			<button
				onclick={() => {
					clientState.Decisions.Products[newProductID].Product.Components.Frame = c[0];
					onProductChange();
					currentComponentSnippet = null;
				}}
				onmouseover={() => {
					currentProductHover.Components.Frame = c[0];
				}}
				onmouseout={resetHoverProduct}
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

{#snippet productBody()}
	<div class="productComponentSelector">
		{#each Object.entries(clientState.productComponents.Body) as c}
			<button
				onclick={() => {
					clientState.Company.Offers[newProductID].Product.Components.Body = c[0];
					clientState.Decisions.Products[newProductID].Product.Components.Body = c[0];
					onProductChange();
					currentComponentSnippet = null;
				}}
				onmouseover={() => {
					currentProductHover.Components.Body = c[0];
				}}
				onmouseout={resetHoverProduct}
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

{#snippet productCoffeeMechanism()}
	<div class="productComponentSelector">
		{#each Object.entries(clientState.productComponents.Mechanism) as c}
			<button
				onclick={() => {
					clientState.Company.Offers[newProductID].Product.Components.Mechanism = c[0];
					clientState.Decisions.Products[newProductID].Product.Components.Mechanism = c[0];
					onProductChange();
					currentComponentSnippet = null;
				}}
				onmouseover={() => {
					currentProductHover.Components.Mechanism = c[0];
				}}
				onmouseout={resetHoverProduct}
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

{#snippet productMisc(/** @type {number} }*/ miscSlot)}
	<div class="productComponentSelector">
		{#each Object.entries(clientState.productComponents.Misc) as c}
			<button
				onclick={() => {
					clientState.Company.Offers[newProductID].Product.Components.Misc[miscSlot] = c[0];
					clientState.Decisions.Products[newProductID].Product.Components.Misc[miscSlot] = c[0];
					onProductChange();
					currentComponentSnippet = null;
				}}
				onmouseover={() => {
					currentProductHover.Components.Misc[miscSlot] = c[0];
				}}
				onmouseout={resetHoverProduct}
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

{#snippet productMisc0()}
	{@render productMisc(0)}
{/snippet}

{#snippet productMisc1()}
	{@render productMisc(1)}
{/snippet}

{#snippet productMisc2()}
	{@render productMisc(2)}
{/snippet}

{#snippet productMisc3()}
	{@render productMisc(3)}
{/snippet}

{#snippet showHoverDifference(
	/** @type {number}*/ mainValue,
	/** @type {number}*/ hoverValue,
	/** @type {(value: number)=>string}*/ format,
	/** @type {boolean?}*/ invert
)}
	{#if mainValue != hoverValue}
		<span
			class={invert
				? hoverValue < mainValue
					? 'green'
					: 'red'
				: hoverValue < mainValue
					? 'red'
					: 'green'}>{format(hoverValue - mainValue)}</span
		>
	{/if}
{/snippet}

<style>
	img.inlineIcon {
		height: 1.5rem;
	}

	span.red {
		color: red;
	}
	span.green {
		color: green;
	}

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
		border-radius: 0;
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

		& > button,
		& > article {
			height: 10rem;
		}
	}

	.marketingIcon {
		height: 1.5rem;
		translate: 0 -0.5rem;
		background-color: transparent;
		background-image: url($lib/images/marketing.svg);
		background-color: transparent;
		background-position: center;
		background-size: 100% 100%;
		border: none;

		&:active {
			border: none !important;
			outline: none !important;
		}
	}
</style>
