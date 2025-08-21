<script>
	import Header from './Header.svelte';
	import '../app.css';
	import '../pico.css';

	/** @type {{children: import('svelte').Snippet}} */
	let { children } = $props();
</script>

<div class="app">
	{@render children()}
</div>

<style>
	.app {
		min-height: 100vh;
		min-width: 100vw;
	}
</style>
