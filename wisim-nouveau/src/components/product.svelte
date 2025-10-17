<script>
	import { format } from '$lib/javascript/format';
	import { calculateProductStats } from '../calculateProduct';
	import storageIcon from '$lib/images/warehouse.svg';
	import Increment from './increment.svelte';
	import Window from './window.svelte';
	import ProductionIcon from '$lib/images/production.svg';
	import MarketingIcon from '$lib/images/marketing.svg';
	import { ignoreError, preventPageReload } from '$lib/helper.svelte';
	import ProductDesigner from './productDesigner.svelte';
	import { createRawSnippet, mount, unmount } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import ConfigurePromotion from './configurePromotion.svelte';
	/** @typedef {Object} Props
	 * @property {import("$lib/javascript/simulation").clientState} clientState,
	 * @property {(Decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions,
	 * @property {(contents: import("svelte").Snippet<[number]>)=>number} newWindow,
	 * @property {(windowId: number)=>void} deleteWindow,
	 * @property {()=>void} openProduction
	 */

	/**
	 * @type {Props}
	 */
	let {
		clientState = $bindable(),
		updateDecisions,
		newWindow,
		deleteWindow,
		openProduction
	} = $props();
	/**
	 * @param {string} productID
	 */
	function calculateMonthlyProduction(productID) {
		let totalProduction = 0;
		for (let m of clientState.Company.Machines) {
			if (m.AssignedProductID == productID) {
				totalProduction += m.ProductionCapacity;
			}
		}

		return totalProduction / clientState.Company.Offers[productID].ProductStats.ProductionCost;
	}

	/**
	 * @param {string} productID
	 */
	function addToDecisionsIfNotPresent(productID) {
		if (clientState.Decisions.Products[productID] == undefined) {
			let offer = clientState.Company.Offers[productID];

			clientState.Decisions.Products[productID] = {
				Name: offer.Product.Name,
				Outdated: offer.Outdated,
				Product: offer.Product,
				Price: offer.Price,
				Promotion: {
					Quantity: offer.Promotion.Quantity,
					Quality: offer.Promotion.StyleQuality,
					Ecology: offer.Promotion.StyleEcology,
					Ethics: offer.Promotion.StyleEthics,
					Durability: offer.Promotion.StyleDurability,
					Price: offer.Promotion.StylePrice
				}
			};
		}
	}

	let showOutdated = $state(false);

	for (let productID of Object.keys(clientState.Company.Offers)) {
		addToDecisionsIfNotPresent(productID);
	}
	updateDecisions(clientState.Decisions);
</script>

<label for="">
	Show Outdated
	<input type="checkbox" bind:checked={showOutdated} />
</label>

