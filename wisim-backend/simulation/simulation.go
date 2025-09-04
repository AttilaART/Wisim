package simulation

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Save_game struct {
	Population []byte // Binary
	Game_state Game_state
}

type Game_state struct {
	Step      int
	Game_name string

	Population              Population
	Employees               Employee_pool
	Companies               []Company
	Current_decisions       []Decisions
	Decisions_submitted     []bool
	Market_sales_statistics []Sales_statistics
	External_factors        External_factors
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
	Ethics_bias                 float32
	Ethics_spread               float32
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
	Base_market_price           float32
	Market_saturation           float32
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
	Id            int
	Name          string
	Balance       float64
	Loans         float64
	Bridge_loans  float64
	employee_pool Employee_pool
	// Global_effects   []Effect
	Decision_history []Decisions

	Reports []Report
	// Research and development
	Global_quality_factor   float32
	Base_marketing_strength float32

	// Product
	Offer  Offer
	Orders int

	// Fulfillment
	Warehouses       []Warehouse
	Items_in_storage int

	// Production
	Machines []Machine
}

// NOTE: Employees themselves keep track of their employers

type Decisions struct {
	Predictions struct {
		Sales_prediction int
	}

	Finances struct {
		Set_bank_loan float64
	}

	Marketing Decisions_marketing

	Employees struct {
		Production_deltas []Delta[Employee]
		Marketing_deltas  []Delta[Employee]

		Severance_pay float32
	}

	Production struct {
		Production_goal int
		Machines        []Delta[Machine]
		Logistics       []Delta[Warehouse]
	}

	Research Decisions_research
}

type Decisions_marketing struct {
	Price float32

	Product Decisions_product

	Promotion struct {
		Quantity         float64
		Style_quality    float32
		Style_ecology    float32
		Style_ethics     float32
		Style_durability float32
	}
}

type Decisions_product struct {
	Materials struct {
		Quality          float32
		Ecology          float32
		Ethical_sourcing float32
	}

	Manufacturing struct {
		Quality             float32
		Ecological_energy   float32
		Material_efficiency float32
		Durability          float32
		Max_durability      int
	}
}

type Decisions_research struct {
	Quality         float32
	Durability      float32
	Ecology         float32
	Promotion       float32
	Production_cost float32
}

type Delta[V any] struct {
	Change int
	Item   V
}

const (
	Delta_New = iota
	Delta_Change
	Delta_Remove
)

type Offer struct {
	Product           Product
	Price             float32
	Promotion_quality float32
	Promotion_goal    struct {
		Quantity         float64
		Style_quality    float32
		Style_ecology    float32
		Style_ethics     float32
		Style_durability float32
	}
}

type Product struct {
	Id              int
	Name            string
	Weight          float32
	Material_use    float32
	Production_cost float32

	Base_material_use    float32
	Base_production_cost float32
	Base_quality         float32
	Base_ecology         float32
	Base_durability      float32

	Ethics_factor  float32 // TODO: Implement Ethicss Factor
	Quality_factor float32
	Ecology_factor float32
	// Coolness_factor float32

	Durabilty int
}

type Machine struct {
	Id int

	Production_capacity  int
	Required_workers     int
	Minimum_workers      int
	Assigned_workers_ids []int
	Energy_use           float32
	Value                float32
	Maintanance_cost     float32 // Monthly
}

func (m Machine) get_id() int {
	return m.Id
}

type Warehouse struct {
	Id              int
	Capacity        int
	Operating_costs float32
	Value           float32
}

func (w Warehouse) get_id() int {
	return w.Id
}

type Employee struct {
	Id       int
	Name     string
	Employer int

	Employee_type Employee_type

	Motivation     float32
	Skill          float32
	Extra_training float32

	// Global_effect *Effect

	Pay           float32
	Bonus         float32
	Working_hours float32
}

type Employee_type int

const (
	Employee_type_production = iota
	Employee_type_marketing
	Employee_type_executive
	Employee_type_all = -1
)

