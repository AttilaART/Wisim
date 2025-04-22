import { writable } from 'svelte/store';

const month_counter = writable(1)
const financial_data_store = writable({ report: undefined, income_statement: undefined, assets: undefined, liabilities: undefined })

export { month_counter, financial_data_store }
