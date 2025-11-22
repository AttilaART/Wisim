<script>
	import { average } from '$lib/helper.svelte';
	import { format } from '$lib/javascript/format';
	import { delta } from '$lib/javascript/simulation';
	import { flip } from 'svelte/animate';
	import { fade, fly } from 'svelte/transition';

	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void}} */
	let { clientState = $bindable(), updateDecisions } = $props();

	/** @type {[string, string]}*/
	let menuState = $state(['all', 'hired']);
	/**
	 * @typedef {Object} Delta
	 * @property {import("$lib/javascript/simulation").Change} Delta.Delta_New
	 * @property {import("$lib/javascript/simulation").Change} Delta.Delta_Change
	 * @property {import("$lib/javascript/simulation").Change} Delta.Delta_Remove
	 */

	function syncDecisionsAndEmployees() {
		/**
		 * @param {import("$lib/javascript/simulation").Delta<import("$lib/javascript/simulation").Employee>} d
		 * @param {import("$lib/javascript/simulation").Employee[]} employees
		 * @param {import("$lib/javascript/simulation").Employee[]} unemployed
		 */
		function handleDeltaNew(d, employees, unemployed) {
			let indexInEmployees = employees.findIndex((e) => {
				return d.Item.ID == e.ID;
			});

			if (indexInEmployees < 0) {
				employees.push(d.Item);
			}

			let indexInUnemployed = unemployed.findIndex((e) => {
				return d.Item.ID == e.ID;
			});

			if (indexInUnemployed >= 0) unemployed.splice(indexInUnemployed, 1);
		}

		/**
		 * @param {import("$lib/javascript/simulation").Delta<import("$lib/javascript/simulation").Employee>} d
		 * @param {import("$lib/javascript/simulation").Employee[]} employees
		 * @param {import("$lib/javascript/simulation").Employee[]} unemployed
		 */
		function handleDeltaRemove(d, employees, unemployed) {
			let indexInEmployees = unemployed.findIndex((e) => {
				return d.Item.ID == e.ID;
			});

			if (indexInEmployees < 0) {
				unemployed.push(d.Item);
			}

			let indexInUnemployed = employees.findIndex((e) => {
				return d.Item.ID == e.ID;
			});

			if (indexInUnemployed >= 0) employees.splice(indexInUnemployed, 1);
		}

		for (let d of clientState.Decisions.Employees.ProductionDeltas) {
			if (d.Change == delta.Delta_New) {
				handleDeltaNew(d, clientState.Employees.production, clientState.Unemployed.production);
			} else if (d.Change == delta.Delta_Remove) {
				handleDeltaRemove(d, clientState.Employees.production, clientState.Unemployed.production);
			}
		}
		for (let d of clientState.Decisions.Employees.MarketingDeltas) {
			if (d.Change == delta.Delta_New) {
				handleDeltaNew(d, clientState.Employees.marketing, clientState.Unemployed.marketing);
			} else if (d.Change == delta.Delta_Remove) {
				handleDeltaRemove(d, clientState.Employees.marketing, clientState.Unemployed.marketing);
			}
		}
	}

	/**
	 * @param {import("$lib/javascript/simulation").Change} change
	 * @param {import("$lib/javascript/simulation").Employee} employee
	 * @param {string} employeeType
	 * @returns {void}
	 */
	function modifyEmployee(change, employee, employeeType) {
		// check if employee is already in deltas
		let deltasList = [];
		employee.Employer = clientState.Company.ID;

		if (employeeType == 'production') {
			deltasList = clientState.Decisions.Employees.ProductionDeltas;
		} else if (employeeType == 'marketing') {
			deltasList = clientState.Decisions.Employees.MarketingDeltas;
		} else throw `invalid employeeType: ${employeeType}`;

		for (/** @type {number} i */ let i in deltasList) {
			if (deltasList[i].Item.ID == employee.ID) {
				if (change == delta.Delta_Change && deltasList[i].Change == delta.Delta_New)
					change = delta.Delta_New;
				deltasList[i].Change = change;
				deltasList[i].Item = employee;

				new Promise(() => {
					syncDecisionsAndEmployees();
				});
				return;
			}
		}

		deltasList.push({
			Change: change,
			Item: employee
		});

		new Promise(() => {
			updateDecisions(clientState.Decisions);
			syncDecisionsAndEmployees();
		});
	}

	/**
	 * @param {import("$lib/javascript/simulation").Employee} e
	 * @param {import("$lib/javascript/simulation").Employee[]} employees
	 */
	function isEmployee(e, employees) {
		return (
			employees.findIndex((em) => {
				return em.ID == e.ID;
			}) >= 0
		);
	}

	/**
	 * @param {string} type
	 * @param {number} index
	 */
	function onModifyEmployee(type, index) {
		// @ts-ignore
		modifyEmployee(delta.Delta_Change, clientState.Employees[type][index], type);
		// console.log($state.snapshot(clientState.Decisions));
		updateDecisions(clientState.Decisions);
	}

	/**
	 * @param {import("$lib/javascript/simulation").Employee} a
	 * @param {import("$lib/javascript/simulation").Employee} b
	 */
	function sortBySkill(a, b) {
		if (a.Skill < b.Skill) return -1;
		if (a.Skill > b.Skill) return 1;
		return 0;
	}

	/** @param {[string, string]} menuState */
	function getEmployeeArray(menuState) {
		/** @type {import("$lib/javascript/simulation").Employee[]} */
		let employeeArray = [];
		if (menuState[0] == 'all' && menuState[1] == 'hired') {
			employeeArray = employeeArray.concat(
				clientState.Employees.marketing,
				clientState.Employees.production
			);
		} else if (menuState[0] == 'production' && menuState[1] == 'hired') {
			employeeArray = clientState.Employees.production;
		} else if (menuState[0] == 'marketing' && menuState[1] == 'hired') {
			employeeArray = clientState.Employees.marketing;
		} else if (menuState[0] == 'all' && menuState[1] == 'prospective') {
			employeeArray = employeeArray.concat(
				clientState.Unemployed.marketing,
				clientState.Unemployed.production
			);
		} else if (menuState[0] == 'production' && menuState[1] == 'prospective') {
			employeeArray = clientState.Unemployed.production;
		} else if (menuState[0] == 'marketing' && menuState[1] == 'prospective') {
			employeeArray = clientState.Unemployed.marketing;
		}

		return employeeArray;
	}

	let employeeArray = $derived(getEmployeeArray(menuState));

	/** @param {number} skill */
	function colorSkill(skill) {
		return `hsl(${Math.round(((skill - 0.5) / 0.8) * 130)}, 85%, 40%)`;
	}

	/** @param {number} motivation */
	function colorMotivation(motivation) {
		return `hsl(clamp(-10, ${Math.round(((motivation - 0.5) / 0.8) * 130)}, 140), 85%, 40%)`;
	}
