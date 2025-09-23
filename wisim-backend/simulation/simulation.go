package simulation

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type SaveGame struct {
	Population []byte // Binary
	GameState  GameState
}

type GameState struct {
	Step     int
	GameName string

	Population            Population
	Employees             Employee_pool
	employeesArray        []Employee
	Companies             []Company
	CurrentDecisions      []Decisions
	DecisionsSubmitted    []bool
	MarketSalesStatistics []Sales_statistics
	ExternalFactors       External_factors
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
	ID           int
	Name         string
	Balance      float64
	Loans        float64
	BridgeLoans  float64
	employeePool Employee_pool
	// Global_effects   []Effect
	DecisionHistory []Decisions

	Reports []Report
	// Research and development
	BaseMarketingStrength float32
	Tech                  TechLevels

	// Product
	Offers map[string]Offer

	// Fulfillment
	Warehouses        []Warehouse
	ProductsInStorage map[string]int

	// Production
	Machines []Machine
}

// NOTE: Employees keep track of their employers, not the companies
type TechLevels struct {
	Quality        float32
	Ecology        float32
	Durability     float32
	ProductionCost float32
	MaterialUse    float32
}

type Decisions struct {
	Predictions struct {
		SalesPrediction int
	}

	Finances struct {
		Set_bank_loan float64
	}

	Products map[string]Decisions_product

	Employees struct {
		Production_deltas []Delta[Employee]
		Marketing_deltas  []Delta[Employee]

		Severance_pay float32
	}

	Production struct {
		Machines  []Delta[Machine]
		Logistics []Delta[Warehouse]
	}

	Research Decisions_research
}

