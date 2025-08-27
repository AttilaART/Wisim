<script>
	import '$lib/javascript/format';
	import { format } from '$lib/javascript/format';
	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void}} */
	let { clientState = $bindable(), updateDecisions } = $props();

	/**
	 * @typedef {Object} Delta
	 * @property {import("$lib/javascript/simulation").Change} Delta.Delta_New
	 * @property {import("$lib/javascript/simulation").Change} Delta.Delta_Change
	 * @property {import("$lib/javascript/simulation").Change} Delta.Delta_Remove
	 */
	/**
	 * @type {Delta}
	 */
	const delta = {
		Delta_New: 0,
		Delta_Change: 1,
		Delta_Remove: 2
	};
	/**
	 * @param {import("$lib/javascript/simulation").Change} delta
	 * @param {import("$lib/javascript/simulation").Employee} employee
	 * @param {string} employeeType
	 * @returns {void}
	 */
	function modifyEmployee(delta, employee, employeeType) {
		/** @type{import("$lib/javascript/simulation").Delta<import("$lib/javascript/simulation").Employee>} */
		let emplyeeDelta = {
			Change: delta,
			Item: employee
		};

		// check if employee is already in deltas

		let deltasList = [];

		if (employeeType == 'production')
			deltasList = clientState.decisions.Employees.Production_deltas;
		else if (employeeType == 'marketing')
			deltasList = clientState.decisions.Employees.Marketing_deltas;
		else throw 'invalid employeeType';

		for (/** @type {number} i */ let i in deltasList) {
			if (clientState.decisions.Employees.Production_deltas[i].Item.Id == employee.Id) {
				clientState.decisions.Employees.Production_deltas[i].Change = delta;
				clientState.decisions.Employees.Production_deltas[i].Item = employee;
				return;
			}
		}

		deltasList.push(emplyeeDelta);
	}
	/**
	 * @param {number} id
	 * @returns {boolean}
	 */
	function isFired(id) {
		for (let e of clientState.decisions.Employees.Production_deltas) {
			if (e.Item.Id == id) {
				if (e.Change == delta.Delta_Remove) {
					return true;
				}
			}
		}
		return false;
	}
</script>

<form
	onsubmit={() => {
		updateDecisions(clientState.decisions);
	}}
>
	<div class="employee-grid">
		{console.warn(clientState.employees.production)}
		{#each clientState.employees.production as e, index (e.Id)}
			{@render employee(index)}
		{/each}
	</div>
</form>

{#snippet employee(/** @type {Number} */ index)}
	<article style="opacity: {isFired(clientState.employees.production[index].Id) ? '0.5' : '1'};">
		<h3>{clientState.employees.production[index].Name}</h3>
		<label for=""
			>Skill
			<progress value={clientState.employees.production[index].Skill - 0.5} max="1"></progress>
		</label>
		<label for=""
			>Motivation
			<progress value={clientState.employees.production[index].Motivation - 0.5} max="1"></progress>
		</label>
		<form
			onchange={() => {
				modifyEmployee(delta.Delta_Change, clientState.employees.production[index], 'production');
				updateDecisions(clientState.decisions);
			}}
		>
			<label for="">
				Working Hours (Per day): {clientState.employees.production[index].Working_hours}h
				<input
					bind:value={clientState.employees.production[index].Working_hours}
					type="range"
					min="1"
					max="12"
				/>
			</label>
			<label for="">
				Monthly Salary
				<input
					bind:value={clientState.employees.production[index].Pay}
					type="number"
					min={clientState.external_factors.Production_minimum_wage}
					step="1000"
				/>
			</label>
		</form>
		{#if !isFired(clientState.employees.production[index].Id)}
			<button
				onclick={() => {
					modifyEmployee(delta.Delta_Remove, clientState.employees.production[index], 'production');
				}}>FIRE</button
			>
		{:else}
			<button
				onclick={() => {
					modifyEmployee(delta.Delta_Change, clientState.employees.production[index], 'production');
				}}>Rehire</button
			>
		{/if}
	</article>
{/snippet}

<style>
	.employee-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--pico-spacing);
	}
</style>
