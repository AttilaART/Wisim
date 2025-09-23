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
		machine.ID = clientState.Decisions.Production.Machines.length;
		clientState.Decisions.Production.Machines.push({ Change: delta.Delta_New, Item: machine });
		clientState.Company.Machines.push(JSON.parse(JSON.stringify(machine)));
		clientState.Company.Balance -= machine.Value;
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
					<select>
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
