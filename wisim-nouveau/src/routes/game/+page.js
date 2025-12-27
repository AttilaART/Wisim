export function load({ url }) {
	let serverAdress = url.searchParams.entries().toArray()[0][0];

	if (serverAdress.endsWith('/')) {
		serverAdress = serverAdress.slice(0, serverAdress.length - 1);
	}

	return {
		serverAdress: serverAdress
	};
}

export const csr = true;
export const ssr = false;