const (
	Employee_employer_none = -1
)

func (e_type Employee_type) to_string() string {
	employee_types := []string{"production", "marketing", "executive"}

	if int(e_type) > len(employee_types)-1 {
		return "unknown"
	} else if e_type < 0 {
		return "unknown"
	}
	return employee_types[e_type]
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

	Financial_report  Financial_report
	Balance_sheet     Balance_sheet
	Personelle_report Personelle_report
	Production_report Production_report
	Sales_report      Sales_report
}

type Financial_report struct {
	// Income
	Income struct {
		Gross_sales   float64
		Other_income  float64
		Cost_of_sales float64
		Gross_profit  float64
	}
	Operating_expenses struct {
		Advertising              float64
		Facilities_and_logistics float64
		Research_and_development float64
	}
	Non_operating_expenses struct {
		Write_offs            float64
		Loan_interest         float64
		Loan_repayment        float64
		Bridge_loan_intrest   float64
		Bridge_loan_repayment float64
		Other                 float64
		Taxes                 float64
	}
	Totals struct {
		Total_operating_expenses     float64
		Total_non_operating_expenses float64
		Income_before_tax            float64
		Net_income                   float64
		Cashflow                     float64
	}
}

func (f *Balance_sheet) add_to_income_statement(name string, group int, info string, cash_cost bool, value float64) *FinanceReportEntry {
	entry := FinanceReportEntry{name, group, info, cash_cost, value}
	f.Invoice_log = append(f.Invoice_log, entry)
	return &f.Invoice_log[len(f.Invoice_log)-1]
}

func (f *Balance_sheet) add_to_equity(name string, group int, info string, cash_cost bool, value float64) *FinanceReportEntry {
	entry := FinanceReportEntry{name, group, info, cash_cost, value}
	f.Assets = append(f.Assets, entry)
	return &f.Assets[len(f.Assets)-1]
}

func (f *Balance_sheet) add_to_liabilities(name string, group int, info string, cash_cost bool, value float64) *FinanceReportEntry {
	entry := FinanceReportEntry{name, group, info, cash_cost, value}
	f.Liabilities = append(f.Liabilities, entry)
	return &f.Liabilities[len(f.Liabilities)-1]
}

type Balance_sheet struct {
	Bank_balance float64

	Invoice_log []FinanceReportEntry

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
	production_personelle
	marketing_personelle
	other_personelle
	facilities
	logistics
	materials
	energy
	research
	employee_training
	loans
	loan_intrest
	bridge_loans
	bridge_loan_intrest
	taxes
	sales
	severance
	predictions
	write_off
	other
)

var AllGroups = []struct {
	Value  int
	TSName string
}{
	{production, "production"},
	{marketing, "marketing"},
	{production_personelle, "prodcution_personelle"},
	{marketing_personelle, "marketing_personelle"},
	{other_personelle, "other_personelle"},
	{facilities, "facilities"},
	{logistics, "logistics"},
	{materials, "materials"},
	{energy, "energy"},
	{research, "product_development"},
	{employee_training, "employee_training"},
	{loans, "loans"},
	{loan_intrest, "loan_intrest"},
	{bridge_loans, "bridge_loans"},
	{bridge_loan_intrest, "bridge_loan_intrest"},
	{taxes, "taxes"},
	{sales, "sales"},
	{severance, "severance"},
	{predictions, "predictions"},
	{write_off, "write_off"},
	{other, "other"},
}

type Personelle_report struct {
	General    Personelle_sub_report
	Marketing  Personelle_sub_report
	Production Personelle_sub_report
}

type Personelle_sub_report struct {
	Number_of_employees  int
	Number_of_hires      int
	Number_of_departures int

	Avg_pay          float32
	Minimum_pay      float32
	Maximum_pay      float32
	Standard_dev_pay float32

	Minimum_skill      float32
	Maximum_skill      float32
	Avg_skill          float32
	Standard_dev_skill float32

	Minimum_motivation      float32
	Maximum_motivation      float32
	Avg_motivation          float32
	Standard_dev_motivation float32

	Minimum_productivity      float32
	Maximum_productivity      float32
	Avg_productivity          float32
	Standard_dev_productivity float32
}

