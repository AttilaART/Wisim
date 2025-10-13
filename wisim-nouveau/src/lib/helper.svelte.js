/** @param {HTMLElement} el */
export const preventPageReload = (el) => {
	el.addEventListener('submit', (/** @type {Event}*/ e) => {
		e.preventDefault();
	});
};
