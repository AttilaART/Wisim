export function load({ url }) {
	let serverAdress = url.searchParams.entries().toArray()[0][0];

	return {
		serverAdress: serverAdress
	};
}

export const csr = true;
export const ssr = false;
