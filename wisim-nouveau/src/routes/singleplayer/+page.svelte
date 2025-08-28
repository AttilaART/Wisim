<script>
	// as@ts-nocheck

	import Canvas from '../../components/canvas.svelte';
	import { baseState, Methods, newConnection } from '$lib/javascript/connection';
	import Product from '../../components/product.svelte';
	import Window from '../../components/window.svelte';
	import '@picocss/pico';
	import '../../app.css';
	import Marketing from '../../components/marketing.svelte';
	import Debt from '../../components/debt.svelte';
	import { format } from '$lib/javascript/format';
	import Employees from '../../components/employees.svelte';
	import Reasearch from '../../components/reasearch.svelte';
	import Finances from '../../components/finances.svelte';
	import FinancialReport from '../../components/financialReport.svelte';

	/** @type {Object.<string, import("svelte").Snippet<[number]>>} */
	let windows = $state({});

	/**
	 *@param {import("svelte").Snippet<[number]>} snippet
	 * @returns {number} window ID
	 */
	function newWindow(snippet) {
		let id = Math.round(Math.random() * 100000);

		if (windows[id] != undefined) {
			return newWindow(snippet);
		}
		windows[id] = snippet;
		return id;
	}

	/**
	 *@param {number} id
	 * @returns {void}
	 */
	function deleteWindow(id) {
		delete windows[id];
	}

	let errorDialogue = $state();
	let companyDialogue = $state();
	let isReady = $state();
	let isSimulating = $state(false);

	/**
	 * @type {number}
	 */
	let companyIDUserFacing = $state(0);
	let companyID = $derived(companyIDUserFacing - 1);

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
					break;
				case Methods.Get_decisions:
					clientState.decisions = dataJSON.Data;
					clientState.decisions.Employees.Production_deltas = [];
					clientState.decisions.Employees.Marketing_deltas = [];
					console.log('decisions updated');
					connection.fProduct_stats(
						clientState.decisions.Marketing.Product,
						clientState.decisions.Research
					);
					connection.sDecisions(clientState.decisions);
					break;
				case Methods.Get_external_factors:
					clientState.external_factors = dataJSON.Data;
					console.log('external_factors updated');
					break;
				case Methods.Func_calculate_product_stats:
					clientState.company.Offer.Product = dataJSON.Data;
					console.log('product Stats updated');
				case Methods.Get_employees:
					if (dataJSON.Data.Type == 'production') {
						clientState.employees.production = dataJSON.Data.Employees;
					} else if (dataJSON.Data.Type == 'marketing') {
						clientState.employees.marketing = dataJSON.Data.Employees;
					}
					console.log('employees updated');
				case Methods.Get_unemployed_employees:
					if (dataJSON.Data.Type == 'production') {
						clientState.unemployed.production = dataJSON.Data.Employees;
					} else if (dataJSON.Data.Type == 'marketing') {
						clientState.unemployed.marketing = dataJSON.Data.Employees;
					}
					console.log('unemployed updated');
			}
		} else {
			switch (dataJSON.Method) {
				case Methods.Sim_starting:
					isReady = true;
					isSimulating = true;
					break;
				case Methods.Sim_done:
					isReady = false;
					isSimulating = false;
					fetchEverything(connection);
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
		connection.gEmployees('production');
		connection.gEmployees('marketing');
		connection.gUnemployedEmployees('production');
		connection.gUnemployedEmployees('marketing');
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
				<input bind:value={companyIDUserFacing} type="number" required placeholder="Company ID" />
				<input type="submit" value="Confirm" />
			</form>
		</article>
	</dialog>

	{#if isSimulating}
		<dialog open>
			<center>
				<h2>Loading...</h2>
				<p>Simulation Simulating</p>
			</center>
		</dialog>
	{/if}

	<div id="ui">
		<Canvas>
			{#each Object.entries(windows) as w}
				{@render w[1](w[0])}
			{/each}
		</Canvas>
		<div id="bottom-menu">
			<button
				onclick={() => {
					newWindow(finances);
				}}
			>
				<strong>
					{format.currency(clientState.company.Balance, true, 2)}
				</strong>
			</button>
			<button
				onclick={() => {
					newWindow(product);
				}}>Product</button
			>
			<button
				onclick={() => {
					newWindow(marketing);
				}}>Marketing</button
			>
			<button
				onclick={() => {
					newWindow(employees);
				}}>Employees</button
			>
			<button>Production</button>
			<button
				onclick={() => {
					newWindow(research);
				}}>Research</button
			>
			<button>Market</button>
			<button>Chat</button>
			{#if !isReady}
				<button
					onclick={() => {
						connection.sReady();
						isReady = true;
					}}>Ready</button
				>{:else}
				<button
					class="secondary"
					onclick={() => {
						connection.sUnready();
						isReady = false;
					}}>Unready</button
				>{/if}
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
				<a href="/">
					<button
						onclick={() => {
							errorDialogue.close();
						}}>Return To Main Menu</button
					>
				</a>
			</footer>
		</article>
	</dialog>
{/await}

{#snippet product(/** @type {Number} id */ id)}
	<Window
		title="Product"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Product
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
				connection.fProduct_stats(decisions.Marketing.Product, decisions.Research);
			}}
			bind:product={clientState.company.Offer.Product}
			externalFactors={clientState.external_factors}
		></Product>
	</Window>
{/snippet}

{#snippet marketing(/** @type {Number} id */ id)}
	<Window
		title="Marketing"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Marketing
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></Marketing>
	</Window>
{/snippet}

{#snippet debt(/** @type {Number} id */ id)}
	<Window
		title="Debt"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Debt
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></Debt>
	</Window>
{/snippet}

{#snippet finances(/** @type {Number} id */ id)}
	<Window
		title="Finances"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Finances
			bind:clientState
			openDebtWindow={() => newWindow(debt)}
			openFinancialReportWindow={() => newWindow(financialReport)}
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></Finances>
	</Window>
{/snippet}

{#snippet financialReport(/** @type {Number} id */ id)}
	<Window
		title="Finance Report"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<FinancialReport
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></FinancialReport>
	</Window>
{/snippet}

{#snippet employees(/** @type {Number} id */ id)}
	<Window
		title="Employees"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Employees
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></Employees>
	</Window>
{/snippet}

{#snippet research(/** @type {Number} id */ id)}
	<Window
		title="Research"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Reasearch
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></Reasearch>
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
		backdrop-filter: blur(10px);
		border-top: 0.5px solid color-mix(in oklab, var(--pico-background-color), transparent 10%);
		z-index: 99;

		button {
			flex: 1 1;
		}
	}
</style>
