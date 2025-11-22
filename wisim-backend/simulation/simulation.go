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
	ExternalFactors       ExternalFactors
	ProductComponents     ProductComponents
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
	SavvynessBias               float32
	SavvynessSpread             float32
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
	ID                int
	Name              string
	Balance           float64
	Loans             float64
	BridgeLoans       float64
	employeePool      Employee_pool
	productComponents *ProductComponents
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
	General struct {
		CompanyName string
	}
	Predictions struct {
		ProductSales map[string]int
		Steps        int
	}

	Finances struct {
		SetBankLoan float64
	}

	Products map[string]Decisions_product

	Employees struct {
		ProductionDeltas []Delta[Employee]
		MarketingDeltas  []Delta[Employee]

		SeverancePay float32
	}

	Production struct {
		Machines                 []Delta[Machine]
		Logistics                []Delta[Warehouse]
		MachineAssignmentPattern AssignmentPattern
	}

	Research Decisions_research
}

type Decisions_product struct {
	Name      string
	Outdated  bool
	Price     float32
	Promotion struct {
		Quantity   float32
		Quality    float32
		Price      float32
		Ecology    float32
		Ethics     float32
		Durability float32
	}

	Product Product
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
	Outdated         bool
	Product          Product
	ProductStats     ProductStats
	Price            float32
	PromotionQuality float32
	Promotion        struct {
		Quantity   float32
		Quality    float32
		Price      float32
		Ecology    float32
		Ethics     float32
		Durability float32
	}
}

type Product struct {
	ID        string
	CompanyID int
	Name      string

	Components struct {
		FormFactor string
		Frame      string
		Body       string
		Mechanism  string
		Misc       []string
	}

	TechLevels TechLevels

	MaterialQuality int
	ExtraDurability int
	ExtraQuality    int
}

type ProductStats struct {
	MiscSlots      int
	ProductionCost float32
	MaterialUse    float32

	Quality    float32
	Ecology    float32
	Ethics     float32
	Durability float32
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
	cash
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

		MaterialUsed float32
		EnergyUsed   float32
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
	ProductsSold  int
	ProductDemand int
	MarketShare   float32

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
	ImpressionCount   int
	// Place
}

