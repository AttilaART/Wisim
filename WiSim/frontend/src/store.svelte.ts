import { writable, Writable } from 'svelte/store';
import { simulation } from '../wailsjs/go/models';

const month = writable(1)
const company = writable(0)
const loading = writable(false)
const decisions = writable(new simulation.Decisions)
const error: Writable<null> | Writable<Error> = writable(null)

type window = {
  id: number,
  z_index: number,
}
const windows: window[] = $state([])

function new_window(): number {
  windows.push({ id: windows.length, z_index: windows.length, })
  return windows[windows.length - 1].id
}
function move_window_to_top(window_id: number) {
  let current_z_index = windows[window_id].z_index
  for (let i in windows) {
    if (windows[i].z_index > current_z_index) {
      windows[i].z_index -= 1
    }
  }
  windows[window_id].z_index = windows.length
}

export { month as month_counter, loading, company, decisions, windows, new_window, move_window_to_top, error }
