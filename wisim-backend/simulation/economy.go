package simulation

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/pehringer/simd"
)

// Economy functions
// purchasingStatisticsMap is passed as argument to reduce unnececary garbage collection
func (population *Population) simulateEconomy(companies *[]Company, externalFactors External_factors, purchasingStatisticsMap map[string]Purchasing_statistics) ([]FinanceReportEntry, error) {
	// Get offers
	var offers []Offer
	var productAvailability []int
	var i int
	for _, c := range *companies {
		for productID, offer := range c.Offers {
			println("1")
			offers = append(offers, offer)
			productAvailability = append(productAvailability, c.ProductsInStorage[productID])
			i++
		}
	}

	if externalFactors.EconomicSituationIndex <= 0 {
		return make([]FinanceReportEntry, len(*companies)),
			errors.New("economic_situation_index cannot be 0")
	}

	// Calculate purchases
	var avgPrice float32
	for _, o := range offers {
		avgPrice += o.Price
	}

	if len(offers) == 0 {
		panic("len(offers) == 0; there are no offers!")
	}

	avgPrice = avgPrice / float32(len(offers))

	tBefore := time.Now()
	purchasingStatistics := make([]Purchasing_statistics, len(offers)+1)

	// Multithreading boilerplate
	var wg sync.WaitGroup
	numThreads := runtime.NumCPU()
	if numThreads <= 0 {
		panic("Thread count <= 0")
	}

	wg.Add(numThreads)

	offerPrices := make([]float32, len(offers))
	offerDurabilities := make([]int, len(offers))
	offersProperties := make([]Properties, len(offers))

	for i, o := range offers {
		offerPrices[i] = o.Price

		offersProperties[i][propertiesQuality] = o.Product.Quality
		offersProperties[i][propertiesEcology] = o.Product.Ecology
		offersProperties[i][propertiesEthics] = o.Product.Ethics
		offersProperties[i][propertiesPrice] = isCheap(o, avgPrice)
		offersProperties[i][propertiesBangForBuck] = o.Product.Quality / o.Price
		if o.Price <= 0 {
			offersProperties[i][propertiesPrice] = 10
			offersProperties[i][propertiesBangForBuck] = 10
		}
		offersProperties[i][propertiesDurability] = float32(o.Product.Durabilty)

		offerDurabilities[i] = int(offersProperties[i][propertiesDurability])
	}
	for _, interval := range split_load(numThreads, len(population.Population)) {
		go simulatePopulationSegment(&wg,
			population.Population[interval.Start:interval.Stop_before],
			population.Preferences[interval.Start:interval.Stop_before],
			offersProperties,
			offerPrices,
			offerDurabilities,
			productAvailability,
			externalFactors,
			purchasingStatistics)
	}

	wg.Wait()

	deltaTime := time.Since(tBefore)
	println("#### Time to calculate: ", deltaTime.String())

	divisionVector := [6]float32{
		float32(len(population.Population)),
		float32(len(population.Population)),
		float32(len(population.Population)),
		float32(len(population.Population)),
		float32(len(population.Population)),
		float32(len(population.Population)),
	}

	for i := range len(purchasingStatistics) - 1 {
		purchasingStatistics[i].AvrDecisionFactor /= float32(len(population.Population))
		purchasingStatistics[i].AvrPurchasingThreshold /= float32(len(population.Population))

		avrPurchasingFactors := purchasingStatistics[i].AvrPurchasingFactors
		simd.DivFloat32(avrPurchasingFactors[:], divisionVector[:], purchasingStatistics[i].AvrPurchasingFactors[:])

		purchasingStatistics[len(purchasingStatistics)-1].ProductsSold += purchasingStatistics[i].ProductsSold
		purchasingStatistics[len(purchasingStatistics)-1].ProductDemand += purchasingStatistics[i].ProductDemand

		purchasingStatistics[len(purchasingStatistics)-1].AvrDecisionFactor += purchasingStatistics[i].AvrDecisionFactor
		purchasingStatistics[len(purchasingStatistics)-1].AvrPurchasingThreshold += purchasingStatistics[i].AvrPurchasingThreshold

		simd.AddFloat32(purchasingStatistics[len(purchasingStatistics)-1].AvrPurchasingFactors[:], purchasingStatistics[i].AvrPurchasingFactors[:], purchasingStatistics[len(purchasingStatistics)-1].AvrPurchasingFactors[:])
	}

	purchasingStatistics[len(purchasingStatistics)-1].AvrDecisionFactor /= float32(len(purchasingStatistics) - 1)
	purchasingStatistics[len(purchasingStatistics)-1].AvrPurchasingThreshold /= float32(len(purchasingStatistics) - 1)

	simd.DivFloat32(purchasingStatistics[len(purchasingStatistics)-1].AvrPurchasingFactors[:], divisionVector[:], purchasingStatistics[len(purchasingStatistics)-1].AvrPurchasingFactors[:])

	results := make([]FinanceReportEntry, len(*companies))
	for i := range results {
		results[i] = FinanceReportEntry{"Products sold in stores", sales, fmt.Sprintf("%d products were sold in strores", purchasingStatistics[i].ProductsSold), true, float64(purchasingStatistics[i].ProductsSold * int(offers[i].Price))}
	}

	for i, offer := range offers {
		purchasingStatisticsMap[offer.Product.ID] = purchasingStatistics[i]
	}

	purchasingStatisticsMap["-1"] = purchasingStatistics[len(purchasingStatistics)-1]

	return results, nil
}

