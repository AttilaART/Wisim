// NOTE: Most of the types defined in this file were generated from Go types through ChatGPT

export const delta = {
	Delta_New: 0,
	Delta_Change: 1,
	Delta_Remove: 2
};

/** @typedef {number} AssignmentPattern */

export const assignmentPatterns = {
	fillMachines: 0,
	distributeWorkers: 1
};

/**
 * @typedef {Object} clientState
 * @property {boolean} predictionMode
 * @property {Company} Company
 * @property {Decisions} Decisions
 * @property {ExternalFactors} ExternalFactors
 * @property {{marketing: Employee[], production: Employee[]}} Employees
 * @property {{marketing: Employee[], production: Employee[]}} Unemployed
 * @property {ProductComponents} productComponents
 */

/**
 * @typedef {Object} Decisions
 * @property {{CompanyName: string, CEO: string}} General
 * @property {{ProductSales: Object.<string, number>, Steps: number}} Predictions
 * @property {{SetBankLoan: number}} Finances
 * @property {Object.<string ,Decisions_product>} Products
 * @property {{ProductionDeltas: Delta<Employee>[], MarketingDeltas: Delta<Employee>[], SeverancePay: number}} Employees
 * @property {{Production_goal: number, Machines: Delta<Machine>[], Logistics: Delta<Warehouse>[], MachineAssignmentPattern: AssignmentPattern}} Production
 * @property {Decisions_research} Research
 */

/**
 * @typedef {Object} ProductComponents
 * @property {Object.<string, Component>} FormFactor
 * @property {Object.<string, Component>} Frame
 * @property {Object.<string, Component>} Body
 * @property {Object.<string, Component>} Mechanism
 * @property {Object.<string, Component>} Misc
 */

/**
 * @typedef {Object} Decisions_product
 * @property {string} Name
 * @property {number} Price
 * @property {boolean} Outdated
 * @property {Object} Promotion
 * @property {number} Promotion.Quantity
 * @property {number} Promotion.Price
 * @property {number} Promotion.Quality
 * @property {number} Promotion.Ecology
 * @property {number} Promotion.Ethics
 * @property {number} Promotion.Durability
 * @property {Product} Product
 */

/**
 * @typedef {Object} Decisions_research
 * @property {number} Quality - Research focus on quality
 * @property {number} Durability - Research focus on durability
 * @property {number} Ecology - Research focus on ecological aspects
 * @property {number} Promotion - Research impact on promotion
 * @property {number} Production_cost - Research impact on production cost
 */

/**
 * @template T
 * @typedef {Object} Delta
 * @property {Change} Change
 * @property {T} Item
 */

/**
 * @typedef {number} Change
 */

/**
 * @typedef {Object} Employee
 * @property {number} ID - Unique identifier of the employee
 * @property {string} Name - Employee's name
 * @property {number} Employer - Reference to the employer ID
 *
 * @property {Employee_type} EmployeeType - Type of the employee
 *
 * @property {number} Motivation - Motivation level
 * @property {number} Skill - Skill level
 * @property {number} ExtraTraining - Extra training factor
 *
 * @property {number} Pay - Base pay
 * @property {number} Bonus - Bonus amount
 * @property {number} WorkingHours - Number of working hours
 */

/**
 * @typedef {number} Employee_type
 */

/**
 * @typedef {Object} Machine
 * @property {number} ID - Unique identifier of the machine
 * @property {number} ProductionCapacity - Maximum production capacity
 * @property {number} RequiredWorkers - Number of workers required
 * @property {number} MinimumWorkers - Minimum number of workers allowed
 * @property {Employee[]} AssignedWorkersIds - Workers assigned to this machine
 * @property {number} EnergyUse - Energy consumption
 * @property {number} Value - Machine value
 * @property {number} MaintananceCost
 * @property {string} AssignedProductID
 */

/**
 * @typedef {Object} Warehouse
 * @property {number} ID - Unique identifier of the warehouse
 * @property {number} Capacity - Storage capacity
 * @property {number} OperatingCosts - Operating costs of the warehouse
 * @property {number} Value - Monetary value of the warehouse
 */

