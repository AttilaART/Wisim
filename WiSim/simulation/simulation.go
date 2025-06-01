package simulation

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Save_game struct {
	Population []byte // Binary
	Game_state Game_state
}

type Game_state struct {
	Step           int
	Step_simulated bool
	Game_name      string

	Population              Population
	Companies               []Company
	Current_decisions       []Decisionsold
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
	Id      int
	Name    string
	Balance float64

	// Global_effects   []Effect
	Decision_history []Decisionsold

	Reports []Report
	// Research and development
	Global_quality_factor float32

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
}

type Decisionsold struct {
	// Planing for budget reports
	Sales_projection int
	Selling_price    float32

	// Quality_of_stores int
	Marketing float32 // Increases probabilty that a potential customer hears of your product

	Quality_development_investment    float32 // Increases Quality
	Ecological_production_investment  float32 // Decreases material use
	Durability_development_investment float32 // increases durabilty
	// Investment_in_coolness_development  float32 // Increases coolness

	Production_target    int
	Purchase_of_machines int
	Selling_of_machines  int

	Material_quality                float32
	Percentage_of_ecological_energy float32

	Purchase_of_warehouses int

	New_hires_in_production             int
	New_hires_in_marketing              int
	Base_pay_for_production             float32
	Base_pay_for_marketing              float32
	Raise_for_production_personelle     float32
	Raise_for_marketing_personelle      float32
	Severance_for_production_personelle float32
	Severance_for_marketing_personelle  float32

	Working_hours_for_production float32
	Working_hours_for_marketing  float32

	Investment_in_production_training float32 // Skilled production staff increase quality & give bonus production
	Investment_in_marketing_training  float32 // Skilled marketing staff increase quality of marketing

	Increase_of_loans float32
	Dividends         float32
}

type Decisions struct {
	Predictions struct {
		Sales_prediction int
	}

	Finances struct {
		set_bank_loan float64
	}

	Marketing struct {
		Price float32

		Product struct {
			Materials struct {
				Quality          float32
				Ecology          float32
				Ethical_sourcing float32
			}

			Manufacturing struct {
				Quality             float32
				Durability          float32
				Ecological_energy   float32
				Material_efficiency float32
			}
		}

		Promotion struct {
			Quantity         float64
			Style_quality    float32
			Style_ecology    float32
			Style_ethics     float32
			Style_durability float32
		}
	}

	Employees struct {
		Production_actions []Employee_action
		Marketing_actions  []Employee_action
	}

	Production struct {
		Production_goal int
		Machines        []Machine
		Logistics       []Warehouse
	}
}

const (
	existing_employee = iota
	new_hire
	fired
)

type Employee_action struct {
	Employee Employee

	Extra_training int
	Pay            float32
	Bonus          float32

	status int
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
	Production_capacity  int
	Required_workers     int
	Minimum_workers      int
	Assigned_workers_ids []int
	Energy_use           float32
	Value                float32
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
	Operating_profit          float64
	Taxes                     float64

	Net_Profit float64
}

type Financial_Reportnew struct {
	// Income
	Income struct {
		Gross_sales   float64
		Cost_of_sales float64
		Gross_profit  float64
	}
	Operating_expenses struct {
		Advertising              float64
		Facilities_and_logistics float64
		Research_and_development float64
		Total_operating_expenses float64
	}
	Non_operating_expenses struct {
		Write_offs                   float64
		Loan_interest                float64
		Loan_repayment               float64
		bridge_loan_intrest          float64
		bridge_loan_repayment        float64
		Total_non_operating_expenses float64
		Taxes                        float64
		Net_income                   float64
		cashflow                     float64
	}
}

func (f *Balance_sheet) add_to_income_statement(name string, group int, info string, cash_cost bool, value float64) {
	f.Invoice_log = append(f.Invoice_log, FinanceReportEntry{name, group, info, cash_cost, value})
}

func (f *Balance_sheet) add_to_equity(name string, group int, info string, cash_cost bool, value float64) {
	f.Assets = append(f.Assets, FinanceReportEntry{name, group, info, cash_cost, value})
}

func (f *Balance_sheet) add_to_liabilities(name string, group int, info string, cash_cost bool, value float64) {
	f.Liabilities = append(f.Liabilities, FinanceReportEntry{name, group, info, cash_cost, value})
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
	severance
	predictions
	other
)

var AllGroups = []struct {
	Value  int
	TSName string
}{
	{production, "production"},
	{marketing, "marketing"},
	{personelle, "personelle"},
	{logistics, "logistics"},
	{materials, "materials"},
	{energy, "energy"},
	{product_development, "product_development"},
	{employee_training, "employee_training"},
	{loans, "loans"},
	{loan_intrest, "loan_intrest"},
	{bridge_loans, "bridge_loans"},
	{bridge_loan_intrest, "bridge_loan_intrest"},
	{taxes, "taxes"},
	{sales, "sales"},
	{severance, "severance"},
	{predictions, "predictions"},
	{other, "other"},
}

type Personelle_report struct {
	General    Personelle_sub_report
	Marketing  Personelle_sub_report
	Production Personelle_sub_report
}

type Personelle_sub_report struct {
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

	avg_productivity     float32
	minimum_productivity float32
	maximum_productivity float32
}

type Production_report struct {
	Machines_purchased       int
	Machines_sold            int
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

	Avr_quality_factor       float32
	Avr_durability_factor    float32
	Avr_ecology_factor       float32
	Avr_price_factor         float32
	Avr_coolness_factor      float32
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

	Avr_quality_factor       float32
	Avr_durability_factor    float32
	Avr_ecology_factor       float32
	Avr_price_factor         float32
	Avr_coolness_factor      float32
	Avr_bang_for_buck_factor float32
}