func simulatePopulationSegment(
	wg *sync.WaitGroup,
	populationSegment []Customer,
	populationSegmentPreferences []Properties,
	offersProperties []Properties,
	offerPrices []float32,
	offerDurabilities []int,
	productAvailability []int,
	externalFactors External_factors,
	purchasingStatistics []Purchasing_statistics,
) {
	productsPurchasingFactors := make([]float32, len(offersProperties))

	for i, customer := range populationSegment {
		customer.Max_price *= 1 + externalFactors.Inflation

		for i := range customer.Owned_products {
			customer.Owned_products[i].Remaining_durabilty -= 1
			if customer.Owned_products[i].Remaining_durabilty > 0 {
				customer.Owned_products = slices.Delete(customer.Owned_products, i, i)
			}
		}

		if customer.Base_need <= len(customer.Owned_products) {
			continue
		}

		for ii := range offersProperties {
			if productAvailability[ii] <= 0 {
				productsPurchasingFactors[ii] = 0
				continue
			} else if offerPrices[ii] > customer.Max_price {
				productsPurchasingFactors[ii] = 0
				continue
			}
			var productPurchasingFactorsComponents Properties
			simd.MulFloat32(populationSegmentPreferences[i][:], offersProperties[ii][:], productPurchasingFactorsComponents[:])
			productPurchasingFactor := float32(int(productPurchasingFactorsComponents[5] +
				productPurchasingFactorsComponents[4] +
				productPurchasingFactorsComponents[3] +
				productPurchasingFactorsComponents[2] +
				productPurchasingFactorsComponents[1] +
				productPurchasingFactorsComponents[0] +
				customer.Brand_loyalty_factor*customer.Loyalties[ii]))

			productPurchasingFactor *= externalFactors.EconomicSituationIndex

			productsPurchasingFactors[ii] = productPurchasingFactor

			simd.AddFloat32(productPurchasingFactorsComponents[:], purchasingStatistics[ii].AvrPurchasingFactors[:], purchasingStatistics[ii].AvrPurchasingFactors[:])

			purchasingStatistics[ii].AvrDecisionFactor += productPurchasingFactor
			purchasingStatistics[ii].AvrPurchasingThreshold += populationSegment[i].Purchashing_threshold
		}

		const noProductChosen = -1

		// Select product using weighted die
		choice := chooseProduct(productsPurchasingFactors, populationSegment[i].Purchashing_threshold)

		if choice != noProductChosen {
			purchasingStatistics[choice].ProductDemand += customer.Base_need - len(customer.Owned_products)
			numberOfProductsPurchased := 0
			for range customer.Base_need - len(customer.Owned_products) {
				if productAvailability[choice] > 0 {
					populationSegment[i].Owned_products = append(customer.Owned_products, Owned_product{choice, offerDurabilities[choice]})
					productAvailability[choice] -= 1
					numberOfProductsPurchased += 1
				} else {
					break
				}
			}
			purchasingStatistics[choice].ProductsSold += numberOfProductsPurchased
		}
		populationSegment[i] = customer
	}
	wg.Done()
}

func isCheap(offer Offer, avrPrice float32) float32 {
	if offer.Price == avrPrice {
		return 0.5
	}

	return 0.5 + (avrPrice - offer.Price)
}

func chooseProduct(decisionFactors []float32, purchasingThreshold float32) int {
	topProductsIndex := []int{}

	for i, p := range decisionFactors {
		if p < purchasingThreshold {
			continue
		} else if len(topProductsIndex) == 0 {
			topProductsIndex = []int{i}
		} else if p > decisionFactors[topProductsIndex[0]] {
			topProductsIndex = []int{i}
		} else if p == decisionFactors[topProductsIndex[0]] {
			topProductsIndex = append(topProductsIndex, i)
		}
	}

	// If only one product is best, choose that one
	switch len(topProductsIndex) {
	case 1:
		return topProductsIndex[0]
	case 0:
		return -1
	default:
		return topProductsIndex[rand.Intn(len(topProductsIndex))]
	}
}
