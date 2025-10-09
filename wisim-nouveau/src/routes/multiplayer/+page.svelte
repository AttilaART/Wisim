<script>
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
				window.location.href = `http://localhost:5173/game?${serverAdress}`;
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
