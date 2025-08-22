<script>
	import '@picocss/pico';
	import '../app.css';

	/** @type {{children: import('svelte').Snippet}} */
	let { children } = $props();
</script>

<svelte:boundary>
	<div class="app">
		{@render children()}
	</div>
	{#snippet pending()}
		<p>loading</p>
	{/snippet}
</svelte:boundary>

<style>
	.app {
		min-height: 100vh;
		min-width: 100vw;
	}
</style>