</script>

<section>
	<div class="tab-selector">
		<button
			class={menuState[0] == 'all' ? '' : 'secondary'}
			onclick={() => {
				menuState[0] = 'all';
			}}>All</button
		>
		<button
			class={menuState[0] == 'production' ? '' : 'secondary'}
			onclick={() => {
				menuState[0] = 'production';
			}}>Production</button
		>
		<button
			class={menuState[0] == 'marketing' ? '' : 'secondary'}
			onclick={() => {
				menuState[0] = 'marketing';
			}}>Marketing</button
		>
	</div>
	<div class="tab-selector secondary">
		<button
			class="outline {menuState[1] == 'hired' ? '' : 'secondary'}"
			onclick={() => {
				menuState[1] = 'hired';
			}}>Hired</button
		>
		<button
			class="outline {menuState[1] == 'prospective' ? '' : 'secondary'}"
			onclick={() => {
				menuState[1] = 'prospective';
			}}>Prospective</button
		>
	</div>
	<div class="employee-grid">
		<table>
			<thead>
				<tr>
					<th> Name </th>
					<th> Skill </th>
					{#if menuState[1] == 'hired'}
						<th> Motivation </th>
						<th> Working hours </th>
						<th> Salary </th>
					{:else}
						<th> Salary expectation </th>
					{/if}
					<th> </th>
				</tr>
			</thead>
			<tbody>
				{#if menuState[0] == 'production'}
					{#if menuState[1] == 'hired'}
						{@render displayProductionHired()}
					{:else}
						{@render displayProductionNotHired()}
					{/if}
				{:else if menuState[0] == 'marketing'}
					{#if menuState[1] == 'hired'}
						{@render displayMarketingHired()}
					{:else}
						{@render displayMarketingNotHired()}
					{/if}
				{:else if menuState[0] == 'all'}
					{#if menuState[1] == 'hired'}
						{@render displayProductionHired()}
						{@render displayMarketingHired()}
					{:else}
						{@render displayProductionNotHired()}
						{@render displayMarketingNotHired()}
					{/if}
				{/if}
			</tbody>
			<tfoot>
				<tr>
					<td>
						<strong>Count: {employeeArray.length}</strong>
					</td>
					<td>
						<progress
							style="--pico-progress-color: {colorSkill(
								average(employeeArray, (e) => {
									return e.Skill;
								})
							)};"
							value={average(employeeArray, (e) => {
								return e.Skill;
							}) - 0.5}
							max="1"
						></progress>
					</td>

					{#if menuState[1] == 'hired'}
						<td>
							<progress
								style="--pico-progress-color: {colorMotivation(
									average(employeeArray, (e) => {
										return e.Motivation;
									})
								)};"
								value={average(employeeArray, (e) => {
									return e.Motivation;
								}) - 0.5}
								max="1"
							></progress>
						</td>

						<td>
							<span style="padding-left: 1rem;">
								{format.number(
									average(employeeArray, (e) => {
										return e.WorkingHours;
									}),
									false,
									1
								)}
							</span>
						</td>
					{/if}
					<td>
						<span style="padding-left: 1rem;">
							{format.currency(
								average(employeeArray, (e) => {
									return e.Pay;
								}),
								false,
								0
							)} / Mon
						</span>
					</td>
					<td></td>
				</tr>
			</tfoot>
		</table>
	</div>
</section>

{#snippet employee(
	/** @type {import("$lib/javascript/simulation").Employee}*/ employee,
	/** @type {import("svelte").Snippet<[]>} */ action,
	/** @type {boolean} */ isEmployee,
	/** @type {string} */ type
)}
	<!--<tr style="opacity: {isFired ? '0.5' : '1'};">-->
	<tr>
		<td>
			{employee.Name}
		</td>
		<td style="min-width: 10rem;">
			<progress
				style="--pico-progress-color: {colorSkill(employee.Skill)};"
				value={employee.Skill - 0.5}
				max="1"
			></progress>
		</td>

		{#if isEmployee}
			<td style="min-width: 10rem;">
				{#if employee.Motivation <= 0}
					<span style="color: red;">BURNT OUT</span>
				{:else}
					<progress
						style="--pico-progress-color: {colorMotivation(employee.Motivation)}"
						value={employee.Motivation - 0.5}
						max="1"
					></progress>
				{/if}
			</td>

			<td>
				<input
					type="number"
					bind:value={employee.WorkingHours}
					onchange={() => {
						if (employee.WorkingHours > 10) employee.WorkingHours = 10;
						else if (employee.WorkingHours < 1) employee.WorkingHours = 1;
						modifyEmployee(delta.Delta_Change, employee, type);
						syncDecisionsAndEmployees();
						updateDecisions(clientState.Decisions);
					}}
					style="margin-bottom: 0;"
				/>
			</td>
		{/if}
		<td>
			{format.currency(employee.Pay, false, 0)} / Mon
		</td>
		<td>
			{@render action()}
		</td>
	</tr>
{/snippet}

{#snippet displayMarketingHired()}
	{#each clientState.Employees.marketing as e, index (e.ID)}
		{#snippet action()}
			<button
				onclick={() => {
					modifyEmployee(delta.Delta_Remove, e, 'marketing');
				}}
			>
				Fire
			</button>
		{/snippet}
		{@render employee(e, action, isEmployee(e, clientState.Employees.marketing), 'marketing')}
	{/each}
{/snippet}

{#snippet displayMarketingNotHired()}
	{#each clientState.Unemployed.marketing.toSorted(sortBySkill) as e, index (e.ID)}
		{#snippet action()}
			<button
				onclick={() => {
					modifyEmployee(delta.Delta_New, e, 'marketing');
				}}>Hire</button
			>
		{/snippet}
		{@render employee(e, action, isEmployee(e, clientState.Employees.marketing), 'marketing')}
	{/each}
{/snippet}

{#snippet displayProductionHired()}
	{#each clientState.Employees.production as e, index (e.ID)}
		{#snippet action()}
			<button
				onclick={() => {
					modifyEmployee(delta.Delta_Remove, e, 'production');
				}}
			>
				Fire
			</button>
		{/snippet}
		{@render employee(e, action, isEmployee(e, clientState.Employees.production), 'production')}
	{/each}
{/snippet}

{#snippet displayProductionNotHired()}
	{#each clientState.Unemployed.production.toSorted(sortBySkill) as e, index (e.ID)}
		{#snippet hireOrCancel()}
			<button
				onclick={() => {
					modifyEmployee(delta.Delta_New, e, 'production');
				}}>Hire</button
			>
		{/snippet}
		{@render employee(
			e,
			hireOrCancel,
			isEmployee(e, clientState.Employees.production),
			'production'
		)}
	{/each}
{/snippet}

<style>
	section {
		min-width: 400px;
	}

	.tab-selector {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
		gap: var(--spacing);
		margin-bottom: 0;
	}

	.tab-selector.secondary {
		grid-template-columns: 1fr 1fr;
	}
</style>
