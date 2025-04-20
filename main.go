package main

import (
	"bytes"
	"cmp"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Save_game struct {
	Population []byte // Binary
	Game_state Game_state
}

type Game_state struct {
	Step           int
	Step_simulated bool
	Game_name      string

	Population        Population
	Companies         []Company
	Current_decisions []Decisions
	External_factors  External_factors
}

type Sim_config struct {
	Default_company Company
	Population_size int

	Min_base_need int
	Max_base_need int

	Quality_bias                float32 // "bias" parameters increase the mean of the normal distributions
	Quality_spread              float32 // "spread" parameters increase the standard deviation of the normal distributions
	Ecology_bias                float32
	Ecology_spread              float32
	Coolness_bias               float32
	Coolness_spread             float32
	Price_bias                  float32
	Price_spread                float32
	Bang_for_buck_bias          float32
	Bang_for_buck_spread        float32
	Durabilty_bias              float32
	Durability_spread           float32
	Purchasing_threshold_bias   float32
	Purchasing_threshold_spread float32
}

// ##########################################################################################
// ##########       _         _                  _                       _         ##########
// ##########      | |       | |                | |                     | |        ##########
// ##########    __| |  __ _ | |_  __ _     ___ | |_  _ __  _   _   ___ | |_  ___  ##########
// ##########  / _` | / _` || __|/ _` |   / __|| __|| '__|| | | | / __|| __|/ __|  ##########
// ########## | (_| || (_| || |_| (_| |   \__ \| |_ | |   | |_| || (__ | |_ \__ \  ##########
// ########## \__,_| \__,_| \__|\__,_|   |___/ \__||_|    \__,_| \___| \__||___/   ##########
// ##########################################################################################

type Company struct {
	// General
	Id               int
	Name             string
	Global_effects   []Effect
	Decision_history []Decisions
	// Product
	Offer                Offer
	Orders               int
	Marketing_personelle []Employee

	// Fulfillment
	Warehouses       []Warehouse
	Items_in_storage int

	// Production
	Machines              []Machine
	Production_personelle []Employee

	// Finances
	Total_income   float64
	Total_expenses float64
	Balance        float64
	Reports        []Report
}

type Decisions struct {
	// Planing for budget reports
	Sales_projection int
	Selling_price    float32

	// Quality_of_stores int
	Marketing float32 // Increases probabilty that a potential customer hears of your product

	Investment_in_quality_development    float32 // Increases Quality
	Investment_in_ecological_production  float32 // Decreases material use
	Investment_in_durability_development float32 // increases durabilty
	// Investment_in_coolness_development  float32 // Increases coolness

	Production_target    int
	Purchase_of_machines int
	Selling_of_machines  int

	Material_quality                float32
	Percentage_of_ecological_energy float32

	Purchase_of_warehouses int

	New_hires_in_production         int
	New_hires_in_marketing          int
	Base_pay_for_production         float32
	Base_pay_for_marketing          float32
	Raise_for_production_personelle float32
	Raise_for_marketing_personelle  float32

	Working_hours_for_production float32
	Working_hours_for_marketing  float32

	Investment_in_production_training float32 // Skilled production staff increase quality & give bonus production
	Investment_in_marketing_training  float32 // Skilled marketing staff increase quality of marketing

	Increase_of_loans float32
	Dividends         float32
}

type Offer struct {
	Product            Product
	Price              float32
	Marketing_quantity float32
	Marketing_strength float32
}

type Product struct {
	Id                         int
	Name                       string
	Weight                     float32
	Material_use               float32
	Material_quality           float32
	Skill_factor               float32
	Quality_development_factor float32
	Quality_factor             float32
	Ecology_factor             float32
	Coolness_factor            float32

	Durabilty int
}

type Machine struct {
	Production_capacity int
	Required_workers    int
	Assigned_workers    []Employee
	Energy_use          float32
	Value               float32
}

type Warehouse struct {
	Id              int
	Capacity        int
	Operating_costs float32
	Value           float32
}

type Employee struct {
	Id            int
	Employee_type int
	Motivation    float32
	Skill         float32
	Global_effect *Effect
	Salary        float32
	Bonus         float32
	Working_hours float32
}

const (
	production_employee = iota
	marketing_employee
	executive_employee
)

func (e Employee) get_type() string {
	return []string{"production", "marketing", "executive"}[e.Employee_type]
}

type Effect struct {
	Id          int
	Name        string
	Description string
	// Effect_function *Effect_Function
}

type Effect_Function func(int)

type Report struct {
	Month int

	Financial_Report  Financial_Report
	Balance_sheet     Balance_sheet
	Personelle_report Personelle_report
	Production_report Production_report
	Sales_report      Sales_report
}

type Financial_Report struct {
	Total_income float64

	Loan_repayments        float64
	Bridge_loan_repayments float64
	New_bridge_loan        float64

	Non_cash_costs            float64
	Cash_costs                float64
	Total_expenses_before_tax float64
	Total_expenses_after_tax  float64
	EBIT                      float64
	Taxes                     float64

	Profit float64
}

func (f *Balance_sheet) add_to_income_statement(name string, group int, info string, cash_cost bool, value float64) {
	f.Income_statement = append(f.Income_statement, FinanceReportEntry{name, group, info, cash_cost, value})
}

func (f *Balance_sheet) add_to_equity(name string, group int, info string, cash_cost bool, value float64) {
	f.Assets = append(f.Assets, FinanceReportEntry{name, group, info, cash_cost, value})
}

func (f *Balance_sheet) add_to_liabilities(name string, group int, info string, cash_cost bool, value float64) {
	f.Liabilities = append(f.Liabilities, FinanceReportEntry{name, group, info, cash_cost, value})
}

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

type Personelle_report struct {
	// General info
	Number_of_employees float32
	Number_of_hires     float32

	Avg_pay          float32
	Minimum_pay      float32
	Maximum_pay      float32
	Standard_dev_pay float32

	Minimum_skill      float32
	Maximum_skill      float32
	Avg_skill          float32
	Standard_dev_skill float32

	top_employees    []Employee
	bottom_employees []Employee

	// Marketing
	Marketing_Number_of_employees float32
	Marketing_Number_of_hires     float32

	Marketing_Avg_pay          float32
	Marketing_Minimum_pay      float32
	Marketing_Maximum_pay      float32
	Marketing_Standard_dev_pay float32

	Marketing_Minimum_skill      float32
	Marketing_Maximum_skill      float32
	Marketing_Avg_skill          float32
	Marketing_Standard_dev_skill float32

	// Production
	Production_Number_of_employees float32
	Production_Number_of_hires     float32

	Production_Avg_pay          float32
	Production_Minimum_pay      float32
	Production_Maximum_pay      float32
	Production_Standard_dev_pay float32

	Production_Minimum_skill      float32
	Production_Maximum_skill      float32
	Production_Avg_skill          float32
	Production_Standard_dev_skill float32

	Production_avg_productivity      float32
	Production_minimum_productivity  float32
	Prouduction_maximum_productivity float32
}

type Production_report struct {
	Machines_purchased       []Machine
	Machines_sold            []Machine
	Worker_surplus           int
	Avg_machine_productivity float32
	// Max_machine_productivity float32
	// Min_machine_productivity float32

	Products_produced int
	Base_production   int
	Bonus_production  int

	Material_used float32
	Energy_used   float32

	Warehouses_bought int
}

type Sales_report struct {
	Products_sold int
}

type Customer struct {
	Base_need      int
	Owned_products []Owned_product

	Quality_preference       float32
	Ecology_preference       float32
	Coolness_preference      float32
	Price_preference         float32
	Bang_for_buck_preference float32 // Price / Quality preference
	Durabilty_preference     float32

	Purchashing_threshold float32
	Satisfaction          []Satisfaction

	Brand_loyalty_factor float32
	Loyalties            []float32
}

type Population struct {
	Population []Customer
}

type Owned_product struct {
	Id                  int
	Remaining_durabilty int
}

type Satisfaction struct {
	Product_id      int
	Decision_factor float32
	Satisfaction    float32
}

type External_factors struct {
	Month int
	// Economy
	Inflation                 float32
	Intrest_rate              float32
	Bridge_loans_intrest_rate float32
	economic_situation_index  float32
	tax_rate                  float32 // as decimal

	// Personelle
	Turnover                float32
	Production_minimum_wage float32
	Marketing_minimum_wage  float32

	// Prdoction
	Machine_on_offer       Machine
	External_storage_price float32 // per item
	Energy_price           float32 // per unit of energy
	Material_price         float32 // per unit of material

	Machine_depreciation_rate float32 // in decimal
}

type Purchasing_statistics struct {
	Products_sold  int
	Product_number int

	Avr_decision_factor      float32
	Avr_purchasing_threshold float32

	Product_quality       float32
	Product_durabilty     int
	Product_ecology       float32
	Product_price_factor  float32
	Product_price         float32
	Product_coolness      float32
	Product_bang_for_buck float32

	Avr_quality_factor       float32
	Avr_durability_factor    float32
	Avr_ecology_factor       float32
	Avr_price_factor         float32
	Avr_coolness_factor      float32
	Avr_bang_for_buck_factor float32
}

// #####################################################################################################
// ##########     _         _              __                      _    _                     ##########
// ##########    | |       | |            / _|                    | |  (_)                    ##########
// ##########  __| |  __ _ | |_  __ _    | |_  _   _  _ __    ___ | |_  _   ___   _ __   ___  ##########
// ########## / _` | / _` || __|/ _` |   |  _|| | | || '_ \  / __|| __|| | / _ \ | '_ \ / __| ##########
// ##########| (_| || (_| || |_| (_| |   | |  | |_| || | | || (__ | |_ | || (_) || | | |\__ \ ##########
// ########## \__,_| \__,_| \__|\__,_|   |_|   \__,_||_| |_| \___| \__||_| \___/ |_| |_||___/ ##########
// #####################################################################################################

func rand_income(mean_income int, standard_dev int) int {
	income := -1
	for income < 1000 {
		income = int(rand.NormFloat64()*float64(standard_dev)) + mean_income
	}
	return income
}

//  func pop(slice []any, index ...int) []any {
//    index = slices.Sorted(index)
//
//    new_slice := slice[0 : index[0]-1]
//    for i, idx := range index[1:] {
//      if i != len(index) {
//        new_slice = append(new_slice, slice[idx+1:index[i+1]-1])
//      } else {
//        new_slice = append(new_slice, slice[idx+1:])
//      }
//    }
//    return new_slice
//  }

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

	for i := range population_size {
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

func generate_employee(base_pay float32, working_hours float32, employee_type int, base_motivation float32) Employee {
	employee := Employee{
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
				external_factors.Production_minimum_wage,
				base_working_hours,
				production_employee,
				1,
			)
			companies[i].Production_personelle[ii].Id = ii
		}
		companies[i].Marketing_personelle = make([]Employee, base_number_of_marketing_personelle)
		for ii := range companies[i].Marketing_personelle {
			companies[i].Marketing_personelle[ii] = generate_employee(
				external_factors.Production_minimum_wage,
				base_working_hours,
				marketing_employee,
				1,
			)
			companies[i].Marketing_personelle[ii].Id = ii
		}

	}
	return companies
}

