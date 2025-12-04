import { dev } from '$app/environment';
import fs from 'fs';
import os from 'os';

export function load({ url }) {
	return {
		serverAdress: url.searchParams.entries().toArray()[0][0]
	};
}

// we don't need any JS on this page, though we'll load
// it in dev so that we get hot module replacement
export const csr = true;
export const ssr = false;
