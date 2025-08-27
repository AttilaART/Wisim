<script>
	import { draggable, controls, ControlFrom } from '@neodrag/svelte';
	let { title, closeWindow, children } = $props();
</script>

<article {@attach draggable([controls({ allow: ControlFrom.selector('header') })])}>
	<header>
		<button onclick={closeWindow} rel="prev" aria-label="Close"></button><center>{title}</center>
	</header>
	<div style="overflow: scroll; max-height: calc(100vh - 136px - 58px - var(--pico-spacing));">
		{@render children()}
	</div>
</article>

<style>
	article {
		max-width: 100vh;
		max-height: calc(100vh - 136px);
		position: fixed;
		top: 0;
		left: 0;
		header {
			button {
				float: right;
				margin: 0;
				margin-left: var(--pico-spacing);
				padding: 0;
				display: block;
				width: 1rem;
				height: 1rem;
				margin-top: calc(var(--pico-spacing) * -1);
				margin-bottom: var(--pico-spacing);
				margin-left: auto;
				border: none;
				background-image: var(--pico-icon-close);
				background-position: center;
				background-size: auto 1rem;
				background-repeat: no-repeat;
				background-color: transparent;
				opacity: 0.5;
				transition: opacity var(--pico-transition);

				margin: 0;
				margin-left: 0px;
				margin-left: var(--pico-spacing);
				padding: 0;
				float: right;
			}
		}
	}
</style>