// Gameloop functions
func (game_state *Game_state) simulate_step() error {
	if len(game_state.Companies) != len(game_state.Current_decisions) {
		return errors.New("amount of decisions does not match number of companies")
	}

	println("=============== Simulating companies ==============")
	for i := range game_state.Companies {
		fmt.Printf("--------------- Simulating company %d -------------- \n", i)
		err := game_state.Companies[i].simulate_company(game_state.Current_decisions[i], game_state.External_factors)
		if err != nil {
			return err
		}
	}
	println("Simulating companies done!")

	println("---------------- Simulatig economy ----------------")
	Results, purchasing_statistcs, err := game_state.Population.simulate_economy(&game_state.Companies, game_state.External_factors)
	if err != nil {
		log.Fatal(err.Error())
	}
	println("Simulatig economy done!")

	println("================ Compiling reports =============== ")
	for i := range game_state.Companies {
		fmt.Printf("Compiling reports for company %d\n", i)
		game_state.Companies[i].compile_reports(game_state.Current_decisions[i], Results[i], game_state.External_factors)

	}
	println("=================================================== ")
	println("               Simulation step done!\n")
	println("===================== RESULTS ===================== ")
	for i, c := range game_state.Companies {
		fmt.Printf("Company %d: %s:\n", i, c.Name)
		fmt.Printf("Products sold: %d\n", c.Reports[len(c.Reports)-1].Sales_report.Products_sold)
		println("")
	}

	total_products_sold := 0
	for _, c := range game_state.Companies {
		total_products_sold += c.Reports[len(c.Reports)-1].Sales_report.Products_sold
	}

	println("Total products sold: ", total_products_sold)

	println("============== Purchasing statistics ============== ")

	for _, p := range purchasing_statistcs {
		s, e := json.MarshalIndent(p, "", "    ")
		if e != nil {
			return e
		}
		println(string(s))
	}
	return nil
}

