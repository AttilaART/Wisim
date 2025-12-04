<script>
	import { goto } from '$app/navigation';

	let serverAdress = $state('');
</script>

<svelte:head>
	<title>Multiplayer</title>
	<meta name="description" content="Multiplayer server selecetion screen" />
</svelte:head>

<dialog open>
	<article>
		<form
			onsubmit={() => {
				console.log('connecting');
				window.location.replace(`/game?${serverAdress}/`);
			}}
		>
			<label for="">
				Server Adress
				<input type="text" placeholder="localhost:8000" required bind:value={serverAdress} />
				<input type="submit" value="Connect" />
			</label>
		</form>
	</article>
</dialog>
