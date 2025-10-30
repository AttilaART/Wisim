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
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Game setup functions

func Load_sim_config(path string) (Sim_config, error) {
	sim_config_file, err := os.ReadFile(path)
	if err != nil {
		return Sim_config{}, errors.New("error loading sim_config at '" + path + "'")
	}

	var sim_config Sim_config
	err = json.Unmarshal(sim_config_file, &sim_config)
	if err != nil {
		return Sim_config{}, errors.New("error in sim_config.json")
	}

	return sim_config, nil
}

func generatePopulation(
	populationSize int,

	minBaseNeed int,
	maxBaseNeed int,

	qualityBias float32, // "bias" parameters increase the mean of the normal distributions
	qualitySpread float32, // "spread" parameters increase the standard deviation of the normal distributions
	ecologyBias float32,
	ecologySpread float32,
	ethicsBias float32,
	ethicsSpread float32,
	// coolnessBias float32,
	// coolnessSpread float32,
	priceBias float32,
	priceSpread float32,
	bangForBuckBias float32,
	bangForBuckSpread float32,
	durabiltyBias float32,
	durabilitySpread float32,
	purchasingThresholdBias float32,
	purchasingThresholdSpread float32,
	savvynessSpread float32,
	savvynessBias float32,
	baseMarketPrice float32,
	marketSaturation float32,

	number_of_companies int,
) ([]Customer, []Properties, error) {
	population := make([]Customer, populationSize)
	populationPreferences := make([]Properties, populationSize)

	// Handle errors
	if minBaseNeed <= 0 {
		return nil, nil, errors.New("min_base_need must be > 0")
	}
	if maxBaseNeed < minBaseNeed {
		return nil, nil, errors.New("max_base_need must be >= min_base_need")
	}

	if bangForBuckBias <= 0 {
		println("bang_for_buck_bias <= 0")
	}
	if bangForBuckSpread <= 0 {
		println("bang_for_buck_spread <= 0")
	}

	avr_max_price := 0.0

	var wg sync.WaitGroup
	for _, interval := range split_load(runtime.NumCPU(), populationSize) {
		wg.Add(1)
		go func(wg *sync.WaitGroup, interval Interval) {
			for i := interval.Start; i < interval.Stop_before; i++ {
				population[i] = Customer{
					Base_need: rand.Intn(maxBaseNeed-minBaseNeed) + minBaseNeed,

					Purchashing_threshold: float32(posNormFloat64())*purchasingThresholdSpread + purchasingThresholdBias,
					Loyalties:             make([]float32, number_of_companies),

					Savyness: float32(posNormFloat64())*savvynessSpread + savvynessBias,
				}
				populationPreferences[i][propertiesQuality] = float32(posNormFloat64())*qualitySpread + qualityBias
				populationPreferences[i][propertiesEcology] = float32(posNormFloat64())*ecologySpread + ecologyBias
				populationPreferences[i][propertiesEthics] = float32(posNormFloat64())*ethicsSpread + ethicsBias
				populationPreferences[i][propertiesPrice] = float32(posNormFloat64())*priceSpread + priceBias
				populationPreferences[i][propertiesBangForBuck] = float32(posNormFloat64())*bangForBuckSpread + bangForBuckBias
				populationPreferences[i][propertiesDurability] = float32(posNormFloat64())*durabilitySpread + durabiltyBias

				number_of_owned_products := int(posNormFloat64()*float64(population[i].Base_need) + float64(marketSaturation))
				for range number_of_owned_products {
					population[i].Owned_products = append(population[i].Owned_products,
						OwnedProduct{-1, int(posNormFloat64()*float64(population[i].Base_need) +
							float64(marketSaturation))})
				}

				population[i].Max_price = ((baseMarketPrice * 1.5) / populationPreferences[i][propertiesPrice]) * float32(posNormFloat64())
				avr_max_price += float64(population[i].Max_price)

				// fmt.Printf("|%6d|%6d|\n", i, customer.income)
			}
			wg.Done()
		}(&wg, interval)
	}
	wg.Wait()

	message.NewPrinter(language.BritishEnglish).Printf("avrg max price: %.2f\n", avr_max_price/float64(len(population)))
	return population, populationPreferences, nil
}

func posNormFloat64() float64 {
	num := rand.NormFloat64()
	if num < 0 {
		return -num
	}
	return num
}

func (g *GameState) Generate_new_employee_id() int {
	randID := rand.Int() % 9999999

	if _, exists := g.Employees[randID]; exists {
		return g.Generate_new_employee_id()
	}

	return randID
}