func (company *Company) simulate_company(decisions Decisions, external_factors External_factors) error {
	company.Reports = append(company.Reports, Report{Month: external_factors.Month})

	// Personelle
	println("Simulatig marketing personelle")
	var err error
	company.Marketing_personelle, err = simulate_employees(
		company.Marketing_personelle,
		decisions.New_hires_in_marketing,
		external_factors,
		decisions.Investment_in_marketing_training,
		decisions.Raise_for_marketing_personelle,
		decisions.Working_hours_for_marketing,
		external_factors.Marketing_minimum_wage,
		decisions.Base_pay_for_marketing,
		&company.Reports[len(company.Reports)-1].Balance_sheet,
	)
	if err != nil {
		return err
	}
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
		"Marketing personelle training",
		employee_training,
		"Training to improve the skill of your employees",
		true,
		float64(-decisions.Investment_in_marketing_training))

	println("Simulatig production personelle")
	company.Production_personelle, err = simulate_employees(
		company.Production_personelle,
		decisions.New_hires_in_production,
		external_factors,
		decisions.Investment_in_production_training,
		decisions.Raise_for_production_personelle,
		decisions.Working_hours_for_production,
		external_factors.Production_minimum_wage,
		decisions.Base_pay_for_production,
		&company.Reports[len(company.Reports)-1].Balance_sheet,
	)
	if err != nil {
		return err
	}
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
		"Production personelle training",
		employee_training,
		"Training to improve the skill of your employees",
		true,
		float64(-decisions.Investment_in_production_training))

	// Offer
	println("Calculating product stats")
	company.Offer.Marketing_strength = marketing_strength(decisions.Marketing, company.Marketing_personelle)

	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
		"Marketing investment",
		marketing,
		"Basic investment in marketing, the quantity of marketing",
		true,
		float64(-decisions.Marketing),
	)
	// Offer: Product

	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Quality development investment", product_development, "", true, float64(-decisions.Investment_in_quality_development))
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Ecology development investment", product_development, "", true, float64(-decisions.Investment_in_ecological_production))
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Durability development investment", product_development, "", true, float64(-decisions.Investment_in_durability_development))

	// calculate avg skill of production personelle -> influences Quality_factor
	var total_production_skill float32 = 0
	for _, e := range company.Production_personelle {
		total_production_skill += e.Skill
	}
	avg_production_skill := total_production_skill / float32(len(company.Production_personelle))

	company.Offer.Product.calculate_quality(
		decisions.Investment_in_quality_development,
		avg_production_skill,
		decisions.Material_quality,
	)

	company.Offer.Product.calculate_ecology(
		1,
		company.Offer.Product.Material_use,
		decisions.Percentage_of_ecological_energy,
	)

	// Production
	println("Calculating production")
	company.calculate_production(decisions, external_factors, &company.Reports[len(company.Reports)-1].Balance_sheet)

	// Logistics
	company.Items_in_storage += company.Reports[len(company.Reports)-1].Production_report.Base_production + company.Reports[len(company.Reports)-1].Production_report.Bonus_production

	println("Calculating logistics")
	Balance_entries, Assets_entries, err := company.calculate_logistics(decisions)
	if err != nil {
		return err
	}

	company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement = append(company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement, Balance_entries...)
	company.Reports[len(company.Reports)-1].Balance_sheet.Assets = append(company.Reports[len(company.Reports)-1].Balance_sheet.Assets, Assets_entries...)
	// Finances

	company.Decision_history = append(company.Decision_history, decisions)

	return nil
}

