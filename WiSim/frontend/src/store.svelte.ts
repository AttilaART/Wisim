import { writable, Writable } from 'svelte/store';
import { simulation } from '../wailsjs/go/models';

export const loading = writable(false)
export const error: Writable<null> | Writable<Error> = writable(null)

export const month = writable(1)

export const company_id = writable(0)
export const decisionsold = writable(new simulation.Decisionsold)


export const canvas = writable({
  x: 0,
  y: 0,
})
