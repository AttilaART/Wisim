import { writable } from 'svelte/store';

const month_counter = writable(1)
const loading = writable(false)

export { month_counter, loading }
