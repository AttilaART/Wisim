# 30/03/2025
Today I finished implementing all the functions necessary to get a single step of the simulation to run.

The file is approaching 1400 lines of GO.

Now I'm making a temporary way of getting the decisions, else the program doesn't run.

As I'm trying to run the simulation for the first time, I'm realising there are so many functions that should error that don't yet, so I'm going around and trying to fix that.

---
I've just gotten the program to run for the first time and I've noticed one major problem:
![](../../attachments/Pasted%20image%2020250820172750.png)
*Graph of my memory usage when starting up the program (32 GB memory + 4 GB swapfile)*

Running the program on my mac, I saw the program use up to 22 GB of memory (mostly swapfile).
![](../../attachments/Pasted%20image%2020250820172824.png)
*Screenshot from activity monitor on my mac*

I think the main problem is that most functions get "Value" arguments instead of pointers. Considering the quantity of data being copied and the latency of the garbage collector, I'm wasting a lot of memory.

For example: the save file (in `.json`) is approx 500 MB, considering the fact that most of the storage space is being used by the `population` array, I found that the whole array was being copied at least 4 times:
```go
gamestate := new_game(...)
gamestate, err = simulate_step(gamestate, ...)

func simulate_step(game_state Gamestate) (Gamestate, error) {
	...
	simulate_economy(game_state.population, ...)
	...
}

func simulate_economy(population []Customer, ...)
```