// Production functions
func (company *Company) calculate_production(decisions Decisions, external_factors External_factors, finance_report *Balance_sheet) {
	var production_report Production_report
	// Purchase Machines
	println("Purchasing machines")
	Machines, production_report, Balance_entries, Assets_entries, err := calculate_machines(
		company.Machines,
		decisions.Purchase_of_machines,
		decisions.Purchase_of_machines,
		external_factors.Machine_on_offer,
		production_report,
		external_factors,
	)
	if err != nil {
		println(err.Error())
		log.Fatal("purchase_machines has errored, fix this properly, Attila!")
	}

	finance_report.Income_statement = append(finance_report.Income_statement, Balance_entries...)
	finance_report.Assets = append(finance_report.Assets, Assets_entries...)
	company.Machines = Machines // doing this to avoid compiler error

	// Produce
	println("Assigning workers")
	company.Machines, production_report.Worker_surplus = assign_workers(company.Machines, company.Production_personelle)

	println("Producing products")
	production_report, Income_entries := produce(company.Machines, company.Offer.Product, production_report, external_factors)
	company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement = append(company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement, Income_entries...)

	company.Reports[len(company.Reports)-1].Production_report = production_report
}

func calculate_machines(
	machines []Machine,
	machines_to_purchase int,
	machines_to_sell int,
	machines_for_sale Machine,
	production_report Production_report,
	external_factors External_factors,
) ([]Machine, Production_report, []FinanceReportEntry, []FinanceReportEntry, error) {
	production_report.Machines_purchased = []Machine{}
	production_report.Machines_sold = []Machine{}
	var Income_entries []FinanceReportEntry
	var Assets_entries []FinanceReportEntry

	for i, m := range machines {
		Income_entries = append(Income_entries, FinanceReportEntry{"Machine depreciation", production, "", false, -float64(m.Value) * float64(external_factors.Machine_depreciation_rate)})
		machines[i].Value = m.Value - m.Value*external_factors.Machine_depreciation_rate
	}

	if machines_to_sell > 0 {
		for range machines_to_sell {
			production_report.Machines_purchased = append(production_report.Machines_purchased, machines[0])
			Income_entries = append(Income_entries, FinanceReportEntry{"Selling of machine", production, "", true, float64(machines[0].Value)})
			machines = machines[1:]
		}
	} else if machines_to_sell < 0 {
		err := errors.New("machines_to_sell is less than 0")
		return nil, production_report, Income_entries, Assets_entries, err
	}

	if machines_to_purchase > 0 {
		for range machines_to_purchase {
			production_report.Machines_purchased = append(production_report.Machines_purchased, machines_for_sale)
			Income_entries = append(Income_entries, FinanceReportEntry{"Purchase of machine", production, "", true, float64(-machines_for_sale.Value)})
			machines = append(machines, machines_for_sale)
		}
	} else if machines_to_purchase < 0 {
		err := errors.New("machines_to_purchase is less than 0")
		return nil, production_report, Income_entries, Assets_entries, err
	}

	for _, m := range machines {
		Assets_entries = append(Assets_entries, FinanceReportEntry{"Machine", production, "Machines, with help of employees, produce your product", false, float64(m.Value)})
	}

	return machines, production_report, Income_entries, Assets_entries, nil
}

// returns machines & number of unassigned wokers (if not enough workers for the machines, it return a negative int)
func assign_workers(machines []Machine, workers []Employee) ([]Machine, int) {
	println("Sorting Employees")

	workers = slices.SortedFunc(func(yield func(Employee) bool) {
		for _, w := range workers {
			if !yield(w) {
				return
			}
		}
	}, func(a, b Employee) int {
		vala := (a.Motivation * a.Skill * a.Working_hours)
		valb := (b.Motivation * b.Skill * a.Working_hours)
		if vala < valb {
			return 1
		} else if vala == valb {
			return 0
		}
		return -1
	})

	machines = slices.SortedFunc(func(yield func(Machine) bool) {
		for _, m := range machines {
			if !yield(m) {
				return
			}
		}
	}, func(a, b Machine) int {
		if a.Production_capacity < b.Production_capacity {
			return 1
		} else if a.Production_capacity == b.Production_capacity {
			return 0
		}
		return -1
	})

	var Worker_surplus int
	ii := 0
	for i := range machines {
		machines[i].Assigned_workers = make([]Employee, 0)

		for range machines[i].Required_workers {
			if ii >= (len(workers) - 1) {
				break
			}
			machines[i].Assigned_workers = append(machines[i].Assigned_workers, workers[ii])
			ii++
		}
		Worker_surplus += len(machines[i].Assigned_workers) - machines[i].Required_workers
	}

	if Worker_surplus >= 0 {
		Worker_surplus = len(workers) - ii - 1
	}

	return machines, Worker_surplus
}

