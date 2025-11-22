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
 * @returns {{WorkerSurplus: number, MachineProduction: number[], MachineWorkerCount: number[]}}
 */

export function calculateProduction(companyID, machines, offers, employees, assignmentPattern) {
	/** @type {{WorkerSurplus: number, MachineProduction: number[], MachineWorkerCount: number[]}}*/

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
 * @param {import("$lib/javascript/simulation").Company} company
 * @param {import("$lib/javascript/simulation").Decisions} decisions
 * @param {import("$lib/javascript/simulation").ExternalFactors} externalFactors
 * @param {{Employees: {marketing: import("$lib/javascript/simulation").Employee[], production: import("$lib/javascript/simulation").Employee[]}, Unemployed: {marketing: import("$lib/javascript/simulation").Employee[], production: import("$lib/javascript/simulation").Employee[]}}} employees
 * @param {()=>void} onEnd
 * @param {number} steps
 * @param {import("$lib/javascript/simulation").ProductComponents} productComponents
 * @returns {import("$lib/javascript/simulation").Company}
 */
export function simulateMockStep(
	company,
	decisions,
	externalFactors,
	employees,
	onEnd,
	steps,
	productComponents
) {
	/** @type {import("$lib/javascript/simulation").Employee[]} */
	let employePool = [];

	employees.Employees.marketing.forEach((e) => {
		e.Employer = company.ID;
	});
	employees.Employees.production.forEach((e) => {
		e.Employer = company.ID;
	});

	employees.Unemployed.marketing.forEach((e) => {
		e.Employer = -1;
	});

	employees.Unemployed.production.forEach((e) => {
		e.Employer = -1;
	});

	employePool = employePool.concat(
		employees.Employees.marketing,
		employees.Employees.production,
		employees.Unemployed.marketing,
		employees.Unemployed.production
	);

	if (typeof steps != typeof 1) {
		steps = 1;
	}

	// @ts-ignore
	let result = SimulateMockStep(
		JSON.stringify(company),
		JSON.stringify(decisions),
		JSON.stringify(externalFactors),
		JSON.stringify(employePool),
		steps,
		JSON.stringify(productComponents)
	);

	onEnd();

	return JSON.parse(result);
}

/**
 * @param {import("$lib/javascript/simulation").Company} company
 * @param {import("$lib/javascript/simulation").Decisions} decisions
 * @param {import("$lib/javascript/simulation").ProductComponents} productComponents
 * @returns {import("$lib/javascript/simulation").Company}
 */
export function syncCompanyWithDecisions(company, decisions, productComponents) {
	// @ts-ignore
	let result = SyncCompanyWithDecisions(
		JSON.stringify(company),
		JSON.stringify(decisions),
		JSON.stringify(productComponents)
	);

	return JSON.parse(result);
}

/**
 * @param {HTMLElement} element
 * @param {echarts.EChartsCoreOption} options
 */
export function chart(element, options) {
	let chart = echarts.init(element, 'dark');
	chart.setOption(options);
}

/**
 * @template T
 * @param {T[]} elements
 * @param {(e: T)=>number} getValue
 * @returns {number}
 */
export function average(elements, getValue) {
	if (elements.length == 0) {
		return 0;
	}

	let total = 0;

	elements.forEach((e) => {
		total += getValue(e);
	});

	return total / elements.length;
}
