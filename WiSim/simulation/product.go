package simulation

// offer functions
func marketing_strength(marketing_investment float32, marketing_personelle []Employee) float32 {
	// Temporary method
	var total_personelle_strength float32 = 1.0
	for _, e := range marketing_personelle {
		total_personelle_strength += e.Motivation * e.Skill * (e.Working_hours / 8.0)
	}

	return marketing_investment * float32(total_personelle_strength)
}

// Product attribute functions
func (product *Product) calculate_ecology(base_ecology float32, material_use float32, percentage_of_ecological_energy float32) {
	product.Ecology_factor = (base_ecology * (material_use/1 +
		percentage_of_ecological_energy/10 +
		product.Quality_factor/10)) / 2
}

func (product *Product) calculate_quality(quality_development_investment float32, production_skill float32, material_quality float32) {
	product.Material_quality = material_quality
	product.Skill_factor = production_skill
	product.Quality_development_factor += quality_development_investment / 1000000

	product.Quality_factor = product.Material_quality * product.Skill_factor * product.Quality_development_factor
}
