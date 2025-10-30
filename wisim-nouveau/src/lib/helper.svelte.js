import * as echarts from 'echarts';

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

/**
 * @param {number} companyID
 * @param {Array.<import("$lib/javascript/simulation").Machine>} machines
 * @param {Object.<string, import("$lib/javascript/simulation").Offer>} offers
 * @param {Array.<import("$lib/javascript/simulation").Employee>} employees
 * @param {import("$lib/javascript/simulation").AssignmentPattern} assignmentPattern
 * @returns {{WorkerSurplus: number, MachineProduction: number[]}}
 */

export function calculateProduction(companyID, machines, offers, employees, assignmentPattern) {
	/** @type {{WorkerSurplus: number, MachineProduction: number[]}}*/

	let result = JSON.parse(
		// @ts-ignore
		CalculateProductionGo(
			companyID,
			JSON.stringify(machines),
			JSON.stringify(offers),
			JSON.stringify(employees),
			assignmentPattern
		)
	);

	return result;
}

/**
 * @param {HTMLElement} element
 * @param {echarts.EChartsCoreOption} options
 */
export function chart(element, options) {
	let chart = echarts.init(element, 'dark');
	chart.setOption(options);
}
