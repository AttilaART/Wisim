/** @param {HTMLElement} el */
export const preventPageReload = (el) => {
	el.addEventListener('submit', (/** @type {Event}*/ e) => {
		e.preventDefault();
	});
};

/** @param {()=>any} func */
export function ignoreError(func) {
	try {
		return func();
	} catch {}
}