type Production_report struct {
	Machines_purchased       int
	Machines_sold            int
	Worker_surplus           int
	Avg_machine_productivity float32
	// Max_machine_productivity float32
	// Min_machine_productivity float32

	Total_production  int
	Base_production   int
	Bonus_production  int
	Excess_production int

	Total_products_produced int
	Base_products_produced  int
	Bonus_products_produced int

	Material_used float32
	Energy_used   float32

	Warehouses_bought int
}

type Sales_report struct {
	Product_statistics       Product_statistics
	Company_sales_statistics Sales_statistics
	Marketing_statistics     Marketing_statistics
}

type Purchasing_statistics struct {
	Product_number               int
	Products_sold                int
	Difference_to_previous_month int
	Product_demand               int

	Avr_decision_factor      float32
	Avr_purchasing_threshold float32

	Avr_quality_factor    float32
	Avr_durability_factor float32
	Avr_ecology_factor    float32
	Avr_price_factor      float32
	Avr_ethics_factor     float32
	// Avr_coolness_factor      float32
	Avr_bang_for_buck_factor float32
}

type Research_statistics struct {
	Quality_development_investment                  float64
	Quality_development_investment_effectiveness    float64
	Durability_development_investment               float64
	Durability_development_investment_effectiveness float64
	Ecological_production_investment                float32 // Decreases material use
	Ecological_production_investment_effectiveness  float32 // Decreases material use
}

type Sales_statistics struct {
	Products_sold                int
	Difference_to_previous_month int
	Product_demand               int
	Market_share                 float32

	Avr_decision_factor      float32
	Avr_purchasing_threshold float32

	Avr_quality_factor    float32
	Avr_durability_factor float32
	Avr_ecology_factor    float32
	Avr_price_factor      float32
	Avr_ethics_factor     float32
	/// Avr_coolness_factor      float32
	Avr_bang_for_buck_factor float32
}

type Marketing_statistics struct {
	Product            Product_statistics
	Price              float64
	Bang_for_buck      float64
	Promotion_quantity float64
	Promotion_quality  float64
	// Place
}

type Product_statistics struct {
	Quality   float32
	Durabilty int
	// Coolness  float32
	Ethics  float32
	Ecology float32
}

