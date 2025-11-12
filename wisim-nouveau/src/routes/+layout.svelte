<script>
	import '../dynamic-pico.css';
	import '../app.css';

	/** @type {{children: import('svelte').Snippet}} */
	let { children } = $props();
</script>

<svelte:boundary>
	<div class="app">
		{@render children()}
	</div>
	{#snippet pending()}
		<dialog open>
			<center>
				<h2>Loading...</h2>
			</center>
		</dialog>
	{/snippet}
</svelte:boundary>

<style>
	.app {
		min-height: 100vh;
		min-width: 100vw;
	}
</style>
