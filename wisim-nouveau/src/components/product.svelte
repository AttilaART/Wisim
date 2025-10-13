<script>
	import { format } from '$lib/javascript/format';
	import { calculateProductStats } from '../calculateProduct';
	import storageIcon from '$lib/images/warehouse.svg';
	import Increment from './increment.svelte';
	import Window from './window.svelte';
	import ProductionIcon from '$lib/images/production.svg';
	import MarketingIcon from '$lib/images/marketing.svg';
	import { preventPageReload } from '$lib/helper.svelte';
	import ProductDesigner from './productDesigner.svelte';
	/** @typedef {Object} Props
	 * @property {import("$lib/javascript/simulation").clientState} clientState,
	 * @property {(Decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions,
	 * @property {(contents: import("svelte").Snippet<[number]>)=>number} newWindow,
	 * @property {(windowId: number)=>void} deleteWindow,
	 */

	/**
	 * @type {Props}
	 */
	let { clientState = $bindable(), updateDecisions, newWindow, deleteWindow } = $props();

	const images = import.meta.glob(['$lib/images/Products/Base_blueprint/*.svg'], {
		eager: true,
		as: 'url'
	});
</script>

<div class="productsGrid">
	<button
		onclick={() => {
			newWindow(newProduct);
		}}
	>
		<div style="line-height: 3.5rem;">
			<span style="font-size: 3rem; font-weight: bolder; vertical-align: middle;">+</span>
			<span style="margin: 0px; vertical-align: middle;"> New Product </span>
		</div>
	</button>
	{#if clientState.Company.Offers}
		{#each Object.entries(clientState.Company.Offers) as offer}
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
					<button class="inlineIcon marketingIcon" aria-label="Marketing"></button>
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
<!--
{#snippet configurePromotion(/** @type {number}*/ windowID)}
	<Window
		title={`Configure ${clientState.Company.Offers[productConfugureWindows[windowID]].Product.Name}`}
		closeWindow={() => {
			deleteWindow(windowID);
			delete productConfugureWindows[windowID];
		}}
	>
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
-->

{#snippet newProduct(/** @type {number}*/ windowID)}
	<Window
		title="Product Designer"
		closeWindow={() => {
			deleteWindow(windowID);
		}}
	>
		<ProductDesigner
			{clientState}
			{updateDecisions}
			closeWindow={() => {
				deleteWindow(windowID);
			}}
			ExistingProduct={null}
		></ProductDesigner>
	</Window>
{/snippet}

<style>
	img.inlineIcon {
		height: 1.5rem;
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
