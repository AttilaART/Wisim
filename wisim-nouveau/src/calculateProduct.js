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

	parts.push(productComponents.FormFactor[`${product.Components.FormFactor}`]);
	parts.push(productComponents.Frame[`${product.Components.Frame}`]);
	parts.push(productComponents.Body[`${product.Components.Body}`]);
	parts.push(productComponents.Mechanism[`${product.Components.Mechanism}`]);

	// console.log(product.Components.Mechanism);
	// console.log(productComponents.Mechanism[`${product.Components.Mechanism}`]);

	for (let component of product.Components.Misc) {
		parts.push(productComponents.Misc[`${component}`]);
	}

	// console.log(parts);
	for (let part of parts) {
		if (part === undefined) {
			continue;
		}
		productStats.ProductionCost += part.ProductionCost;
		productStats.MaterialUse += part.MaterialUse;

		productStats.Quality += part.Quality;
		productStats.Ecology += part.Ecology;
		productStats.Ethics += part.Ethics;
		productStats.Durability += part.Durability;
		productionLineCost += part.ProductionLineCost;
	}

	productStats.Durability += product.ExtraDurability;
	productStats.ProductionCost += 5 * product.ExtraDurability;

	productStats.Quality += product.ExtraQuality;
	productStats.ProductionCost += 5 * product.ExtraQuality;

	productStats.Quality += product.MaterialQuality;
	productStats.MaterialUse += 5 * product.MaterialQuality;

	productStats.Durability = Math.round(productStats.Durability);
	// console.log(productStats);

	return { productStats, productionLineCost };
}
