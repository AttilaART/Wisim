import { dev } from '$app/environment';
export const prerender = true;

export function load({ url }) {
	let serverAdress = url.searchParams.entries().toArray()[0][0];

	return {
		serverAdress: serverAdress
	};
}

// we don't need any JS on this page, though we'll load
// it in dev so that we get hot module replacement
export const csr = dev;
