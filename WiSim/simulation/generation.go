package simulation

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"slices"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Game setup functions

func generate_population(
	population_size int,

	min_base_need int,
	max_base_need int,

	quality_bias float32, // "bias" parameters increase the mean of the normal distributions
	quality_spread float32, // "spread" parameters increase the standard deviation of the normal distributions
	ecology_bias float32,
	ecology_spread float32,
	coolness_bias float32,
	coolness_spread float32,
	price_bias float32,
	price_spread float32,
	bang_for_buck_bias float32,
	bang_for_buck_spread float32,
	durabilty_bias float32,
	durability_spread float32,
	purchasing_threshold_bias float32,
	purchasing_threshold_spread float32,
	base_market_price float32,

	number_of_companies int,
) ([]Customer, error) {
	population := make([]Customer, population_size)

	// Handle errors
	if min_base_need <= 0 {
		return nil, errors.New("min_base_need must be > 0")
	}
	if max_base_need < min_base_need {
		return nil, errors.New("max_base_need must be >= min_base_need")
	}

	avr_max_price := 0.0

	var wg sync.WaitGroup
	for _, interval := range split_load(runtime.NumCPU(), population_size) {
		wg.Add(1)
		go func(wg *sync.WaitGroup, interval Interval) {
			for i := interval.Start; i < interval.Stop_before; i++ {
				population[i] = Customer{
					Base_need: rand.Intn(max_base_need-min_base_need) + min_base_need,

					Quality_preference:       float32(PosNormFloat64())*quality_spread + quality_bias,
					Ecology_preference:       float32(PosNormFloat64())*ecology_spread + ecology_bias,
					Coolness_preference:      float32(PosNormFloat64())*coolness_spread + coolness_bias,
					Price_preference:         float32(PosNormFloat64())*price_spread + price_bias,
					Bang_for_buck_preference: float32(PosNormFloat64())*bang_for_buck_spread + bang_for_buck_bias,
					Durabilty_preference:     float32(PosNormFloat64())*durability_spread + durabilty_bias,
					Purchashing_threshold:    float32(PosNormFloat64())*purchasing_threshold_spread + purchasing_threshold_bias,
					Loyalties:                make([]float32, number_of_companies),
				}

				population[i].Max_price = ((base_market_price * 1.1) / population[i].Price_preference) * float32(PosNormFloat64()*100)
				avr_max_price += float64(population[i].Max_price)

				// fmt.Printf("|%6d|%6d|\n", i, customer.income)
			}
			wg.Done()
		}(&wg, interval)
	}
	wg.Wait()

	message.NewPrinter(language.BritishEnglish).Printf("avrg max price: %.2f\n", avr_max_price/float64(len(population)))
	return population, nil
}

func PosNormFloat64() float64 {
	var num float64
	num = rand.NormFloat64()
	if num < 0 {
		return -num
	}
	return num
}

func (g *Game_state) Generate_new_employee_id() int {
	// Find taken IDs
	var taken_ids []int

	for _, e := range g.Employees {
		taken_ids = append(taken_ids, e.Id)
	}

	// Generate lowest non-taken ID

	gen_id := 0
	for range g.Employees {
		if !slices.Contains(taken_ids, gen_id) {
			return gen_id
		}
		gen_id++
	}

	return len(g.Employees)
}

func (g *Game_state) Generate_employee(base_pay float32, working_hours float32, employee_type int, base_motivation float32) *Employee {
	employee := Employee{
		Id:            g.Generate_new_employee_id(),
		Employee_type: employee_type,
		Motivation:    base_motivation,
		Skill:         float32(rand.NormFloat64()*0.1 + 1),
		Pay:           base_pay,
		Working_hours: working_hours,
	}
	g.Employees = append(g.Employees, employee)

	return &g.Employees[len(g.Employees)-1]
}

func (g *Game_state) generate_companies(
	default_company Company,
	number_of_companies int,
	external_factors External_factors,
	base_working_hours float32,
	base_number_of_marketing_personelle int,
) []Company {
	default_company.employee_pool = &g.Employees
	// Make each company according to defaults & preferences
	companies := make([]Company, number_of_companies)

	for i := range number_of_companies {
		companies[i] = default_company
		companies[i].Id = i
		companies[i].Name = "Unnamed Company"

		required_production_personelle := 0
		for _, m := range companies[i].Machines {
			required_production_personelle += m.Required_workers
		}

		companies[i].Production_personelle = make([]*Employee, required_production_personelle)
		for ii := range companies[i].Production_personelle {
			companies[i].Production_personelle[ii] = g.Generate_employee(
				external_factors.Production_minimum_wage,
				base_working_hours,
				Production_employee,
				1,
			)
		}
		companies[i].Marketing_personelle = make([]*Employee, base_number_of_marketing_personelle)
		for ii := range companies[i].Marketing_personelle {
			companies[i].Marketing_personelle[ii] = g.Generate_employee(
				external_factors.Production_minimum_wage,
				base_working_hours,
				Marketing_employee,
				1,
			)
		}

	}
	return companies
}

