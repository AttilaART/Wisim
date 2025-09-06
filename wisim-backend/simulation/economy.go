package simulation

import (
	"errors"
	"fmt"
	"math"
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
	product_availability := make([]int, len(*companies))
	for i, c := range *companies {
		offers[i] = c.Offer
		product_availability[i] = c.Items_in_storage
	}

	if external_factors.Economic_situation_index <= 0 {
		return make([]FinanceReportEntry, len(*companies)),
			make([]Purchasing_statistics, len(offers)+1),
			errors.New("economic_situation_index cannot be 0")
	}

	// Calculate purchases
	var avg_price float32
	for _, o := range offers {
		avg_price += o.Price
	}

	if len(offers) == 0 {
		panic("len(offers) == 0; there are no offers!")
	}

	avg_price = avg_price / float32(len(offers))

	tBefore := time.Now()
	purchasingStatistics := make([]Purchasing_statistics, len(offers)+1)

	// Multithreading boilerplate
	var wg sync.WaitGroup
	numThreads := runtime.NumCPU() * 2

	wg.Add(numThreads)

	offerPrices := make([]float32, len(offers))
	offerDurabilities := make([]int, len(offers))
	offersProperties := make([]Properties, len(offers))
	for i, o := range offers {
		offerPrices[i] = o.Price

		offersProperties[i][properties_quality] = o.Product.Quality_factor
		offersProperties[i][properties_ecology] = o.Product.Ecology_factor
		offersProperties[i][properties_ethics] = o.Product.Ethics_factor
		offersProperties[i][properties_price] = is_cheap(o, avg_price)
		offersProperties[i][properties_bang_for_buck] = o.Product.Quality_factor / o.Price
		if o.Price <= 0 {
			offersProperties[i][properties_price] = 10
			offersProperties[i][properties_bang_for_buck] = 10
		}
		offersProperties[i][properties_durability] = float32(o.Product.Durabilty)

		offerDurabilities[i] = int(offersProperties[i][properties_durability])
	}
	for id, interval := range split_load(numThreads, len(population.Population)) {
		go func(wg *sync.WaitGroup, populationSegment []Customer, id int,
		) {
			allCustomerPreferences := make([]float32, len(populationSegment)*len(Properties{}))
			for i := range populationSegment {
				for ii := range populationSegment[i].Preferences {
					allCustomerPreferences[i*6+ii] = populationSegment[i].Preferences[ii]
				}
			}

			allCustomerAllProductDecisionFactorConstituents := make([][]float32, len(offersProperties))
			for i, offerProperties := range offersProperties {
				correspondingProductProperties := make([]float32, len(populationSegment)*len(Properties{}))
				for ii := range correspondingProductProperties {
					correspondingProductProperties[ii] = offerProperties[ii%len(offerProperties)]
				}
				allCustomerProductDecisionFactorConstituents := make([]float32, len(population.Population)*len(Properties{}))
				simd.MulFloat32(allCustomerPreferences, correspondingProductProperties, allCustomerProductDecisionFactorConstituents)
				allCustomerAllProductDecisionFactorConstituents[i] = allCustomerProductDecisionFactorConstituents
			}

			allPurchasingFactors := make([][]float32, len(populationSegment))
			for i := range allPurchasingFactors {
				allPurchasingFactors[i] = make([]float32, len(offers))
				for ii := range allPurchasingFactors[i] {
					propertyI := i * len(Properties{})
					productDecisionFactorConstituents := allCustomerAllProductDecisionFactorConstituents[ii][propertyI : propertyI+6]
					allPurchasingFactors[i][ii] = productDecisionFactorConstituents[5] +
						productDecisionFactorConstituents[4] +
						productDecisionFactorConstituents[3] +
						productDecisionFactorConstituents[2] +
						productDecisionFactorConstituents[1] +
						productDecisionFactorConstituents[0]
				}
			}

			for currentCustomerIndex := range populationSegment {
				populationSegment[currentCustomerIndex].Max_price *= 1 + external_factors.Inflation

				populationSegment[currentCustomerIndex] = calculate_purchase(
					simulate_deterioration(populationSegment[currentCustomerIndex]),
					allPurchasingFactors[currentCustomerIndex],
					offerPrices,
					offerDurabilities,
					external_factors,
					&product_availability,
					&purchasingStatistics)

			}
			wg.Done()
		}(&wg, population.Population[interval.Start:interval.Stop_before], id)
	}

	wg.Wait()

	//for i, c := range population.Population {
	//	customer := calcualte_durability(c)
	//	product_purchased, quanity, customer, individual_purchasing_statistics := calculate_purchase(customer, offers, avg_price, external_factors, product_availability)
	//	population.Population[i] = customer
	//	purchases[product_purchased] += quanity
	//
	//	for i, s := range individual_purchasing_statistics {
	//		purchasing_statistics[i].Products_sold += s.Products_sold
	//		purchasing_statistics[i].Product_demand += s.Product_demand
	//		purchasing_statistics[i].Product_number = s.Product_number
	//		purchasing_statistics[i].Avr_decision_factor += s.Avr_decision_factor
	//		purchasing_statistics[i].Avr_purchasing_threshold += s.Avr_purchasing_threshold
	//
	//		purchasing_statistics[i].Avr_quality_factor += s.Avr_quality_factor
	//		purchasing_statistics[i].Avr_durability_factor += s.Avr_durability_factor
	//		purchasing_statistics[i].Avr_ecology_factor += s.Avr_ecology_factor
	//		purchasing_statistics[i].Avr_price_factor += s.Avr_price_factor
	//		purchasing_statistics[i].Avr_coolness_factor += s.Avr_coolness_factor
	//	}
	//}
	delta_time := time.Since(tBefore)
	println("#### Time to calculate: ", delta_time.String())

	div_vector := [6]float32{
		float32(len(population.Population)),
		float32(len(population.Population)),
		float32(len(population.Population)),
		float32(len(population.Population)),
		float32(len(population.Population)),
		float32(len(population.Population)),
	}

	if len(purchasingStatistics[0].Avr_purchasing_factors[:]) != len(div_vector[:]) {
		panic("div vector not same length as Properties, expect undefined behavior")
	}

	for i := range len(purchasingStatistics) - 1 {
		purchasingStatistics[i].Avr_decision_factor /= float32(len(population.Population))
		purchasingStatistics[i].Avr_purchasing_threshold /= float32(len(population.Population))

		Avr_purchasing_factors := purchasingStatistics[i].Avr_purchasing_factors
		simd.DivFloat32(Avr_purchasing_factors[:], div_vector[:], purchasingStatistics[i].Avr_purchasing_factors[:])

		purchasingStatistics[len(purchasingStatistics)-1].Products_sold += purchasingStatistics[i].Products_sold
		purchasingStatistics[len(purchasingStatistics)-1].Product_demand += purchasingStatistics[i].Product_demand

		purchasingStatistics[len(purchasingStatistics)-1].Avr_decision_factor += purchasingStatistics[i].Avr_decision_factor
		purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_threshold += purchasingStatistics[i].Avr_purchasing_threshold

		general_avr_purchasing_factors := purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_factors
		simd.AddFloat32(general_avr_purchasing_factors[:], purchasingStatistics[i].Avr_purchasing_factors[:], purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_factors[:])
	}

	purchasingStatistics[len(purchasingStatistics)-1].Avr_decision_factor /= float32(len(purchasingStatistics) - 1)
	purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_threshold /= float32(len(purchasingStatistics) - 1)

	general_avr_purchasing_factors := purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_factors
	simd.DivFloat32(general_avr_purchasing_factors[:], div_vector[:], purchasingStatistics[len(purchasingStatistics)-1].Avr_purchasing_factors[:])

	results := make([]FinanceReportEntry, len(*companies))
	for i := range results {
		results[i] = FinanceReportEntry{"Products sold in stores", sales, fmt.Sprintf("%d products were sold in strores", purchasingStatistics[i].Products_sold), true, float64(purchasingStatistics[i].Products_sold * int(offers[i].Price))}
	}

	return results, purchasingStatistics, nil
}

