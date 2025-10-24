package simulation

import (
	"fmt"
	"math"
)

func (c *Company) newProduct(productID string, companyID int, productName string, productDecisions Decisions_product) Offer {
	productDecisions.Product.TechLevels = c.Tech
	productStats, productionLineCost := CalculateProductStats(productDecisions.Product, *c.productComponents)
	offer := Offer{Outdated: productDecisions.Outdated, Product: productDecisions.Product, ProductStats: productStats}
	offer.Product.ID = productID
	offer.Product.CompanyID = companyID
	offer.Product.Name = productName

	c.Reports[len(c.Reports)-1].BalanceSheet.add_to_income_statement(fmt.Sprintf("Production line cost for product %s", productName), facilities, "The cost of setting up the production line for the product.", true, -productionLineCost)

	return offer
}

// Calcualtes base product stats
func CalculateProductStats(
	product Product, productComponents ProductComponents,
) (productStats ProductStats, productionLineCost float64) { // Calcualtes product stats without side effects

	parts := []Component{}

	parts = append(parts, productComponents.FormFactor[product.Components.FormFactor])
	parts = append(parts, productComponents.Frame[product.Components.Frame])
	parts = append(parts, productComponents.Body[product.Components.Body])
	parts = append(parts, productComponents.Mechanism[product.Components.Mechanism])

	productStats.MiscSlots += productComponents.FormFactor[product.Components.FormFactor].MiscSlots
	productStats.MiscSlots += productComponents.Frame[product.Components.Frame].MiscSlots
	productStats.MiscSlots += productComponents.Body[product.Components.Body].MiscSlots
	productStats.MiscSlots += productComponents.Mechanism[product.Components.Mechanism].MiscSlots

	for i, component := range product.Components.Misc {
		if i < productStats.MiscSlots {
			parts = append(parts, productComponents.Misc[component])
		}
	}

	for _, part := range parts {
		productStats.ProductionCost += part.ProductionCost
		productStats.MaterialUse += part.MaterialUse

		productStats.Quality += part.Quality
		productStats.Ecology += part.Ecology
		productStats.Ethics += part.Ethics
		productStats.Durability += part.Durability
		productionLineCost += float64(part.ProductionLineCost)
	}

	productStats.Durability = float32(math.Round(float64(productStats.Durability)))

	productStats.Durability += float32(product.ExtraDurability)
	productStats.ProductionCost += float32(5 * product.ExtraDurability)

	productStats.Quality += float32(product.ExtraQuality)
	productStats.ProductionCost += float32(5 * product.ExtraQuality)

	productStats.Quality += float32(product.MaterialQuality)
	productStats.MaterialUse += float32(5 * product.MaterialQuality)

	productStats.Durability *= product.TechLevels.Durability
	productStats.Ecology *= product.TechLevels.Ecology
	productStats.Quality *= product.TechLevels.Quality

	if product.TechLevels.MaterialUse != 0 {
		productStats.MaterialUse /= product.TechLevels.MaterialUse
	}
	if product.TechLevels.ProductionCost != 0 {
		productStats.ProductionCost /= product.TechLevels.ProductionCost
	}

	return productStats, productionLineCost
}

func (c *Company) calculatePromotion(decisions Decisions) {
	c.BaseMarketingStrength += decisions.Research.Promotion / 1000 * c.BaseMarketingStrength

	promotionQuality := promotionQuality(c.employeePool, c.BaseMarketingStrength, c.employeePool.Get_employees_of_company(c.ID, Employee_type_marketing))

	for productID, offer := range c.Offers {
		offer.PromotionQuality = promotionQuality
		offer.Promotion.Quantity = decisions.Products[productID].Promotion.Quantity
		offer.Price = decisions.Products[productID].Price
		offer.Product.Name = decisions.Products[productID].Name

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
