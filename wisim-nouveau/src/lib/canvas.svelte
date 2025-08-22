<script>
	/** @type {{children: import("svelte").Snippet?}} */
	let { children } = $props();
</script>

<div>
	{@render children?.()}
</div>

<style>
	div {
		flex: 1 1;
	}
</style>
