package simulation

import (
	"fmt"
	"math"
)

func (c *Company) Calculate_product(
	product_decisions Decisions_product, research Decisions_research,
) (product Product, error error) { // Calcualtes product stats without side effects
	product = c.Offer.Product

	product.calcualte_material_use(research.Ecology, product_decisions.Manufacturing.Material_efficiency, product_decisions.Materials.Quality)
	product.calculate_quality(research.Quality, c.employee_pool.Get_avr_skill(c.Id, Employee_type_production), product_decisions.Materials.Quality, product_decisions.Manufacturing.Quality)
	product.calculate_durability(research.Durability, product_decisions.Manufacturing.Max_durability, product_decisions.Manufacturing.Durability)
	product.calculate_ecology(research.Ecology, product_decisions.Manufacturing.Ecological_energy, product_decisions.Materials.Ecology)
	product.calculate_production_cost(research.Production_cost, product_decisions.Manufacturing)
	product.Ethics_factor = 1

	err := check_product(product)
	if err != nil {
		error = fmt.Errorf("%w /n Product: %#+v", err, product)
	}
	return product, error
}

// offer functions
func promotion_quality(base_marketing_strength float32, marketing_personelle []*Employee) float32 {
	// Temporary method
	var total_personelle_strength float32 = 1.0
	for _, e := range marketing_personelle {
		total_personelle_strength += e.Motivation * e.Skill * (e.Working_hours / 8.0)
	}

	return base_marketing_strength *
		(total_personelle_strength / float32(len(marketing_personelle))) *
		float32((1 + math.Log(float64(len(marketing_personelle)))))
}

// Product attribute functions

// Calculates Material_use of product w/ side effects and returns value
func (product *Product) calcualte_material_use(ecology_research float32, material_efficiency float32, quality float32) float32 {
	product.Base_material_use *= 1 + ecology_research/2000
	product.Material_use = product.Base_material_use / material_efficiency * quality
	return product.Material_use
}

// Calculates Ecology_factor of product w/ side effects and returns value
func (product *Product) calculate_ecology(
	ecology_research float32,
	percentage_of_ecological_energy float32,
	material_ecology float32,
) float32 {
	product.Base_ecology *= 1 + ecology_research/1000
	product.Ecology_factor = product.Base_ecology*(material_ecology*product.Material_use) + percentage_of_ecological_energy/20 + float32(product.Durabilty)/5
	return product.Ecology_factor
}

// Calculates Quality_factor of product w/ side effects and returns value
func (product *Product) calculate_quality(quality_research, production_skill, material_quality, manufacturing_quality float32) float32 {
	product.Base_quality *= 1 + quality_research/1000
	product.Quality_factor = material_quality * production_skill * product.Base_quality * manufacturing_quality
	return product.Quality_factor
}

// Calculates Durabilty of product w/ side effects and returns value
func (product *Product) calculate_durability(durability_research float32, max_durability int, durability_focus float32) int {
	product.Base_durability *= 1 + durability_research/1000
	product.Durabilty = int(math.Round(float64(product.Base_durability*product.Quality_factor + durability_focus*2)))

	if product.Durabilty > max_durability {
		product.Durabilty = max_durability
	}

	return product.Durabilty
}

func (product *Product) calculate_production_cost(
	production_speed_research float32,
	manufacturing struct {
		Quality             float32
		Ecological_energy   float32
		Material_efficiency float32
		Durability          float32
		Max_durability      int
	},
) (production_speed float32) {
	product.Base_production_cost *= 1 + production_speed_research/1000
	product.Production_cost = product.Base_production_cost*manufacturing.Quality + product.Base_production_cost*manufacturing.Material_efficiency*manufacturing.Durability/10
	return production_speed
}