func produce(machines []Machine, product Product, production_report Production_report, external_factors External_factors) (Production_report, []FinanceReportEntry) {
	Income_entries := make([]FinanceReportEntry, 2)
	base_production := 0
	bonus_production := 0

	energy_use := 0.0
	for _, m := range machines {
		base_prod_of_machine, bonus_prod_of_machine := production_of_machine(m)
		base_production += base_prod_of_machine
		bonus_production += bonus_prod_of_machine

		energy_use += float64(m.Energy_use)
	}

	production_report.Products_produced = base_production + bonus_production
	production_report.Base_production = base_production
	production_report.Bonus_production = bonus_production

	production_report.Material_used = product.Material_use * float32(production_report.Products_produced)
	production_report.Energy_used = float32(energy_use)
	Income_entries[0] = FinanceReportEntry{"Material costs", 2, "The cost of materials used in your products", true, -round(float64(external_factors.Material_price)*float64(production_report.Material_used), 2)}
	Income_entries[1] = FinanceReportEntry{"Energy costs", 2, "The cost of energy used by machines in production", true, -round(float64(external_factors.Energy_price)*float64(production_report.Energy_used), 2)}

	production_report.Avg_machine_productivity = float32(production_report.Products_produced) / float32(len(machines))

	return production_report, Income_entries
}

// return (base production, bonus production)
func production_of_machine(machine Machine) (int, int) {
	// calculate averages
	var skill float32 = 0
	var motivation float32 = 0
	var working_hours float32 = 0

	for _, e := range machine.Assigned_workers {
		skill += e.Skill
		motivation += e.Motivation
		working_hours += e.Working_hours
	}

	num_assigned_workers := len(machine.Assigned_workers)
	skill = skill / float32(num_assigned_workers)
	motivation = motivation / float32(num_assigned_workers)
	working_hours = working_hours / float32(num_assigned_workers)

	base_production := int(float32(machine.Production_capacity) * (working_hours / 8))
	bonus_production := int(float32(base_production)*skill*motivation) - base_production

	return base_production, bonus_production
}

// Logistics
func (company *Company) calculate_logistics(decisions Decisions) ([]FinanceReportEntry, []FinanceReportEntry, error) {
	var Income_entries []FinanceReportEntry
	var Assets_entries []FinanceReportEntry

	if decisions.Purchase_of_warehouses > 0 {
		existing_warehouses := len(company.Warehouses)

		for i := range decisions.Purchase_of_warehouses {
			company.Warehouses = append(company.Warehouses, Warehouse{
				existing_warehouses + i,
				20000,
				20000,
				100000,
			})
			Income_entries = append(Income_entries, FinanceReportEntry{"Purchase of warehouse", logistics, "", false, -100000})
		}
	} else if (-decisions.Purchase_of_warehouses) == len(company.Warehouses) {
		company.Reports[len(company.Reports)].Production_report.Warehouses_bought = len(company.Warehouses)
		company.Warehouses = []Warehouse{}
	} else if decisions.Purchase_of_warehouses < 0 {
		if -decisions.Purchase_of_warehouses > len(company.Warehouses) {
			return Income_entries, Assets_entries, errors.New("can't sell more warehouses than you have")
		}
		for range -decisions.Purchase_of_warehouses {
			Income_entries = append(Income_entries, FinanceReportEntry{
				"Selling of warehouse",
				logistics,
				"",
				true,
				float64(company.Warehouses[len(company.Warehouses)-1].Value),
			})
			company.Warehouses = company.Warehouses[1:len(company.Warehouses)]
		}

		for _, w := range company.Warehouses {
			Assets_entries = append(Assets_entries, FinanceReportEntry{"Warehouse as asset", logistics, "", false, float64(w.Value)})
		}
	}

	company.Reports[len(company.Reports)-1].Production_report.Warehouses_bought = decisions.Purchase_of_warehouses
	return Income_entries, Assets_entries, nil
}

// offer functions
func marketing_strength(marketing_investment float32, marketing_personelle []Employee) float32 {
	// Temporary method
	var total_personelle_strength float32 = 1.0
	for _, e := range marketing_personelle {
		total_personelle_strength += e.Motivation * e.Skill * (e.Working_hours / 8.0)
	}

	return marketing_investment * float32(total_personelle_strength)
}

// Product attribute functions
func (product *Product) calculate_ecology(base_ecology float32, material_use float32, percentage_of_ecological_energy float32) {
	product.Ecology_factor = (base_ecology * (material_use/1 +
		percentage_of_ecological_energy/10 +
		product.Quality_factor/10)) / 2
}

func (product *Product) calculate_quality(quality_development_investment float32, production_skill float32, material_quality float32) {
	product.Material_quality = material_quality
	product.Skill_factor = production_skill
	product.Quality_development_factor += quality_development_investment / 1000000

	product.Quality_factor = product.Material_quality * product.Skill_factor * product.Quality_development_factor
}

// Finances

