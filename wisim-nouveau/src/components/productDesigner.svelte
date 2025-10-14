<script>
	import { format } from '$lib/javascript/format';
	import { calculateProductStats } from '../calculateProduct';
	import Increment from './increment.svelte';

	/** @typedef {Object} Props
	 * @property {import("$lib/javascript/simulation").clientState} clientState,
	 * @property {(Decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions,
	 * @property {()=>void} closeWindow,
	 * @property {import("$lib/javascript/simulation").Product?} existingProduct
	 */

	/** @type {Props} */

	let { clientState = $bindable(), updateDecisions, closeWindow, existingProduct } = $props();

	/** @type {import("$lib/javascript/simulation").Decisions_product} */
	let productDecisions = $state({
		Price: 150,
		Name: 'Unnamed Product',
		Promotion: {
			Quantity: 10000,
			Quality: 0.2,
			Durability: 0.2,
			Price: 0.2,
			Ecology: 0.2,
			Ethics: 0.2
		},
		Product: {
			ID: `${Math.trunc(Math.random() * 100000000)}`,
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
		}
	});

	/** @type {number} */
	/** @type {{productStats: import("$lib/javascript/simulation").ProductStats , productionLineCost: number}} */
	let { productStats, productionLineCost } = $derived(
		calculateProductStats(productDecisions.Product, clientState.productComponents)
	);

	/** @type {import("$lib/javascript/simulation").Decisions_product} */
	let hoverProductDecisions = $state(JSON.parse(JSON.stringify(productDecisions)));

	$effect(() => {
		hoverProductDecisions = JSON.parse(JSON.stringify(productDecisions));
	});

	/** @type {{productStats: import("$lib/javascript/simulation").ProductStats , productionLineCost: number}} */
	let { productStats: hoverProductStats, productionLineCost: hoverProductionLineCost } = $derived(
		calculateProductStats(hoverProductDecisions.Product, clientState.productComponents)
	);

	/** @type {import("svelte").Snippet} */
	let currentDesignerSnippet = $state(selectPart);

	const images = import.meta.glob(['$lib/images/Products/Base_blueprint/*.svg'], {
		eager: true,
		as: 'url'
	});

	let mousePosition = $state({ x: 0, y: 0 });

	window.addEventListener('mousemove', (event) => {
		mousePosition.x = event.clientX;
		mousePosition.y = event.clientY;
	});

	/** @param {HTMLElement} el*/
	function followMouse(el) {
		$effect(() => {
			el.setAttribute(
				'style',
				`
        position: absolute;
        top: 0;
        left: 0;
        z-index: 99;
        translate: ${mousePosition.x}px ${mousePosition.y}px;`
			);
		});
	}

	if (existingProduct != null) {
		let id = productDecisions.Product.ID;
		productDecisions.Product = JSON.parse(JSON.stringify(existingProduct));
		productDecisions.Product.ID = id;
	}
</script>

