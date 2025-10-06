<script>
	/**
	 * @typedef {Object} props
	 * @property {string} label
	 * @property {number} value
	 * @property {number | undefined} min
	 * @property {number | undefined} max
	 * @property {()=>void | undefined} onclick
	 */
	/** @type {props} */
	let { label, value = $bindable(), min, max, onclick } = $props();
</script>

<div
	style="display: inline-flex; flex-direction: row; gap: var(--pico-spacing);text-align: center; margin: 1rem; margin-top: 2rem; justify-items: center;"
>
	<button
		onclick={() => {
			if (min !== undefined && value > min) {
				value -= 1;
				onclick();
			}
		}}
		aria-label="-"
		class="minus"
	>
	</button>
	<span style="position: relative; min-width: 2rem; vertical-align: middle;">
		<div
			style="position: absolute; left: 50%; top: -1.5rem; translate: -50%; width: 10rem; text-align: center;"
		>
			{label}
		</div>
		<div style="height: 100%; line-height: normal; position: relative;">
			<span style="font-size: 1rem; position: absolute; top: 50%; left: 50%; translate: -50% -50%;">
				{value}
			</span>
		</div>
	</span>
	<button
		onclick={() => {
			if (max !== undefined && value < max) {
				value += 1;
				onclick();
			}
		}}
		aria-label="-"
		class="plus"
	>
	</button>
</div>

<style>
	button {
		border-radius: 50px;
		aspect-ratio: 1 / 1;
		padding: 0;
		height: 1.5rem;
		font-weight: bold;
		background-color: transparent;
		background-position: center;
		background-size: 100% 100%;
		border: none;
	}
	.minus {
		background-image: url('$lib/images/minusIcon.svg');
		&:active {
			background-image: url('$lib/images/minusIconDepressed.svg');
		}
	}

	.plus {
		background-image: url('$lib/images/plusIcon.svg');
		&:active {
			background-image: url('$lib/images/plusIconDepressed.svg');
		}
	}
</style>
