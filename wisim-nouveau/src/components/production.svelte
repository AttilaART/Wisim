<script>
	import { format } from '$lib/javascript/format';
	import { delta, financeReportCategories } from '$lib/javascript/simulation';

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
		machine = JSON.parse(JSON.stringify(machine));
		machine.ID = clientState.Company.Machines.length;
		clientState.Company.Machines.push(machine);
		clientState.Company.Balance -= machine.Value;
		updateMachineDecision(machine, delta.Delta_New);
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
</script>

<table>
	<thead>
		<tr>
			<th>Name </th>
			<th> Production capacity</th>
			<th>Required Workers</th>
			<th>Assigned Product</th>
		</tr>
	</thead>
	<tbody>
		{#each clientState.Company.Machines as m}
			<tr>
				<td>
					Machine {m.ID}
				</td>
				<td>
					<h2>{m.ProductionCapacity}</h2>
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
			</tr>
		{/each}
	</tbody>
</table>

<center>
	<button
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