<div style="min-width: 50rem;">
	<div class="main-grid">
		<div>
			<input type="text" autocomplete="off" bind:value={productDecisions.Product.Name} />

			{@render currentDesignerSnippet()}

			<div class="grid">
				<center>
					<Increment
						label="Durability"
						bind:value={productDecisions.Product.ExtraDurability}
						min={0}
						max={5}
						onclick={() => {}}
					></Increment>
				</center>
				<center>
					<Increment
						label="Quality"
						bind:value={productDecisions.Product.ExtraQuality}
						min={0}
						max={5}
						onclick={() => {}}
					></Increment>
				</center>
			</div>
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
							>{format.number(productStats.ProductionCost, false, 1)}
							{@render showHoverDifference(
								productStats.ProductionCost,
								hoverProductStats.ProductionCost,
								(value) => {
									return format.number(value, true, 1);
								},
								true
							)}
						</td>
						<td>Quality: </td>
						<td
							>{format.number(productStats.Quality, false, 1)}
							{@render showHoverDifference(
								productStats.Quality,
								hoverProductStats.Quality,
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
							>{format.number(productStats.MaterialUse, false, 1)}
							{@render showHoverDifference(
								productStats.MaterialUse,
								hoverProductStats.MaterialUse,
								(value) => {
									return format.number(value, true, 1);
								},
								false
							)}
						</td>
						<td>Ecology: </td>
						<td
							>{format.number(productStats.Ecology, false, 1)}
							{@render showHoverDifference(
								productStats.Ecology,
								hoverProductStats.Ecology,
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
								productStats.MaterialUse * clientState.ExternalFactors.MaterialPrice,
								false,
								1
							)}
							{@render showHoverDifference(
								productStats.MaterialUse * clientState.ExternalFactors.MaterialPrice,
								hoverProductStats.MaterialUse * clientState.ExternalFactors.MaterialPrice,
								(value) => {
									return format.number(value, true, 1);
								},
								false
							)}
						</td>
						<td>Ethics: </td>
						<td
							>{format.number(productStats.Ethics, false, 1)}
							{@render showHoverDifference(
								productStats.Ethics,
								hoverProductStats.Ethics,
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
							>{productStats.Durability}
							{@render showHoverDifference(
								productStats.Durability,
								hoverProductStats.Durability,
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
			<fieldset role="group">
				<input type="number" placeholder="price" bind:value={productDecisions.Price} />
				<input type="text" disabled value="CHF" style="width: 5rem;" />
			</fieldset>
		</div>
	</div>
	<footer class="grid">
		<button onclick={closeWindow}>Cancel</button>
		<button
			onclick={() => {
				productDecisions.Name = productDecisions.Product.Name;
				clientState.Decisions.Products[productDecisions.Product.ID] = productDecisions;
				clientState.Company.Offers[productDecisions.Product.ID] = {
					Price: productDecisions.Price,
					productStats: productStats,
					PromotionQuality: 0,
					Promotion: {
						Quantity: productDecisions.Promotion.Quantity,
						StyleQuality: productDecisions.Promotion.Quality,
						StyleEcology: productDecisions.Promotion.Ecology,
						StyleEthics: productDecisions.Promotion.Ethics,
						StyleDurability: productDecisions.Promotion.Durability
					},
					Product: productDecisions.Product
				};
				updateDecisions(clientState.Decisions);
				closeWindow();
			}}
			disabled={(() => {
				if (!productDecisions.Product.Components.FormFactor) return true;
				if (!productDecisions.Product.Components.Body) return true;
				if (!productDecisions.Product.Components.Mechanism) return true;
				if (!productDecisions.Product.Components.Frame) return true;
				return false;
			})()}>Confirm</button
		>
	</footer>
</div>

{#snippet selectPart()}
	<div class="parts-grid">
		<button
			onclick={() => {
				currentDesignerSnippet = formFactor;
			}}
		>
			{@render showImageOrPlus(
				clientState.productComponents.FormFactor,
				productDecisions.Product.Components.FormFactor
			)}
		</button>
		<button
			onclick={() => {
				currentDesignerSnippet = frame;
			}}
		>
			{@render showImageOrPlus(
				clientState.productComponents.Frame,
				productDecisions.Product.Components.Frame
			)}
		</button>
		<button
			onclick={() => {
				currentDesignerSnippet = body;
			}}
		>
			{@render showImageOrPlus(
				clientState.productComponents.Body,
				productDecisions.Product.Components.Body
			)}
		</button>
		<button
			onclick={() => {
				currentDesignerSnippet = mechanism;
			}}
		>
			{@render showImageOrPlus(
				clientState.productComponents.Mechanism,
				productDecisions.Product.Components.Mechanism
			)}
		</button>

		<div>
			<img
				src={'/src/lib/images/' +
					clientState.productComponents.FormFactor[
						`${productDecisions.Product.Components.FormFactor}`
					]?.Image}
				alt=""
			/>
		</div>

		<button
			onclick={() => {
				currentDesignerSnippet = misc1;
			}}
			disabled={productStats.MiscSlots < 1}
		>
			{@render showImageOrPlus(
				clientState.productComponents.Misc,
				productDecisions.Product.Components.Misc[0]
			)}
		</button>
		<button
			onclick={() => {
				currentDesignerSnippet = misc2;
			}}
			disabled={productStats.MiscSlots < 2}
		>
			{@render showImageOrPlus(
				clientState.productComponents.Misc,
				productDecisions.Product.Components.Misc[1]
			)}
		</button>
		<button
			onclick={() => {
				currentDesignerSnippet = misc3;
			}}
			disabled={productStats.MiscSlots < 3}
		>
			{@render showImageOrPlus(
				clientState.productComponents.Misc,
				productDecisions.Product.Components.Misc[2]
			)}
		</button>
		<button
			onclick={() => {
				currentDesignerSnippet = misc4;
			}}
			disabled={productStats.MiscSlots < 4}
		>
			{@render showImageOrPlus(
				clientState.productComponents.Misc,
				productDecisions.Product.Components.Misc[3]
			)}
		</button>
	</div>
{/snippet}

{#snippet formFactor()}
	<div class="grid">
		{#each Object.entries(clientState.productComponents.FormFactor) as c}
			{@render renderComponent('FormFactor', c)}
		{/each}
	</div>
{/snippet}

{#snippet frame()}
	<div class="grid">
		{#each Object.entries(clientState.productComponents.Frame) as c}
			{@render renderComponent('Frame', c)}
		{/each}
	</div>
{/snippet}

{#snippet body()}
	<div class="grid">
		{#each Object.entries(clientState.productComponents.Body) as c}
			{@render renderComponent('Body', c)}
		{/each}
	</div>
{/snippet}

{#snippet mechanism()}
	<div class="grid">
		{#each Object.entries(clientState.productComponents.Mechanism) as c}
			{@render renderComponent('Mechanism', c)}
		{/each}
	</div>
{/snippet}

{#snippet misc1()}
	{@render misc(0)}
{/snippet}
{#snippet misc2()}
	{@render misc(1)}
{/snippet}
{#snippet misc3()}
	{@render misc(2)}
{/snippet}
{#snippet misc4()}
	{@render misc(3)}
{/snippet}

{#snippet misc(/** @type {number} */ slot)}
	<div class="grid">
		{#each Object.entries(clientState.productComponents.Misc) as c}
			<button
				class="component-button"
				onclick={() => {
					productDecisions.Product.Components.Misc[slot] = c[0];
					currentDesignerSnippet = selectPart;
				}}
				onmouseenter={(_) => {
					hoverProductDecisions.Product.Components.Misc[slot] = c[0];
				}}
				onmouseleave={(_) => {
					hoverProductDecisions = JSON.parse(JSON.stringify(productDecisions));
				}}
			>
				<img src={'/src/lib/images/' + c[1].Image} alt="" style="mix-blend-mode: lighten;" />
				{@render componentTooltip(c)}
			</button>
		{/each}
	</div>
{/snippet}

{#snippet renderComponent(
	/** @type {string}*/ part,
	/** @type {[string, import("$lib/javascript/simulation").Component]}*/ c
)}
	<button
		class="component-button"
		onclick={() => {
			productDecisions.Product.Components[part] = c[0];
			currentDesignerSnippet = selectPart;
		}}
		onmouseenter={(_) => {
			hoverProductDecisions.Product.Components[part] = c[0];
		}}
		onmouseleave={(_) => {
			hoverProductDecisions = JSON.parse(JSON.stringify(productDecisions));
		}}
	>
		<img src={'/src/lib/images/' + c[1].Image} alt="" style="mix-blend-mode: lighten;" />
		{@render componentTooltip(c)}
	</button>
{/snippet}

{#snippet componentTooltip(
	/** @type {[string, import("$lib/javascript/simulation").Component]}*/ c
)}
	<article use:followMouse class="tooltip">
		<h4>{c[1].Name}</h4>
		<ul>
			{#each [{ value: c[1].MiscSlots, name: 'Miscelaneous slots' }, { value: c[1].ProductionCost, name: 'ProductionCost', invert: true }, { value: c[1].MaterialUse, name: 'MaterialUse', invert: true }, { value: c[1].Ecology, name: 'Ecology' }, { value: c[1].Ethics, name: 'Ethics' }, { value: c[1].Quality, name: 'Quality' }, { value: c[1].Durability, name: 'Durability' }, { value: c[1].ProductionLineCost, name: 'ProductionLineCost', invert: true }] as p}
				{#if p.value != 0}
					<li>
						{@render colorProperty(p.value, p.invert)}
						{p.name}
					</li>
				{/if}
			{/each}
		</ul>
	</article>
{/snippet}

{#snippet colorProperty(/** @type {number} */ value, /** @type {boolean | undefined} */ invert)}
	<span class={value * (invert ? -1 : 1) > 0 ? 'green' : 'red'}
		>{format.number(value, true, 0)}</span
	>
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

{#snippet showImageOrPlus(
	/** @type {Object.<string, import("$lib/javascript/simulation").Component>}*/ imageSource,
	/** @type {string?}*/ field
)}
	{#if !field}
		+
	{:else}
		<img src={'/src/lib/images/' + `${imageSource[field]?.Image}`} alt="" />
	{/if}
{/snippet}

<style>
	span.red {
		color: red;
	}
	span.green {
		color: green;
	}

	.tooltip {
		display: none;
		text-align: left;
	}

	@keyframes showTooltip {
		start {
			opacity: 0%;
		}
		end {
			opacity: 100%;
		}
	}

	*:hover > .tooltip {
		display: unset;
		animation: showTooltip 1s 1s;
	}

	.main-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}

	.component-button,
	.parts-grid > button {
		background-color: transparent;
		border-radius: 0;
		border-color: silver;

		box-shadow: inset 0px 0px 20px black;

		transition: box-shadow 0.5s;

		&:hover {
			box-shadow: inset 0px 0px 10px color-mix(in oklab, gray, transparent 30%);
		}

		img {
			content-fit: contain;
			height: 2rem;
		}
	}

	.parts-grid {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;

		button {
			font-size: 2rem;
			padding: 0 1rem;
		}

		div {
			grid-column: 1 / span 4;
			aspect-ratio: 16 / 9;
			border-right: solid 1px silver;
			border-left: solid 1px silver;

			box-shadow: inset 0px 0px 20px black;
			background-color: #bdbdbd;

			img {
				padding: 1rem;
				height: 100%;
				width: 100%;

				object-fit: contain;
				mix-blend-mode: lighten;
			}
		}
	}
</style>
