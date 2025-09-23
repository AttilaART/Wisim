// NOTE: This may get out of sync with the Go function of the same name
/**
 * @param {import("$lib/javascript/simulation").Company} c
 * @param {import("$lib/javascript/simulation").Product} product
 * @param {import("$lib/javascript/simulation").Decisions_product} productDecisions
 * @param {import("$lib/javascript/simulation").Employee[]} productionEmployees
 * @returns {import("$lib/javascript/simulation").Product}
 */
export function Calculate_product(c, product, productDecisions, productionEmployees) {
	console.log('calculating product stats');
	// Calcualtes product stats without side effects
	c = JSON.parse(JSON.stringify(c));

	product.Name = productDecisions.Name;

	product = calcualteMaterialUse(
		product,
		c.Tech.Ecology,
		productDecisions.Manufacturing.MaterialEfficiency,
		productDecisions.Materials.Quality,
		c.Tech
	);

	product = calculate_quality(
		product,
		avr_skill(productionEmployees),
		productDecisions.Materials.Quality,
		productDecisions.Manufacturing.Quality,
		c.Tech
	);

	product = calculate_durability(
		product,
		productDecisions.Manufacturing.MaxDurability,
		productDecisions.Manufacturing.Durability,
		c.Tech
	);
	product = calculate_ecology(
		product,
		productDecisions.Manufacturing.EcologicalEnergy,
		productDecisions.Materials.Ecology,
		c.Tech
	);
	product = calculate_production_cost(product, productDecisions.Manufacturing, c.Tech);
	product.Ethics = 1;
	product.Weight = 1;

	console.log(product);

	return product;
}

/**
 * @param {import("$lib/javascript/simulation").Employee[]} employees
 * @returns {number}
 */
function avr_skill(employees) {
	if (employees.length <= 0) return 0;

	let totalSkill = 0;
	for (let e of employees) {
		totalSkill += e.Skill;
	}
	return totalSkill / employees.length;
}

// Product attribute functions

/**
 * Calculates Material_use of product w/ side effects and returns value
 * @param {import("$lib/javascript/simulation").Product} product
 * @param {number} ecology_research
 * @param {number} material_efficiency
 * @param {number} quality
 * @param {import("$lib/javascript/simulation").TechLevels} techLevels
 */
function calcualteMaterialUse(product, ecology_research, material_efficiency, quality, techLevels) {
	techLevels.MaterialUse *= 1 + ecology_research / 2000;
	product.MaterialUse =
		(techLevels.MaterialUse * quality * 0.25) / Math.sqrt(0.1 * Math.max(material_efficiency, 0.1));

	return product;
}

/**
 * Calculates Ecology_factor of product w/ side effects and returns value
 * @param {import("$lib/javascript/simulation").Product} product
 * @param {number} percentageOfEcologicalEnergy
 * @param {number} materialEcology
 * @param {import("$lib/javascript/simulation").TechLevels} techLevels
 */
function calculate_ecology(product, percentageOfEcologicalEnergy, materialEcology, techLevels) {
	product.Ecology =
		2 *
		Math.log(
			techLevels.Ecology * (materialEcology * product.MaterialUse) +
				percentageOfEcologicalEnergy / 20 +
				product.Durabilty / 5
		);
	product.Ecology = Math.min(product.Ecology, Number.MAX_VALUE);
	product.Ecology = Math.max(product.Ecology, 0);

	return product;
}

/** Calculates Quality_factor of product w/ side effects and returns value
 * @param {import("$lib/javascript/simulation").Product} product
 * @param {number} production_skill
 * @param {number} material_quality
 * @param {number} manufacturing_quality
 * @param {import("$lib/javascript/simulation").TechLevels} techLevels
 */

function calculate_quality(
	product,
	production_skill,
	material_quality,
	manufacturing_quality,
	techLevels
) {
	product.Quality =
		2 * Math.sqrt(material_quality * techLevels.Quality * manufacturing_quality) +
		Math.sqrt(production_skill);
	product.Quality = Math.min(product.Quality, Number.MAX_VALUE);
	product.Quality = Math.max(product.Quality, 0);

	product.Weight = 1;

	return product;
}

/**
 * Calculates Durabilty of product w/ side effects and returns value
 * @param {import("$lib/javascript/simulation").Product} product
 * @param {number} max_durability
 * @param {number} durability_focus
 * @param {import("$lib/javascript/simulation").TechLevels} techLevels
 */
function calculate_durability(product, max_durability, durability_focus, techLevels) {
	product.Durabilty = Math.round(techLevels.Durability * product.Quality + durability_focus * 2);

	if (product.Durabilty > max_durability) {
		product.Durabilty = max_durability;
	}

	product.Durabilty = Math.min(product.Durabilty, Number.MAX_SAFE_INTEGER);

	return product;
}

/**
 * @param {import("$lib/javascript/simulation").Product} product
 * @param {{Quality: number, EcologicalEnergy: number, MaterialEfficiency: number, Durability: number, MaxDurability: number}} manufacturing
 * @param {import("$lib/javascript/simulation").TechLevels} techLevels
 */

function calculate_production_cost(product, manufacturing, techLevels) {
	product.ProductionCost =
		techLevels.ProductionCost *
			(Math.pow(1.1, manufacturing.Quality) * 0.1 +
				Math.pow(1.1, manufacturing.MaterialEfficiency) * 0.01) +
		manufacturing.Durability / 10;
	product.ProductionCost = Math.max(product.ProductionCost, 0.5);

	return product;
}