type Customer struct {
	Base_need      int
	Owned_products []Owned_product

	Quality_preference float32
	Ecology_preference float32
	// Coolness_preference      float32
	Ethics_preference        float32
	Price_preference         float32
	Bang_for_buck_preference float32 // Price / Quality preference
	Durabilty_preference     float32

	Purchashing_threshold float32
	Max_price             float32
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
	Economic_situation_index  float32
	Tax_rate                  float32 // as decimal

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

// #####################################################################################################
// ##########     _         _              __                      _    _                     ##########
// ##########    | |       | |            / _|                    | |  (_)                    ##########
// ##########  __| |  __ _ | |_  __ _    | |_  _   _  _ __    ___ | |_  _   ___   _ __   ___  ##########
// ########## / _` | / _` || __|/ _` |   |  _|| | | || '_ \  / __|| __|| | / _ \ | '_ \ / __| ##########
// ##########| (_| || (_| || |_| (_| |   | |  | |_| || | | || (__ | |_ | || (_) || | | |\__ \ ##########
// ########## \__,_| \__,_| \__|\__,_|   |_|   \__,_||_| |_| \___| \__||_| \___/ |_| |_||___/ ##########
// #####################################################################################################

func (c Company) Mock_simulate_step(decisions Decisions, external_factors External_factors) Report {
	results := FinanceReportEntry{"Predicted sales", predictions, "The amount of you predict you'll make", true, float64(decisions.Predictions.Sales_prediction) * float64(decisions.Marketing.Price)}
	c.compile_reports(
		decisions,
		results,
		Purchasing_statistics{Products_sold: decisions.Predictions.Sales_prediction},
		&Sales_statistics{},
		external_factors,
	)
	return c.Reports[len(c.Reports)-1]
}

// Gameloop functions
func (game_state *Game_state) Simulate_step() error {
	game_state.Step += 1

	game_state.External_factors.Month = game_state.Step

	if len(game_state.Companies) != len(game_state.Current_decisions) {
		return errors.New("amount of decisions does not match number of companies")
	}
	println("============ Simulating Hiring / Firing ===========")
	game_state.handleEmployeeDeltas()

	println("=============== Simulating companies ==============")
	for i := range game_state.Companies {
		game_state.Companies[i].employee_pool = game_state.Employees
		fmt.Printf("--------------- Simulating company %d -------------- \n", i)
		err := game_state.Companies[i].simulate_company(game_state.Current_decisions[i], game_state.External_factors)
		if err != nil {
			return err
		}
	}

	println("Simulating companies done!")

	println("---------------- Simulatig economy ----------------")
	Results, purchasing_statistics, err := game_state.Population.simulate_economy(
		&game_state.Companies,
		game_state.External_factors,
	)
	if err != nil {
		return err
	}
	println("Simulatig economy done!")

	println("================ Compiling reports =============== ")

	game_state.Market_sales_statistics = append(game_state.Market_sales_statistics, Sales_statistics{})
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Products_sold = purchasing_statistics[len(purchasing_statistics)-1].Products_sold
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Product_demand = purchasing_statistics[len(purchasing_statistics)-1].Product_demand
	if len(game_state.Market_sales_statistics) >= 2 {
		game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Difference_to_previous_month = purchasing_statistics[len(purchasing_statistics)-1].Products_sold - game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-2].Products_sold
	} else {
		game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Difference_to_previous_month = purchasing_statistics[len(purchasing_statistics)-1].Products_sold
	}
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Market_share = 100

	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_decision_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_decision_factor
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_purchasing_threshold = purchasing_statistics[len(purchasing_statistics)-1].Avr_purchasing_threshold

	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_quality_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_quality_factor
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_durability_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_durability_factor
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_ecology_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_ecology_factor
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_price_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_price_factor
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_ethics_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_ethics_factor
	// game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_coolness_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_coolness_factor
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_bang_for_buck_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_bang_for_buck_factor

	//for range game_state.Market_sales_statistics {
	//	println("----------")
	//	fmt.Printf("%f \n", game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_decision_factor)
	//	fmt.Printf("%f \n", game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_purchasing_threshold)
	//	fmt.Printf("%f \n", game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_quality_factor)
	//	fmt.Printf("%f \n", game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_durability_factor)
	//	fmt.Printf("%f \n", game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_ecology_factor)
	//	fmt.Printf("%f \n", game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_price_factor)
	//	fmt.Printf("%f \n", game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_ethics_factor)
	//	fmt.Printf("%f \n", game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_bang_for_buck_factor)
	//	println("----------")
	//}

	for i := range game_state.Companies {
		fmt.Printf("Compiling reports for company %d\n", i)
		game_state.Companies[i].compile_reports(
			game_state.Current_decisions[i],
			Results[i],
			purchasing_statistics[i],
			&game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1],
			game_state.External_factors,
		)

	}

	for i := range game_state.Decisions_submitted {
		game_state.Decisions_submitted[i] = false
	}
	println("=================================================== ")
	println("               Simulation step done!\n")
	println("===================== RESULTS ===================== ")
	println("Month: ", game_state.Step)

	printer := message.NewPrinter(language.Swedish)
	for i, c := range game_state.Companies {
		printer.Printf("Company %d: %s:\n", i, c.Name)
		printer.Printf("Total Production: %d\n", c.Reports[len(c.Reports)-1].Production_report.Total_production)
		printer.Printf("Excess Production: %d\n", c.Reports[len(c.Reports)-1].Production_report.Excess_production)
		printer.Printf("Total Products Produced: %d\n", c.Reports[len(c.Reports)-1].Production_report.Total_products_produced)
		printer.Printf("Products sold: %d\n", c.Reports[len(c.Reports)-1].Sales_report.Company_sales_statistics.Products_sold)
		printer.Printf("--> Net profit: %.2f", c.Reports[len(c.Reports)-1].Financial_report.Totals.Net_income)
		println("")
		printer.Printf("Number of employees: %d\n", c.Reports[len(c.Reports)-1].Personelle_report.General.Number_of_employees)
		printer.Printf("Number of production employees: %d\n", c.Reports[len(c.Reports)-1].Personelle_report.Production.Number_of_employees)
		printer.Printf("Number of marketing employees: %d\n", c.Reports[len(c.Reports)-1].Personelle_report.Marketing.Number_of_employees)
		printer.Printf("Number of hires: %d\n", c.Reports[len(c.Reports)-1].Personelle_report.General.Number_of_hires)
		printer.Printf("Number of departures: %d\n", c.Reports[len(c.Reports)-1].Personelle_report.General.Number_of_departures)
		printer.Printf("Avr Pay: %.2f\n", c.Reports[len(c.Reports)-1].Personelle_report.General.Avg_pay)
		printer.Printf("Avr Skill: %.2f\n", c.Reports[len(c.Reports)-1].Personelle_report.General.Avg_skill)
		printer.Printf("Avr Motivation: %.2f\n", c.Reports[len(c.Reports)-1].Personelle_report.General.Avg_motivation)
		printer.Printf("Avr Productivity: %.2f\n", c.Reports[len(c.Reports)-1].Personelle_report.General.Avg_productivity)
		println("")
	}

	total_products_sold := 0
	for _, c := range game_state.Companies {
		total_products_sold += c.Reports[len(c.Reports)-1].Sales_report.Company_sales_statistics.Products_sold
	}

	println("Total products sold: ", total_products_sold)

	missing_products := 0
	for _, p := range game_state.Population.Population {
		missing_products += p.Base_need - len(p.Owned_products)
	}
	avr_missing_products := float64(missing_products) / float64(len(game_state.Population.Population))
	fmt.Printf("avr missing products: %f\n", avr_missing_products)

	// println("============== Purchasing statistics ============== ")
	//
	//	for _, p := range purchasing_statistcs {
	//		s, e := json.MarshalIndent(p, "", "    ")
	//		if e != nil {
	//			return e
	//		}
	//		println(string(s))
	//	}
	return nil
}

