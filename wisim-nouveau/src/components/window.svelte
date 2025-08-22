<script>
	import { draggable } from '@neodrag/svelte';
	let { children } = $props();
</script>

<article use:draggable>
	<header>Window</header>
	{@render children()}
</article>

<style>
	article {
		max-width: 100vh;
		max-height: 100vh;
		overflow: scroll;
	}
</style>
