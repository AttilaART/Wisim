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
	import Invoices from '../../components/invoices.svelte';
	import Production from '../../components/production.svelte';
	import { preventPageReload } from '$lib/helper.svelte';
	import MonthlyOverview from '../../components/monthlyOverview.svelte';

	/** handle wasm import */
	import { wasm_exec } from '$lib/wasm_exec.js';

	wasm_exec();

	// @ts-ignore
	const go = new Go();
	go.run((await WebAssembly.instantiateStreaming(fetch('/main.wasm'), go.importObject)).instance);

	/**/

	/** @type {{data: {serverAdress: string}}}*/
	let { data } = $props();

	/** @type {Object.<string, (id: number)=>ReturnType<import("svelte").Snippet>>} */
	let windows = $state({});

	/**
	 *@param {(id: number)=>ReturnType<import("svelte").Snippet>} snippet
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
	let overviewWindowOpen = $state(false);

	/**
	 * @type {number}
	 */
	let companyIDUserFacing = $state(0);
	let companyID = $derived(companyIDUserFacing - 1);

	/** @type {Promise<{connection: import("$lib/javascript/connection").Connection, clientState: import("$lib/javascript/simulation").clientState}>} */
	let connectionPromise = $state(
		newConnection(`ws://${data.serverAdress}`, handleConnection, onClose, onError)
	);
	/** @type {import("$lib/javascript/connection").Connection} */
	let connection = $state(null);
	/** @type {import("$lib/javascript/simulation").clientState} */
	let clientState = $state(JSON.parse(JSON.stringify(baseState)));
	/** @type {{Message: string, From: string}[]} */
	let chats = $state([]);
	/** @type {string} */
	let chatMessage = $state('');

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
		connectionPromise = newConnection(
			`ws://${data.serverAdress}`,
			handleConnection,
			onClose,
			onError
		);

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
					clientState.Company = dataJSON.Data;
					console.log('company updated');
					break;
				case Methods.Get_decisions:
					clientState.Decisions = dataJSON.Data;
					clientState.Decisions.Employees.ProductionDeltas = [];
					clientState.Decisions.Employees.MarketingDeltas = [];
					console.log('decisions updated');
					connection.sDecisions(clientState.Decisions);
					break;
				case Methods.Get_external_factors:
					clientState.ExternalFactors = dataJSON.Data;
					console.log('external_factors updated');
					break;
				case Methods.Get_employees:
					if (dataJSON.Data.Type == 'production') {
						clientState.Employees.production = dataJSON.Data.Employees;
					} else if (dataJSON.Data.Type == 'marketing') {
						clientState.Employees.marketing = dataJSON.Data.Employees;
					}
					console.log('employees updated');
				case Methods.Get_unemployed_employees:
					if (dataJSON.Data.Type == 'production') {
						clientState.Unemployed.production = dataJSON.Data.Employees;
					} else if (dataJSON.Data.Type == 'marketing') {
						clientState.Unemployed.marketing = dataJSON.Data.Employees;
					}
					console.log('unemployed updated');
				case Methods.Get_product_components:
					clientState.productComponents = dataJSON.Data;
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
					if (!overviewWindowOpen) {
						newWindow(monthlyOverview);
					}
					break;
				case Methods.Broadcast_chat:
					chats.push(dataJSON.Data);
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
		connection.gProductComponents();
	}
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
				use:preventPageReload
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
			{#each Object.entries(windows) as w (w[0])}
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
					{format.currency(clientState.Company.Balance, true, 2)}
				</strong>
			</button>
			<button
				onclick={() => {
					newWindow(product);
				}}>Product</button
			>
			<button
				onclick={() => {
					newWindow(employees);
				}}>Employees</button
			>
			<button
				onclick={() => {
					newWindow(production);
				}}>Production</button
			>
			<button
				onclick={() => {
					newWindow(research);
				}}>Research</button
			>
			<button>Market</button>
			<button
				onclick={() => {
					newWindow(chat);
				}}>Chat</button
			>
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
			}}
			{newWindow}
			deleteWindow={(id) => {
				deleteWindow(id);
			}}
			openProduction={() => {
				newWindow(production);
			}}
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
			openInvoicesWindow={() => newWindow(invoiceLog)}
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

{#snippet invoiceLog(/** @type {Number} id */ id)}
	<Window
		title="Finance Report"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Invoices
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></Invoices>
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

{#snippet production(/** @type {Number} id */ id)}
	<Window
		title="Production"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Production
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
		></Production>
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

{#snippet chat(/** @type {Number} id */ id)}
	<Window
		title="Chat"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<div style="display: grid; grid-template-columns: auto 1fr;">
			{#each chats as c}
				<span>
					{c.From}:
				</span>
				<span style="margin-left: 10px;">
					{c.Message}
				</span>
			{/each}
		</div>
		<form
			onsubmit={() => {
				connection.bChat(chatMessage);
				chatMessage = '';
			}}
		>
			<input bind:value={chatMessage} type="text" />
			<input type="submit" style="display: none" />
		</form>
	</Window>
{/snippet}

{#snippet monthlyOverview(/** @type {Number} id */ id)}
	<Window
		title="Monthly Report"
		closeWindow={() => {
			overviewWindowOpen = false;
			deleteWindow(id);
		}}
	>
		<span hidden>{(overviewWindowOpen = true)}</span>
		<MonthlyOverview bind:clientState></MonthlyOverview>
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