type Customer struct {
	Base_need      int
	Owned_products []OwnedProduct
	KnownProducts  []string

	// Preferences moved to sepeate array in gamestate

	Purchashing_threshold float32
	Max_price             float32
	Savyness              float32
	brandSatisfaction     []Satisfaction

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

type OwnedProduct struct {
	ID                 int
	RemainingDurabilty int
}

type Satisfaction struct {
	ProductID      int
	DecisionFactor float32
	Satisfaction   float32
}

type ExternalFactors struct {
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

type ProductComponents struct {
	FormFactor map[string]Component
	Frame      map[string]Component
	Body       map[string]Component
	Mechanism  map[string]Component
	Misc       map[string]Component
}

type Component struct {
	Name               string
	MiscSlots          int
	ProductionCost     float32
	MaterialUse        float32
	Ecology            float32
	Ethics             float32
	Quality            float32
	Durability         float32
	ProductionLineCost float32
	Image              string
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
		game_state.Companies[i].prepareCompany(game_state.CurrentDecisions[i], game_state.ExternalFactors, game_state.Employees)
	}

	println("Simulating companies done!")

	println("---------------- Simulatig economy ----------------")

	purchasingStatistics := make(map[string]Purchasing_statistics)
	promotionImpression := make(map[string]int)
	err := game_state.Population.simulateEconomy(
		game_state.Companies,
		game_state.ExternalFactors,
		purchasingStatistics,
		promotionImpression,
	)
	if err != nil {
		return err
	}
	println("Simulatig economy done!")

	println("================ Compiling reports =============== ")

	game_state.MarketSalesStatistics = append(game_state.MarketSalesStatistics, Sales_statistics{})
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].ProductsSold = purchasingStatistics["-1"].ProductsSold
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].ProductDemand = purchasingStatistics["-1"].ProductDemand
	game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1].MarketShare = 100

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
			purchasingStatistics,
			promotionImpression,
			game_state.MarketSalesStatistics[len(game_state.MarketSalesStatistics)-1],

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

		totalProduction := 0
		totalProductsProducted := 0

		for _, r := range c.Reports[len(c.Reports)-1].ProductionReport.ProductSpecificReport {
			totalProduction += r.TotalProduction
			totalProductsProducted += r.TotalProductsProduced
		}

		printer.Printf("Total producton: %d\n", totalProduction)
		printer.Printf("Total products produced: %d\n", totalProductsProducted)
	}

	totalProductsSold := 0
	for _, c := range game_state.Companies {
		for _, salesReport := range c.Reports[len(c.Reports)-1].SalesReport {
			totalProductsSold += salesReport.ProductSalesStatistics.ProductsSold
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

	game_state.resetCurrentDecisions()

	return nil
}

func SimulateMockStep(company Company, decisions Decisions, externalFactors ExternalFactors, employeePool Employee_pool, steps int) Company {
	for range steps {
		company.prepareCompany(decisions, externalFactors, employeePool)

		purchasingStatistics := map[string]Purchasing_statistics{}

		for i := range decisions.Predictions.ProductSales {

			productsSold := min(decisions.Predictions.ProductSales[i], company.ProductsInStorage[i])

			company.Reports[len(company.Reports)-1].BalanceSheet.add_to_income_statement(
				"Sales of "+company.Offers[i].Product.Name,
				sales,
				fmt.Sprintf("%d %ss were sold in strores", productsSold, company.Offers[i].Product.Name),
				true,
				float64(productsSold*int(company.Offers[i].Price)))

			purchasingStatistics[i] = Purchasing_statistics{ProductDemand: decisions.Predictions.ProductSales[i], ProductsSold: productsSold}
			println(company.Offers[i].Product.Name, productsSold)
		}

		company.compileReports(
			decisions,
			purchasingStatistics,
			map[string]int{},
			Sales_statistics{},

			externalFactors,
		)

		decisions.resetDecisions()

	}

	assignCompanySlicesAndMaps(&company)

	return company
}

func (c *Company) prepareCompany(decisions Decisions, externalFactors ExternalFactors, employeePool Employee_pool) {
	c.Reports = append(c.Reports, Report{Month: externalFactors.Month})
	c.DecisionHistory = append(c.DecisionHistory, decisions)

	// renaming company
	c.Name = decisions.General.CompanyName

	fmt.Printf("--------------- Simulating company %d -------------- \n", c.ID)

	c.employeePool = employeePool

	println("Overworking employees...")
	c.simulateEmployees(externalFactors)

	// Add new products
	println("Developing Products...")

	c.addNewProducts(decisions)

	println("Posting advertisments...")
	c.calculatePromotion(decisions)

	println("Researching...")
	c.research(decisions)

	println("Producing...")
	c.calculateProduction(decisions, externalFactors)

	for productID := range c.Offers {
		c.ProductsInStorage[productID] += c.Reports[len(c.Reports)-1].ProductionReport.ProductSpecificReport[productID].TotalProductsProduced
	}

	println("Shipping...")
	c.calculate_logistics(decisions)
}

func (c *Company) addNewProducts(decisions Decisions) {
	for ID, productDecisions := range decisions.Products {
		if _, exists := c.Offers[ID]; !exists {
			c.Offers[ID] = c.newProduct(ID, c.ID, productDecisions.Name, decisions.Products[ID])
		} else {
			o := c.Offers[ID]
			o.Outdated = productDecisions.Outdated
			o.Product.Name = productDecisions.Product.Name
			o.Price = productDecisions.Price
			c.Offers[ID] = o
		}
	}
}

func (company *Company) compileReports(
	decisions Decisions,
	companyPurchasingStatistcs map[string]Purchasing_statistics,
	promotionImpressions map[string]int,
	marketPurchasingStatistics Sales_statistics,
	externalFactors ExternalFactors,
) error {
	company.Reports[len(company.Reports)-1].SalesReport = make(map[string]Sales_report)

	for productID := range company.Offers {
		company.compileSalesReport(
			companyPurchasingStatistcs,
			promotionImpressions,
			marketPurchasingStatistics.ProductsSold,
			company.Reports[len(company.Reports)-1].SalesReport,
		)

		company.ProductsInStorage[productID] -= company.Reports[len(company.Reports)-1].SalesReport[productID].ProductSalesStatistics.ProductsSold
	}

	// Finance
	company.calculateBudget(decisions, externalFactors)
	company.Reports[len(company.Reports)-1].BalanceSheet.Assets = cleanUpFinanceReportEntries(company.Reports[len(company.Reports)-1].BalanceSheet.Assets)
	company.Reports[len(company.Reports)-1].BalanceSheet.InvoiceLog = cleanUpFinanceReportEntries(company.Reports[len(company.Reports)-1].BalanceSheet.InvoiceLog)
	company.Reports[len(company.Reports)-1].BalanceSheet.Liabilities = cleanUpFinanceReportEntries(company.Reports[len(company.Reports)-1].BalanceSheet.Liabilities)

	// Personelle
	company.Reports[len(company.Reports)-1].PersonelleReport = company.compilePersonelleReport(decisions)

	return nil
}

func (c *Company) compilePersonelleReport(decisions Decisions) Personelle_report {
	personelleReport := Personelle_report{}

	personelleReport.General = c.compilePersonelleSubreport(decisions, Employee_type_all)
	personelleReport.Marketing = c.compilePersonelleSubreport(decisions, Employee_type_marketing)
	personelleReport.Production = c.compilePersonelleSubreport(decisions, Employee_type_production)

	return personelleReport
}

func (c *Company) compilePersonelleSubreport(decisions Decisions, employeeType Employee_type) Personelle_sub_report {
	var subReport Personelle_sub_report
	employeeIDs := c.employeePool.Get_employees_of_company(c.ID, employeeType)

	subReport.NumberOfEmployees = len(employeeIDs)

	if subReport.NumberOfEmployees == 0 {
		return subReport
	}

	var employeeDeltas []Delta[Employee] // We can trust that the employees exist because we checked this when "simulating employees"
	switch employeeType {
	case Employee_type_marketing:
		employeeDeltas = decisions.Employees.MarketingDeltas
	case Employee_type_production:
		employeeDeltas = decisions.Employees.ProductionDeltas
	case Employee_type_all:
		employeeDeltas = decisions.Employees.ProductionDeltas
		employeeDeltas = append(employeeDeltas, decisions.Employees.MarketingDeltas...)
	}

	for _, eDelta := range employeeDeltas {
		switch eDelta.Change {
		case Delta_New:
			subReport.NumberOfHires += 1
		case Delta_Remove:
			subReport.NumberOfDepartures += 1
		}
	}

	pay := make([]float32, len(employeeIDs))
	skill := make([]float32, len(employeeIDs))
	motivation := make([]float32, len(employeeIDs))
	productivity := make([]float32, len(employeeIDs))
	for i, e := range employeeIDs {
		pay[i] = c.employeePool[e].Pay
		skill[i] = c.employeePool[e].Skill
		motivation[i] = c.employeePool[e].Motivation
		productivity[i] = c.employeePool[e].Motivation * c.employeePool[e].Skill * c.employeePool[e].WorkingHours // TODO: Make sure this is actually accurate
	}

	subReport.AvgPay = avr(pay)
	subReport.MaximumPay = slices.Max(pay)
	subReport.MinimumPay = slices.Min(pay)
	subReport.StandardDevPay = std_dev(pay...)

	subReport.AvgSkill = avr(skill)
	subReport.MaximumSkill = slices.Max(skill)
	subReport.MinimumSkill = slices.Min(skill)
	subReport.StandardDevSkill = std_dev(skill...)

	subReport.AvgMotivation = avr(motivation)
	subReport.MaximumMotivation = slices.Max(motivation)
	subReport.MinimumMotivation = slices.Min(motivation)
	subReport.StandardDevMotivation = std_dev(motivation...)

	subReport.AvgProductivity = avr(productivity)
	subReport.MaximumProductivity = slices.Max(productivity)
	subReport.MinimumProductivity = slices.Min(productivity)
	subReport.StandardDevProductivity = std_dev(productivity...)

	return subReport
}

func (c *Company) compileSalesReport(
	purchasingStatiscs map[string]Purchasing_statistics,
	impressions map[string]int,
	MarketProductsSold int,
	salesReportsMap map[string]Sales_report,
) {
	// ----------------

	for productID, productSpecificPurchasingStatiscs := range purchasingStatiscs {
		if _, isCompaniesProducts := c.Offers[productID]; !isCompaniesProducts {
			continue
		}

		salesStatistics := Sales_statistics{}

		salesStatistics.ProductsSold = productSpecificPurchasingStatiscs.ProductsSold
		println("salesStatistics.ProductsSold: ", salesStatistics.ProductsSold, productID)
		salesStatistics.ProductDemand = productSpecificPurchasingStatiscs.ProductDemand
		if MarketProductsSold != 0 {
			salesStatistics.MarketShare = (float32(productSpecificPurchasingStatiscs.ProductsSold) / float32(MarketProductsSold))
		} else {
			salesStatistics.MarketShare = 0
		}

		salesStatistics.AvrDecisionFactor = productSpecificPurchasingStatiscs.AvrDecisionFactor
		salesStatistics.AvrPurchasingThreshold = productSpecificPurchasingStatiscs.AvrPurchasingThreshold

		salesStatistics.AvrQualityFactor = productSpecificPurchasingStatiscs.AvrPurchasingFactors[propertiesQuality]
		salesStatistics.AvrDurabilityFactor = productSpecificPurchasingStatiscs.AvrPurchasingFactors[propertiesDurability]
		salesStatistics.AvrEcologyFactor = productSpecificPurchasingStatiscs.AvrPurchasingFactors[propertiesEcology]
		salesStatistics.AvrPriceFactor = productSpecificPurchasingStatiscs.AvrPurchasingFactors[propertiesPrice]
		salesStatistics.AvrEthicsFactor = productSpecificPurchasingStatiscs.AvrPurchasingFactors[propertiesEthics]
		// salesReport.Company_sales_statistics.Avr_coolness_factor = productSpecificPurchasing_statiscs.Avr_coolness_factor
		salesStatistics.AvrBangForBuckFactor = productSpecificPurchasingStatiscs.AvrPurchasingFactors[propertiesBangForBuck]

		marketingStatistics := Marketing_statistics{}
		marketingStatistics.Quality = c.Offers[productID].ProductStats.Quality
		marketingStatistics.Durabilty = int(c.Offers[productID].ProductStats.Durability)
		marketingStatistics.Ethics = c.Offers[productID].ProductStats.Ethics
		// reportMarketingStatistics.Coolness = offer.Product.Coolness_factor
		marketingStatistics.Ecology = c.Offers[productID].ProductStats.Ecology

		marketingStatistics.Price = float64(c.Offers[productID].Price)
		marketingStatistics.PromotionQuantity = float64(c.Offers[productID].Promotion.Quantity)
		marketingStatistics.PromotionQuality = float64(c.Offers[productID].PromotionQuality)
		marketingStatistics.ImpressionCount = impressions[productID]

		salesReportsMap[productID] = Sales_report{salesStatistics, marketingStatistics}
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