func (g *GameState) Generate_employee(basePay float32, workingHours float32, employeeType Employee_type, baseMotivation float32) (int, *Employee) {
	employeeeID := g.Generate_new_employee_id()

	g.employeesArray = append(g.employeesArray, Employee{
		ID:           employeeeID,
		Name:         randomName(employeeeID),
		EmployeeType: employeeType,
		Motivation:   baseMotivation,
		Skill:        float32(rand.NormFloat64()*0.1 + 1),
		Pay:          basePay,
		WorkingHours: workingHours,
		Employer:     Employee_employer_none,
	},
	)
	g.Employees[employeeeID] = &g.employeesArray[len(g.employeesArray)-1]

	return employeeeID, &g.employeesArray[len(g.employeesArray)-1]
}

func (g *GameState) RefillUnemployed(wantedNumberOfEmployees int, basePay float32, workingHours float32, employeeType Employee_type, baseMotivation float32) {
	numEmployees := 0
	for _, e := range g.Employees {
		if e.EmployeeType == employeeType && e.Employer == Employee_employer_none {
			numEmployees += 1
		}
	}

	if wantedNumberOfEmployees-numEmployees <= 0 {
		return
	}

	for range wantedNumberOfEmployees - numEmployees {
		g.Generate_employee(basePay, workingHours, employeeType, baseMotivation)
	}
}

func randomName(seed int) string {
	firstNames := []string{
		"Borg",
		"Tim",
		"Trim",
		"Andrin",
		"Attila",
		"Ali",
		"Brahim",
		"Fatima",
		"Bob",
		"Felix",
		"Yäl",
		"Jael",
		"Lizza",
		"Carmelon",
		"El-Ias",
		"Keira",
		"Kimberly",
		"John",
		"Heldegard",
		"Tinish",
		"Chasable",
		"Cecily",
		"Earnest",
		"Gwendilyn",
		"Patrik",
	}
	lastNames := []string{
		"Smith",
		"Brugger",
		"Schaginhaufen",
		"Gehfehler",
		"McChicken",
		"Ouegouri",
		"Fullagar",
		"Keller",
		"Acharya",
		"Mozarella",
		"Stift",
		"Pluss",
		"The Great",
		"Bunbury",
	}

	randomiser := rand.New(rand.NewSource(int64(seed)))

	return firstNames[randomiser.Int()%len(firstNames)] + " " + lastNames[randomiser.Int()%len(lastNames)]
}

func (g *GameState) generateCompanies(
	defaultCompany Company,
	numberOfCompanies int,
	externalFactors ExternalFactors,
	baseWorkingHours float32,
	baseNumberOfMarketingPersonelle int,
) []Company {
	// Make each company according to defaults & preferences
	companies := make([]Company, numberOfCompanies)

	for i := range numberOfCompanies {
		companies[i] = defaultCompany
		companies[i].ID = i
		companies[i].Name = "Unnamed Company"
		companies[i].Reports = make([]Report, 0)
		companies[i].DecisionHistory = make([]Decisions, 0)
		companies[i].employeePool = g.Employees
		companies[i].productComponents = &g.ProductComponents
		companies[i].Offers = make(map[string]Offer)
		companies[i].ProductsInStorage = make(map[string]int)
		if len(defaultCompany.Machines) == 0 {
			companies[i].Machines = make([]Machine, 0)
		}

		requiredProductionPersonelle := 0
		for _, m := range companies[i].Machines {
			requiredProductionPersonelle += m.RequiredWorkers
		}

		for range requiredProductionPersonelle {
			_, e := g.Generate_employee(
				externalFactors.ProductionMinimumWage,
				baseWorkingHours,
				Employee_type_production,
				1,
			)
			e.Employer = i
		}
		for range baseNumberOfMarketingPersonelle {
			_, e := g.Generate_employee(
				externalFactors.MarketingMinimumWage,
				baseWorkingHours,
				Employee_type_marketing,
				1,
			)
			e.Employer = i
		}

	}
	return companies
}

