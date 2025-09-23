package simulation

import (
	"fmt"
	"log"
	"math"
)

func (c *Company) newProduct(productID, productName string, productDecisions Decisions_product) Offer {
	offer := Offer{Status: "current", Product: c.Calculate_product(Product{ID: productID, Name: productName}, productDecisions)}

	return offer
}

func (c Company) Calculate_product(
	product Product, productDecisions Decisions_product,
) Product { // Calcualtes product stats without side effects

	product.calcualteMaterialUse(c.Tech.Ecology,
		productDecisions.Manufacturing.MaterialEfficiency,
		productDecisions.Materials.Quality,
		c.Tech)

	product.calculateQuality(c.employeePool.Get_avr_skill(c.ID, Employee_type_production),
		productDecisions.Materials.Quality,
		productDecisions.Manufacturing.Quality, c.Tech)

	product.calculate_durability(productDecisions.Manufacturing.MaxDurability, productDecisions.Manufacturing.Durability, c.Tech)
	product.calculate_ecology(productDecisions.Manufacturing.EcologicalEnergy, productDecisions.Materials.Ecology, c.Tech)
	product.calculate_production_cost(productDecisions.Manufacturing, c.Tech)
	product.Ethics = 1
	product.Weight = 1

	errString := check_product(product)
	if errString != "" {
		log.Panicln(errString)
	}
	return product
}

// offer functions
func promotionQuality(employee_pool Employee_pool, baseMarketingStrength float32, marketingPersonelleIds []int) float32 {
	// Temporary method
	var totalPersonelleStrength float32 = 1.0
	for _, id := range marketingPersonelleIds {
		totalPersonelleStrength += employee_pool[id].Motivation * employee_pool[id].Skill * (employee_pool[id].WorkingHours / 8.0)
	}

	return baseMarketingStrength *
		(totalPersonelleStrength / float32(len(marketingPersonelleIds))) *
		float32(1+math.Log(float64(len(marketingPersonelleIds))))
}

// Product attribute functions

// Calculates Material_use of product w/ side effects and returns value
func (product *Product) calcualteMaterialUse(ecology_research float32, material_efficiency float32, quality float32, techLevels TechLevels) float32 {
	techLevels.MaterialUse *= 1 + ecology_research/2000
	product.MaterialUse = (techLevels.MaterialUse * quality * 0.25) / float32(math.Sqrt(float64(0.1*max(material_efficiency, 0.1))))
	return product.MaterialUse
}

// Calculates Ecology_factor of product w/ side effects and returns value
func (product *Product) calculate_ecology(
	percentageOfEcologicalEnergy float32,
	materialEcology float32,
	techLevels TechLevels,
) float32 {
	product.Ecology = float32(2 * math.Log(float64(techLevels.Ecology*(materialEcology*product.MaterialUse)+percentageOfEcologicalEnergy/20+float32(product.Durabilty)/5)))
	product.Ecology = min(product.Ecology, math.MaxFloat32)
	product.Ecology = max(product.Ecology, 0)
	return product.Ecology
}

// Calculates Quality_factor of product w/ side effects and returns value
func (product *Product) calculateQuality(productionSkill, materialQuality, manufacturingQuality float32, techLevels TechLevels) float32 {
	product.Quality = float32(2*math.Sqrt(float64(materialQuality*techLevels.Quality*manufacturingQuality))) + float32(math.Sqrt(float64(productionSkill)))
	product.Quality = min(product.Quality, math.MaxFloat32)
	product.Quality = max(product.Quality, 0)
	return product.Quality
}

// Calculates Durabilty of product w/ side effects and returns value
func (product *Product) calculate_durability(max_durability int, durability_focus float32, techLevels TechLevels) int {
	product.Durabilty = int(math.Round(float64(techLevels.Durability*product.Quality + durability_focus*2)))

	if product.Durabilty > max_durability {
		product.Durabilty = max_durability
	}

	product.Durabilty = min(product.Durabilty, math.MaxInt)

	return product.Durabilty
}

func (product *Product) calculate_production_cost(
	manufacturing struct {
		Quality            float32
		EcologicalEnergy   float32
		MaterialEfficiency float32
		Durability         float32
		MaxDurability      int
	}, techLevels TechLevels,
) (production_speed float32) {
	product.ProductionCost = techLevels.ProductionCost * (float32(exponential(1.1, float64(manufacturing.Quality), 0.1)) + float32(exponential(1.1, float64(manufacturing.MaterialEfficiency), 0.01)) + manufacturing.Durability/10)
	product.ProductionCost = max(product.ProductionCost, 0.5)
	return production_speed
}

func (c *Company) calculatePromotion(decisions Decisions) {
	c.BaseMarketingStrength += decisions.Research.Promotion / 1000 * c.BaseMarketingStrength

	promotionQuality := promotionQuality(c.employeePool, c.BaseMarketingStrength, c.employeePool.Get_employees_of_company(c.ID, Employee_type_marketing))

	for productID, offer := range c.Offers {
		offer.PromotionQuality = promotionQuality
		offer.Price = decisions.Products[productID].Price

		c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(
			"Advertisement costs",
			marketing,
			fmt.Sprintf("Cost of your ads for product %d (equals promotion quantity)", productID),
			true,
			float64(-decisions.Products[productID].Promotion.Quantity),
		)

		c.Offers[productID] = offer
	}
}