func (company *Company) calculate_budget(decisions Decisions, external_factors External_factors) Financial_Report {
	var financial_report Financial_Report

	local_storage_capacity := 0
	local_storage_costs := 0.0
	for _, w := range company.Warehouses {
		local_storage_costs += float64(-w.Operating_costs)
		local_storage_capacity += w.Capacity
	}
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
		"Internal storage costs",
		logistics,
		"The operating costs of our own warehouses",
		true,
		local_storage_costs)
	items_in_external_storage := (company.Items_in_storage - local_storage_capacity)
	if items_in_external_storage > 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
			"External storage",
			logistics,
			"The cost of storing products in external warehouses, this happens when our own are full",
			true,
			-(float64(external_factors.External_storage_price) * float64(company.Items_in_storage)))
	}

	// Loans

	// Get loans from last year
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
		if e.Group == loans || e.Group == bridge_loans {
			company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities = append(company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities, e)
		}
	}

	// Calculate interest
	intrest := 0.0
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
		if e.Group == loans {
			intrest += e.Value * float64(external_factors.Intrest_rate)
		}
	}

	if intrest > 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Intrest payments", loan_intrest, "", true, -intrest)
	}
	// Calculate bride loan intrest
	bl_intrest := 0.0
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
		if e.Group == bridge_loans {
			bl_intrest += e.Value * float64(external_factors.Bridge_loans_intrest_rate)
		}
	}

	if bl_intrest > 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Bridge loan intrest payments", bridge_loans, "", true, -intrest)
	}

	// Increase or decrease loans
	if decisions.Increase_of_loans > 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_liabilities("Bank loan", loans, "", true, float64(decisions.Increase_of_loans))
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Income from bank loan", loans, "", true, float64(decisions.Increase_of_loans))
	} else if decisions.Increase_of_loans < 0 {
		money_remaining := float64(-decisions.Increase_of_loans)
		var loans_to_delete []int
		for i, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
			if e.Group == 1 {
				if money_remaining >= e.Value {
					loans_to_delete = append(loans_to_delete, i)
					money_remaining -= e.Value
					company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Payement of loan", loans, "", true, -e.Value)
				} else if money_remaining < e.Value {
					company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities[i].Value -= money_remaining
					company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Payement of loan", loans, "", true, -money_remaining)
					money_remaining = 0
					break
				}
			}
		}

		seq := func(yield func(int) bool) {
			for i := range loans_to_delete {
				if !yield(i) {
					return
				}
			}
		}
		loans_to_delete = slices.SortedFunc(seq, func(a, b int) int { return cmp.Compare(b, a) })

		for _, i := range loans_to_delete {
			company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities = slices.Delete(company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities, i, i)
		}
	}

	// Totals
	financial_report.Total_income = 0
	financial_report.Cash_costs = 0
	financial_report.Non_cash_costs = 0
	financial_report.Total_expenses_before_tax = 0
	financial_report.EBIT = 0
	financial_report.Taxes = 0
	financial_report.Bridge_loan_repayments = 0
	financial_report.Loan_repayments = 0
	financial_report.Total_expenses_after_tax = 0
	financial_report.Profit = financial_report.Total_expenses_after_tax + financial_report.Total_income

	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement {
		financial_report.Profit += e.Value
		if e.Value > 0 {
			financial_report.Total_income += e.Value
		} else {
			financial_report.Total_expenses_after_tax += e.Value
			financial_report.Total_expenses_before_tax += e.Value
			if e.Cash_cost {
				financial_report.Cash_costs += e.Value
			} else {
				financial_report.Non_cash_costs += e.Value
			}
		}
		if e.Group != loan_intrest && e.Group != bridge_loan_intrest {
			financial_report.EBIT += e.Value
		}
		if e.Group == loans {
			financial_report.Loan_repayments += e.Value
		}
		if e.Group == bridge_loans {
			financial_report.Bridge_loan_repayments += e.Value
		}
		if e.Cash_cost {
			company.Balance += e.Value
		}
	}

	financial_report.Taxes = profit_taxes(financial_report.EBIT, external_factors)
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
		"Taxes",
		taxes,
		"Taxes paid on our profit",
		true,
		financial_report.Taxes)

	financial_report.Profit += financial_report.Taxes
	financial_report.Total_expenses_after_tax += financial_report.Taxes
	financial_report.Cash_costs += financial_report.Taxes
	company.Balance += financial_report.Taxes

	// Calculate bridge loans

	// try to pay off existing bridge loans
	if company.Balance > 0 {
		var loans_to_delete []int
		for i, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
			if e.Group != bridge_loans {
				continue
			}
			if company.Balance >= e.Value {
				company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Payment of bridge loan", bridge_loans, "", true, -e.Value)
				loans_to_delete = append(loans_to_delete, i)
				company.Balance -= e.Value
				financial_report.Cash_costs -= e.Value
				financial_report.Total_expenses_before_tax -= e.Value
				financial_report.Total_expenses_after_tax -= e.Value
				financial_report.Profit -= e.Value
				financial_report.Bridge_loan_repayments -= e.Value
			} else if company.Balance < e.Value {
				company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities[i].Value = e.Value - company.Balance
				company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement("Payement of bridge loan", bridge_loans, "", true, -company.Balance)
				financial_report.Cash_costs -= company.Balance
				financial_report.Total_expenses_before_tax -= company.Balance
				financial_report.Total_expenses_after_tax -= company.Balance
				financial_report.Profit -= company.Balance
				financial_report.Bridge_loan_repayments -= company.Balance
				company.Balance = 0
				break
			}
		}
		seq := func(yield func(int) bool) {
			for i := range loans_to_delete {
				if !yield(i) {
					return
				}
			}
		}
		loans_to_delete = slices.SortedFunc(seq, func(a, b int) int { return cmp.Compare(b, a) })

		for _, i := range loans_to_delete {
			company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities = slices.Delete(company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities, i, i)
		}
	} else if company.Balance < 0 {
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_income_statement(
			"Bridge loan",
			bridge_loans,
			"You are automatically lent out bridge loans when your balance goes beneath 0",
			true,
			-company.Balance)
		company.Reports[len(company.Reports)-1].Balance_sheet.add_to_liabilities(
			"Bridge loan",
			bridge_loans,
			"",
			true,
			-company.Balance)
		// financial_report.Total_income += -company.Balance
		// financial_report.EBIT += -company.Balance
		// financial_report.Profit += -company.Balance
		financial_report.New_bridge_loan += -company.Balance
		company.Balance = 0
	}

	// calculate Liabilities
	total_assets := 0.0
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Assets {
		total_assets += e.Value
	}
	total_liabilities := 0.0
	for _, e := range company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities {
		total_liabilities += e.Value
	}

	equity := total_assets - total_liabilities
	company.Reports[len(company.Reports)-1].Balance_sheet.add_to_liabilities(
		"Private equity",
		other,
		"The amount of money that is owned exclusively by the company",
		false,
		equity)

	return financial_report
}

