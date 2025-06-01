import { writable, Writable, get } from 'svelte/store';
import { simulation } from '../wailsjs/go/models';
import { Get_Decisions, Get_External_Factors } from '../wailsjs/go/main/App';

// hide_tabs: "after_hover || always || never"
export const preferences = $state({ hide_tabs: "after_hover" })

export const loading = writable(false)
export const error: Writable<null> | Writable<Error> = writable(null)

export const month = writable(1)

export const company_id = writable(0)
export const decisionsold = writable(new simulation.Decisionsold)

export const decisions: simulation.Decisions = $state(await Get_Decisions(0, 0))
export const external_factors: simulation.External_factors = $state(await Get_External_Factors(0, 0))

console.log($state.snapshot(external_factors))
console.log($state.snapshot(decisions))


export const canvas = writable({
  x: 0,
  y: 0,
})
