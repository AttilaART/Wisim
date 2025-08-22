<script>
	// as@ts-nocheck

	import Canvas from '../../components/canvas.svelte';
	import { baseState, Methods, newConnection } from '$lib/javascript/connection';
	import Product from '../../components/product.svelte';
	import Window from '../../components/window.svelte';
	import '@picocss/pico';
	import '../../app.css';
	import Marketing from '../../components/marketing.svelte';

	/** @type {import("svelte").Snippet[]} */
	let windows = $state([]);

	let errorDialogue = $state();
	let companyDialogue = $state();
	let isReady = $state();

	/**
	 * @type {number}
	 */
	let companyID = $state(-1);

	/** @type {Promise<{connection: import("$lib/javascript/connection").Connection, clientState: import("$lib/javascript/simulation").clientState}>} */
	let connectionPromise = $state(
		newConnection('ws://localhost:8000', handleConnection, onClose, onError)
	);
	/** @type {import("$lib/javascript/connection").Connection} */
	let connection = $state(null);
	/** @type {import("$lib/javascript/simulation").clientState} */
	let clientState = $state(JSON.parse(JSON.stringify(baseState)));

	$effect(() => {
		connectionPromise.then((value) => {
			clientState = value.clientState;
			connection = value.connection;
		});
	});

	/**
	 * @param {CloseEvent?} event
	 */
	async function onClose(event) {
		console.log(`WebSocket Closed. Reconnecting \n ${event}`);
		connectionPromise = newConnection('ws://localhost:8000', handleConnection, onClose, onError);

		({ connection, clientState } = await connectionPromise);
		connection.sCompany(companyID);
	}

	/**
	 * @param {Event} event
	 */
	function onError(event) {
		console.log(`WebSocket Closed. \n ${event}`);
	}

	/**
	 * @param {MessageEvent} event
	 */

	function handleConnection(event) {
		let dataJSON = JSON.parse(event.data);
		console.log(dataJSON);

		if (dataJSON.IsResponse) {
			switch (dataJSON.Method) {
				case Methods.Set_company:
					if (dataJSON.Data == true) {
						companyDialogue.close();
						fetchEverything(connection);
					} else {
						console.log('invalid company id set');
					}
					break;
				case Methods.Get_company:
					clientState.company = dataJSON.Data;
					console.log('company updated');
					console.log(clientState.company);
					break;
				case Methods.Get_decisions:
					clientState.decisions = dataJSON.Data;
					console.log('decisions updated');
					console.log(clientState.decisions);
					break;
				case Methods.Get_external_factors:
					clientState.external_factors = dataJSON.Data;
					console.log('external_factors updated');
					console.log(clientState.external_factors);
					break;
			}
		}
	}

	/**
	 * @param {import("$lib/javascript/connection").Connection} connection
	 */
	function fetchEverything(connection) {
		console.log('fetching everything');
		connection.gCompany();
		connection.gDecisions();
		connection.gExternal_factors();
	}

	/**
	 * @type {string}
	 */
	let chat = $state('');
</script>

<svelte:head>
	<title>Singlplayer</title>
	<meta name="description" content="About this app" />
</svelte:head>

{#await connectionPromise}
	<dialog open>
		<center>
			<h2>Loading...</h2>
			<p>Connecting to server</p>
		</center>
	</dialog>
{:then}
	<dialog bind:this={companyDialogue} open>
		<article>
			<header>
				<p>Choose Your Company</p>
			</header>
			<form
				onsubmit={() => {
					connection.sCompany(companyID);
				}}
			>
				<input bind:value={companyID} type="number" required placeholder="Company ID" />
				<input type="submit" value="Confirm" />
			</form>
		</article>
	</dialog>
	<div id="ui">
		<Canvas>
			{#each windows as w}
				{@render w()}
			{/each}
		</Canvas>
		<div id="bottom-menu">
			<button>{clientState?.company.Balance} $</button>
			<button
				onclick={() => {
					windows.push(product);
				}}>Product</button
			>
			<button
				onclick={() => {
					windows.push(marketing);
				}}>Marketing</button
			>
			<button>Employees</button>
			<button>Production</button>
			<button>Research</button>
			<button>Market</button>
			<button>Chat</button>
			<button
				style="display: {isReady ? 'none' : 'unset'};"
				onclick={() => {
					connection.sReady();
					isReady = true;
				}}>Ready</button
			>
			<button
				class="secondary"
				style="display: {!isReady ? 'none' : 'unset'};"
				onclick={() => {
					connection.sUnready();
					isReady = false;
				}}>Unready</button
			>
		</div>
	</div>
{:catch error}
	<script>
	</script>
	<dialog bind:this={errorDialogue} open>
		<article>
			<header>
				<button
					aria-label="Close"
					rel="prev"
					onclick={() => {
						errorDialogue.close();
					}}
				></button>
				<p>Connection Failed</p>
			</header>
			<div>
				<p>Failed to connect to server</p>
				<small>{error}</small>
			</div>
			<footer>
				<button
					class="secondary"
					onclick={() => {
						onClose(null);
					}}>Retry Connection</button
				>
				<button
					onclick={() => {
						errorDialogue.close();
					}}>Return To Main Menu</button
				>
			</footer>
		</article>
	</dialog>
{/await}

{#snippet product()}
	<Window title="Product" closeWindow={() => {}}>
		<Product
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></Product>
	</Window>
{/snippet}

{#snippet marketing()}
	<Window title="Marketing" closeWindow={() => {}}>
		<Marketing
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></Marketing>
	</Window>
{/snippet}

<style>
	#ui {
		display: flex;
		flex-direction: column;
		height: 100vh;
		width: 100%;
	}

	#bottom-menu {
		flex: 0 0;
		display: flex;
		flex-direction: row;
		gap: var(--pico-spacing);
		padding: var(--pico-spacing);

		button {
			flex: 1 1;
		}
	}
</style>
