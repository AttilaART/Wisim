<script>
	import { draggable, controls, ControlFrom, Compartment, position, events } from '@neodrag/svelte';
	import { render } from 'svelte/server';
	let { title, closeWindow, children } = $props();
	/** @type {{x: number, y: number}}*/
	let viewportSize = $derived({
		x: document.documentElement.clientWidth,
		y: document.documentElement.clientHeight
	});
	let windowPos = $state({ x: 20, y: 20 });

	const positionComp = Compartment.of(() => position({ current: windowPos }));

	const eventsComp = Compartment.of(() =>
		events({
			onDrag: (data) => {
				windowPos = validateWindowPos(data.offset);
				savePosition(windowPos);
			}
		})
	);

	/** @param {{x: number, y: number}} pos*/
	function savePosition(pos) {
		localStorage.setItem('windowPos', JSON.stringify(pos));
	}

	/** @param {{x: number, y: number}} pos*/
	function validateWindowPos(pos) {
		if (pos.y < 32) pos.y = 32;
		else if (pos.y > viewportSize.y - 200) pos.y = viewportSize.y - 200;
		if (pos.x < 0) pos.x = 0;
		else if (pos.x > viewportSize.x - 10) pos.x = viewportSize.x - 10;

		return pos;
	}
</script>

<article
	{@attach draggable(() => [
		controls({ allow: ControlFrom.selector('.windowHeader') }),
		positionComp,
		eventsComp
	])}
>
	<header class="windowHeader">
		<button class="windowHeader" onclick={closeWindow} rel="prev" aria-label="Close"></button>
		<span class="windowHeader">{title}</span>
	</header>
	<div style="overflow: auto; max-height: calc(100vh - 136px - 58px - var(--pico-spacing)*2);">
		{@render children()}
	</div>
</article>

<style>
	article {
		max-width: 100vw;
		max-height: calc(100vh - 136px);
		position: fixed;
		top: 0;
		left: 0;
		header {
			height: 3rem;

			user-select: none;
			position: relative;
			span {
				pointer-events: none;

				position: absolute;
				left: 50%;
				top: 50%;
				translate: -50% -50%;
				text-align: center;
				white-space: nowrap;
			}
		}
		div {
			resize: both;
		}
	}
</style>
