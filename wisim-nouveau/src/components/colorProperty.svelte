<script>
	/** @type {{value: number, invert: undefined | boolean, formatter: (v: number)=>string}}*/
	let { value, invert, formatter } = $props();
</script>

<span class={value * (invert ? -1 : 1) >= 0 ? 'green' : 'red'}>{formatter(value)}</span>