func simulate_deterioration(customer Customer) Customer {
	var new_owned_products []Owned_product
	for _, p := range customer.Owned_products {
		p.Remaining_durabilty -= 1
		if p.Remaining_durabilty > 0 {
			new_owned_products = append(new_owned_products, p)
		}
	}

	customer.Owned_products = new_owned_products
	return customer
}

func calculate_purchase(customer Customer, purchasing_factors []float32, offer_prices []float32, offer_durabilities []int, external_factors External_factors, product_availability *[]int, purchasing_statistics *[]Purchasing_statistics) Customer {
	decision_factors := make([]float32, len(offer_prices))
	for i := range offer_prices {
		if (*product_availability)[i] <= 0 {
			continue
		} // idk if this is good but it saves 2s of processing

		if offer_prices[i] > customer.Max_price {
			continue
		}

		brand_loyalty_factor := clamp(customer.Brand_loyalty_factor*customer.Loyalties[i], customer.Brand_loyalty_factor*10)

		decision_factors[i] = (purchasing_factors[i] + brand_loyalty_factor) * external_factors.Economic_situation_index

		(*purchasing_statistics)[i].Avr_decision_factor += decision_factors[i]
		(*purchasing_statistics)[i].Avr_purchasing_threshold += customer.Purchashing_threshold
		Avr_purchasing_factors := (*purchasing_statistics)[i].Avr_purchasing_factors
		simd.AddFloat32(purchasing_factors[:], Avr_purchasing_factors[:], (*purchasing_statistics)[i].Avr_purchasing_factors[:])
	}

	// Select product using weighted die
	choice := choose_product(decision_factors, customer.Purchashing_threshold)

	if choice != -1 {
		(*purchasing_statistics)[choice].Product_demand += customer.Base_need - len(customer.Owned_products)
		number_of_products_purchased := 0
		for range customer.Base_need - len(customer.Owned_products) {
			if (*product_availability)[choice] > 0 {
				customer.Owned_products = append(customer.Owned_products, Owned_product{choice, offer_durabilities[choice]})
				(*product_availability)[choice] -= 1
				number_of_products_purchased += 1
			} else {
				break
			}
		}
		(*purchasing_statistics)[choice].Products_sold += number_of_products_purchased
		return customer
	}
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

	// println("----------")
	// Round decision_factors
	for i, f := range decision_factors {
		decision_factors[i] = float32(math.Round(float64(f)))
		// fmt.Printf("decision_factors[%d]: %f\n", i, decision_factors[i])
	}
	// fmt.Printf("purchasing_threshold: %f\n", purchasing_threshold)

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
