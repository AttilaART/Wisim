<script>
	import { calculateProduction } from '$lib/helper.svelte';
	import { format } from '$lib/javascript/format';
	import { assignmentPatterns, delta, financeReportCategories } from '$lib/javascript/simulation';

	/**
	 * @typedef {Object} props
	 * @property {import("$lib/javascript/simulation").clientState} clientState
	 * @property {(decisions: import("$lib/javascript/simulation").Decisions)=>void} updateDecisions
	 */
	/** @type {props} */
	let { clientState = $bindable(), updateDecisions } = $props();

	/** @param {import("$lib/javascript/simulation").Machine} machine */
	function newMachine(machine) {
		console.log(clientState);

		if (clientState.Company.Machines == null) {
			clientState.Company.Machines = [];
		}

		const newMachine = JSON.parse(JSON.stringify(machine));
		newMachine.ID = clientState.Company.Machines.length + 1;
		clientState.Company.Machines.push(newMachine);
		clientState.Company.Balance -= newMachine.Value;
		updateMachineDecision(newMachine, delta.Delta_New);
	}

	/** @param {import("$lib/javascript/simulation").Machine} machine */
	function sellMachine(machine) {
		// console.log(clientState);
		let index = clientState.Company.Machines.findIndex((m) => {
			return m.ID == machine.ID;
		});
		clientState.Company.Machines.splice(index, 1);
		clientState.Company.Balance += machine.Value;
		updateMachineDecision(machine, delta.Delta_Remove);
	}

	/**
	 * @param {import("$lib/javascript/simulation").Machine} machine
	 * @param {number} change
	 */
	function updateMachineDecision(machine, change) {
		let deltaIndex = -1;
		try {
			for (let i in clientState.Decisions.Production.Machines) {
				if (clientState.Decisions.Production.Machines[i].Item.ID == machine.ID) {
					deltaIndex = Number(i);
					break;
				}
			}
		} catch (e) {
			if (e instanceof TypeError) {
				clientState.Decisions.Production.Machines = [];
			} else {
				throw e;
			}
		}

		if (deltaIndex != -1) {
			if (clientState.Decisions.Production.Machines[deltaIndex].Change > change) {
				// Follow change higherarchy (new -> change -> remove)
				clientState.Decisions.Production.Machines[deltaIndex].Change = change;
			}
			clientState.Decisions.Production.Machines[deltaIndex].Item = machine;
		} else {
			clientState.Decisions.Production.Machines.push({ Change: change, Item: machine });
		}

		updateDecisions(clientState.Decisions);
	}

	let { MachineProduction, WorkerSurplus: workerSurplus } = $derived.by(() => {
		return calculateProduction(
			clientState.Company.ID,
			clientState.Company.Machines,
			clientState.Company.Offers,
			clientState.Employees.production,
			clientState.Decisions.Production.MachineAssignmentPattern
		);
	});

	/**
	 * @param {import("$lib/javascript/simulation").Machine[]} machines
	 * @returns {import("$lib/javascript/simulation").Machine[]}
	 */
	function removeDuplicates(machines) {
		/** @type {import("$lib/javascript/simulation").Machine[]} */
		let returnArray = [];
		for (let m of machines) {
			if (returnArray.findIndex((e) => e.ID == m.ID) == -1) {
				returnArray.push(m);
			}
		}

		return returnArray;
	}
</script>

<div>
	<article class="grid">
		<span style="line-height: 200%;">
			{#if workerSurplus > 0}
				You have <strong>{workerSurplus}</strong> production employees too many.
			{:else if workerSurplus < 0}
				You are missing <strong>{-workerSurplus}</strong> production employees!
			{/if}
		</span>
		<fieldset role="group" style="margin: 0;">
			<input type="text" value="Worker Distribution" disabled />
			<select
				disabled={clientState.predictionMode}
				bind:value={clientState.Decisions.Production.MachineAssignmentPattern}
				onchange={() => {
					updateDecisions(clientState.Decisions);
				}}
			>
				<option value={assignmentPatterns.fillMachines}>Fill Machines</option>
				<option value={assignmentPatterns.distributeWorkers}>Distribute Evenly</option>
			</select>
		</fieldset>
	</article>
</div>
<table style="max-height: calc(100vw - 100px);">
	<thead>
		<tr>
			<th>Name </th>
			<th>Production capacity</th>
			<th>Minimum Worker Count</th>
			<th>Optimal Worker Count</th>
			<th>Assigned Product</th>
			<th></th>
		</tr>
	</thead>
	<tbody>
		{#each removeDuplicates(clientState.Company.Machines) as m, i (m.ID)}
			<tr>
				<td>
					Machine {m.ID}
				</td>
				<td>
					<h2 style={MachineProduction[i] >= m.ProductionCapacity ? '' : 'color: red;'}>
						{MachineProduction[i]}
					</h2>
				</td>
				<td>
					<h2>{m.MinimumWorkers}</h2>
				</td>
				<td>
					<h2>{m.RequiredWorkers}</h2>
				</td>
				<td>
					<select
						bind:value={m.AssignedProductID}
						onchange={() => {
							updateMachineDecision(m, delta.Delta_Change);
						}}
					>
						{#if m.AssignedProductID && clientState.Company.Offers[m.AssignedProductID]}{clientState
								.Company.Offers[m.AssignedProductID].Product.Name}{:else}None{/if}
						{#each Object.entries(clientState.Company.Offers) as offer}
							<option value={offer[0]}>{offer[1].Product.Name}</option>
						{/each}
					</select>
				</td>
				<td>
					<button
						disabled={clientState.predictionMode}
						onclick={() => {
							sellMachine(m);
						}}
					>
						Sell
					</button>
				</td>
			</tr>
		{/each}
	</tbody>
</table>

<center>
	<button
		disabled={clientState.predictionMode}
		onclick={() => {
			newMachine(clientState.ExternalFactors.MachineOnOffer);
		}}
		>Buy Machines ({format.currency(
			clientState.ExternalFactors.MachineOnOffer.Value,
			false,
			2
		)})</button
	>
</center>
