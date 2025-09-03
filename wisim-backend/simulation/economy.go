package simulation

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
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

	t_before := time.Now()
	purchasing_statistics := make([]Purchasing_statistics, len(offers)+1)

	// Multithreading boilerplate
	var wg sync.WaitGroup
	num_threads := runtime.NumCPU()

	wg.Add(num_threads)
	for id, interval := range split_load(num_threads, len(population.Population)) {
		go func(wg *sync.WaitGroup, population_range Interval, id int,
		) {
			for current_customer_index := population_range.Start; current_customer_index < population_range.Stop_before; current_customer_index++ {

				if current_customer_index >= len(population.Population) {
					log.Panic("current_customer_index higher than len")
					break
				}

				population.Population[current_customer_index] = calculate_purchase(
					simulate_deterioration(population.Population[current_customer_index]),
					offers, avg_price, external_factors,
					&product_availability,
					&purchasing_statistics)

				population.Population[current_customer_index].Max_price *= 1 + external_factors.Inflation
			}

			wg.Done()
		}(&wg, interval, id)
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
	delta_time := time.Since(t_before)
	println("#### Time to calculate: ", delta_time.String())

	for i := range len(purchasing_statistics) - 1 {
		purchasing_statistics[i].Avr_decision_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_purchasing_threshold /= float32(len(population.Population))

		purchasing_statistics[i].Avr_quality_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_durability_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_ecology_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_price_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_ethics_factor /= float32(len(population.Population))
		// purchasing_statistics[i].Avr_coolness_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_bang_for_buck_factor /= float32(len(population.Population))

		purchasing_statistics[len(purchasing_statistics)-1].Products_sold += purchasing_statistics[i].Products_sold
		purchasing_statistics[len(purchasing_statistics)-1].Product_demand += purchasing_statistics[i].Product_demand

		purchasing_statistics[len(purchasing_statistics)-1].Avr_decision_factor += purchasing_statistics[i].Avr_decision_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_purchasing_threshold += purchasing_statistics[i].Avr_purchasing_threshold

		purchasing_statistics[len(purchasing_statistics)-1].Avr_quality_factor += purchasing_statistics[i].Avr_quality_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_durability_factor += purchasing_statistics[i].Avr_durability_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_ecology_factor += purchasing_statistics[i].Avr_ecology_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_price_factor += purchasing_statistics[i].Avr_price_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_ethics_factor += purchasing_statistics[i].Avr_ethics_factor
		// purchasing_statistics[len(purchasing_statistics)-1].Avr_coolness_factor += purchasing_statistics[i].Avr_coolness_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_bang_for_buck_factor += purchasing_statistics[i].Avr_bang_for_buck_factor
	}

	purchasing_statistics[len(purchasing_statistics)-1].Avr_decision_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_purchasing_threshold /= float32(len(purchasing_statistics) - 1)

	purchasing_statistics[len(purchasing_statistics)-1].Avr_quality_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_durability_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_ecology_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_price_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_ethics_factor /= float32(len(purchasing_statistics) - 1)
	// purchasing_statistics[len(purchasing_statistics)-1].Avr_coolness_factor /= float32(len(purchasing_statistics) - 1)

	results := make([]FinanceReportEntry, len(*companies))
	for i := range results {
		results[i] = FinanceReportEntry{"Products sold in stores", sales, fmt.Sprintf("%d products were sold in strores", purchasing_statistics[i].Products_sold), true, float64(purchasing_statistics[i].Products_sold * int(offers[i].Price))}
	}

	return results, purchasing_statistics, nil
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

func calculate_purchase(customer Customer, offers []Offer, avg_price float32, external_factors External_factors, product_availability *[]int, purchasing_statistics *[]Purchasing_statistics) Customer {
	decision_factors := make([]float32, len(offers))
	for i, o := range offers {
		if (*product_availability)[i] <= 0 {
			continue
		} // idk if this is good but it saves 2s of processing

		if o.Price > customer.Max_price {
			continue
		}

		quality_factor := clamp(customer.Quality_preference*o.Product.Quality_factor, customer.Quality_preference*10)
		durability_factor := clamp(customer.Durabilty_preference*float32(o.Product.Durabilty)*3, customer.Durabilty_preference*10)
		ecology_factor := clamp(customer.Ecology_preference*o.Product.Ecology_factor, customer.Ecology_preference*10)
		ethics_factor := clamp(customer.Ethics_preference*o.Product.Ethics_factor, customer.Ethics_preference*10)
		price_factor := clamp(customer.Price_preference*is_cheap(o, avg_price), customer.Price_preference*10)
		if o.Price <= 0 {
			price_factor = customer.Price_preference * 10
		}
		bang_for_buck_factor := clamp(customer.Bang_for_buck_preference*(o.Product.Quality_factor/o.Price), customer.Bang_for_buck_preference*10)
		if o.Price <= 0 {
			bang_for_buck_factor = customer.Bang_for_buck_preference * 10
		}
		brand_loyalty_factor := clamp(customer.Brand_loyalty_factor*customer.Loyalties[i], customer.Brand_loyalty_factor*10)

		decision_factors[i] = (quality_factor + ecology_factor + ethics_factor + price_factor + bang_for_buck_factor + brand_loyalty_factor) * external_factors.Economic_situation_index

		(*purchasing_statistics)[i].Avr_decision_factor += decision_factors[i]
		(*purchasing_statistics)[i].Avr_purchasing_threshold += customer.Purchashing_threshold
		(*purchasing_statistics)[i].Avr_quality_factor += quality_factor
		(*purchasing_statistics)[i].Avr_durability_factor += durability_factor
		(*purchasing_statistics)[i].Avr_ecology_factor += ecology_factor
		(*purchasing_statistics)[i].Avr_price_factor += price_factor
		(*purchasing_statistics)[i].Avr_ethics_factor += ethics_factor
		// (*purchasing_statistics)[i].Avr_coolness_factor += customer.Coolness_preference * offers[i].Product.Coolness_factor
		(*purchasing_statistics)[i].Avr_bang_for_buck_factor += bang_for_buck_factor
	}

	// Select product using weighted die
	choice := choose_product(decision_factors, customer.Purchashing_threshold)

	if choice != -1 {
		(*purchasing_statistics)[choice].Product_demand += customer.Base_need - len(customer.Owned_products)
		number_of_products_purchased := 0
		for range customer.Base_need - len(customer.Owned_products) {
			if (*product_availability)[choice] > 0 {
				customer.Owned_products = append(customer.Owned_products, Owned_product{choice, offers[choice].Product.Durabilty})
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