func New_game(sim_config Sim_config, number_of_companies int, game_name string) Game_state {
	var game_state Game_state

	game_state.Step = -1
	game_state.Step_simulated = false
	game_state.Game_name = game_name

	game_state.External_factors = External_factors{
		Inflation:                 0.005,
		Economic_situation_index:  1,
		Tax_rate:                  0.147,
		Material_price:            10,
		Energy_price:              96.2,
		Machine_depreciation_rate: 0.1,

		Turnover:                0.08,
		Production_minimum_wage: 60000,
		Marketing_minimum_wage:  80000,

		Machine_on_offer: Machine{
			Production_capacity: 15000,
			Required_workers:    5,
			Minimum_workers:     3,
			Energy_use:          0.5,
			Value:               100000,
		},
	}

	// "Population_size": 1000000,
	// "min_base_need": 0,
	// "max_base_need": 4,
	// "quality_bias": 1,
	// "quality_spread": 1,
	// "ecology_bias": 1,
	// "ecology_spread": 1,
	// "coolness_bias": 1,
	// "coolness_spread": 1,
	// "price_bias": 1,
	// "price_spread": 1,
	// "bang_for_buck_bias": 1,
	// "bang_for_buck_spread": 1,
	// "durabilty_bias": 1,
	// "durability_spread": 1,
	// "purchasing_threshold_bias": 1,
	// "purchasing_threshold_spread": 1

	var err error
	game_state.Population.Population, err = generate_population(
		sim_config.Population_size,
		sim_config.Min_base_need,
		sim_config.Max_base_need,
		sim_config.Quality_bias,
		sim_config.Quality_spread,
		sim_config.Ecology_bias,
		sim_config.Ecology_spread,
		sim_config.Coolness_bias,
		sim_config.Coolness_spread,
		sim_config.Price_bias,
		sim_config.Price_spread,
		sim_config.Bang_for_buck_bias,
		sim_config.Bang_for_buck_spread,
		sim_config.Durabilty_bias,
		sim_config.Durability_spread,
		sim_config.Purchasing_threshold_bias,
		sim_config.Purchasing_threshold_spread,
		sim_config.Base_market_price,

		number_of_companies,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	game_state.Companies = game_state.generate_companies(
		sim_config.Default_company,
		number_of_companies,
		game_state.External_factors,
		8,
		5,
	)

	game_state.Current_decisions = make([]Decisions, number_of_companies)
	game_state.Decisions_submitted = make([]bool, number_of_companies)

	return game_state
}

func Load_game(path string) (Game_state, error) {
	println("Loading game")

	var save_file []byte

	if path[len(path)-4:] == ".zip" {
		println("Decompressing save")
		r, err := zip.OpenReader(path)
		if err != nil {
			return Game_state{}, err
		}
		defer r.Close()

		save_file_reader, err := r.File[0].Open()
		if err != nil {
			return Game_state{}, err
		}
		defer save_file_reader.Close()

		save_file, err = io.ReadAll(save_file_reader)
		if err != nil {
			return Game_state{}, err
		}

	} else {
		println("Opening file")
		file, err := os.Open(path)
		if err != nil {
			return Game_state{}, err
		}
		save_file, err = io.ReadAll(file)
		if err != nil {
			return Game_state{}, err
		}
	}

	var save Save_game
	err := json.Unmarshal(save_file, &save)
	if err != nil {
		return Game_state{}, err
	}

	var game_state Game_state
	population_buffer := bytes.NewBuffer(save.Population)
	decoder := gob.NewDecoder(population_buffer)

	var population Population
	err = decoder.Decode(&population.Population)
	if err != nil {
		return Game_state{}, err
	}

	game_state = save.Game_state

	for i := range game_state.Companies { // fix employee pointer stuff
		game_state.Companies[i].employee_pool = &game_state.Employees

		for ii, e := range game_state.Companies[i].Production_personelle {
			var err error
			game_state.Companies[i].Production_personelle[ii], err = game_state.Employees.find_employee_by_id(e.Id)
			if err != nil {
				panic(err)
			}
		}

		for ii, e := range game_state.Companies[i].Marketing_personelle {
			var err error
			game_state.Companies[i].Marketing_personelle[ii], err = game_state.Employees.find_employee_by_id(e.Id)
			if err != nil {
				panic(err)
			}
		}

		for ii := range game_state.Companies[i].Machines {
			for iii, e := range game_state.Companies[i].Machines[ii].Assigned_workers_ptr {
				var err error
				game_state.Companies[i].Machines[ii].Assigned_workers_ptr[iii], err = game_state.Employees.find_employee_by_id(e.Id)
				if err != nil {
					panic(err)
				}
			}
		}

		if len(game_state.Companies[i].Decision_history) >= 1 {
			for ii := range game_state.Companies[i].Decision_history {
				for iii := range game_state.Companies[i].Decision_history[ii].Employees.Marketing_actions {

					game_state.Companies[i].Decision_history[ii].Employees.Marketing_actions[iii].employee, err = game_state.Employees.find_employee_by_id(game_state.Companies[i].Decision_history[ii].Employees.Marketing_actions[iii].Employee_id)
					if err != nil {
						panic(err)
					}
				}
				for iii := range game_state.Companies[i].Decision_history[ii].Employees.Production_actions {

					game_state.Companies[i].Decision_history[ii].Employees.Production_actions[iii].employee, err = game_state.Employees.find_employee_by_id(game_state.Companies[i].Decision_history[ii].Employees.Production_actions[iii].Employee_id)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}

	game_state.Population = population

	if len(game_state.Population.Population) == 0 {
		return game_state, errors.New("Failed to load population")
	}

	println("Successfully opened ", game_state.Game_name)

	s, err := json.MarshalIndent(game_state.External_factors, "", "    ")
	if err != nil {
		return game_state, err
	}

	println(string(s))

	return game_state, nil
}
