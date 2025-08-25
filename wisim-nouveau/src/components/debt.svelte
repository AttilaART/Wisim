<script>
	import '$lib/javascript/format';
	import { format } from '$lib/javascript/format';
	/** @type {{clientState: import("$lib/javascript/simulation").clientState, updateDecisions: (decisions: import("$lib/javascript/simulation").Decisions)=>void}} */
	let { clientState = $bindable(), updateDecisions } = $props();

	/** @type {number}*/
	let newAmount = $state(clientState.decisions.Finances.Set_bank_loan);
	let creditLimit = $state(10000000); //TODO: Implement Credit limit
</script>

<form
	onsubmit={() => {
		clientState.company.Balance =
			clientState.company.Balance - clientState.decisions.Finances.Set_bank_loan + newAmount;
		clientState.decisions.Finances.Set_bank_loan = newAmount;
		updateDecisions(clientState.decisions);
	}}
>
	<label for="debt">
		<h2>Bank Loan</h2>
		<input
			id="debt"
			bind:value={newAmount}
			type="range"
			min="0"
			max={creditLimit - clientState.company.Bridge_loans}
		/>
	</label>
	<div class="grid">
		<div>
			<p style="margin-bottom: calc(var(--pico-spacing) / 2)">Loan Amount</p>
			<h3 style="margin-top: 0">{format.currency(newAmount, false, 0)}</h3>
		</div>
		{#if clientState.company.Bridge_loans > 0}
			<div>
				<p style="margin-bottom: calc(var(--pico-spacing) / 2)">Bridge Loan</p>
				<h3 style="margin-top: 0; color: red">
					{format.currency(clientState.company.Bridge_loans, false, 0)}
				</h3>
			</div>
		{/if}
	</div>
	<table>
		<tbody>
			<tr>
				<td> Current balance: </td>
				<td style="text-align: right;">
					{format.currency(clientState.company.Balance, true, 0)}
				</td>
				<td></td>
			</tr>
			<tr>
				<td> New balance </td>
				<td style="text-align: right;">
					{format.currency(
						clientState.company.Balance - clientState.decisions.Finances.Set_bank_loan + newAmount,
						true,
						0
					)}
				</td>
				<td style="text-align: left;">
					<small
						>{format.currency(
							newAmount - clientState.decisions.Finances.Set_bank_loan,
							true,
							0
						)}</small
					>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="grid">
		<input
			type="button"
			value="Cancel"
			class="secondary"
			onclick={() => {
				newAmount = clientState.decisions.Finances.Set_bank_loan;
			}}
		/>
		<input type="submit" value="Confirm Loan" />
	</div>
</form>

<style>
	form {
		min-width: 30rem;
	}
</style>