type Decisions_product struct {
	Price          float32
	Name           string
	ProductionGoal int

	Materials struct {
		Quality         float32
		Ecology         float32
		EthicalSourcing float32
	}

	Manufacturing struct {
		Quality            float32
		EcologicalEnergy   float32
		MaterialEfficiency float32
		Durability         float32
		MaxDurability      int
	}

	Promotion struct {
		Quantity        float64
		StyleQuality    float32
		StyleEcology    float32
		StyleEthics     float32
		StyleDurability float32
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
	Status           string
	Product          Product
	Price            float32
	PromotionQuality float32
	PromotionGoal    struct {
		Quantity        float64
		StyleQuality    float32
		StyleEcology    float32
		StyleEthics     float32
		StyleDurability float32
	}
}

type Product struct {
	ID             string
	Name           string
	Weight         float32
	MaterialUse    float32
	ProductionCost float32

	Ethics  float32 // TODO: Implement Ethicss Factor
	Quality float32
	Ecology float32
	// Coolness_factor float32

	Durabilty int
}

type Machine struct {
	ID int

	ProductionCapacity int
	RequiredWorkers    int
	MinimumWorkers     int
	AssignedWorkersIDs []int
	EnergyUse          float32
	Value              float32
	MaintananceCost    float32 // Monthly
	AssignedProductID  string
}

func (m Machine) get_id() int {
	return m.ID
}

type Warehouse struct {
	ID             int
	Capacity       int
	OperatingCosts float32
	Value          float32
}

func (w Warehouse) get_id() int {
	return w.ID
}

type Employee struct {
	ID              int
	Name            string
	Employer        int
	MonthsAtCompany int

	EmployeeType Employee_type

	Motivation    float32
	Skill         float32
	ExtraTraining float32

	// Global_effect *Effect

	Pay          float32
	Bonus        float32
	WorkingHours float32
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

type Effect struct {
	ID          int
	Name        string
	Description string
	// Effect_function *Effect_Function
}

type Effect_Function func(int)

type Report struct {
	Month int

	FinancialReport  Financial_report
	BalanceSheet     Balance_sheet
	PersonelleReport Personelle_report
	ProductionReport Production_report
	SalesReport      map[string]Sales_report
}

type Financial_report struct {
	// Income
	Income struct {
		GrossSales  float64
		OtherIncome float64
		CostOfSales float64
		GrossProfit float64
	}
	OperatingExpenses struct {
		Advertising            float64
		FacilitiesAndLogistics float64
		ResearchAndDevelopment float64
	}
	NonOperatingExpenses struct {
		WriteOffs           float64
		LoanInterest        float64
		LoanRepayment       float64
		BridgeLoanIntrest   float64
		BridgeLoanRepayment float64
		Other               float64
		Taxes               float64
	}
	Totals struct {
		TotalOperatingExpenses    float64
		TotalNonOperatingExpenses float64
		IncomeBeforeTax           float64
		NetIncome                 float64
		Cashflow                  float64
	}
}

func (f *Balance_sheet) add_to_income_statement(name string, group int, info string, cash_cost bool, value float64) *FinanceReportEntry {
	entry := FinanceReportEntry{name, group, info, cash_cost, value}
	f.InvoiceLog = append(f.InvoiceLog, entry)
	return &f.InvoiceLog[len(f.InvoiceLog)-1]
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

	InvoiceLog []FinanceReportEntry

	Assets      []FinanceReportEntry
	Liabilities []FinanceReportEntry
}

type FinanceReportEntry struct {
	Name     string
	Group    int
	Info     string
	CashCost bool

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
	NumberOfEmployees  int
	NumberOfHires      int
	NumberOfDepartures int

	AvgPay         float32
	MinimumPay     float32
	MaximumPay     float32
	StandardDevPay float32

	MinimumSkill     float32
	MaximumSkill     float32
	AvgSkill         float32
	StandardDevSkill float32

	MinimumMotivation     float32
	MaximumMotivation     float32
	AvgMotivation         float32
	StandardDevMotivation float32

	MinimumProductivity     float32
	MaximumProductivity     float32
	AvgProductivity         float32
	StandardDevProductivity float32
}

type Production_report struct {
	MachinesPurchased      int
	MachinesSold           int
	WorkerSurplus          int
	AvgMachineProductivity float32
	// Max_machine_productivity float32
	// Min_machine_productivity float32

	ProductSpecificReport map[string]struct {
		TotalProduction  int
		BaseProduction   int
		BonusProduction  int
		ExcessProduction int

		TotalProductsProduced int
		BaseProductsProduced  int
		BonusProductsProduced int
	}

	MaterialUsed float32
	EnergyUsed   float32

	WarehousesBought int
}

type Sales_report struct {
	ProductSalesStatistics Sales_statistics
	MarketingStatistics    Marketing_statistics
}

type Purchasing_statistics struct {
	ProductNumber int
	ProductsSold  int
	ProductDemand int

	AvrDecisionFactor      float32
	AvrPurchasingThreshold float32

	AvrPurchasingFactors Properties
}

type Research_statistics struct {
	QualityDevelopmentInvestment                 float64
	QualityDevelopmentInvestmentEffectiveness    float64
	DurabilityDevelopmentInvestment              float64
	DurabilityDevelopmentInvestmentEffectiveness float64
	EcologicalProductionInvestment               float32 // Decreases material use
	EcologicalProductionInvestmentEffectiveness  float32 // Decreases material use
}

type Sales_statistics struct {
	Products_sold  int
	Product_demand int
	Market_share   float32

	AvrDecisionFactor      float32
	AvrPurchasingThreshold float32

	AvrQualityFactor    float32
	AvrDurabilityFactor float32
	AvrEcologyFactor    float32
	AvrPriceFactor      float32
	AvrEthicsFactor     float32
	/// Avr_coolness_factor      float32
	AvrBangForBuckFactor float32
}

type Marketing_statistics struct {
	Quality   float32
	Durabilty int
	// Coolness  float32
	Ethics            float32
	Ecology           float32
	Price             float64
	BangForBuck       float64
	PromotionQuantity float64
	PromotionQuality  float64
	// Place
}

type Customer struct {
	Base_need      int
	Owned_products []Owned_product

	// Preferences moved to sepeate array in gamestate

	Purchashing_threshold float32
	Max_price             float32
	Satisfaction          []Satisfaction

	Brand_loyalty_factor float32
	Loyalties            []float32
}

type Properties [8]float32 // last 2 slots are left empty (for optimisation)

const (
	propertiesQuality = iota
	propertiesEcology
	propertiesEthics
	propertiesPrice
	propertiesBangForBuck
	propertiesDurability
)

type Population struct {
	Population  []Customer
	Preferences []Properties
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
	Inflation              float32
	IntrestRate            float32
	BridgeLoansIntrestRate float32
	EconomicSituationIndex float32
	TaxRate                float32 // as decimal

	// Personelle
	ProductionMinimumWage float32
	MarketingMinimumWage  float32

	// Prdoction
	MachineOnOffer       Machine
	ExternalStoragePrice float32 // per item
	EnergyPrice          float32 // per unit of energy
	MaterialPrice        float32 // per unit of material

	MachineDepreciationRate float32 // in decimal
}

// #####################################################################################################
// ##########     _         _              __                      _    _                     ##########
// ##########    | |       | |            / _|                    | |  (_)                    ##########
// ##########  __| |  __ _ | |_  __ _    | |_  _   _  _ __    ___ | |_  _   ___   _ __   ___  ##########
// ########## / _` | / _` || __|/ _` |   |  _|| | | || '_ \  / __|| __|| | / _ \ | '_ \ / __| ##########
// ##########| (_| || (_| || |_| (_| |   | |  | |_| || | | || (__ | |_ | || (_) || | | |\__ \ ##########
// ########## \__,_| \__,_| \__|\__,_|   |_|   \__,_||_| |_| \___| \__||_| \___/ |_| |_||___/ ##########
// #####################################################################################################

/*
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
*/

// Gameloop functions

func (game_state *GameState) SimulateStep() error {
	game_state.Step += 1

	game_state.ExternalFactors.Month = game_state.Step

	if len(game_state.Companies) != len(game_state.CurrentDecisions) {
		return errors.New("amount of decisions does not match number of companies")
	}
	println("============ Simulating Hiring / Firing ===========")
	game_state.handleEmployeeDeltas()

	println("=============== Simulating companies ==============")
	for i := range game_state.Companies {
		c := &game_state.Companies[i]
		c.Reports = append(c.Reports, Report{Month: game_state.ExternalFactors.Month})
		c.DecisionHistory = append(c.DecisionHistory, game_state.CurrentDecisions[i])

		fmt.Printf("--------------- Simulating company %d -------------- \n", i)

		c.employeePool = game_state.Employees

		// Add new products
		println("Developing Products...")

		b, _ := json.MarshalIndent(c.Offers, "", "    ")
		fmt.Printf("%s\n", b)

		b, _ = json.MarshalIndent(game_state.CurrentDecisions[i].Products, "", "    ")
		fmt.Printf("%s\n", b)
		for ID, decisions := range game_state.CurrentDecisions[i].Products {
			if _, exists := c.Offers[ID]; !exists {
				c.Offers[ID] = c.newProduct(ID, decisions.Name, game_state.CurrentDecisions[i].Products[ID])
			}
		}

		println("Overworking employees...")
		c.simulateEmployees(game_state.ExternalFactors)

		println("Posting advertisments...")
		c.calculatePromotion(game_state.CurrentDecisions[i])

		println("Researching...")
		c.research(game_state.CurrentDecisions[i])

		println("Producing...")
		c.calculateProduction(game_state.CurrentDecisions[i], game_state.ExternalFactors)
		for productID := range c.Offers {
			c.ProductsInStorage[productID] += c.Reports[len(c.Reports)-1].ProductionReport.ProductSpecificReport[productID].TotalProductsProduced
		}

		println("Shipping...")
		c.calculate_logistics(game_state.CurrentDecisions[i])
	}

	println("Simulating companies done!")

	println("---------------- Simulatig economy ----------------")

	purchasingStatistics := make(map[string]Purchasing_statistics)
	Results, err := game_state.Population.simulateEconomy(
		&game_state.Companies,
		game_state.ExternalFactors,
		purchasingStatistics,
	)
	if err != nil {
		return err
	}
	println("Simulatig economy done!")

	println("================ Compiling reports =============== ")

	game_state.MarketSalesStatistics = append(game_state.MarketSalesStatistics, Sales_statistics{})
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].Products_sold = purchasingStatistics["-1"].ProductsSold
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].Product_demand = purchasingStatistics["-1"].ProductDemand
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].Market_share = 100

	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].AvrDecisionFactor = purchasingStatistics["-1"].AvrDecisionFactor
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].AvrPurchasingThreshold = purchasingStatistics["-1"].AvrPurchasingThreshold

	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].AvrQualityFactor = purchasingStatistics["-1"].AvrPurchasingFactors[propertiesQuality]
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].AvrDurabilityFactor = purchasingStatistics["-1"].AvrPurchasingFactors[propertiesDurability]
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].AvrEcologyFactor = purchasingStatistics["-1"].AvrPurchasingFactors[propertiesEcology]
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].AvrPriceFactor = purchasingStatistics["-1"].AvrPurchasingFactors[propertiesPrice]
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].AvrEthicsFactor = purchasingStatistics["-1"].AvrPurchasingFactors[propertiesEthics]
	// game_state.Market_sales_statistics[len(game_state.Market_sales_statistics)-1].Avr_coolness_factor = purchasingStatistics["-1"].Avr_coolness_factor
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].AvrBangForBuckFactor = purchasingStatistics["-1"].AvrPurchasingFactors[propertiesBangForBuck]

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
		game_state.Companies[i].compileReports(
			game_state.CurrentDecisions[i],
			Results[i],
			purchasingStatistics,
			&game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1],
			game_state.ExternalFactors,
		)

	}

	for i := range game_state.DecisionsSubmitted {
		game_state.DecisionsSubmitted[i] = false
	}
	println("=================================================== ")
	println("               Simulation step done!\n")
	println("===================== RESULTS ===================== ")
	println("Month: ", game_state.Step)

	printer := message.NewPrinter(language.Swedish)
	for _, c := range game_state.Companies {
		/*
			printer.Printf("Company %d: %s:\n", i, c.Name)
			printer.Printf("Total Production: %d\n", c.Reports[len(c.Reports)-1].Production_report..Total_production)
			printer.Printf("Excess Production: %d\n", c.Reports[len(c.Reports)-1].Production_report.Excess_production)
			printer.Printf("Total Products Produced: %d\n", c.Reports[len(c.Reports)-1].Production_report.Total_products_produced)
			printer.Printf("Products sold: %d\n", c.Reports[len(c.Reports)-1].Sales_report.Product_sales_statistics.Products_sold)
		*/
		printer.Printf("--> Net profit: %.2f", c.Reports[len(c.Reports)-1].FinancialReport.Totals.NetIncome)
		println("")
		printer.Printf("Number of employees: %d\n", c.Reports[len(c.Reports)-1].PersonelleReport.General.NumberOfEmployees)
		printer.Printf("Number of production employees: %d\n", c.Reports[len(c.Reports)-1].PersonelleReport.Production.NumberOfEmployees)
		printer.Printf("Number of marketing employees: %d\n", c.Reports[len(c.Reports)-1].PersonelleReport.Marketing.NumberOfEmployees)
		printer.Printf("Number of hires: %d\n", c.Reports[len(c.Reports)-1].PersonelleReport.General.NumberOfHires)
		printer.Printf("Number of departures: %d\n", c.Reports[len(c.Reports)-1].PersonelleReport.General.NumberOfDepartures)
		printer.Printf("Avr Pay: %.2f\n", c.Reports[len(c.Reports)-1].PersonelleReport.General.AvgPay)
		printer.Printf("Avr Skill: %.2f\n", c.Reports[len(c.Reports)-1].PersonelleReport.General.AvgSkill)
		printer.Printf("Avr Motivation: %.2f\n", c.Reports[len(c.Reports)-1].PersonelleReport.General.AvgMotivation)
		printer.Printf("Avr Productivity: %.2f\n", c.Reports[len(c.Reports)-1].PersonelleReport.General.AvgProductivity)
		println("")
	}

	totalProductsSold := 0
	for _, c := range game_state.Companies {
		for _, salesReport := range c.Reports[len(c.Reports)-1].SalesReport {
			totalProductsSold += salesReport.ProductSalesStatistics.Products_sold
		}
	}

	println("Total products sold: ", totalProductsSold)

	missingProducts := 0
	for _, p := range game_state.Population.Population {
		missingProducts += p.Base_need - len(p.Owned_products)
	}
	avrMissingProducts := float64(missingProducts) / float64(len(game_state.Population.Population))
	fmt.Printf("avr missing products: %f\n", avrMissingProducts)

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

