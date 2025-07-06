package simulation

import "math"

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
func (product *Product) calculate_ecology(
	material_use float32,
	percentage_of_ecological_energy float32,
	material_ecology float32,
) {
	product.Ecology_factor = (product.Base_ecology * material_ecology * (material_use/1 +
		percentage_of_ecological_energy/10 +
		product.Quality_factor/10)) / 2
}

func (product *Product) calculate_quality(production_skill float32, material_quality float32, manufacturing_quality float32) {
	product.Quality_factor = material_quality * production_skill * product.Base_quality * manufacturing_quality
}

func (product *Product) calculate_durability(max_durability int, durability_focus float32) {
	product.Durabilty = int(math.Round(float64(clamp(
		float32(product.Base_durability)*product.Quality_factor*durability_focus,
		float32(max_durability),
	))))
}

func calculate_production_speed(
	manufacturing struct {
		Quality             float32
		Ecological_energy   float32
		Material_efficiency float32
		Durability          float32
		Max_durability      int
	}, base_production_speed float32,
) (production_speed float32) {
	production_speed = base_production_speed / manufacturing.Quality / manufacturing.Material_efficiency / manufacturing.Durability

	if math.IsInf(float64(production_speed), 0) {
		panic("production_speed is infinite")
	}
	return production_speed
}
