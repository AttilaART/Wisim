import { writable, Writable, get } from 'svelte/store';
import { simulation } from '../wailsjs/go/models';
import { newConnection } from './api.svelte';

// hide_tabs: "after_hover || always || never"
export const preferences = $state({ hide_tabs: "after_hover" })

export let connection = newConnection("ws://localhost:8000")

export const loading = writable(false)
export const error: Writable<null> | Writable<Error> = writable(null)

export const company_id: Writable<number> = writable(0)

export const current_decisions: simulation.Decisions = $state({})
export const company: simulation.Company = $state(new simulation.Company)
export const latest_reports: simulation.Report = $state({})
export const external_factors: simulation.External_factors = $state(new simulation.External_factors)

console.log($state.snapshot(external_factors))
console.log($state.snapshot(current_decisions))

export function update_external_factors(new_value: simulation.External_factors) {
  for (let field of Object.keys(new_value)) {
    external_factors[field] = new_value[field]
  }
}

export function update_decisions(new_value: simulation.Decisions) {
  for (let field of Object.keys(new_value)) {
    current_decisions[field] = new_value[field]
  }
}

export function update_reports(new_value: simulation.Report) {
  for (let field of Object.keys(new_value)) {
    latest_reports[field] = new_value[field]
  }
}

export function update_company(new_value: simulation.Company) {
  for (let field of Object.keys(new_value)) {
    company[field] = new_value[field]
  }
}


export const canvas = writable({
  x: 0,
  y: 0,
})