func (company *Company) compileReports(
	decisions Decisions,
	results FinanceReportEntry,
	companyPurchasingStatistcs map[string]Purchasing_statistics,
	marketPurchasingStatistics *Sales_statistics,
	externalFactors External_factors,
) error {
	company.Reports[len(company.Reports)-1].SalesReport = make(map[string]Sales_report)

	for productID := range company.Offers {
		company.compile_sales_report(
			companyPurchasingStatistcs,
			marketPurchasingStatistics.Products_sold,
			company.Reports[len(company.Reports)-1].SalesReport,
		)

		company.ProductsInStorage[productID] -= marketPurchasingStatistics.Products_sold
	}

	company.Reports[len(company.Reports)-1].BalanceSheet.InvoiceLog = append(company.Reports[len(company.Reports)-1].BalanceSheet.InvoiceLog, results)

	// Finance
	company.calculate_budget(decisions, externalFactors)
	company.Reports[len(company.Reports)-1].BalanceSheet.Assets = clean_up_financeReportEntries(company.Reports[len(company.Reports)-1].BalanceSheet.Assets)
	company.Reports[len(company.Reports)-1].BalanceSheet.InvoiceLog = clean_up_financeReportEntries(company.Reports[len(company.Reports)-1].BalanceSheet.InvoiceLog)
	company.Reports[len(company.Reports)-1].BalanceSheet.Liabilities = clean_up_financeReportEntries(company.Reports[len(company.Reports)-1].BalanceSheet.Liabilities)

	// Personelle
	company.Reports[len(company.Reports)-1].PersonelleReport = company.compile_personelle_report(decisions)

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
	employee_ids := c.employeePool.Get_employees_of_company(c.ID, employee_type)

	sub_report.NumberOfEmployees = len(employee_ids)

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
			sub_report.NumberOfHires += 1
		} else if e_delta.Change == Delta_Remove {
			sub_report.NumberOfDepartures += 1
		}
	}

	pay := make([]float32, len(employee_ids))
	skill := make([]float32, len(employee_ids))
	motivation := make([]float32, len(employee_ids))
	productivity := make([]float32, len(employee_ids))
	for i, e := range employee_ids {
		pay[i] = c.employeePool[e].Pay
		skill[i] = c.employeePool[e].Skill
		motivation[i] = c.employeePool[e].Motivation
		productivity[i] = c.employeePool[e].Motivation * c.employeePool[e].Skill * c.employeePool[e].WorkingHours // TODO: Make sure this is actually accurate
	}

	sub_report.AvgPay = avr(pay)
	sub_report.MaximumPay = slices.Max(pay)
	sub_report.MinimumPay = slices.Min(pay)
	sub_report.StandardDevPay = std_dev(pay...)

	sub_report.AvgSkill = avr(skill)
	sub_report.MaximumSkill = slices.Max(skill)
	sub_report.MinimumSkill = slices.Min(skill)
	sub_report.StandardDevSkill = std_dev(skill...)

	sub_report.AvgMotivation = avr(motivation)
	sub_report.MaximumMotivation = slices.Max(motivation)
	sub_report.MinimumMotivation = slices.Min(motivation)
	sub_report.StandardDevMotivation = std_dev(motivation...)

	sub_report.AvgProductivity = avr(productivity)
	sub_report.MaximumProductivity = slices.Max(productivity)
	sub_report.MinimumProductivity = slices.Min(productivity)
	sub_report.StandardDevProductivity = std_dev(productivity...)

	return sub_report
}