/**
 * @typedef {Object} Company
 *
 * @property {number} ID
 * @property {string} CEO
 * @property {string} Name
 * @property {number} Balance - Current financial balance
 * @property {number} Loans - Outstanding loans
 * @property {number} BridgeLoans - Outstanding bridge loans
 * @property {Decisions[]} DecisionHistory - History of company decisions
 * @property {Report[]} Reports - Reports related to the company
 *
 * @property {number} BaseMarketingStrength - Base marketing strength from R&D
 * @property {TechLevels} Tech - Base marketing strength from R&D
 *
 * @property {Object.<string, Offer>} Offers - Product offer
 *
 * @property {Warehouse[]} Warehouses - Warehouses owned by the company
 * @property {Object.<string, number>} ProductsInStorage - Number of items in storage
 *
 * @property {Machine[]} Machines - Machines owned by the company
 */

/**
 * @typedef {{Quality: number, Ecology: number, Durability: number, ProductionCost: number, MaterialUse: number}} TechLevels
 */

/**
 * @typedef {Object} Product
	* @property {string} ID
	* @property {number} CompanyID
	* @property {string} Name

	* @property {object} Components
  * @property {string?} Components.FormFactor
  * @property {string?} Components.Frame
  * @property {string?} Components.Body
  * @property {string?} Components.Mechanism
  * @property {Array.<string?>} Components.Misc
  *
  * @property {TechLevels} TechLevels
  *
	* @property {number} MaterialQuality
	* @property {number} ExtraDurability
	* @property {number} ExtraQuality
}
 */

/**
 * @typedef {Object} ProductStats
  * @property {number} MiscSlots
  * @property {number} ProductionCost
  * @property {number} MaterialUse

  * @property {number} Quality
  * @property {number} Ecology
  * @property {number} Ethics
  * @property {number} Durability
*/

/**
 * @typedef {Object} Component
 * @property {string} Name
 * @property {number} MiscSlots
 * @property {number} ProductionCost
 * @property {number} MaterialUse
 * @property {number} Ecology
 * @property {number} Ethics
 * @property {number} Quality
 * @property {number} Durability
 * @property {number} ProductionLineCost
 * @property {string} Image
 */

/**
 * @typedef {Object} Offer
 * @property {Product} Product - The product being offered
 * @property {boolean} Outdated
 * @property {number} Price
 * @property {ProductStats} ProductStats
 * @property {number} PromotionQuality
 * @property {Object} Promotion - Promotion goals
 * @property {number} Promotion.Quantity
 * @property {number} Promotion.Quality
 * @property {number} Promotion.Ecology
 * @property {number} Promotion.Ethics
 * @property {number} Promotion.Durability
 * @property {number} Promotion.Price
 */

/**
 * @typedef {Object} Report
 * @property {number} Month - Reporting month
 * @property {Financial_report} FinancialReport - Financial report
 * @property {Balance_sheet} BalanceSheet - Balance sheet
 * @property {Personelle_report} PersonelleReport - Personnel report
 * @property {Production_report} ProductionReport - Production report
 * @property {Object.<string, Sales_report>} SalesReport - Sales report
 */

/**
 * @typedef {Object} FinanceReportEntry
 * @property {string} Name - Entry name
 * @property {number} Group - Group identifier
 * @property {string} Info - Additional information
 * @property {boolean} CashCost - Whether the entry is a cash cost
 * @property {number} Value - Monetary value
 */

/**
 * @typedef {Object} Balance_sheet
 * @property {number} Bank_balance - Current bank balance
 * @property {FinanceReportEntry[]} InvoiceLog - Log of invoices
 * @property {FinanceReportEntry[]} Assets - List of assets
 * @property {FinanceReportEntry[]} Liabilities - List of liabilities
 */