func NewGame(simConfig Sim_config, numberOfCompanies int, gameName string) GameState {
	var gameState GameState

	gameState.Step = 1
	gameState.GameName = gameName
	gameState.Employees = make(Employee_pool)

	gameState.ExternalFactors = ExternalFactors{
		Inflation:               0.005,
		EconomicSituationIndex:  1,
		TaxRate:                 0.147,
		MaterialPrice:           3.5,
		EnergyPrice:             96.2,
		MachineDepreciationRate: 0.1,

		IntrestRate:            0.04,
		BridgeLoansIntrestRate: 0.08,

		ProductionMinimumWage: 60000 / 12,
		MarketingMinimumWage:  80000 / 12,

		MachineOnOffer: Machine{
			ProductionCapacity: 1500,
			RequiredWorkers:    3,
			MinimumWorkers:     1,
			EnergyUse:          0.01,
			Value:              10000,
			MaintananceCost:    100,
		},
	}

	// Load product components
	bComponents, err := os.ReadFile("../../common/components.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bComponents, &gameState.ProductComponents)
	if err != nil {
		log.Fatal(err)
	}

	gameState.Population.Population, gameState.Population.Preferences, err = generatePopulation(
		simConfig.Population_size,
		simConfig.Min_base_need,
		simConfig.Max_base_need,
		simConfig.Quality_bias,
		simConfig.Quality_spread,
		simConfig.Ecology_bias,
		simConfig.Ecology_spread,
		simConfig.Ethics_bias,
		simConfig.Ethics_spread,
		simConfig.Price_bias,
		simConfig.Price_spread,
		simConfig.Bang_for_buck_bias,
		simConfig.Bang_for_buck_spread,
		simConfig.Durabilty_bias,
		simConfig.Durability_spread,
		simConfig.Purchasing_threshold_bias,
		simConfig.Purchasing_threshold_spread,
		simConfig.SavvynessSpread,
		simConfig.SavvynessBias,
		simConfig.Base_market_price,
		simConfig.Market_saturation,

		numberOfCompanies,
	)
	if err != nil {
		log.Fatal(err)
	}

	gameState.Companies = gameState.generateCompanies(
		simConfig.Default_company,
		numberOfCompanies,
		gameState.ExternalFactors,
		8,
		1,
	)

	gameState.CurrentDecisions = make([]Decisions, numberOfCompanies)

	for i := range gameState.CurrentDecisions {
		gameState.CurrentDecisions[i] = Decisions{
			General: struct{ CompanyName string }{
				CompanyName: gameState.Companies[i].Name,
			},
			Products: make(map[string]Decisions_product),
			Research: Decisions_research{
				Quality:         1000,
				Durability:      1000,
				Ecology:         1000,
				Promotion:       1000,
				Production_cost: 1000,
			},
			Production: struct {
				Machines                 []Delta[Machine]
				Logistics                []Delta[Warehouse]
				MachineAssignmentPattern AssignmentPattern
			}{
				Logistics:                make([]Delta[Warehouse], 0),
				Machines:                 make([]Delta[Machine], 0),
				MachineAssignmentPattern: FillMachines,
			},
			Employees: struct {
				ProductionDeltas []Delta[Employee]
				MarketingDeltas  []Delta[Employee]
				SeverancePay     float32
			}{
				ProductionDeltas: make([]Delta[Employee], 0),
				MarketingDeltas:  make([]Delta[Employee], 0),
				SeverancePay:     10000,
			},
		}
	}

	gameState.DecisionsSubmitted = make([]bool, numberOfCompanies)

	return gameState
}

func Load_game(path string) (GameState, error) {
	println("Loading game")

	var save_file []byte

	if path[len(path)-4:] == ".zip" {
		println("Decompressing save")
		r, err := zip.OpenReader(path)
		if err != nil {
			return GameState{}, err
		}
		defer r.Close()

		save_file_reader, err := r.File[0].Open()
		if err != nil {
			return GameState{}, err
		}
		defer save_file_reader.Close()

		save_file, err = io.ReadAll(save_file_reader)
		if err != nil {
			return GameState{}, err
		}

	} else {
		println("Opening file")
		file, err := os.Open(path)
		if err != nil {
			return GameState{}, err
		}
		save_file, err = io.ReadAll(file)
		if err != nil {
			return GameState{}, err
		}
	}

	var save SaveGame
	err := json.Unmarshal(save_file, &save)
	if err != nil {
		return GameState{}, err
	}

	var game_state GameState
	population_buffer := bytes.NewBuffer(save.Population)
	decoder := gob.NewDecoder(population_buffer)

	var population Population
	err = decoder.Decode(&population.Population)
	if err != nil {
		return GameState{}, err
	}

	game_state = save.GameState

	for i := range game_state.Companies {
		// fix employee pointer stuff
		game_state.Companies[i].employeePool = game_state.Employees
	}

	game_state.Population = population

	if len(game_state.Population.Population) == 0 {
		return game_state, errors.New("Failed to load population")
	}

	println("Successfully opened ", game_state.GameName)

	s, err := json.MarshalIndent(game_state.ExternalFactors, "", "    ")
	if err != nil {
		return game_state, err
	}

	println(string(s))

	return game_state, nil
}