func (c *Company) compile_sales_report(purchasing_statiscs map[string]Purchasing_statistics, Market_products_sold int, salesReportsMap map[string]Sales_report) {
	// ----------------

	for productId, productSpecificPurchasing_statiscs := range purchasing_statiscs {
		salesReport := Sales_statistics{}

		salesReport.Products_sold = productSpecificPurchasing_statiscs.ProductsSold
		salesReport.Product_demand = productSpecificPurchasing_statiscs.ProductDemand
		if Market_products_sold != 0 {
			salesReport.Market_share = (float32(productSpecificPurchasing_statiscs.ProductsSold) / float32(Market_products_sold))
		} else {
			salesReport.Market_share = 0
		}

		salesReport.AvrDecisionFactor = productSpecificPurchasing_statiscs.AvrDecisionFactor
		salesReport.AvrPurchasingThreshold = productSpecificPurchasing_statiscs.AvrPurchasingThreshold

		salesReport.AvrQualityFactor = productSpecificPurchasing_statiscs.AvrPurchasingFactors[propertiesQuality]
		salesReport.AvrDurabilityFactor = productSpecificPurchasing_statiscs.AvrPurchasingFactors[propertiesDurability]
		salesReport.AvrEcologyFactor = productSpecificPurchasing_statiscs.AvrPurchasingFactors[propertiesEcology]
		salesReport.AvrPriceFactor = productSpecificPurchasing_statiscs.AvrPurchasingFactors[propertiesPrice]
		salesReport.AvrEthicsFactor = productSpecificPurchasing_statiscs.AvrPurchasingFactors[propertiesEthics]
		// salesReport.Company_sales_statistics.Avr_coolness_factor = productSpecificPurchasing_statiscs.Avr_coolness_factor
		salesReport.AvrBangForBuckFactor = productSpecificPurchasing_statiscs.AvrPurchasingFactors[propertiesBangForBuck]

		marketingStatistics := Marketing_statistics{}
		marketingStatistics.Quality = c.Offers[productId].Product.Quality
		marketingStatistics.Durabilty = c.Offers[productId].Product.Durabilty
		marketingStatistics.Ethics = c.Offers[productId].Product.Ethics
		// reportMarketingStatistics.Coolness = offer.Product.Coolness_factor
		marketingStatistics.Ecology = c.Offers[productId].Product.Ecology

		marketingStatistics.Price = float64(c.Offers[productId].Price)
		marketingStatistics.PromotionQuantity = c.Offers[productId].PromotionGoal.Quantity
		marketingStatistics.PromotionQuality = float64(c.Offers[productId].PromotionQuality)

		salesReportsMap[productId] = Sales_report{salesReport, marketingStatistics}
	}
}

