package simulation

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
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
		return make([]FinanceReportEntry, len(*companies)), make([]Purchasing_statistics, len(offers)+1), errors.New("economic_situation_index cannot be 0")
	}

	purchases := make([]int, len(*companies))
	// Calculate purchases
	var avg_price float32
	for _, o := range offers {
		avg_price += o.Price
	}
	avg_price = avg_price / float32(len(offers))

	purchasing_statistics := make([]Purchasing_statistics, len(offers)+1)
	for i, c := range population.Population {
		customer := calcualte_durability(c)
		product_purchased, quanity, customer, individual_purchasing_statistics := calculate_purchase(customer, offers, avg_price, external_factors, product_availability)
		population.Population[i] = customer
		purchases[product_purchased] += quanity

		for i, s := range individual_purchasing_statistics {
			purchasing_statistics[i].Products_sold += s.Products_sold
			purchasing_statistics[i].Product_number = s.Product_number
			purchasing_statistics[i].Avr_decision_factor += s.Avr_decision_factor
			purchasing_statistics[i].Avr_purchasing_threshold += s.Avr_purchasing_threshold

			purchasing_statistics[i].Product_quality = s.Product_quality
			purchasing_statistics[i].Product_durabilty = s.Product_durabilty
			purchasing_statistics[i].Product_ecology = s.Product_ecology
			purchasing_statistics[i].Product_price_factor = s.Product_price_factor
			purchasing_statistics[i].Product_price = s.Product_price
			purchasing_statistics[i].Product_coolness = s.Product_coolness
			purchasing_statistics[i].Avr_bang_for_buck_factor = s.Avr_bang_for_buck_factor

			purchasing_statistics[i].Avr_quality_factor += s.Avr_quality_factor
			purchasing_statistics[i].Avr_durability_factor += s.Avr_durability_factor
			purchasing_statistics[i].Avr_ecology_factor += s.Avr_ecology_factor
			purchasing_statistics[i].Avr_price_factor += s.Avr_price_factor
			purchasing_statistics[i].Avr_coolness_factor += s.Avr_coolness_factor
			purchasing_statistics[i].Avr_bang_for_buck_factor += s.Avr_bang_for_buck_factor
		}
	}

	for i := range len(purchasing_statistics) - 1 {
		purchasing_statistics[i].Avr_decision_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_purchasing_threshold /= float32(len(population.Population))

		purchasing_statistics[i].Avr_quality_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_durability_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_ecology_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_price_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_coolness_factor /= float32(len(population.Population))
		purchasing_statistics[i].Avr_bang_for_buck_factor /= float32(len(population.Population))

		purchasing_statistics[len(purchasing_statistics)-1].Products_sold += purchasing_statistics[i].Products_sold

		purchasing_statistics[len(purchasing_statistics)-1].Avr_decision_factor += purchasing_statistics[i].Avr_decision_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_purchasing_threshold += purchasing_statistics[i].Avr_purchasing_threshold
		purchasing_statistics[len(purchasing_statistics)-1].Product_quality += purchasing_statistics[i].Product_quality
		purchasing_statistics[len(purchasing_statistics)-1].Product_durabilty += purchasing_statistics[i].Product_durabilty
		purchasing_statistics[len(purchasing_statistics)-1].Product_ecology += purchasing_statistics[i].Product_ecology
		purchasing_statistics[len(purchasing_statistics)-1].Product_price_factor += purchasing_statistics[i].Product_price_factor
		purchasing_statistics[len(purchasing_statistics)-1].Product_price += purchasing_statistics[i].Product_price
		purchasing_statistics[len(purchasing_statistics)-1].Product_coolness += purchasing_statistics[i].Product_coolness

		purchasing_statistics[len(purchasing_statistics)-1].Avr_quality_factor += purchasing_statistics[i].Avr_quality_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_durability_factor += purchasing_statistics[i].Avr_durability_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_ecology_factor += purchasing_statistics[i].Avr_ecology_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_price_factor += purchasing_statistics[i].Avr_price_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_coolness_factor += purchasing_statistics[i].Avr_coolness_factor
		purchasing_statistics[len(purchasing_statistics)-1].Avr_bang_for_buck_factor += purchasing_statistics[i].Avr_bang_for_buck_factor
	}

	purchasing_statistics[len(purchasing_statistics)-1].Avr_decision_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_purchasing_threshold /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Product_quality /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Product_durabilty /= len(purchasing_statistics)
	purchasing_statistics[len(purchasing_statistics)-1].Product_ecology /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Product_price_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Product_price /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Product_coolness /= float32(len(purchasing_statistics) - 1)

	purchasing_statistics[len(purchasing_statistics)-1].Avr_quality_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_durability_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_ecology_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_price_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_coolness_factor /= float32(len(purchasing_statistics) - 1)
	purchasing_statistics[len(purchasing_statistics)-1].Avr_bang_for_buck_factor += float32(len(purchasing_statistics) - 1)

	results := make([]FinanceReportEntry, len(*companies))
	for i := range results {
		results[i] = FinanceReportEntry{"Products sold in stores", sales, fmt.Sprintf("%d products were sold in strores", purchases[i]), true, float64(purchases[i] * int(offers[i].Price))}
	}

	return results, purchasing_statistics, nil
}

