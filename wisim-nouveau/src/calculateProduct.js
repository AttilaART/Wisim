// NOTE: This may get out of sync with the Go function of the same name
/**
 * @param {import("$lib/javascript/simulation").Product} product
 * @param {import("$lib/javascript/simulation").ProductComponents} productComponents
 * @returns {{productStats: import("$lib/javascript/simulation").ProductStats, productionLineCost:number}}
 */
// Calcualtes base product stats
export function calculateProductStats(product, productComponents) {
	// Calcualtes product stats without side effects
	/** @type {import("$lib/javascript/simulation").ProductStats} */
	let productStats = {
		ProductionCost: 0,
		MaterialUse: 0,

		Quality: 0,
		Ecology: 0,
		Ethics: 0,
		Durability: 0
	};

	let productionLineCost = 0;

	/** @type {import("$lib/javascript/simulation").Component[]} */
	let parts = [];

	product.Components.FormFactor
		? parts.push(productComponents.FormFactor[product.Components.FormFactor])
		: undefined;
	product.Components.Frame
		? parts.push(productComponents.FormFactor[product.Components.Frame])
		: undefined;
	product.Components.Body
		? parts.push(productComponents.FormFactor[product.Components.Body])
		: undefined;
	product.Components.Mechanism
		? parts.push(productComponents.FormFactor[product.Components.Mechanism])
		: undefined;

	for (let component of product.Components.Misc) {
		parts.push(productComponents.Misc[component]);
	}

	for (let part of parts) {
		productStats.ProductionCost += part.ProductionCost;
		productStats.MaterialUse += part.MaterialUse;

		productStats.Quality += part.Quality;
		productStats.Ecology += part.Ecology;
		productStats.Ethics += part.Ethics;
		productStats.Durability += part.Durability;
		productionLineCost += part.ProductionLineCost;
	}

	productStats.Durability = Math.round(productStats.Durability);

	return { productStats, productionLineCost };
}
