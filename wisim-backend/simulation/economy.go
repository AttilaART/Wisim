package simulation

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/pehringer/simd"
)

// Economy functions
func (population *Population) simulate_economy(companies *[]Company, external_factors External_factors) ([]FinanceReportEntry, []Purchasing_statistics, error) {
	// Get offers
	offers := make([]Offer, len(*companies))

	if external_factors.Economic_situation_index <= 0 {
		return make([]FinanceReportEntry, len(*companies)),
			make([]Purchasing_statistics, len(offers)+1),
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

	wg.Add(numThreads)

	offerPrices := make([]float32, len(offers))
	offerDurabilities := make([]int, len(offers))
	offersProperties := make([]Properties, len(offers))

	productAvailability := make([]int, len(offersProperties))
	for i, c := range *companies {
		offers[i] = c.Offer
		productAvailability[i] = c.Items_in_storage
	}
	for i, o := range offers {
		offerPrices[i] = o.Price

		offersProperties[i][propertiesQuality] = o.Product.Quality_factor
		offersProperties[i][propertiesEcology] = o.Product.Ecology_factor
		offersProperties[i][propertiesEthics] = o.Product.Ethics_factor
		offersProperties[i][propertiesPrice] = is_cheap(o, avgPrice)
		offersProperties[i][propertiesBangForBuck] = o.Product.Quality_factor / o.Price
		if o.Price <= 0 {
			offersProperties[i][propertiesPrice] = 10
			offersProperties[i][propertiesBangForBuck] = 10
		}
		offersProperties[i][propertiesDurability] = float32(o.Product.Durabilty)

		offerDurabilities[i] = int(offersProperties[i][propertiesDurability])
	}
	for _, interval := range split_load(numThreads, len(population.Population)) {
		go func(wg *sync.WaitGroup, populationSegment []Customer, populationSegmentPreferences []Properties,
		) {
			for i := range populationSegment {
				populationSegment[i].Max_price *= 1 + external_factors.Inflation
				populationSegment[i] = simulateDeterioration(populationSegment[i])
			}

			productsPurchasingFactors := make([]float32, len(offersProperties))

			for i := range populationSegmentPreferences {
				customer := populationSegment[i]

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

					productPurchasingFactor *= external_factors.Economic_situation_index

					productsPurchasingFactors[ii] = productPurchasingFactor

					simd.AddFloat32(productPurchasingFactorsComponents[:], purchasingStatistics[ii].Avr_purchasing_factors[:], purchasingStatistics[ii].Avr_purchasing_factors[:])

					purchasingStatistics[ii].Avr_decision_factor += productPurchasingFactor
					purchasingStatistics[ii].Avr_purchasing_threshold += populationSegment[i].Purchashing_threshold
				}

				const noProductChosen = -1

				// Select product using weighted die
				choice := choose_product(productsPurchasingFactors, populationSegment[i].Purchashing_threshold)

				if choice != noProductChosen {
					purchasingStatistics[choice].Product_demand += customer.Base_need - len(customer.Owned_products)
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
					purchasingStatistics[choice].Products_sold += numberOfProductsPurchased
				}
			}
			wg.Done()
		}(&wg, population.Population[interval.Start:interval.Stop_before], population.Preferences[interval.Start:interval.Stop_before])
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
		purchasingStatistics[i].Avr_decision_factor /= float32(len(population.Population))
		purchasingStatistics[i].Avr_purchasing_threshold /= float32(len(population.Population))

		avrPurchasingFactors := purchasingStatistics[i].Avr_purchasing_factors
		simd.DivFloat32(avrPurchasingFactors[:], divisionVector[:], purchasingStatistics[i].Avr_purchasing_factors[:])

		purchasingStatistics[len(purchasingStatistics)-1].Products_sold += purchasingStatistics[i].Products_sold
		purchasingStatistics[len(purchasingStatistics)-1].Product_demand += purchasingStatistics[i].Product_demand

		purchasingStatistics[len(purchasingStatistics)-1].Avr_decision_factor += purchasingStatistics[i].Avr_decision_factor
		purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_threshold += purchasingStatistics[i].Avr_purchasing_threshold

		simd.AddFloat32(purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_factors[:], purchasingStatistics[i].Avr_purchasing_factors[:], purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_factors[:])
	}

	purchasingStatistics[len(purchasingStatistics)-1].Avr_decision_factor /= float32(len(purchasingStatistics) - 1)
	purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_threshold /= float32(len(purchasingStatistics) - 1)

	simd.DivFloat32(purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_factors[:], divisionVector[:], purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_factors[:])

	results := make([]FinanceReportEntry, len(*companies))
	for i := range results {
		results[i] = FinanceReportEntry{"Products sold in stores", sales, fmt.Sprintf("%d products were sold in strores", purchasingStatistics[i].Products_sold), true, float64(purchasingStatistics[i].Products_sold * int(offers[i].Price))}
	}

	return results, purchasingStatistics, nil
}

func simulateDeterioration(customer Customer) Customer {
	var newOwnedProducts []Owned_product
	for _, p := range customer.Owned_products {
		p.Remaining_durabilty -= 1
		if p.Remaining_durabilty > 0 {
			newOwnedProducts = append(newOwnedProducts, p)
		}
	}

	customer.Owned_products = newOwnedProducts
	return customer
}

func is_cheap(offer Offer, avr_price float32) float32 {
	if offer.Price == avr_price {
		return 0.5
	}

	return 0.5 + (avr_price - offer.Price)
}

func choose_product(decision_factors []float32, purchasing_threshold float32) int {
	top_products_index := []int{}

	for i, p := range decision_factors {
		if p < purchasing_threshold {
			continue
		} else if len(top_products_index) == 0 {
			top_products_index = []int{i}
		} else if p > decision_factors[top_products_index[0]] {
			top_products_index = []int{i}
		} else if p == decision_factors[top_products_index[0]] {
			top_products_index = append(top_products_index, i)
		}
	}

	// If only one product is best, choose that one
	switch len(top_products_index) {
	case 1:
		return top_products_index[0]
	case 0:
		return -1
	default:
		return top_products_index[rand.Intn(len(top_products_index))]
	}
}