# 02/04/2025
After editing all most functions to use pointers, I discovered the source of the insane memory usage: the `sort_employees()` function. The function implements [[selction sort]]. In short selection sort selects the lowest entry in the array, appends it to a new array and deletes it from the original array. My function had 2 major mistakes (I think I wrote it way too late at night…):
1. The condition controlling the behavior of the `for` loop was set incorrectly. (I, for some strange reason, wanted to iterate through the array backwards and made the loop never stop iterating.)
2. At every iteration of the loop, I would append to the `slice`[^1] instead of assigning by index (I had actually [allocated a slice](https://go.dev/tour/moretypes/13) of the correct length before running the loop).
These 2 mistakes made it so that the array would infinitely increase in size. (very bad).

---
After some more basic debugging I finally got the simulation to run a complete iteration.

I went to handle the problem of saving the state to disk. The easiest way was to simply save export the whole `game_state` struct into a `.json`file using the `json` module from the standard library.

So I went and did that!

The result: a 6 GB file\*!

\*when saving 10 million people...

The simulation needs to simulate up to multiple millions of individual customers, however simply dumping a 10 million entry long array into json, rather predictably, results in a LOT of duplication...

Instead, to remedy this, I have settled with using the `gob`binary encoder provided in the standard library and just replacing the `population` field with a binary representation of the data, instead of json objects.

I'm considering completely overhauling the storing of the population, perhaps by calculating a seed for the generation of preferences, instead of the preferences themselves, however for the moment it doesn't seem necessary.

---

**Anyway the simulation runs!**

Now begins the second phase of designing the simulation: fleshing out purchasing habits, balancing and otherwise going into detail about how things really work and how i'm going to implement them.

For example the product purchasing algorithm needs to use a weighted die to choose which product to buy (to prevent that if 2 products have the same stats, the first one will be chosen every time).
![](../../attachments/Pasted%20image%2020250820172934.png)
*In my test runs every company was identical, however only the first sold any products*

I also want to do some research into what makes products *cool*.

# 18.04.2025
After fixing the previous bug that caused people only to buy from the first company (By using a weighted dice roll instead of just choosing the best option), I went to tackle finances.

I spent approx a week designing an over-complicated system, where a single function running at the end of each simulation step would calculate every cost and input it into a massive struct that looked like this:

```go
type Financial_Report struct {
	// Finance:
	// Finance: Income Statement
	// Finance: Income Statement: Income
	Income_from_sales         float64
	Increase_of_bank_loan     float64
	Income_from_other_sources float64
	Total_income              float64

	// Finance: Costs

	// Finance: Costs: personelle
	Base_production_personelle_costs float64
	Production_training_costs        float64
	Production_bonuses               float64

	Base_marketing_personelle_costs float64
	Marketing_training_costs        float64
	Marketing_bonuses               float64

	//Base_product_development_personelle_costs float64
	//Product_development_training_costs        float64
	//Product_development_bonuses               float64
	//
	//Base_adminstrative_costs float64
	//Adminstrative_costs      float64
	//Adminstrative_bonuses    float64
	//
	//Base_management_costs float64
	//Management_costs      float64
	//Management_bonuses    float64

	Total_personelle_costs float64

	// Finance: Costs: Marketing
	Non_personelle_marketing_costs float64

	// Finance: Costs: Production
	Material_costs                       float64
	Energy_costs                         float64
	Write_offs                           float64
	Investment_in_quality_development    float64 // Increases Quality
	Investment_in_ecological_production  float64 // Decreases material use
	Investment_in_durability_development float64

	// Finance: Costs: Distribution & Fulfillment
	Distribution_costs     float64
	Local_storage_cost     float64
	External_storage_costs float64

	// Finacne: Costs: Investments
	Purchase_of_warehouses    float64
	Purchase_of_factory_space float64
	Purchase_of_real_estate   float64

	// Finance: financial costs
	Loan_repayments        float64
	Bridge_loan_repayments float64

	// Finance: Totals
	Non_cash_costs float64
	Cash_costs     float64
	Total_expenses float64
	EBIT           float64
	Taxes          float64

	Profit float64

	// Finance: Balance Sheet
	// Finance: Balance Sheet: Assets
	Stock float64

	Machines          float64
	Warehouses        float64
	Factory_space     float64
	Total_Real_estate float64

	Bank_balance               float64
	Number_of_stocks           float64
	Stocks_value               float64
	Interest_value             float64
	Bridge_loans_intrest_value float64

	Total_assets float64

	// Fincance: Balance Sheet: Liabilities
	Equity float64

	Bridge_loans   float64
	Base_bank_loan float64

	Total_credit float64
	Investment   float64

	Total_liabilities float64
}
```

At some point I realised this was overly complicated an inflexible. The problem is that the structure had a fixed number of slots, and if I wanted to add slots, it would break existing saves. Also calculating everything at once made it really easy to simply *forget* to calculate certain values. As well as the other problems, there was a confusing mix of totals and simple values. All in all it was very complicated...

Instead I switched to a simple model. I would use a struct containing 3 arrays, each representing a part of the balance sheet, liabilities, assets and the income statement. That way whenever a function did something that impacted finances, it, itself, would just add on a new entry.

```go
type Balance_sheet struct {
	Bank_balance float64

	Income_statement []FinanceReportEntry

	Assets      []FinanceReportEntry
	Liabilities []FinanceReportEntry
}

type FinanceReportEntry struct {
	Name      string
	Group     int
	Info      string
	Cash_cost bool

	Value float64
}

const (
	production = iota
	marketing
	personelle
	logistics
	materials
	energy
	product_development
	employee_training
	loans
	loan_intrest
	bridge_loans
	bridge_loan_intrest
	taxes
	sales
	other
)
```

The finance report would now be shortened to just a few totals that would be important for the player to know.

```go
type Financial_Report struct {
	Total_income float64

	Loan_repayments        float64
	Bridge_loan_repayments float64

	Non_cash_costs            float64
	Cash_costs                float64
	Total_expenses_before_tax float64
	Total_expenses_after_tax  float64
	EBIT                      float64
	Taxes                     float64

	Profit float64
}
```

This change also made the finances function (`calcualate_budget()`) much shorter and simpler.

Also loans would now be handled individually instead of just one big loan, this actually has no impact on Gameplay, I just think that's interesting.

---

Next up I modified the logic for choosing products, to actually implement the `purchasing_threashold`, so that products have to be good enough for people to buy them. This led to a rather interesting situation:

```zsh
=================================================== 
               Simulation step done!

===================== RESULTS ===================== 
Company 0: Unnamed Company:
Products sold: 0

Company 1: Unnamed Company:
Products sold: 0

Company 2: Unnamed Company:
Products sold: 0

Company 3: Unnamed Company:
Products sold: 0

Total products sold:  0
Saving file

```

No-one is buying anything. I told the program to print the average `decision_factor` and `purchasing_threashold` and it seems to point to the same problem.

```zsh
---------------- Simulatig economy ----------------
Average decision_factor: 0.253219
Average purchasing threshold: 1.796335
Simulatig economy done!
```

Yup. This means I need to find out how the products are absolutely terrible.

# 19.04.2025

It turns out products weren't bad, in fact they were extremely good, perhaps too good (unbalanced). The problem was the `choose_product()`. What it does is it gets a list of `decision factors`, normalises them, so their sum is equal to one, puts them on a number line, one after the other and then picks a random number. If the random number is in the range of a specific decision factor, that product gets chosen. (for a better explanation, see illustration below)

![](Development/Choose_product.md#^group=weL-8LOsqiR89INYyzb2b|1000)

The problem was that I was comparing the `purchasing_threashold`with the normalised `decision factors`instead of the original `decision factors`. This caused the normalised decision factors to always be less than the purchasing theashold, which meant that no-one was buying anything.

---
Although the simulation isn't done, my next priority is to get the basic UI done, so that I can actually play the game and test, how things are going. Because at the moment it's too annoying to open up save files and be inundated with information that *is* important for the simulation, but not for me right now. Also editing values atm is too tedious. I might develop an AI to try to optimise the initial values.
# 01/06/2025
## MULTITHREADING
I just updated the `simulate_economy` function to use goroutines to calculate the purchasing of products and...

### Hardware
The tests were run on a AMD Ryzen 7 5700X with 8 cores/16 threads & 32 GB of Memory.

### Setup
All implementations were run with 50'000'000 cims, with only the main loop in `simulate_economy` being timed.

### Results
Single-threaded:
![](../../attachments/Pasted%20image%2020250820173044.png)

Naive Multi-threaded:
![](../../attachments/Pasted%20image%2020250820173046.png)

Slightly smarter implementation:
![](../../attachments/Pasted%20image%2020250820173049.png)

### Implementation
#### Single threaded:
This is the base implementation. It just runs a loop through all the elements in `population.Population`

```go
func (population *Population) simulate_economy(companies *[]Company, external_factors External_factors) ([]FinanceReportEntry, []Purchasing_statistics, error) {
	// Get offers
	...

	// Main loop
	for i, c := range population.Population {
		customer := calcualte_durability(c)
		product_purchased, quanity, customer, individual_purchasing_statistics := calculate_purchase(customer, offers, avg_price, external_factors, product_availability)
		population.Population[i] = customer
		purchases[product_purchased] += quanity

		for i, s := range individual_purchasing_statistics {
			purchasing_statistics[i].Products_sold += s.Products_sold
			purchasing_statistics[i].Product_demand += s.Product_demand
			purchasing_statistics[i].Product_number = s.Product_number
			purchasing_statistics[i].Avr_decision_factor += s.Avr_decision_factor
			purchasing_statistics[i].Avr_purchasing_threshold += s.Avr_purchasing_threshold

			purchasing_statistics[i].Avr_quality_factor += s.Avr_quality_factor
			purchasing_statistics[i].Avr_durability_factor += s.Avr_durability_factor
			purchasing_statistics[i].Avr_ecology_factor += s.Avr_ecology_factor
			purchasing_statistics[i].Avr_price_factor += s.Avr_price_factor
			purchasing_statistics[i].Avr_coolness_factor += s.Avr_coolness_factor
		}
	}
	delta_time := time.Since(t_before)
	println("#### Time to calculate: ", delta_time.String())

	// More stuff
	...

```

#### Naive multi-threading:

This is almost the same as the original except it has multiple threads doing the calculations in a loop whose index is protected by a mutex (to prevent multiple threads from doing the same calculations multiple times). This makes it quite simple to implement as I don't need to make a load spreading system.

I think this still has a lot of room for improvement since the mutex forces each thread to wait for a long time for the mutex to be unlocked.

```go
func (population *Population) simulate_economy(companies *[]Company, external_factors External_factors) ([]FinanceReportEntry, []Purchasing_statistics, error) {
	// Get offers
	...

	// Main loop
	// Multithreading boilerplate
	var wg sync.WaitGroup
	latest_customer_to_handle := struct {
		index int
		mu    sync.Mutex
	}{}
	num_threads := runtime.NumCPU()
	
	wg.Add(num_threads)
	for id := range num_threads {
		go func(wg *sync.WaitGroup, latest_customer_to_handle *struct {
			index int
			mu    sync.Mutex
		}, id int,
		) {
			for {
				latest_customer_to_handle.mu.Lock()
				current_customer_index := latest_customer_to_handle.index
				// fmt.Printf("current_customer_index: %d, id: %d\n", current_customer_index, id)
				latest_customer_to_handle.index++
				latest_customer_to_handle.mu.Unlock()
	
				if current_customer_index >= len(population.Population) {
					break
				}
				current_customer := population.Population[current_customer_index]
	
				current_customer = calcualte_durability(current_customer)
				product_purchased, quanity, customer, individual_purchasing_statistics := calculate_purchase(current_customer, offers, avg_price, external_factors, product_availability)
				population.Population[current_customer_index] = customer
				purchases[product_purchased] += quanity
	
				for i, s := range individual_purchasing_statistics {
					purchasing_statistics[i].Products_sold += s.Products_sold
					purchasing_statistics[i].Product_demand += s.Product_demand
					purchasing_statistics[i].Product_number = s.Product_number
					purchasing_statistics[i].Avr_decision_factor += s.Avr_decision_factor
					purchasing_statistics[i].Avr_purchasing_threshold += s.Avr_purchasing_threshold
	
					purchasing_statistics[i].Avr_quality_factor += s.Avr_quality_factor
					purchasing_statistics[i].Avr_durability_factor += s.Avr_durability_factor
					purchasing_statistics[i].Avr_ecology_factor += s.Avr_ecology_factor
					purchasing_statistics[i].Avr_price_factor += s.Avr_price_factor
					purchasing_statistics[i].Avr_coolness_factor += s.Avr_coolness_factor
				}
			}
			wg.Done()
		}(&wg, &latest_customer_to_handle, id)
	}
	
	wg.Wait()
	delta_time := time.Since(t_before)
	println("#### Time to calculate: ", delta_time.String())
	
	// More stuff
	...

```

#### Slightly smarter multi-threading:

Instead of using a mutex, each thread was given a interval of indexes that it would handle. This makes each completely independent.

```go
func (population *Population) simulate_economy(companies *[]Company, external_factors External_factors) ([]FinanceReportEntry, []Purchasing_statistics, error) {
	// Get offers
	...

	// Main loop
	// Multithreading boilerplate
	var wg sync.WaitGroup
	num_threads := runtime.NumCPU()

	type interval struct {
		Start       int
		Stop_before int
	}

	thread_people_range := make([]interval, num_threads)

	count_per_thread := len(population.Population) / num_threads
	remainder := len(population.Population) % num_threads
	offset := 0

	for i := range thread_people_range {
		thread_people_range[i].Start = offset
		thread_people_range[i].Stop_before = offset + count_per_thread
		offset += count_per_thread

		if remainder > 0 {
			thread_people_range[i].Stop_before += 1
			remainder -= 1
			offset += 1
		}
	}

	wg.Add(num_threads)
	for id := range num_threads {
		go func(wg *sync.WaitGroup, population_range interval, id int,
		) {
			for current_customer_index := population_range.Start; current_customer_index < population_range.Stop_before; current_customer_index++ {

				if current_customer_index >= len(population.Population) {
					log.Panic("current_customer_index higher than len")
					break
				}
				current_customer := population.Population[current_customer_index]

				current_customer = calcualte_durability(current_customer)
				product_purchased, quanity, customer, individual_purchasing_statistics := calculate_purchase(current_customer, offers, avg_price, external_factors, product_availability)
				population.Population[current_customer_index] = customer
				purchases[product_purchased] += quanity

				for i, s := range individual_purchasing_statistics {
					purchasing_statistics[i].Products_sold += s.Products_sold
					purchasing_statistics[i].Product_demand += s.Product_demand
					purchasing_statistics[i].Product_number = s.Product_number
					purchasing_statistics[i].Avr_decision_factor += s.Avr_decision_factor
					purchasing_statistics[i].Avr_purchasing_threshold += s.Avr_purchasing_threshold

					purchasing_statistics[i].Avr_quality_factor += s.Avr_quality_factor
					purchasing_statistics[i].Avr_durability_factor += s.Avr_durability_factor
					purchasing_statistics[i].Avr_ecology_factor += s.Avr_ecology_factor
					purchasing_statistics[i].Avr_price_factor += s.Avr_price_factor
					purchasing_statistics[i].Avr_coolness_factor += s.Avr_coolness_factor
				}
			}
			wg.Done()
		}(&wg, thread_people_range[id], id)
	}
	wg.Wait()

	delta_time := time.Since(t_before)
	println("#### Time to calculate: ", delta_time.String())
	
	// More stuff
	...

```
## What's next?
As we saw, multi-threading is extremely powerful. I'm going to implement it in other similar functions like the `generate_population` function, as it's currently one of the slowest (painfully so).

Maybe if I'm really performance constrained I could try to leverage the GPU to compute these large amounts of data.

# 02/06/2026
## Further performance improvements
By trimming down the amount of tasks necessary (eg. reducing number of useless in-between steps, checking if a calculation is necessary, etc.). This made my the simulation run the benchmark of 50 million cims in 760 ms![^2]

![](../../attachments/Pasted%20image%2020250820173133.png)
---

---

[^1]: In `Go` slices function similarly to dynamically sized arrays like in Javascript.
[^2]: Note after the fact: a lot of this speed up was because calculations were being skipped, when no more products were available.
 