// ##################################################
// ##########                    _         ##########
// ##########                   (_)        ##########
// ##########  _ __ ___    __ _  _  _ __   ##########
// ########## | '_ ` _ \  / _` || || '_ \  ##########
// ########## | | | | | || (_| || || | | | ##########
// ########## |_| |_| |_| \__,_||_||_| |_| ##########
// ##################################################

func Get_decisions_from_file(saveLocation string, numberOfCompanies int) ([]Decisions, error) {
	decisions := make([]Decisions, numberOfCompanies)
	for i := range decisions {
		decisions_json, err := os.ReadFile(fmt.Sprintf("%s/decisions_company_%d.json", saveLocation, i))
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

func (game_state GameState) Save_game(location string, compress bool) error {
	filename := fmt.Sprintf(
		"%s-%d.json",
		game_state.GameName,
		game_state.Step,
	)

	var save SaveGame

	var populationBuffer bytes.Buffer
	encoder := gob.NewEncoder(&populationBuffer)
	err := encoder.Encode(game_state.Population.Population)
	if err != nil {
		return err
	}
	save.Population = populationBuffer.Bytes()
	game_state.Population = Population{}

	save.GameState = game_state

	saveFile, err := json.MarshalIndent(save, "", "    ")
	if err != nil {
		return err
	}

	// Turning file into zip
	// (IDK what's happening)

	if compress {

		zipFileBuffer := new(bytes.Buffer)

		w := zip.NewWriter(zipFileBuffer)

		file, err := w.Create(filename)
		if err != nil {
			return err
		}

		_, err = file.Write(saveFile)
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		err = os.WriteFile(fmt.Sprint(location, "/", filename, ".zip"), zipFileBuffer.Bytes(), 0644)
		if err != nil {
			return err
		}

		return nil
	} // else if !compress {

	err = os.WriteFile(fmt.Sprint(location, "/", filename), saveFile, 0644)
	if err != nil {
		return err
	}

	return nil
}