func (company *Company) compile_reports(
	decisions Decisions,
	results FinanceReportEntry,
	company_purchasing_statistcs Purchasing_statistics,
	market_purchasing_statistics *Sales_statistics,
	external_factors External_factors,
) error {
	products_sold := company_purchasing_statistcs.Products_sold

	company.Reports[len(company.Reports)-1].Sales_report = company.compile_sales_report(
		company_purchasing_statistcs,
		market_purchasing_statistics.Products_sold,
	)

	company.Items_in_storage -= products_sold
	company.Reports[len(company.Reports)-1].Balance_sheet.Invoice_log = append(company.Reports[len(company.Reports)-1].Balance_sheet.Invoice_log, results)

	// Finance
	company.calculate_budget(decisions, external_factors)
	company.Reports[len(company.Reports)-1].Balance_sheet.Assets = clean_up_financeReportEntries(company.Reports[len(company.Reports)-1].Balance_sheet.Assets)
	company.Reports[len(company.Reports)-1].Balance_sheet.Invoice_log = clean_up_financeReportEntries(company.Reports[len(company.Reports)-1].Balance_sheet.Invoice_log)
	company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities = clean_up_financeReportEntries(company.Reports[len(company.Reports)-1].Balance_sheet.Liabilities)

	// Personelle
	company.Reports[len(company.Reports)-1].Personelle_report = company.compile_personelle_report(decisions)

	return nil
}

