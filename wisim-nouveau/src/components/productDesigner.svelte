<script>
	import ProductionIcon from '$lib/images/production.svg';
	import { format } from '$lib/javascript/format';
	import noIcon from '$lib/images/noIcon.svg';
	import Increment from './increment.svelte';
	import ColorProperty from './colorProperty.svelte';

	/** @typedef {Object} Props
	 * @property {import("$lib/javascript/simulation").clientState} clientState,
	 * @property {(Decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions,
	 * @property {()=>void} closeWindow,
	 * @property {import("$lib/javascript/simulation").Product?} existingProduct
	 * @property {boolean} viewOnly
	 * @property {()=>void} openProduction
	 */

	/** @type {Props} */

	let {
		clientState = $bindable(),
		updateDecisions,
		closeWindow,
		existingProduct,
		viewOnly,
		openProduction
	} = $props();

	/** @type {HTMLDialogElement} */
	let machinesDialogue;

	/** @type {import("$lib/javascript/simulation").Decisions_product} */
	let productDecisions = $state({
		Price: 150,
		Name: '',
		Outdated: false,
		Promotion: {
			Quantity: 1000,
			Quality: 0.4472135955,
			Durability: 0.4472135955,
			Price: 0.4472135955,
			Ecology: 0.4472135955,
			Ethics: 0.4472135955
		},
		Product: {
			ID: `${Math.trunc(Math.random() * 100000000)}`,
			CompanyID: clientState.Company.ID,
			Name: '',

			Components: {
				FormFactor: 'FormFactorMedium',
				Frame: null,
				Body: null,
				Mechanism: null,
				Misc: []
			},

			TechLevels: clientState.Company.Tech,

			MaterialQuality: 0,
			ExtraDurability: 0,
			ExtraQuality: 0
		}
	});

	/** @type {number} */
	/** @type {{ProductStats: import("$lib/javascript/simulation").ProductStats , ProductionLineCost: number}} */
	let { ProductStats: productStats, ProductionLineCost: productionLineCost } = $derived.by(() => {
		// @ts-ignore
		return JSON.parse(
			CalculateProductStatsGo(
				JSON.stringify(productDecisions.Product),
				JSON.stringify(clientState.productComponents)
			)
		);
	});

	/** @type {import("$lib/javascript/simulation").Decisions_product} */
	let hoverProductDecisions = $state(JSON.parse(JSON.stringify(productDecisions)));

	$effect(() => {
		hoverProductDecisions = JSON.parse(JSON.stringify(productDecisions));
	});

	/** @type {{ProductStats: import("$lib/javascript/simulation").ProductStats , ProductionLineCost: number}} */
	let { ProductStats: hoverProductStats, ProductionLineCost: hoverProductionLineCost } =
		$derived.by(() => {
			// @ts-ignore
			return JSON.parse(
				CalculateProductStatsGo(
					JSON.stringify(hoverProductDecisions.Product),
					JSON.stringify(clientState.productComponents)
				)
			);
		});

	/** @type {import("svelte").Snippet} */
	let currentDesignerSnippet = $state(selectPart);

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
		productDecisions.Price = clientState.Company.Offers[existingProduct.ID].Price;
		productDecisions.Product = JSON.parse(JSON.stringify(existingProduct));
		if (!viewOnly) {
			productDecisions.Product.ID = id;
			productDecisions.Product.TechLevels = clientState.Company.Tech;
		}

		let promotion = clientState.Company.Offers[existingProduct.ID].Promotion;
		// console.log(promotion);

		productDecisions.Promotion.Quantity = promotion.Quantity;

		productDecisions.Promotion.Ecology = promotion.Ecology;
		productDecisions.Promotion.Quality = promotion.Quality;
		productDecisions.Promotion.Durability = promotion.Durability;
		productDecisions.Promotion.Ethics = promotion.Ethics;
		productDecisions.Promotion.Price = promotion.Price;
	}

	/**
	 * @param {HTMLElement} el
	 */
	function focusOnMount(el) {
		el.focus();
	}
</script>

