import { writable } from 'svelte/store';
import { simulation } from '../wailsjs/go/models';

const month_counter = writable(1)
const company = writable(0)
const loading = writable(false)
const decisions = writable(new simulation.Decisions)

export { month_counter, loading, company, decisions }
