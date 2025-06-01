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

	var wg sync.WaitGroup
	for _, interval := range split_load(runtime.NumCPU(), population_size) {
		wg.Add(1)
		go func(wg *sync.WaitGroup, interval Interval) {
			for i := interval.Start; i < interval.Stop_before; i++ {
				customer := Customer{
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
				population[i] = customer
				// fmt.Printf("|%6d|%6d|\n", i, customer.income)
			}
			wg.Done()
		}(&wg, interval)
	}
	wg.Wait()

	return population, nil
}

func PosNormFloat64() float64 {
	var num float64
	for {
		num = rand.NormFloat64()
		if num >= 0 {
			return num
		}
	}
}

func generate_new_employee_id(employees *[]Employee) int {
	// Find taken IDs
	var taken_ids []int

	for _, e := range *employees {
		taken_ids = append(taken_ids, e.Id)
	}

	// Generate lowest non-taken ID

	gen_id := 0
	for range *employees {
		if !slices.Contains(taken_ids, gen_id) {
			return gen_id
		}
		gen_id++
	}

	return len(*employees)
}

func generate_employee(id int, base_pay float32, working_hours float32, employee_type int, base_motivation float32) Employee {
	employee := Employee{
		Id:            id,
		Employee_type: employee_type,
		Motivation:    base_motivation,
		Skill:         float32(rand.NormFloat64()*0.1 + 1),
		Salary:        base_pay,
		Working_hours: working_hours,
	}
	return employee
}

func generate_companies(
	default_company Company,
	number_of_companies int,
	external_factors External_factors,
	base_working_hours float32,
	base_number_of_marketing_personelle int,
) []Company {
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

		companies[i].Production_personelle = make([]Employee, required_production_personelle)
		for ii := range companies[i].Production_personelle {
			companies[i].Production_personelle[ii] = generate_employee(
				generate_new_employee_id(&companies[i].Production_personelle),
				external_factors.Production_minimum_wage,
				base_working_hours,
				production_employee,
				1,
			)
		}
		companies[i].Marketing_personelle = make([]Employee, base_number_of_marketing_personelle)
		for ii := range companies[i].Marketing_personelle {
			companies[i].Marketing_personelle[ii] = generate_employee(
				generate_new_employee_id(&companies[i].Marketing_personelle),
				external_factors.Production_minimum_wage,
				base_working_hours,
				marketing_employee,
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

		number_of_companies,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	game_state.Companies = generate_companies(
		sim_config.Default_company,
		number_of_companies,
		game_state.External_factors,
		8,
		5,
	)

	game_state.Current_decisions = make([]Decisionsold, number_of_companies)
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