<div>
	<div class="main-grid">
		<div>
			<input
				use:focusOnMount
				type="text"
				autocomplete="off"
				placeholder="Product Name"
				bind:value={productDecisions.Product.Name}
			/>

			{@render currentDesignerSnippet()}

			<div class="grid">
				<center>
					<Increment
						label="Durability"
						bind:value={productDecisions.Product.ExtraDurability}
						min={0}
						max={5}
						onclick={() => {}}
						disabled={viewOnly}
					></Increment>
				</center>
				<center>
					<Increment
						label="Quality"
						bind:value={productDecisions.Product.ExtraQuality}
						min={0}
						max={5}
						onclick={() => {}}
						disabled={viewOnly}
					></Increment>
				</center>
			</div>
		</div>
		<div style="width: 100%;">
			<table style="width: 100%;">
				<thead>
					<tr>
						<th colspan="2"> <h3>Manufacturing Stats</h3> </th>
						<th colspan="2"> <h3>Product Stats</h3> </th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td> Production Cost:</td>
						<td>
							<img
								class="inlineIcon"
								style="height: 1.2rem; translate: 0 0.1rem ;"
								src={ProductionIcon}
								alt=""
							/>
							{format.number(productStats.ProductionCost, false, 1)}
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
							>{format.currency(
								productStats.MaterialUse * clientState.ExternalFactors.MaterialPrice,
								false,
								2
							)}
							{@render showHoverDifference(
								productStats.MaterialUse * clientState.ExternalFactors.MaterialPrice,
								hoverProductStats.MaterialUse * clientState.ExternalFactors.MaterialPrice,
								(value) => {
									return format.currency(value, true, 2);
								},
								true
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
						<td data-tooltip="The cost to develop the product and production line tooling"
							>Development cost:</td
						>
						<td
							>{format.currency(hoverProductionLineCost, false, 0)}
							{@render showHoverDifference(
								productionLineCost,
								hoverProductionLineCost,
								(value) => {
									return format.currency(value, true, 0);
								},
								false
							)}
						</td>

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
				</tbody>
			</table>
			<!-- svelte-ignore a11y_no_redundant_roles -->
			<fieldset role="group">
				<input type="number" placeholder="price" bind:value={productDecisions.Price} />
				<input type="text" disabled value="CHF" style="width: 5rem;" />
			</fieldset>
		</div>
	</div>
	<footer class="grid">
		<button class="secondary" onclick={closeWindow}>Cancel</button>
		<button
			onclick={() => {
				if (productDecisions.Product.Name == '') {
					productDecisions.Product.Name = 'Unnamed Product';
				}
				productDecisions.Name = productDecisions.Product.Name;
				clientState.Decisions.Products[productDecisions.Product.ID] = productDecisions;
				clientState.Company.Offers[productDecisions.Product.ID] = {
					Price: productDecisions.Price,
					ProductStats: productStats,
					PromotionQuality: 0,
					Outdated: productDecisions.Outdated,
					Promotion: {
						Quantity: productDecisions.Promotion.Quantity,
						Quality: productDecisions.Promotion.Quality,
						Ecology: productDecisions.Promotion.Ecology,
						Ethics: productDecisions.Promotion.Ethics,
						Durability: productDecisions.Promotion.Durability,
						Price: productDecisions.Promotion.Price
					},
					Product: productDecisions.Product
				};

				updateDecisions(clientState.Decisions);
				if (viewOnly) {
					closeWindow();
				} else {
					clientState.Company.Balance -= productionLineCost;
					machinesDialogue.show();
				}
			}}
			disabled={(() => {
				if (!productDecisions.Product.Components.FormFactor) return true;
				if (!productDecisions.Product.Components.Body) return true;
				if (!productDecisions.Product.Components.Mechanism) return true;
				if (!productDecisions.Product.Components.Frame) return true;
				return false;
			})()}
			>Confirm <span data-tooltip="The cost to develop the product and production line tooling"
				>{!viewOnly ? format.currency(-hoverProductionLineCost, true, 2) : ''}</span
			></button
		>
		<dialog bind:this={machinesDialogue}>
			<article>
				<p>Would you like to also assign some production machines machines to your product?</p>
				<footer>
					<button
						onclick={() => {
							machinesDialogue.close();
							closeWindow();
						}}
						class="secondary outline">Later</button
					>
					<button
						onclick={() => {
							machinesDialogue.close();
							closeWindow();
							openProduction();
						}}>Yes</button
					>
				</footer>
			</article>
		</dialog>
	</footer>
</div>

{#snippet selectPart()}
	<div class="parts-grid">
		<button
			onclick={() => {
				currentDesignerSnippet = formFactor;
			}}
			disabled={viewOnly}
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
			disabled={viewOnly}
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
			disabled={viewOnly}
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
			disabled={viewOnly}
		>
			{@render showImageOrPlus(
				clientState.productComponents.Mechanism,
				productDecisions.Product.Components.Mechanism
			)}
		</button>

		<div>
			<img
				src={'/' +
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
			disabled={productStats.MiscSlots < 1 || viewOnly}
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
			disabled={productStats.MiscSlots < 2 || viewOnly}
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
			disabled={productStats.MiscSlots < 3 || viewOnly}
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
			disabled={productStats.MiscSlots < 4 || viewOnly}
		>
			{@render showImageOrPlus(
				clientState.productComponents.Misc,
				productDecisions.Product.Components.Misc[3]
			)}
		</button>
	</div>
{/snippet}

{#snippet formFactor()}
	<div class="component-grid">
		<center>
			<h3>Form Factor</h3>
		</center>
		<div>
			{#each Object.entries(clientState.productComponents.FormFactor) as c}
				{@render renderComponent('FormFactor', c)}
			{/each}
		</div>
	</div>
{/snippet}

{#snippet frame()}
	<div class="component-grid">
		<center>
			<h3>Frame</h3>
		</center>
		<div>
			{#each Object.entries(clientState.productComponents.Frame) as c}
				{@render renderComponent('Frame', c)}
			{/each}
		</div>
	</div>
{/snippet}

{#snippet body()}
	<div class="component-grid">
		<center>
			<h3>Body</h3>
		</center>
		<div>
			{#each Object.entries(clientState.productComponents.Body) as c}
				{@render renderComponent('Body', c)}
			{/each}
		</div>
	</div>
{/snippet}

{#snippet mechanism()}
	<div class="component-grid">
		<center>
			<h3>Coffee Mechanism</h3>
		</center>
		<div>
			{#each Object.entries(clientState.productComponents.Mechanism) as c}
				{@render renderComponent('Mechanism', c)}
			{/each}
		</div>
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
	<div class="component-grid">
		<center>
			<h3>Miscelaneous</h3>
		</center>
		<div>
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
					{#if c[1].Image != ''}
						<img src={'/' + c[1].Image} alt="" style="mix-blend-mode: lighten;" />
					{:else}
						<small>{c[1].Name}</small>
					{/if}
					{@render componentTooltip(c)}
				</button>
			{/each}

			<button
				class="component-button"
				onclick={() => {
					productDecisions.Product.Components.Misc[slot] = null;
					currentDesignerSnippet = selectPart;
				}}
				onmouseenter={(_) => {
					hoverProductDecisions.Product.Components.Misc[slot] = null;
				}}
				onmouseleave={(_) => {
					hoverProductDecisions = JSON.parse(JSON.stringify(productDecisions));
				}}
			>
				<img src={noIcon} alt="" style="mix-blend-mode: lighten;" />
			</button>
		</div>
	</div>
{/snippet}

{#snippet renderComponent(
	/** @type {string}*/ part,
	/** @type {[string, import("$lib/javascript/simulation").Component]}*/ c
)}
	<button
		class="component-button"
		onclick={() => {
			// @ts-ignore
			productDecisions.Product.Components[part] = c[0];
			currentDesignerSnippet = selectPart;
		}}
		onmouseenter={(_) => {
			// @ts-ignore
			hoverProductDecisions.Product.Components[part] = c[0];
		}}
		onmouseleave={(_) => {
			hoverProductDecisions = JSON.parse(JSON.stringify(productDecisions));
		}}
	>
		{#if c[1].Image != ''}
			<img src={'/' + c[1].Image} alt="" style="mix-blend-mode: lighten;" />
		{:else}
			<small>{c[1].Name}</small>
		{/if}
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
						<ColorProperty
							value={p.value}
							invert={p.invert}
							formatter={(v) => format.number(v, true, 0)}
						></ColorProperty>
						{p.name}
					</li>
				{/if}
			{/each}
		</ul>
	</article>
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
		<span style="font-size: 2rem;">+</span>
	{:else if imageSource[field]?.Image != ''}
		<img src={'/' + `${imageSource[field]?.Image}`} alt="" style="mix-blend-mode: lighten;" />
	{:else}
		<small>{imageSource[field]?.Name}</small>
	{/if}
{/snippet}

<style>
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
		poimter-events: none;
		animation: showTooltip 1s 1s;
	}

	.main-grid {
		width: 60rem;
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}

	.parts-grid {
		height: 22.5rem;
	}

	.component-button,
	.parts-grid > button {
		margin: 0;

		background-color: transparent;
		border-radius: 0;
		border: 1px solid silver;

		box-shadow: inset 0px 0px 20px black;

		transition: box-shadow 0.5s;

		&:hover,
		&:focus {
			box-shadow: inset 0px 0px 10px color-mix(in oklab, gray, transparent 30%);
		}

		img {
			content-fit: contain;
			padding: 5px;
			height: 3rem;
		}
	}

	.parts-grid {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		margin-bottom: 1.5rem;

		button {
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
				padding: 0.5rem;
				height: 100%;
				width: 100%;

				object-fit: contain;
				mix-blend-mode: lighten;
			}
		}
	}

	.component-grid > div {
		display: grid;
		gap: 1rem;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		grid-template-rows: auto;
	}

	.component-grid {
		height: 22.5rem;
	}

	.tooltip {
		pointer-events: none;
	}
</style>
