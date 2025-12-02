<script>
	import { goto } from '$app/navigation';
	import { preventPageReload } from '$lib/helper.svelte';
	/** @type {{data: {serverAdress: string}}}*/
	let { data } = $props();

	/** @type {FileList | undefined}*/
	let fileList = $derived(undefined);

	/** @type {File | undefined}*/
	let saveFile = $derived(fileList !== undefined ? fileList[0] : undefined);
	let errorText = $state('');
</script>

<section class="container">
	<form
		style="padding: 8rem  4rem;"
		use:preventPageReload
		onsubmit={async () => {
			if (saveFile === undefined) {
				errorText = 'No file given';
				return;
			}
			console.log(saveFile);

			if (!saveFile.name.endsWith('.notajson')) {
				errorText = "incorrect file format: expected '.notajson'";
				return;
			}

			let response = await fetch(`http://${data.serverAdress}/save`, {
				method: 'POST',
				// @ts-ignore
				body: await saveFile.bytes()
			});

			let responseText = await response.text();

			console.log(responseText);

			if (responseText.includes('success')) {
				goto(`/game?${data.serverAdress}/`);
			} else {
				errorText = responseText;
			}
		}}
	>
		<h1>Load save file</h1>
		<input type="file" bind:files={fileList} />

		{#if errorText != ''}<strong style="color: red;">Error: {errorText}</strong>{/if}

		<input type="submit" disabled={saveFile == undefined} />
	</form>
</section>