<div class="products-grid">
	{#snippet newProductNoExistingProduct(/** @type {number}*/ windowID)}
		{@render newProduct(windowID, null, false)}
	{/snippet}
	<button
		onclick={() => {
			newWindow(newProductNoExistingProduct);
		}}
	>
		<div style="line-height: 3.5rem;">
			<span style="font-size: 3rem; font-weight: bolder; vertical-align: middle;">+</span>
			<span style="margin: 0px; vertical-align: middle;"> New Product </span>
		</div>
	</button>
	{#if clientState.Company.Offers}
		{#each Object.entries(clientState.Company.Offers) as offer (offer[0])}
			{#if !clientState.Decisions.Products[offer[0]].Outdated || showOutdated}
				{#snippet viewProduct(/** @type {number}*/ windowID)}
					{@render newProduct(windowID, offer[1].Product, true)}
				{/snippet}
				<article
					in:fly={{ y: -100 }}
					onclick={() => {
						newWindow(viewProduct);
					}}
				>
					<img
						src={'/src/lib/images/' +
							ignoreError(() => {
								return clientState.productComponents.FormFactor[
									`${offer[1].Product.Components.FormFactor}`
								]
									? clientState.productComponents.FormFactor[
											`${offer[1].Product.Components.FormFactor}`
										].Image
									: '';
							})}
						alt=""
						style="position: absolute; pointer-events: none; left: 50%; top: 50%; transform: translate(-50%, calc(-50% + 0.5rem)); height: 6rem; mix-blend-mode: lighten;"
					/>
					<div>
						{#snippet configurePromotionOfProduct(/** @type {number}*/ windowID)}
							{@render configurePromotion(windowID, offer[0])}
						{/snippet}
						<button
							onclick={(e) => {
								e.stopPropagation();
								newWindow(configurePromotionOfProduct);
							}}
							class="inlineIcon marketingIcon"
							aria-label="Marketing"
						></button>
						<h4 style="display: inline;">
							{offer[1].Product.Name}
						</h4>
					</div>
					<div style="text-align: right;">
						<h4>{format.currency(offer[1].Price, false, 0)}</h4>
						<!--<div>
						{clientState.Company.ProductsInStorage[offer[0]]
							? format.number(clientState.Company.ProductsInStorage[offer[0]], false, 0)
							: '0'}
						<img
							class="inlineIcon"
							style="height: 1.2rem; translate: 0 -0.1rem ;"
							src={storageIcon}
							alt=""
						/>
					</div>-->
						<div>
							<img
								class="inlineIcon"
								style="height: 1.2rem; translate: 0 -0.1rem ;"
								src={ProductionIcon}
								alt=""
							/>
							<strong>
								{format.number(offer[1].ProductStats.ProductionCost, false, 1)}
							</strong>
						</div>
						<small
							>{format.number(
								calculateMonthlyProduction(offer[1].Product.ID),
								false,
								1
							)}/month</small
						>
					</div>
					<label
						onclick={(e) => {
							e.stopPropagation();
							updateDecisions(clientState.Decisions);
						}}
						for="outdated{offer[1].Product.ID}"
						style="position: absolute; bottom: 0.5rem; left: 0.5rem;"
					>
						<input
							id="outdated{offer[1].Product.ID}"
							bind:checked={clientState.Decisions.Products[offer[0]].Outdated}
							type="checkbox"
						/>
						Mark as outdated
					</label>

					{#snippet newProductWithExistingProduct(/** @type {number}*/ windowID)}
						{@render newProduct(windowID, offer[1].Product, false)}
					{/snippet}
					<button
						onclick={(e) => {
							e.stopPropagation();
							newWindow(newProductWithExistingProduct);
						}}
						class="contrast outline"
						style="position: absolute; bottom: 0.5rem; right: 0.5rem; padding: 0.5rem; line-height: 1rem;"
					>
						Make a copy
					</button>
				</article>
			{/if}
		{/each}
	{/if}
</div>

{#snippet configurePromotion(/** @type {number}*/ windowID, /** @type {string}*/ productID)}
	<Window
		title={'Configure Promotion'}
		closeWindow={() => {
			deleteWindow(windowID);
		}}
	>
		<ConfigurePromotion bind:clientState {updateDecisions} {productID}></ConfigurePromotion>
	</Window>
{/snippet}

{#snippet newProduct(
	/** @type {number}*/ windowID,
	/** @type {import("$lib/javascript/simulation").Product?}*/ existingProduct,
	/** @type {boolean}*/ viewOnly
)}
	<Window
		title="Product Designer"
		closeWindow={() => {
			deleteWindow(windowID);
		}}
	>
		<ProductDesigner
			bind:clientState
			{updateDecisions}
			closeWindow={() => {
				deleteWindow(windowID);
			}}
			{existingProduct}
			{viewOnly}
			{openProduction}
		></ProductDesigner>
	</Window>
{/snippet}

<style>
	@property --top-color {
		syntax: '<color>';
		initial-value: black;
		inherits: false;
	}

	.products-grid {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		min-width: 30rem;

		& > article,
		& > button {
			min-height: 9rem;
		}

		article {
			position: relative;
			display: grid;
			grid-template-columns: auto 30%;
			grid-template-rows: auto auto auto;
			gap: 0.5rem;

			box-shadow: inset 0px 0px 40px black;
			border: solid 1px silver;
			padding: 0.5rem;
			margin: 0;

			transition:
				--top-color 0.5s,
				box-shadow 0.25s;

			background-image: linear-gradient(to bottom, var(--top-color), transparent);

			&:hover {
				--top-color: gray;
			}

			&:active {
				box-shadow: inset 0px 0px 20px black;
			}
		}
	}
	img.inlineIcon,
	button.inlineIcon {
		min-height: 0 !important;
		height: 1.5rem !important;
	}

	.marketingIcon {
		min-height: 0 !important;
		height: 1.5rem !important;
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
