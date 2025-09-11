<script>
	import { format } from '$lib/javascript/format';
	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void}} */
	let { clientState = $bindable(), updateDecisions } = $props();

	let menuState = $state(['all', 'hired']);
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
		else throw `invalid employeeType: ${employeeType}`;

		for (/** @type {number} i */ let i in deltasList) {
			if (deltasList[i].Item.Id == employee.Id) {
				deltasList[i].Change = delta;
				deltasList[i].Item = employee;
				return;
			}
		}

		deltasList.push(emplyeeDelta);
	}
	/**
	 * @param {number} id
	 * @param {string} type
	 * @returns {boolean}
	 */
	function isFired(id, type) {
		for (let e of clientState.decisions.Employees[format.capitaliseFirstLetter(type) + '_deltas']) {
			if (e.Item.Id == id) {
				if (e.Change == delta.Delta_Remove) {
					return true;
				}
			}
		}
		return false;
	}

	/**
	 * @param {number} id
	 * @param {string} type
	 * @returns {boolean}
	 */
	function isHired(id, type) {
		for (let e of clientState.decisions.Employees[format.capitaliseFirstLetter(type) + '_deltas']) {
			if (e.Item.Id == id) {
				if (e.Change == delta.Delta_New) {
					return true;
				}
			}
		}
		return false;
	}

	/**
	 * @param {string} type
	 * @param {number} index
	 */
	function onModifyEmployee(type, index) {
		modifyEmployee(delta.Delta_Change, clientState.employees[type][index], type);
		console.log($state.snapshot(clientState.decisions));
		updateDecisions(clientState.decisions);
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
						{#each clientState.employees.production as e, index (e.Id)}
							{@render employee(index, 'production')}
						{/each}
					{:else}
						{#each clientState.unemployed.production as e, index (e.Id)}
							{@render prospective(index, 'production')}
						{/each}
					{/if}
				{:else if menuState[0] == 'marketing'}
					{#if menuState[1] == 'hired'}
						{#each clientState.employees.marketing as e, index (e.Id)}
							{@render employee(index, 'marketing')}
						{/each}
					{:else}
						{#each clientState.unemployed.marketing as e, index (e.Id)}
							{@render prospective(index, 'marketing')}
						{/each}
					{/if}
				{:else if menuState[0] == 'all'}
					{#if menuState[1] == 'hired'}
						{#each clientState.employees.marketing as e, index (e.Id)}
							{@render employee(index, 'marketing')}
						{/each}
						{#each clientState.employees.production as e, index (e.Id)}
							{@render employee(index, 'production')}
						{/each}
					{:else}
						{#each clientState.unemployed.marketing as e, index (e.Id)}
							{@render prospective(index, 'marketing')}
						{/each}
						{#each clientState.unemployed.production as e, index (e.Id)}
							{@render prospective(index, 'production')}
						{/each}
					{/if}
				{/if}
			</tbody>
		</table>
	</div>
</section>

{#snippet employee(/** @type {Number} */ index, /** @type {String} */ type)}
	<tr style="opacity: {isFired(clientState.employees[type][index].Id, type) ? '0.5' : '1'};">
		<td>
			{clientState.employees[type][index].Name}
		</td>
		<td>
			<progress value={clientState.employees[type][index].Skill - 0.5} max="1"></progress>
		</td>
		<td>
			<progress value={clientState.employees[type][index].Motivation - 0.5} max="1"></progress>
		</td>
		<td>
			<div style="display: inline-block;">
				{clientState.employees[type][index].Working_hours}h
			</div>
			<input
				style="display: inline-block;"
				bind:value={clientState.employees[type][index].Working_hours}
				type="range"
				min="1"
				max="12"
				onchange={() => {
					onModifyEmployee(type, index);
				}}
			/>
		</td>
		<td>
			<label
				for=""
				onchange={() => {
					onModifyEmployee(type, index);
				}}
			>
				<input
					bind:value={clientState.employees[type][index].Pay}
					type="number"
					min={clientState.external_factors[format.capitaliseFirstLetter(type) + '_minimum_wage']}
					step="1000"
				/>
			</label>
		</td>
		<td>
			{#if !isFired(clientState.employees[type][index].Id, type)}
				<button
					onclick={() => {
						modifyEmployee(delta.Delta_Remove, clientState.employees[type][index], type);
						updateDecisions(clientState.decisions);
					}}>FIRE</button
				>
			{:else}
				<button
					onclick={() => {
						modifyEmployee(delta.Delta_Change, clientState.employees[type][index], type);
						updateDecisions(clientState.decisions);
					}}>Rehire</button
				>
			{/if}
		</td>
	</tr>
{/snippet}

{#snippet prospective(/** @type {Number} */ index, /** @type {String} */ type)}
	<tr style="opacity: {isFired(clientState.employees[type][index].Id, type) ? '0.5' : '1'};">
		<td>
			{clientState.unemployed[type][index].Name}
		</td>
		<td>
			<progress value={clientState.unemployed[type][index].Skill - 0.5} max="1"></progress>
		</td>
		<td>
			{format.currency(clientState.unemployed[type][index].Pay, false, 0)} / Mon
		</td>
		<td>
			{#if !isHired(clientState.unemployed[type][index].Id, type)}
				<button
					onclick={() => {
						modifyEmployee(delta.Delta_New, clientState.unemployed[type][index], type);
						updateDecisions(clientState.decisions);
					}}>Send Offer</button
				>
			{:else}
				<button
					class="secondary"
					onclick={() => {
						modifyEmployee(delta.Delta_Remove, clientState.unemployed[type][index], type);
						updateDecisions(clientState.decisions);
					}}>Rescind Offer</button
				>
			{/if}
		</td>
	</tr>
{/snippet}

<style>
	section {
		min-width: 400px;
	}

	.tab-selector {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
		gap: var(--pico-spacing);
		margin-bottom: var(--pico-spacing);
	}

	.tab-selector.secondary {
		grid-template-columns: 1fr 1fr;
	}
</style>