func (c *Company) compile_personelle_report(decisions Decisions) Personelle_report {
	personelle_report := Personelle_report{}

	personelle_report.General = c.compile_personelle_subreport(decisions, Employee_type_all)
	personelle_report.Marketing = c.compile_personelle_subreport(decisions, Employee_type_marketing)
	personelle_report.Production = c.compile_personelle_subreport(decisions, Employee_type_production)

	return personelle_report
}

func (c *Company) compile_personelle_subreport(decisions Decisions, employee_type Employee_type) Personelle_sub_report {
	var sub_report Personelle_sub_report
	employee_ids := c.employee_pool.Get_employees_of_company(c.Id, employee_type)

	sub_report.Number_of_employees = len(employee_ids)

	var employee_deltas []Delta[Employee] // We can trust that the employees exist because we checked this when "simulating employees"
	if employee_type == Employee_type_marketing {
		employee_deltas = decisions.Employees.Marketing_deltas
	} else if employee_type == Employee_type_production {
		employee_deltas = decisions.Employees.Production_deltas
	} else if employee_type == Employee_type_all {
		employee_deltas = decisions.Employees.Production_deltas
		employee_deltas = append(employee_deltas, decisions.Employees.Marketing_deltas...)
	}

	for _, e_delta := range employee_deltas {
		if e_delta.Change == Delta_New {
			sub_report.Number_of_hires += 1
		} else if e_delta.Change == Delta_Remove {
			sub_report.Number_of_departures += 1
		}
	}

	pay := make([]float32, len(employee_ids))
	skill := make([]float32, len(employee_ids))
	motivation := make([]float32, len(employee_ids))
	productivity := make([]float32, len(employee_ids))
	for i, e := range employee_ids {
		pay[i] = c.employee_pool[e].Pay
		skill[i] = c.employee_pool[e].Skill
		motivation[i] = c.employee_pool[e].Motivation
		productivity[i] = c.employee_pool[e].Motivation * c.employee_pool[e].Skill * c.employee_pool[e].Working_hours // TODO: Make sure this is actually accurate
	}

	sub_report.Avg_pay = avr(pay)
	sub_report.Maximum_pay = max(pay...)
	sub_report.Minimum_pay = min(pay...)
	sub_report.Standard_dev_pay = std_dev(pay...)

	sub_report.Avg_skill = avr(skill)
	sub_report.Maximum_skill = max(skill...)
	sub_report.Minimum_skill = min(skill...)
	sub_report.Standard_dev_skill = std_dev(skill...)

	sub_report.Avg_motivation = avr(motivation)
	sub_report.Maximum_motivation = max(motivation...)
	sub_report.Minimum_motivation = min(motivation...)
	sub_report.Standard_dev_motivation = std_dev(motivation...)

	sub_report.Avg_productivity = avr(productivity)
	sub_report.Maximum_productivity = max(productivity...)
	sub_report.Minimum_productivity = min(productivity...)
	sub_report.Standard_dev_productivity = std_dev(productivity...)

	return sub_report
}