func calcualte_durability(customer Customer) Customer {
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

func calculate_purchase(customer Customer, offers []Offer, avg_price float32, external_factors External_factors, product_availability []int) (int, int, Customer, []Purchasing_statistics) { // returns (index of product purchased, quality purchased, customer in question)

	decision_factors := make([]float32, len(offers))
	purchasing_statistics := make([]Purchasing_statistics, len(offers))
	for i, o := range offers {
		decision_factors[i] = (customer.Quality_preference*o.Product.Quality_factor +
			customer.Ecology_preference*o.Product.Ecology_factor +
			customer.Coolness_preference*o.Product.Coolness_factor +
			customer.Price_preference*is_cheap(o, avg_price) +
			customer.Bang_for_buck_preference*(o.Product.Quality_factor/o.Price) +
			customer.Brand_loyalty_factor*customer.Loyalties[i]) * external_factors.Economic_situation_index

		purchasing_statistics[i] = Purchasing_statistics{
			Products_sold:  0,
			Product_number: i,

			Avr_decision_factor:      decision_factors[i],
			Avr_purchasing_threshold: customer.Purchashing_threshold,

			Product_quality:       offers[i].Product.Quality_factor,
			Product_durabilty:     offers[i].Product.Durabilty,
			Product_ecology:       offers[i].Product.Ecology_factor,
			Product_price_factor:  is_cheap(offers[i], avg_price),
			Product_price:         offers[i].Price,
			Product_coolness:      offers[i].Product.Coolness_factor,
			Product_bang_for_buck: offers[i].Product.Quality_factor / offers[i].Price,

			Avr_quality_factor:       customer.Quality_preference * offers[i].Product.Quality_factor,
			Avr_durability_factor:    customer.Durabilty_preference * float32(offers[i].Product.Durabilty),
			Avr_ecology_factor:       customer.Ecology_preference * offers[i].Product.Ecology_factor,
			Avr_price_factor:         customer.Price_preference * is_cheap(offers[i], avg_price),
			Avr_coolness_factor:      customer.Coolness_preference * offers[i].Product.Coolness_factor,
			Avr_bang_for_buck_factor: customer.Bang_for_buck_preference * (offers[i].Product.Quality_factor / offers[i].Price),
		}
	}

	// Select product using weighted die
	choice, err := choose_product(decision_factors, customer.Purchashing_threshold)
	if err != nil {
		log.Fatalf("%s at customer %+v with offer", err.Error(), customer)
	}

	if choice != -1 {
		number_of_products_purchased := 0
		for range customer.Base_need - len(customer.Owned_products) {
			if product_availability[choice] > 0 {
				customer.Owned_products = append(customer.Owned_products, Owned_product{choice, offers[choice].Product.Durabilty})
				product_availability[choice] -= 1
				number_of_products_purchased += 1
			} else {
				break
			}
		}
		purchasing_statistics[choice].Products_sold = number_of_products_purchased
		return choice, number_of_products_purchased, customer, purchasing_statistics
	}
	return 0, 0, customer, purchasing_statistics
}

func is_cheap(offer Offer, avr_price float32) float32 {
	difference_from_avg := (avr_price - offer.Price) / avr_price * 10

	return difference_from_avg
}

func choose_product(decision_factors []float32, purchasing_threshold float32) (int, error) {
	og_decision_factors := make([]float32, len(decision_factors))
	copy(og_decision_factors, decision_factors)

	// make all sections have a sum of 1
	var len_of_weights float32 // get max
	for _, s := range decision_factors {
		len_of_weights += s
	}

	for i := range decision_factors {
		decision_factors[i] /= len_of_weights
	}

	sections := make([]float32, len(decision_factors))

	var len_of_sections float32 = 0
	for i, w := range decision_factors {
		sections[i] = len_of_sections + w
		len_of_sections += w
	}
	sections[len(sections)-1] = 1

	// check the total of sections
	//len_of_sections = 0
	//for _, s := range sections {
	//	if s > len_of_sections {
	//		len_of_sections = s
	//	}
	//}
	// println(len_of_sections)

	random_number := rand.Float32()

	choice := -2
	for i := range sections {
		if i == 0 {
			if random_number <= sections[i] {
				choice = i
			}
		} else if random_number > sections[i-1] && random_number <= sections[i] {
			choice = i
		}
	}
	if og_decision_factors[choice] >= purchasing_threshold {
		return choice, nil
	} else if decision_factors[choice] < purchasing_threshold {
		return -1, nil
	}

	return -2, errors.New("failed to select element")
}