func profit_taxes(EBIT float64, external_factors External_factors) float64 {
	if EBIT > 0 {
		return -round(EBIT*float64(external_factors.tax_rate), 2)
	}
	return 0
}

// Employee functions

func simulate_employees(
	employees []Employee,
	hires int,
	external_factors External_factors,
	training_investment float32,
	raises float32,
	working_hours float32,
	minimum_wage float32,
	base_pay float32,
	financial_report *Balance_sheet,
) ([]Employee, error) {
	// Calculate Turnover
	number_of_employees_leaving := (int(math.Round(float64(len(employees)) * float64(external_factors.Turnover))))
	if hires < 0 {
		number_of_employees_leaving += hires
	}

	// Detirmine who leaves (fired / voluntary)
	for range number_of_employees_leaving {
		if len(employees) <= 0 {
			break
		}
		delete_index := rand.Intn(len(employees))
		employees = slices.Delete(employees, delete_index, delete_index)
	}

	// Geneate new employees
	if hires > 0 {
		for range hires {
			employees = append(employees, generate_employee(base_pay, working_hours, employees[0].Employee_type, 1.0))
		}
	}

	// Calculate Payroll & working hours

	old_employees := make([]Employee, len(employees))
	copy(old_employees, employees)

	employees, pay_entries := calculate_payroll(employees, raises)
	financial_report.Income_statement = append(financial_report.Income_statement, pay_entries...)

	var err error
	employees, err = set_working_hours(employees, working_hours)
	if err != nil {
		return employees, err
	}

	// Calculate training
	employees = train_employees(employees, training_investment, 0.1)

	// Calculate Motivation
	employees = calculate_motivation(employees, old_employees, minimum_wage, 1.0)

	return employees, nil
}

func train_employees(employees []Employee, investment float32, passive_training float32) []Employee {
	investment_per_employee := investment / float32(len(employees))
	for i, e := range employees {
		employees[i].Skill += e.Skill * investment_per_employee / 1000.0
		employees[i].Skill += passive_training
	}
	return employees
}

func calculate_payroll(employees []Employee, pay_raise float32) ([]Employee, []FinanceReportEntry) {
	employees_pay := make([]FinanceReportEntry, len(employees))
	for i := range employees {
		employees[i].Salary += employees[i].Salary * pay_raise
		employees_pay[i] = FinanceReportEntry{
			fmt.Sprintf("Pay for %s employee %d", employees[i].get_type(), employees[i].Id),
			personelle,
			"",
			true,
			float64(-employees[i].Salary),
		}
	}
	return employees, employees_pay
}

func set_working_hours(employees []Employee, working_hours float32) ([]Employee, error) {
	if working_hours <= 0 {
		return employees, errors.New("working hours must be > 0h")
	}

	for i := range employees {
		employees[i].Working_hours = working_hours
	}

	return employees, nil
}

func calculate_motivation(employees_after_calculations []Employee, employees_before_calculations []Employee, minimum_wage float32, base_motivation float32) []Employee {
	// REDO MOTIVATION
	for i, e := range employees_after_calculations {
		pay_factor := e.Salary / (minimum_wage * 1.2)
		raise_factor := e.Salary / employees_before_calculations[i].Salary
		working_hours_factor := e.Working_hours / 8
		time_off_factor := e.Working_hours / employees_before_calculations[i].Working_hours
		training_factor := e.Skill / employees_before_calculations[i].Skill

		employees_after_calculations[i].Motivation = ((base_motivation*2 +
			pay_factor +
			raise_factor*2 +
			working_hours_factor +
			time_off_factor*2 +
			training_factor) /
			9)
	}
	return employees_after_calculations
}

// Economy functions