func (c *Company) compile_sales_report(purchasing_statiscs Purchasing_statistics, Market_products_sold int) Sales_report {
	report := Sales_report{}

	report.Product_statistics.Quality = c.Offer.Product.Quality_factor
	report.Product_statistics.Durabilty = c.Offer.Product.Durabilty
	report.Product_statistics.Ecology = c.Offer.Product.Ecology_factor
	report.Product_statistics.Ethics = c.Offer.Product.Ethics_factor
	// report.Product_statistics.Coolness = c.Offer.Product.Coolness_factor

	// ----------------

	report.Company_sales_statistics.Products_sold = purchasing_statiscs.Products_sold
	if len(c.Reports) >= 2 {
		report.Company_sales_statistics.Difference_to_previous_month = purchasing_statiscs.Products_sold - c.Reports[len(c.Reports)-2].Sales_report.Company_sales_statistics.Products_sold
	} else {
		report.Company_sales_statistics.Difference_to_previous_month = purchasing_statiscs.Products_sold
	}
	report.Company_sales_statistics.Product_demand = purchasing_statiscs.Product_demand
	if Market_products_sold != 0 {
		report.Company_sales_statistics.Market_share = (float32(purchasing_statiscs.Products_sold) / float32(Market_products_sold))
	} else {
		report.Company_sales_statistics.Market_share = 0
	}

	report.Company_sales_statistics.Avr_decision_factor = purchasing_statiscs.Avr_decision_factor
	report.Company_sales_statistics.Avr_purchasing_threshold = purchasing_statiscs.Avr_purchasing_threshold

	report.Company_sales_statistics.Avr_quality_factor = purchasing_statiscs.Avr_quality_factor
	report.Company_sales_statistics.Avr_durability_factor = purchasing_statiscs.Avr_durability_factor
	report.Company_sales_statistics.Avr_ecology_factor = purchasing_statiscs.Avr_ecology_factor
	report.Company_sales_statistics.Avr_price_factor = purchasing_statiscs.Avr_price_factor
	report.Company_sales_statistics.Avr_ethics_factor = purchasing_statiscs.Avr_ethics_factor
	// report.Company_sales_statistics.Avr_coolness_factor = purchasing_statiscs.Avr_coolness_factor
	report.Company_sales_statistics.Avr_bang_for_buck_factor = purchasing_statiscs.Avr_bang_for_buck_factor

	// ----------------

	report.Marketing_statistics.Product.Quality = c.Offer.Product.Quality_factor
	report.Marketing_statistics.Product.Durabilty = c.Offer.Product.Durabilty
	report.Marketing_statistics.Product.Ethics = c.Offer.Product.Ethics_factor
	// report.Marketing_statistics.Product.Coolness = c.Offer.Product.Coolness_factor
	report.Marketing_statistics.Product.Ecology = c.Offer.Product.Ecology_factor

	report.Marketing_statistics.Price = float64(c.Offer.Price)
	report.Marketing_statistics.Promotion_quantity = c.Offer.Promotion_goal.Quantity
	report.Marketing_statistics.Promotion_quality = float64(c.Offer.Promotion_quality)

	return report
}

// ##################################################
// ##########                    _         ##########
// ##########                   (_)        ##########
// ##########  _ __ ___    __ _  _  _ __   ##########
// ########## | '_ ` _ \  / _` || || '_ \  ##########
// ########## | | | | | || (_| || || | | | ##########
// ########## |_| |_| |_| \__,_||_||_| |_| ##########
// ##################################################

func Get_decisions_from_file(save_location string, number_of_companies int) ([]Decisions, error) {
	decisions := make([]Decisions, number_of_companies)
	for i := range decisions {
		decisions_json, err := os.ReadFile(fmt.Sprintf("%s/decisions_company_%d.json", save_location, i))
		// println(string(decisions_json))
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

func (game_state Game_state) Save_game(location string, compress bool) error {
	filename := fmt.Sprintf(
		"%s-%d.json",
		game_state.Game_name,
		game_state.Step,
	)

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

	// Turning file into zip
	// (IDK what's happening)

	if compress {

		zip_file_buffer := new(bytes.Buffer)

		w := zip.NewWriter(zip_file_buffer)

		file, err := w.Create(filename)
		if err != nil {
			return err
		}

		_, err = file.Write(save_file)
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		err = os.WriteFile(fmt.Sprint(location, "/", filename, ".zip"), zip_file_buffer.Bytes(), 0644)
		if err != nil {
			return err
		}

		return nil
	} // else if !compress {

	err = os.WriteFile(fmt.Sprint(location, "/", filename), save_file, 0644)
	if err != nil {
		return err
	}

	return nil
}
