<script>
	// as@ts-nocheck

	import Canvas from '../../components/canvas.svelte';
	import { baseState, Methods, newConnection } from '$lib/javascript/connection';
	import Product from '../../components/product.svelte';
	import Window from '../../components/window.svelte';
	import Marketing from '../../components/marketing.svelte';
	import Debt from '../../components/debt.svelte';
	import { format } from '$lib/javascript/format';
	import Employees from '../../components/employees.svelte';
	import Reasearch from '../../components/reasearch.svelte';
	import Production from '../../components/production.svelte';
	import { beforeNavigate } from '$app/navigation';
	import {
		preventPageReload,
		simulateMockStep,
		syncCompanyWithDecisions
	} from '$lib/helper.svelte';
	import MonthlyOverview from '../../components/monthlyOverview.svelte';

	/** handle wasm import */
	import { wasm_exec } from '$lib/wasm_exec.js';
	import Reports from '../../components/reports.svelte';
	import Market from '../../components/market.svelte';
	import ColorProperty from '../../components/colorProperty.svelte';
	import { on } from 'svelte/events';
	import { redirect } from '@sveltejs/kit';
	import { number } from 'echarts/core';

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
	let isSimulatingMock = $state(false);
	let mainMenuOpen = $state(false);
	let quitDialogueOpen = $derived(!mainMenuOpen);
	let host = $state('');

	/** @param {KeyboardEvent} event */
	function onKeyUp(event) {
		if (event.key == 'Escape') {
			mainMenuOpen = !mainMenuOpen;
		}
	}

	/**
	 * @type {number}
	 */
	let companyIDUserFacing = $state(0);
	let companyID = $derived(companyIDUserFacing - 1);

	function handleMockSimulation() {
		clientState = JSON.parse(JSON.stringify(clientStateBase));

		isSimulatingMock = true;

		clientState.Company = simulateMockStep(
			clientState.Company,
			clientState.Decisions,
			clientState.ExternalFactors,
			{ Employees: clientState.Employees, Unemployed: clientState.Unemployed },
			() => {
				isSimulatingMock = false;
				clientState.predictionMode = true;
			},
			clientState.Decisions.Predictions.Steps,
			clientState.productComponents
		);

		clientStateBase.Decisions = clientState.Decisions;
	}

	const hooks = {
		sDecisions: () => {
			if (clientState.predictionMode) {
				handleMockSimulation();
			}
		}
	};

	/** @type {Promise<{connection: import("$lib/javascript/connection").Connection, clientState: import("$lib/javascript/simulation").clientState}>} */
	let connectionPromise = $state(
		newConnection(`ws://${data.serverAdress}`, handleConnection, onClose, onError, hooks)
	);
	/** @type {import("$lib/javascript/connection").Connection} */
	let connection = $state(null);
	/** @type {import("$lib/javascript/simulation").clientState} */
	let clientState = $state(JSON.parse(JSON.stringify(baseState)));

	/** @type {import("$lib/javascript/simulation").clientState} */
	let clientStateBase = $state(JSON.parse(JSON.stringify(clientState)));

	function togglePredictionMode() {
		if (clientState.predictionMode) {
			clientStateBase.Decisions = clientState.Decisions;
			clientStateBase.Employees = clientState.Employees;
			clientStateBase.Unemployed = clientState.Unemployed;

			clientState = clientStateBase;

			clientState.Company = syncCompanyWithDecisions(
				clientState.Company,
				clientState.Decisions,
				clientState.productComponents
			);

			clientState.predictionMode = false;
		} else {
			clientStateBase = JSON.parse(JSON.stringify(clientState));

			clientState.predictionMode = true;

			handleMockSimulation();
		}
	}

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
		if (event?.code == 1000) {
			return;
		}
		console.log(`WebSocket Closed. Reconnecting \n ${event}`);
		connectionPromise = newConnection(
			`ws://${data.serverAdress}`,
			handleConnection,
			onClose,
			onError,
			hooks
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

	let invalidCompany = $state(false);
	let renamedCompany = $state(false);
	let namedCEO = $state(false);

	/** @type {{ID: number,Name: string,Balance: number,Taken: boolean}[]} */
	let gameCompanies = $state([]);

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
					/*
					if (!overviewWindowOpen) {
						newWindow(monthlyOverview);
					}
          */
					break;
				case Methods.Broadcast_chat:
					chats.push(dataJSON.Data);
				case Methods.Broadcast_server_address:
					host = dataJSON.Data;
				case Methods.Broadcast_game_companies:
					gameCompanies = dataJSON.Data;
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

	beforeNavigate(() => {
		connection.socket.close(1000);
	});

	function totalProductsSold() {
		let total = 0;
		for (let r of clientState.Company.Reports) {
			for (let p of Object.values(r.SalesReport)) {
				total += p.ProductSalesStatistics.ProductsSold;
			}
		}
		return total;
	}
</script>

<svelte:head>
	<title>Singlplayer</title>
	<meta name="description" content="About this app" />
</svelte:head>

<svelte:window onkeyup={onKeyUp} />

{#await connectionPromise}
	<dialog open>
		<center>
			<h2>Loading...</h2>
			<p>Connecting to server</p>
		</center>
	</dialog>
{:then}
	<dialog bind:this={companyDialogue} open>
		<article style="min-height: 30rem; max-height: 60%; overflow-y: auto; width: 60vw;;">
			<header>
				<p>Choose Your Company</p>
			</header>

			{#if gameCompanies == []}
				<h4>Loading ...</h4>
			{:else}
				<div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: var(--spacing); ">
					{#each gameCompanies as c}
						<button
							class="contrast outline"
							style="{c.Taken ? 'opacity: 0.5;' : ''} text-align: left;"
							onclick={() => {
								connection.sCompany(c.ID);
							}}
						>
							<h2>{c.Name}</h2>
							<small>Balance: {format.currency(c.Balance, true, 0)}</small>
						</button>
					{/each}
				</div>
				{#if invalidCompany}
					<small style="color: red;">Invalid company chosen</small>
				{/if}
			{/if}
		</article>
	</dialog>

	{#if !renamedCompany && clientState.Company.Name.includes('\r')}
		<dialog open>
			<article>
				<h2>Name Company</h2>
				<form
					use:preventPageReload
					onsubmit={() => {
						clientState.Company.Name = clientState.Decisions.General.CompanyName;
						renamedCompany = true;

						connection.sDecisions(clientState.Decisions);
					}}
				>
					<input
						type="text"
						onfocus={(event) => event.target.select()}
						bind:value={clientState.Decisions.General.CompanyName}
					/>
					<input type="submit" value="Set Company Name" />
				</form>
			</article>
		</dialog>
	{:else if !namedCEO && clientState.Company.CEO == '' && clientState.Company.Name.length > 0}
		<dialog open>
			<article>
				<h2>Name CEO</h2>
				<form
					use:preventPageReload
					onsubmit={() => {
						clientState.Company.CEO = clientState.Decisions.General.CEO;
						renamedCompany = true;

						connection.sDecisions(clientState.Decisions);
					}}
				>
					<input
						type="text"
						onfocus={(event) => event.target.select()}
						bind:value={clientState.Decisions.General.CEO}
					/>
					<small>You cannot change this later</small>
					<input type="submit" value="Set CEO Name" />
				</form>
			</article>
		</dialog>
	{/if}

	{#if isSimulating}
		<dialog open>
			<center>
				<h2>Loading...</h2>
				<p>Simulation Simulating</p>
			</center>
		</dialog>
	{/if}

	{#if isSimulatingMock}
		<dialog open>
			<center>
				<h2>Loading...</h2>
				<p>Calculating Budget</p>
			</center>
		</dialog>
	{/if}

	{#if mainMenuOpen}
		<dialog open>
			<article>
				<h1>Main Menu</h1>
				<div style="display: flex; flex-direction: column;">
					<button
						onclick={() => {
							mainMenuOpen = !mainMenuOpen;
						}}>Back to Game</button
					>
					<a href="http://{data.serverAdress}/save/" target="_blank"><button>Save</button></a>
					<button
						class="outline contrast"
						onclick={() => {
							quitDialogueOpen = !quitDialogueOpen;
						}}>Quit</button
					>
				</div>

				{#if host != ''}<small>Server address: {host}</small>{/if}
			</article>
		</dialog>

		{#if quitDialogueOpen}
			<dialog open>
				<article>
					<h1>Are you sure you want to quit</h1>
					<p>Any unsaved data will be lost</p>
					<footer class="grid">
						<a
							href="/"
							onclick={() => {
								if (
									data.serverAdress.includes('localhost') ||
									data.serverAdress.includes('127.0.0.1')
								)
									fetch(`http://${data.serverAdress}/quit/`);
							}}><button class="outline">Quit</button></a
						>
						<button
							onclick={() => {
								quitDialogueOpen = !quitDialogueOpen;
							}}>Go back</button
						>
					</footer>
				</article>
			</dialog>
		{/if}
	{/if}

	<div id="ui" data-scheme={clientState.predictionMode ? 'prediction' : 'normal'}>
		<div id="top-bar">
			<span>{clientState.Company.Name}</span>
			<span
				>Balance: <ColorProperty
					value={clientState.Company.Balance}
					invert={false}
					formatter={(v) => format.currency(v, true, 0)}
				></ColorProperty></span
			>
			<span
				>Monthly Cashflow: <ColorProperty
					value={clientState.Company.Reports.length >= 1
						? clientState.Company.Reports[clientState.Company.Reports.length - 1].FinancialReport
								.Totals.Cashflow
						: 0}
					invert={false}
					formatter={(/** @type {number}*/ v) => {
						return format.currency(v, true, 0);
					}}
				></ColorProperty></span
			>
			<span>
				Employees: {format.number(
					clientState.Employees.marketing.length + clientState.Employees.production.length,
					false,
					0
				)}
			</span>
			<span>
				Total Sales: {clientState.Company.Reports.length >= 1
					? format.number(totalProductsSold(), false, 0)
					: 0}
			</span>
		</div>

		<button
			class="contrast outline"
			style="height: 2.5rem; width: 2.5rem; margin: 0.5rem; border-radius: 100px; z-index: 99; text-align: center; font-size: 2rem; position: absolute; top: 0.5rem; right: 0.5rem;"
			onclick={() => {
				mainMenuOpen = !mainMenuOpen;
			}}
		>
			<span
				style="position: absolute; top: 50%; left: 50%; transform: translate(calc(-50% + 0.01rem), calc(-50% - 0.05rem));"
				>≡</span
			></button
		>

		<Canvas>
			{#each Object.entries(windows) as w (w[0])}
				{@render w[1](Number(w[0]))}
			{/each}
		</Canvas>
		<div id="bottom-menu">
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
			<button
				onclick={() => {
					newWindow(market);
				}}>Market</button
			>
			<button
				onclick={() => {
					newWindow(reports);
				}}
			>
				Reports
			</button>

			{#snippet predictionMenu()}
				<article id="prediction">
					{#if clientState.Company.Offers != undefined && clientState.Decisions.Predictions.ProductSales != undefined}
						{#each Object.entries(clientState.Company.Offers) as o}
							{#if clientState.Decisions.Products[o[0]] != undefined}
								{#if !clientState.Decisions.Products[o[0]].Outdated}
									<label for="">
										<strong>{o[1].Product.Name}</strong> demand
										<input
											bind:value={clientState.Decisions.Predictions.ProductSales[o[0]]}
											onchange={() => {
												connection.sDecisions(clientState.Decisions);
											}}
											type="number"
											step="100"
											min="0"
										/>
										{#if clientState.predictionMode && clientState.Decisions.Predictions.ProductSales[o[0]] > clientState.Company.Reports[clientState.Company.Reports.length - 1].ProductionReport.ProductSpecificReport[o[0]].TotalProductsProduced}
											<small style="color: orange; display: block; margin-top: -0.5rem;"
												>Predicted Demand is more than production.
											</small>
										{/if}
									</label>
								{/if}
							{/if}
						{/each}
					{/if}
					<fieldset role="group">
						<input disabled type="text" value="Steps to calculate" />
						<input
							type="number"
							bind:value={clientState.Decisions.Predictions.Steps}
							min="1"
							step="1"
							onchange={() => {
								connection.sDecisions(clientState.Decisions);
							}}
						/>
					</fieldset>
					<button
						style="float: right;"
						onclick={() => {
							togglePredictionMode();
						}}
					>
						{#if clientState.predictionMode}
							Exit Budget Mode
						{:else}
							Enter Budget Mode
						{/if}
					</button>
				</article>
			{/snippet}

			<div id="ready-div">
				{#if clientState.predictionMode}
					<button
						style="float: right;"
						onclick={() => {
							togglePredictionMode();
						}}>Exit Budget Mode</button
					>
					{@render predictionMenu()}
				{:else if !isReady}
					<button
						id="ready"
						onclick={() => {
							connection.sReady();
							isReady = true;
						}}>Ready</button
					>
					{@render predictionMenu()}
				{:else}
					<button
						class="secondary"
						onclick={() => {
							connection.sUnready();
							isReady = false;
						}}>Unready</button
					>
				{/if}
			</div>
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
			{deleteWindow}
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

{#snippet reports(/** @type {Number} id */ id)}
	<Window
		title="Reports"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Reports
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
			{newWindow}
			{deleteWindow}
		></Reports>
	</Window>
{/snippet}

{#snippet market(/** @type {Number} id */ id)}
	<Window
		title="Market"
		closeWindow={() => {
			deleteWindow(id);
		}}
	>
		<Market
			bind:clientState
			updateDecisions={(decisions) => {
				connection.sDecisions(decisions);
			}}
			{newWindow}
			{deleteWindow}
			serverAdress={data.serverAdress}
		></Market>
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

<!--
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
			use:preventPageReload
			onsubmit={() => {
				connection.bChat(chatMessage);
				chatMessage = '';
			}}
		>
			<input bind:value={chatMessage} type="text" />
			<input type="submit" style="display: none" />
		</form>
	</Window>
{/snippet}/
-->

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
		gap: var(--spacing);
		padding: calc(var(--spacing) * 0.5);
		backdrop-filter: blur(10px);
		border-top: 0.5px solid color-mix(in oklab, var(--background), transparent 30%);
		z-index: 99;

		button,
		div {
			flex: 1 1;
		}
	}

	#top-bar {
		padding-top: 0.5rem;
		display: flex;
		backdrop-filter: blur(10px);
		border-bottom: 0.5px solid color-mix(in oklab, var(--background), transparent 30%);
		z-index: 99;
		* {
			flex: 1 0;
			text-align: center;
		}
	}

	#ready-div {
		position: relative;
		button {
			width: 100%;
		}
	}

	#ready {
		width: 100%;
	}

	#prediction {
		position: absolute;
		bottom: 0;
		transform: scaleY(0);
		right: 0;
		opacity: 0%;
		width: 10rem;

		&:has(:first-child) {
			width: 30rem;
		}

		&:not(:has(:first-child)) {
			display: none;
		}

		transition:
			opacity 0.5s,
			bottom 0.75s 0.5s,
			transform 0s 1s;
	}

	#ready-div:hover #prediction {
		bottom: calc(100% + 1rem);
		opacity: 100%;
		transform: scaleY(1);
		transition:
			opacity 0.5s 0.25s,
			bottom 0.75s;
	}
</style>