func (population *Population) simulate_economy(companies *[]Company, external_factors External_factors) ([]FinanceReportEntry, []Purchasing_statistics, error) {
	// Get offers
	offers := make([]Offer, len(*companies))
	product_availability := make([]int, len(*companies))
	for i, c := range *companies {
		offers[i] = c.Offer
		product_availability[i] = c.Items_in_storage
	}

	if external_factors.economic_situation_index <= 0 {
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
		product_purchased, quanity, customer, individual_purchasing_statistics := calculate_purchase(c, offers, avg_price, external_factors, product_availability)
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

func calculate_purchase(customer Customer, offers []Offer, avg_price float32, external_factors External_factors, product_availability []int) (int, int, Customer, []Purchasing_statistics) { // returns (index of product purchased, quality purchased, customer in question)
	decision_factors := make([]float32, len(offers))
	purchasing_statistics := make([]Purchasing_statistics, len(offers))
	for i, o := range offers {
		decision_factors[i] = (customer.Quality_preference*o.Product.Quality_factor +
			customer.Ecology_preference*o.Product.Ecology_factor +
			customer.Coolness_preference*o.Product.Coolness_factor +
			customer.Price_preference*is_cheap(o, avg_price) +
			customer.Bang_for_buck_preference*(o.Product.Quality_factor/o.Price) +
			customer.Brand_loyalty_factor*customer.Loyalties[i]) * external_factors.economic_situation_index

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

func (company *Company) compile_reports(decisions Decisions, Results FinanceReportEntry, external_factors External_factors) error {
	products_sold, err := strconv.Atoi(strings.Split(Results.Info, " ")[0])
	if err != nil {
		return err
	}
	company.Reports[len(company.Reports)-1].Sales_report.Products_sold = products_sold
	company.Items_in_storage -= company.Reports[len(company.Reports)-1].Sales_report.Products_sold
	company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement = append(company.Reports[len(company.Reports)-1].Balance_sheet.Income_statement, Results)

	// Finance
	company.Reports[len(company.Reports)-1].Financial_Report = company.calculate_budget(decisions, external_factors)

	return nil
}

func round(num float64, decimal_place int) float64 {
	num = num * math.Pow(10, (float64(decimal_place)))
	num = math.Round(num)
	num = num / math.Pow(10, (float64(decimal_place)))
	return num
}

// ##################################################
// ##########                    _         ##########
// ##########                   (_)        ##########
// ##########  _ __ ___    __ _  _  _ __   ##########
// ########## | '_ ` _ \  / _` || || '_ \  ##########
// ########## | | | | | || (_| || || | | | ##########
// ########## |_| |_| |_| \__,_||_||_| |_| ##########
// ##################################################

func new_game(sim_config Sim_config, number_of_companies int, game_name string) Game_state {
	var game_state Game_state

	game_state.Step = 0
	game_state.Step_simulated = false
	game_state.Game_name = game_name

	game_state.External_factors = External_factors{
		Inflation:                 0.005,
		economic_situation_index:  1,
		tax_rate:                  0.147,
		Material_price:            10,
		Energy_price:              96.2 / 10,
		Machine_depreciation_rate: 0.1,

		Turnover:                0.08,
		Production_minimum_wage: 60000,
		Marketing_minimum_wage:  80000,

		Machine_on_offer: Machine{
			Production_capacity: 15000,
			Required_workers:    5,
			Energy_use:          75000,
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

	return game_state
}

func get_decisions(save_location string, number_of_companies int) ([]Decisions, error) {
	decisions := make([]Decisions, number_of_companies)
	for i := range decisions {
		decisions_json, err := os.ReadFile(fmt.Sprintf("%s/decisions_company_%d.json", save_location, i))

		println(string(decisions_json))
		if err != nil {
			return decisions, err
		}

		err = json.Unmarshal(decisions_json, &decisions[i])
		if err != nil {
			return nil, err
		}
	}

	return decisions, nil
}

func (game_state Game_state) save_game(location string) error {
	var save Save_game

	var population_buffer bytes.Buffer
	encoder := gob.NewEncoder(&population_buffer)
	err := encoder.Encode(game_state.Population.Population)
	if err != nil {
		return err
	}
	save.Population = population_buffer.Bytes()
	game_state.Population = Population{}

	save.Game_state = game_state

	save_file, err := json.MarshalIndent(save, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf(
		"%s/%s-%d.json",
		location,
		game_state.Game_name,
		game_state.Step,
	), save_file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Load settings
	sim_config_file, err := os.ReadFile("Config/sim_config.json")
	if err != nil {
		println("Error loading sim_config.json")
		os.Exit(1)
	}

	var sim_config Sim_config
	err = json.Unmarshal(sim_config_file, &sim_config)
	if err != nil {
		println("Error in sim_config.json")
		os.Exit(2)
	}

	// Setup game
	game_state := new_game(sim_config, 4, "Test_game")

	game_state.Current_decisions, err = get_decisions(
		fmt.Sprintf("Saves/%s-%d/Decisions", game_state.Game_name, game_state.Step),
		len(game_state.Companies),
	)
	if err != nil {
		println(err.Error())
		log.Fatal("Failed to get current decisions")
	}

	// fmt.Printf("%+#v\n", game_state.Current_decisions[0])

	for range 1 {
		err = game_state.simulate_step()
		if err != nil {
			println(err.Error())
			log.Fatal("An error has occured")
		}
	}
	//for _, c := range game_state.Companies {
	//	company_as_json, err := json.MarshalIndent(c, "", "  ")
	//	if err != nil {
	//		log.Fatal("Error occored printing companies")
	//	}
	//
	//	 println(string(company_as_json))
	//}

	println("Saving file")
	err = game_state.save_game("Saves")
	if err != nil {
		log.Fatal(err.Error())
	}
}