type Marketing_statistics struct {
	Product       Product_statistics
	Price         float64
	Bang_for_buck float64
	Promotion     float64
	// Place
}

type Product_statistics struct {
	Quality   float32
	Durabilty int
	Coolness  float32
	Ecology   float32
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

func (c Company) Mock_simulate_step(decisions Decisionsold, external_factors External_factors) Report {
	results := FinanceReportEntry{"Predicted sales", predictions, "The amount of you predict you'll make", true, float64(decisions.Sales_projection) * float64(decisions.Selling_price)}
	c.compile_reports(decisions, results, Purchasing_statistics{Products_sold: decisions.Sales_projection}, &Sales_statistics{}, external_factors)
	return c.Reports[len(c.Reports)-1]
}

// Gameloop functions
func (game_state *Game_state) Simulate_step() error {
	game_state.Step_simulated = false
	game_state.Step += 1

	game_state.External_factors.Month = game_state.Step

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
	Results, purchasing_statistics, err := game_state.Population.simulate_economy(&game_state.Companies, game_state.External_factors)
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
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_coolness_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_coolness_factor
	game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_bang_for_buck_factor = purchasing_statistics[len(purchasing_statistics)-1].Avr_bang_for_buck_factor

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

	game_state.Step_simulated = true

	for i := range game_state.Decisions_submitted {
		game_state.Decisions_submitted[i] = false
	}
	println("=================================================== ")
	println("               Simulation step done!\n")
	println("===================== RESULTS ===================== ")
	println("Month: ", game_state.Step)
	for i, c := range game_state.Companies {
		fmt.Printf("Company %d: %s:\n", i, c.Name)
		fmt.Printf("Products sold: %d\n", c.Reports[len(c.Reports)-1].Sales_report.Company_sales_statistics.Products_sold)
		println("")
	}

	total_products_sold := 0
	for _, c := range game_state.Companies {
		total_products_sold += c.Reports[len(c.Reports)-1].Sales_report.Company_sales_statistics.Products_sold
	}

	println("Total products sold: ", total_products_sold)

	missing_products := 0.0
	for _, p := range game_state.Population.Population {
		missing_products += float64(p.Base_need - len(p.Owned_products))
	}
	missing_products /= float64(len(game_state.Population.Population))
	fmt.Printf("avr missing products: %f\n", missing_products)

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
	decisions Decisionsold,
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
	company.Reports[len(company.Reports)-1].Financial_Report = company.calculate_budget(decisions, external_factors)

	return nil
}

func (c *Company) compile_sales_report(purchasing_statiscs Purchasing_statistics, Market_products_sold int) Sales_report {
	report := Sales_report{}

	report.Product_statistics.Quality = c.Offer.Product.Quality_factor
	report.Product_statistics.Durabilty = c.Offer.Product.Durabilty
	report.Product_statistics.Ecology = c.Offer.Product.Ecology_factor
	report.Product_statistics.Coolness = c.Offer.Product.Coolness_factor

	// ----------------

	report.Company_sales_statistics.Products_sold = purchasing_statiscs.Products_sold
	if len(c.Reports) >= 2 {
		report.Company_sales_statistics.Difference_to_previous_month = purchasing_statiscs.Products_sold - c.Reports[len(c.Reports)-2].Sales_report.Company_sales_statistics.Products_sold
	} else {
		report.Company_sales_statistics.Difference_to_previous_month = purchasing_statiscs.Products_sold
	}
	report.Company_sales_statistics.Product_demand = purchasing_statiscs.Product_demand
	report.Company_sales_statistics.Market_share = (float32(purchasing_statiscs.Products_sold) / float32(Market_products_sold)) * 100

	report.Company_sales_statistics.Avr_decision_factor = purchasing_statiscs.Avr_decision_factor
	report.Company_sales_statistics.Avr_purchasing_threshold = purchasing_statiscs.Avr_purchasing_threshold

	report.Company_sales_statistics.Avr_quality_factor = purchasing_statiscs.Avr_quality_factor
	report.Company_sales_statistics.Avr_durability_factor = purchasing_statiscs.Avr_durability_factor
	report.Company_sales_statistics.Avr_ecology_factor = purchasing_statiscs.Avr_ecology_factor
	report.Company_sales_statistics.Avr_price_factor = purchasing_statiscs.Avr_price_factor
	report.Company_sales_statistics.Avr_coolness_factor = purchasing_statiscs.Avr_coolness_factor
	report.Company_sales_statistics.Avr_bang_for_buck_factor = purchasing_statiscs.Avr_bang_for_buck_factor

	// ----------------

	report.Marketing_statistics.Product.Quality = c.Offer.Product.Quality_factor
	report.Marketing_statistics.Product.Durabilty = c.Offer.Product.Durabilty
	report.Marketing_statistics.Product.Coolness = c.Offer.Product.Coolness_factor
	report.Marketing_statistics.Product.Ecology = c.Offer.Product.Ecology_factor

	report.Marketing_statistics.Price = float64(c.Offer.Price)
	report.Marketing_statistics.Promotion = float64(c.Decision_history[len(c.Decision_history)-1].Marketing)

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

func Get_decisions(save_location string, number_of_companies int) ([]Decisionsold, error) {
	decisions := make([]Decisionsold, number_of_companies)
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