export const production = 0,
	marketing = 1,
	production_personelle = 2,
	marketing_personelle = 3,
	other_personelle = 4,
	facilities = 5,
	logistics = 6,
	materials = 7,
	energy = 8,
	research = 9,
	employee_training = 10,
	loans = 11,
	loan_intrest = 12,
	bridge_loans = 13,
	bridge_loan_intrest = 14,
	taxes = 15,
	sales = 16,
	severance = 17,
	predictions = 18,
	write_off = 19,
	other = 20;

/** @type {Object<string, string>} */
export const financeReportCategories = {
	0: 'production',
	1: 'marketing',
	2: 'production_personelle',
	3: 'marketing_personelle',
	4: 'other_personelle',
	5: 'facilities',
	6: 'logistics',
	7: 'materials',
	8: 'energy',
	9: 'research',
	10: 'employee_training',
	11: 'loans',
	12: 'loan_intrest',
	13: 'bridge_loans',
	14: 'bridge_loan_intrest',
	15: 'taxes',
	16: 'sales',
	17: 'severance',
	18: 'predictions',
	19: 'write_off',
	20: 'other'
};

/**
 * @typedef {Object} Financial_report
 *
 * @property {Object} Income - Income-related values
 * @property {number} Income.GrossSales - Gross sales
 * @property {number} Income.OtherIncome - Other income
 * @property {number} Income.CostOfSales - Cost of sales
 * @property {number} Income.GrossProfit - Gross profit
 *
 * @property {Object} OperatingExpenses - Operating expenses
 * @property {number} OperatingExpenses.Advertising - Advertising expenses
 * @property {number} OperatingExpenses.FacilitiesAndLogistics - Facilities and logistics costs
 * @property {number} OperatingExpenses.ResearchAndDevelopment - Research and development expenses
 *
 * @property {Object} NonOperatingExpenses - Non-operating expenses
 * @property {number} NonOperatingExpenses.WriteOffs - Write-offs
 * @property {number} NonOperatingExpenses.LoanInterest - Loan interest
 * @property {number} NonOperatingExpenses.LoanRepayment - Loan repayment
 * @property {number} NonOperatingExpenses.BridgeLoanIntrest - Bridge loan interest
 * @property {number} NonOperatingExpenses.BridgeLoanRepayment - Bridge loan repayment
 * @property {number} NonOperatingExpenses.Other - Other non-operating expenses
 * @property {number} NonOperatingExpenses.Taxes - Taxes
 *
 * @property {Object} Totals
 * @property {number} Totals.TotalOperatingExpenses - Total operating expenses
 * @property {number} Totals.TotalNonOperatingExpenses - Total non-operating expenses
 * @property {number} Totals.IncomeBeforeTax - Income before tax
 * @property {number} Totals.Cashflow - Cashflow
 * @property {number} Totals.NetIncome - Net income
 */

/**
 * @typedef {Object} Personelle_sub_report
 * @property {number} NumberOfEmployees
 * @property {number} NumberOfHires
 * @property {number} NumberOfDepartures
 *
 * @property {number} AvgPay
 * @property {number} MinimumPay
 * @property {number} MaximumPay
 * @property {number} StandardDevPay
 *
 * @property {number} MinimumSkill
 * @property {number} MaximumSkill
 * @property {number} AvgSkill
 * @property {number} StandardDevSkill
 *
 * @property {number} MinimumMotivation
 * @property {number} MaximumMotivation
 * @property {number} AvgMotivation
 * @property {number} StandardDevMotivation
 *
 * @property {number} MinimumProductivity
 * @property {number} MaximumProductivity
 * @property {number} AvgProductivity
 * @property {number} StandardDevProductivity
 */

/**
 * @typedef {Object} Personelle_report
 * @property {Personelle_sub_report} General
 * @property {Personelle_sub_report} Marketing
 * @property {Personelle_sub_report} Production
 */

/**
 * @typedef {Object} Production_report
 * @property {number} MachinesPurchased
 * @property {number} MachinesSold
 * @property {number} WorkerSurplus
 * @property {number} AvgMachineProductivity
 *
 * @property {Object.<string, ProductSpecificReport>} ProductSpecificReport
 *
 * @property {number} MaterialUsed
 * @property {number} EnergyUsed
 *
 * @property {number} WarehousesBought
 */

