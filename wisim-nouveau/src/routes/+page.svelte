<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section id="main-menu-container" class="container">
	<article id="main-menu">
		<h1>Wisim</h1>
		<div id="buttons">
			<button disabled>Continue</button>
			<a href="/multiplayer"><button>Multiplayer</button></a>
			<a href="/game?localhost:8000"><button>Singleplayer</button></a>
			<button class="secondary">Exit</button>
		</div>
	</article>
</section>

<style>
	#main-menu-container {
		display: flex;
		align-items: center;
		justify-content: center;

		min-height: 100vh;
		min-width: 100%;

		margin: 0;
	}

	#buttons {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--pico-spacing);

		a {
			width: 100%;
			button {
				width: 100%;
			}
		}
	}

	#main-menu {
		text-align: center;
		min-width: 20rem;
	}
</style>
