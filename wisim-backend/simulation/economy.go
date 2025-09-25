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
func (population *Population) simulateEconomy(companies []Company, externalFactors ExternalFactors, purchasingStatisticsMap map[string]Purchasing_statistics) error {
	// Get offers
	var offers []Offer
	var offerIDs []struct {
		ID      string
		Company int
	}
	var productAvailability []int
	for _, c := range companies {
		for productID, offer := range c.Offers {
			println(offer.Product.Name)
			offers = append(offers, offer)
			offerIDs = append(offerIDs, struct {
				ID      string
				Company int
			}{offer.Product.ID, c.ID})
			productAvailability = append(productAvailability, c.ProductsInStorage[productID])
		}
	}

	if externalFactors.EconomicSituationIndex <= 0 {
		return errors.New("economic_situation_index cannot be 0")
	}

	// Calculate purchases
	var avgPrice float32
	for _, o := range offers {
		avgPrice += o.Price
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
	offerMutexes := make([]sync.Mutex, len(offers))

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
			offers,
			offerPrices,
			offerDurabilities,
			offerMutexes,
			productAvailability,
			externalFactors,
			purchasingStatistics,
			len(population.Population))
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

	for i, o := range offerIDs {
		companies[o.Company].Reports[len(companies[o.Company].Reports)-1].BalanceSheet.add_to_income_statement("Sales of "+offers[i].Product.Name, sales, fmt.Sprintf("%d %ss were sold in strores", purchasingStatistics[i].ProductsSold, offers[i].Product.Name), true, float64(purchasingStatistics[i].ProductsSold*int(offers[i].Price)))
	}

	for i, offer := range offers {
		purchasingStatisticsMap[offer.Product.ID] = purchasingStatistics[i]
	}

	purchasingStatisticsMap["-1"] = purchasingStatistics[len(purchasingStatistics)-1]

	return nil
}

func simulatePopulationSegment(
	wg *sync.WaitGroup,
	populationSegment []Customer,
	populationSegmentPreferences []Properties,
	offersProperties []Properties,
	offers []Offer,
	offerPrices []float32,
	offerDurabilities []int,
	offerMutexes []sync.Mutex,
	productAvailability []int,
	externalFactors ExternalFactors,
	purchasingStatistics []Purchasing_statistics,
	populationSize int,
) {
	for i, customer := range populationSegment {
		productsPurchasingFactors := make([]float32, len(offersProperties))
		customer.Max_price *= 1 + externalFactors.Inflation

		for i := range customer.Owned_products {
			customer.Owned_products[i].RemainingDurabilty -= 1
			if customer.Owned_products[i].RemainingDurabilty > 0 {
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
				customer.Brand_loyalty_factor*customer.Loyalties[offers[ii].Product.CompanyID]))

			productPurchasingFactor *= externalFactors.EconomicSituationIndex

			productsPurchasingFactors[ii] = productPurchasingFactor

			simd.AddFloat32(productPurchasingFactorsComponents[:], purchasingStatistics[ii].AvrPurchasingFactors[:], purchasingStatistics[ii].AvrPurchasingFactors[:])

			purchasingStatistics[ii].AvrDecisionFactor += productPurchasingFactor
			purchasingStatistics[ii].AvrPurchasingThreshold += populationSegment[i].Purchashing_threshold
		}

		// calculate promotion
		for ii, o := range offers {
			marketingReach := (o.Promotion.Quantity / float32(populationSize)) * rand.Float32() * 10
			debugString := "marketingReach %s (%s) %.4f "
			if marketingReach > customer.Savyness {
				debugString += "> "
			} else {
				debugString += "< "
			}
			debugString += "customer.Savyness (%.2f); threshold(%.2f): %t (%.2f)\n"

			productsPurchasingFactors[ii] += productsPurchasingFactors[ii] * (customer.Loyalties[o.Product.CompanyID] * customer.Brand_loyalty_factor)

			if customer.Savyness <= marketingReach {
				marketingStickingFactor := (populationSegmentPreferences[i][propertiesQuality]*o.Promotion.StyleQuality +
					populationSegmentPreferences[i][propertiesEcology]*o.Promotion.StyleEcology +
					populationSegmentPreferences[i][propertiesEthics]*o.Promotion.StyleEthics +
					populationSegmentPreferences[i][propertiesPrice]*o.Promotion.StylePrice +
					populationSegmentPreferences[i][propertiesDurability]*o.Promotion.StyleDurability) *
					o.Promotion.Quality

				productsPurchasingFactors[ii] += productsPurchasingFactors[ii] + marketingStickingFactor
			} else if customer.Loyalties[o.Product.CompanyID] < 0.5 {
				productsPurchasingFactors[ii] *= 0.75
			}

			// fmt.Printf(debugString, o.Product.Name, o.Product.ID, marketingReach, customer.Savyness, customer.Purchashing_threshold, productsPurchasingFactors[ii] >= customer.Purchashing_threshold, productsPurchasingFactors[ii])

		}

		const noProductChosen = -1

		// Select product using weighted die
		choice := chooseProduct(productsPurchasingFactors, populationSegment[i].Purchashing_threshold)

		if choice != noProductChosen {
			purchasingStatistics[choice].ProductDemand += customer.Base_need - len(customer.Owned_products)
			populationSegment[i].Loyalties[offers[choice].Product.CompanyID] += productsPurchasingFactors[choice] / 10
			for range customer.Base_need - len(customer.Owned_products) {
				offerMutexes[choice].Lock()
				if productAvailability[choice] > 0 {
					productAvailability[choice] -= 1
					purchasingStatistics[choice].ProductsSold++
					offerMutexes[choice].Unlock()
					populationSegment[i].Owned_products = append(customer.Owned_products, OwnedProduct{choice, offerDurabilities[choice]})
				} else {
					offerMutexes[choice].Unlock()
					break
				}
			}

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