/** 
 * @typedef {Object} ProductSpecificReport
 * @property {number} TotalProduction 
 * @property {number} BaseProduction  
 * @property {number} BonusProduction 
 * @property {number} ExcessProduction

 * @property {number} TotalProductsProduced
 * @property {number} BaseProductsProduced 
 * @property {number} BonusProductsProduced
 */

/**
 * @typedef {Object} Product_statistics
 * @property {number} Quality
 * @property {number} Durabilty
 * @property {number} Ethics
 * @property {number} Ecology
 */

/**
 * @typedef {Object} Purchasing_statistics
 * @property {number} ProductNumber
 * @property {number} ProductsSold
 * @property {number} ProductDemand
 *
 * @property {number} AvrDecisionFactor
 * @property {number} AvrPurchasingThreshold
 *
 * @property {number} AvrPurchasingFactors
 * @property {[number, number, number, number, number, number, number, number,]} AvrPurchasingFactors
 */

/**
 * @typedef {Object} Research_statistics
 * @property {number} QualityDevelopmentInvestment
 * @property {number} QualityDevelopmentInvestmentEffectiveness
 * @property {number} DurabilityDevelopmentInvestment
 * @property {number} DurabilityDevelopmentInvestmentEffectiveness
 * @property {number} EcologicalProductionInvestment
 * @property {number} EcologicalProductionInvestmentEffectiveness
 */

/**
 * @typedef {Object} Sales_statistics
 * @property {number} ProductsSold
 * @property {number} ProductDemand
 * @property {number} MarketShare
 *
 * @property {number} AvrDecisionFactor
 * @property {number} AvrPurchasingThreshold
 *
 * @property {number} AvrQualityFactor
 * @property {number} AvrDurabilityFactor
 * @property {number} AvrEcologyFactor
 * @property {number} AvrPriceFactor
 * @property {number} AvrEthicsFactor
 * @property {number} AvrBangForBuckFactor
 */

/**
 * @typedef {Object} Marketing_statistics
 * @property {string} Name
 *
 * @property {number} Quality
 * @property {number} Durabilty

 * @property {number} Ethics
 * @property {number} Ecology
 * @property {number} Price
 * @property {number} BangForBuck
 * @property {number} PromotionQuantity
 * @property {number} PromotionQuality
 * @property {number} ImpressionCount
 */

/**
 * @typedef {Object} Sales_report
 * @property {Sales_statistics} ProductSalesStatistics
 * @property {Marketing_statistics} MarketingStatistics
 */

/**
@typedef {Object} CompanyMarketStatistics 
	@property {number} CompanyID
	@property {number} Step
	@property {string} CEO
	@property {string} Name
	@property {ArrayBuffer} Logo              
	@property {number} EmployeeCount
	@property {number} QuartalyNetIncome
	@property {number} Assets
	@property {number} Value
	@property {number} MarketShare
	@property {number} MonthlySales
  @property {number} ValueOfMonthlySales
	@property {number} TotalSales
	@property {number} TotalValueOfSales

	@property { Object.<string, CompanyMarketStatisticsProduct> } Products 
*/

/**
@typedef {Object} CompanyMarketStatisticsProduct  
	@property {string} ID

	@property {Marketing_statistics} Marketing_statistics
	@property {number} MarketShare
	@property {number} MonthlySales
  @property {number} ValueOfMonthlySales
	@property {number} TotalSales
	@property {number} TotalValueOfSales
*/

/**
 * @typedef {Object} ExternalFactors
 * @property {number} Month
 * @property {number} Inflation
 * @property {number} IntrestRate
 * @property {number} BridgeLoansIntrestRate
 * @property {number} EconomicSituationIndex
 * @property {number} TaxRate

 * @property {number} ProductionMinimumWage
 * @property {number} MarketingMinimumWage

 * @property {Machine[]} MachinesOnOffer
 * @property {number} ExternalStoragePrice
 * @property {number} EnergyPrice
 * @property {number} MaterialPrice

 * @property {number} MachineDepreciationRate
 */

export {};
