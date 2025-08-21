// NOTE: Most of the types defined in this file were generated from Go types through ChatGPT

/**
 * @typedef {{company: Company, decisions: Decisions, external_factors: External_factors}} clientState
 */

/**
 * @typedef {Object} Decisions
 * @property {{Sales_prediction: number}} Predictions
 * @property {{Set_bank_loan: number}} Finances
 * @property {Decisions_marketing} Marketing
 * @property {{Production_deltas: Delta<Employee>[], Marketing_deltas: Delta<Employee>[], Severance_pay: number}} Employees
 * @property {{Production_goal: number, Machines: Delta<Machine>[], Logistics: Delta<Warehouse>[]}} Production
 * @property {Decisions_research} Research
 */

/**
 * @typedef {Object} Decisions_marketing
 * @property {number} Price - Price setting
 * @property {Decisions_product | undefined} Product - Product-related decisions
 * @property {Object} Promotion - Promotion decisions
 * @property {number} Promotion.Quantity - Promotion quantity
 * @property {number} Promotion.Style_quality - Quality style emphasis
 * @property {number} Promotion.Style_ecology - Ecology style emphasis
 * @property {number} Promotion.Style_ethics - Ethics style emphasis
 * @property {number} Promotion.Style_durability - Durability style emphasis
 */

/**
 * @typedef {Object} Decisions_product
 * @property {Object} Materials - Material-related decisions
 * @property {number} Materials.quality - Material quality
 * @property {number} Materials.ecology - Ecological impact of materials
 * @property {number} Materials.ethical_sourcing - Ethical sourcing emphasis
 *
 * @property {Object} Manufacturing - Manufacturing-related decisions
 * @property {number} Manufacturing.Quality - Manufacturing quality
 * @property {number} Manufacturing.Ecological_energy - Use of ecological energy
 * @property {number} Manufacturing.Material_efficiency - Efficiency of material usage
 * @property {number} Manufacturing.Durability - Product durability
 * @property {number} Manufacturing.Max_durability - Maximum possible durability
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
 * @property {number} Id - Unique identifier of the employee
 * @property {string} Name - Employee's name
 * @property {number} Employer - Reference to the employer ID
 * @property {Employee_type} Employee_type - Type of the employee
 * @property {number} Motivation - Motivation level
 * @property {number} Skill - Skill level
 * @property {number} Extra_training - Extra training factor
 * @property {number} Pay - Base pay
 * @property {number} Bonus - Bonus amount
 * @property {number} Working_hours - Number of working hours
 */

/**
 * @typedef {number} Employee_type
 */

/**
 * @typedef {Object} Machine
 * @property {number} Id - Unique identifier of the machine
 * @property {number} Production_capacity - Maximum production capacity
 * @property {number} Required_workers - Number of workers required
 * @property {number} Minimum_workers - Minimum number of workers allowed
 * @property {Employee[]} Assigned_workers_ptr - Workers assigned to this machine
 * @property {number} Energy_use - Energy consumption
 * @property {number} Value - Machine value
 */

/**
 * @typedef {Object} Warehouse
 * @property {number} Id - Unique identifier of the warehouse
 * @property {number} Capacity - Storage capacity
 * @property {number} Operating_costs - Operating costs of the warehouse
 * @property {number} Value - Monetary value of the warehouse
 */

/**
 * @typedef {Object} Company
 *
 * @property {number} Id - Unique identifier of the company
 * @property {string} Name - Name of the company
 * @property {number} Balance - Current financial balance
 * @property {number} Loans - Outstanding loans
 * @property {number} Bridge_loans - Outstanding bridge loans
 * @property {Decisions[]} Decision_history - History of company decisions
 * @property {Report[]} Reports - Reports related to the company
 *
 * @property {number} Global_quality_factor - Global quality factor from R&D
 * @property {number} Base_marketing_strength - Base marketing strength from R&D
 *
 * @property {Offer | undefined} Offer - Product offer
 * @property {number} Orders - Number of orders
 *
 * @property {Warehouse[]} Warehouses - Warehouses owned by the company
 * @property {number} Items_in_storage - Number of items in storage
 *
 * @property {Machine[]} Machines - Machines owned by the company
 */

/**
 * @typedef {Object} Product
 * @property {number} Id
 * @property {string} Name
 * @property {number} Weight
 * @property {number} Material_use
 * @property {number} Production_cost
 *
 * @property {number} Base_material_use
 * @property {number} Base_production_cost
 * @property {number} Base_quality
 * @property {number} Base_ecology
 * @property {number} Base_durability
 *
 * @property {number} Ethics_factor
 * @property {number} Quality_factor
 * @property {number} Ecology_factor
 * @property {number} Durabilty
 */

/**
 * @typedef {Object} Offer
 * @property {Product} Product - The product being offered
 * @property {number} Price
 * @property {number} Promotion_quality
 * @property {Object} Promotion_goal - Promotion goals
 * @property {number} Promotion_goal.Quantity
 * @property {number} Promotion_goal.Style_quality
 * @property {number} Promotion_goal.Style_ecology
 * @property {number} Promotion_goal.Style_ethics
 * @property {number} Promotion_goal.Style_durability
 */

/**
 * @typedef {Object} Report
 * @property {number} Month - Reporting month
 * @property {Financial_report} Financial_Report - Financial report
 * @property {Balance_sheet} Balance_sheet - Balance sheet
 * @property {Personelle_report} Personelle_report - Personnel report
 * @property {Production_report} Production_report - Production report
 * @property {Sales_report} Sales_report - Sales report
 */

/**
 * @typedef {Object} FinanceReportEntry
 * @property {string} Name - Entry name
 * @property {number} Group - Group identifier
 * @property {string} Info - Additional information
 * @property {boolean} Cash_cost - Whether the entry is a cash cost
 * @property {number} Value - Monetary value
 */

/**
 * @typedef {Object} Balance_sheet
 * @property {number} Bank_balance - Current bank balance
 * @property {FinanceReportEntry[]} Invoice_log - Log of invoices
 * @property {FinanceReportEntry[]} Assets - List of assets
 * @property {FinanceReportEntry[]} Liabilities - List of liabilities
 */

/**
 * @typedef {Object} Financial_report
 *
 * @property {Object} Income - Income-related values
 * @property {number} Income.Gross_sales - Gross sales
 * @property {number} Income.Other_income - Other income
 * @property {number} Income.Cost_of_sales - Cost of sales
 * @property {number} Income.Gross_profit - Gross profit
 *
 * @property {Object} Operating_expenses - Operating expenses
 * @property {number} Operating_expenses.Advertising - Advertising expenses
 * @property {number} Operating_expenses.Facilities_and_logistics - Facilities and logistics costs
 * @property {number} Operating_expenses.Research_and_development - Research and development expenses
 * @property {number} Operating_expenses.Total_operating_expenses - Total operating expenses
 *
 * @property {Object} Non_operating_expenses - Non-operating expenses
 * @property {number} Non_operating_expenses.Write_offs - Write-offs
 * @property {number} Non_operating_expenses.Loan_interest - Loan interest
 * @property {number} Non_operating_expenses.Loan_repayment - Loan repayment
 * @property {number} Non_operating_expenses.Bridge_loan_intrest - Bridge loan interest
 * @property {number} Non_operating_expenses.Bridge_loan_repayment - Bridge loan repayment
 * @property {number} Non_operating_expenses.Other - Other non-operating expenses
 * @property {number} Non_operating_expenses.Total_non_operating_expenses - Total non-operating expenses
 * @property {number} Non_operating_expenses.Income_before_tax - Income before tax
 * @property {number} Non_operating_expenses.Taxes - Taxes
 * @property {number} Non_operating_expenses.Net_income - Net income
 * @property {number} Non_operating_expenses.Cashflow - Cashflow
 */

/**
 * @typedef {Object} Personelle_sub_report
 * @property {number} Number_of_employees
 * @property {number} Number_of_hires
 * @property {number} Number_of_departures
 *
 * @property {number} Avg_pay
 * @property {number} Minimum_pay
 * @property {number} Maximum_pay
 * @property {number} Standard_dev_pay
 *
 * @property {number} Minimum_skill
 * @property {number} Maximum_skill
 * @property {number} Avg_skill
 * @property {number} Standard_dev_skill
 *
 * @property {number} Minimum_motivation
 * @property {number} Maximum_motivation
 * @property {number} Avg_motivation
 * @property {number} Standard_dev_motivation
 *
 * @property {number} Minimum_productivity
 * @property {number} Maximum_productivity
 * @property {number} Avg_productivity
 * @property {number} Standard_dev_productivity
 */

/**
 * @typedef {Object} Personelle_report
 * @property {Personelle_sub_report} General
 * @property {Personelle_sub_report} Marketing
 * @property {Personelle_sub_report} Production
 */

/**
 * @typedef {Object} Production_report
 * @property {number} Machines_purchased
 * @property {number} Machines_sold
 * @property {number} Worker_surplus
 * @property {number} Avg_machine_productivity
 * @property {number} Products_produced
 * @property {number} Base_production
 * @property {number} Bonus_production
 * @property {number} Material_used
 * @property {number} Energy_used
 * @property {number} Warehouses_bought
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
 * @property {number} Product_number
 * @property {number} Products_sold
 * @property {number} Difference_to_previous_month
 * @property {number} Product_demand
 *
 * @property {number} Avr_decision_factor
 * @property {number} Avr_purchasing_threshold
 *
 * @property {number} Avr_quality_factor
 * @property {number} Avr_durability_factor
 * @property {number} Avr_ecology_factor
 * @property {number} Avr_price_factor
 * @property {number} Avr_ethics_factor
 * @property {number} Avr_bang_for_buck_factor
 */

/**
 * @typedef {Object} Research_statistics
 * @property {number} Quality_development_investment
 * @property {number} Quality_development_investment_effectiveness
 * @property {number} Durability_development_investment
 * @property {number} Durability_development_investment_effectiveness
 * @property {number} Ecological_production_investment
 * @property {number} Ecological_production_investment_effectiveness
 */

/**
 * @typedef {Object} Sales_statistics
 * @property {number} Products_sold
 * @property {number} Difference_to_previous_month
 * @property {number} Product_demand
 * @property {number} Market_share
 *
 * @property {number} Avr_decision_factor
 * @property {number} Avr_purchasing_threshold
 *
 * @property {number} Avr_quality_factor
 * @property {number} Avr_durability_factor
 * @property {number} Avr_ecology_factor
 * @property {number} Avr_price_factor
 * @property {number} Avr_ethics_factor
 * @property {number} Avr_bang_for_buck_factor
 */

/**
 * @typedef {Object} Marketing_statistics
 * @property {Product_statistics} Product
 * @property {number} Price
 * @property {number} Bang_for_buck
 * @property {number} Promotion_quantity
 * @property {number} Promotion_quality
 */

/**
 * @typedef {Object} Sales_report
 * @property {Product_statistics} Product_statistics
 * @property {Sales_statistics} Company_sales_statistics
 * @property {Marketing_statistics} Marketing_statistics
 */

/**
 * @typedef {Object} External_factors
 *
 * @property {number} Month
 *
 * // Economy
 * @property {number} Inflation
 * @property {number} Intrest_rate
 * @property {number} Bridge_loans_intrest_rate
 * @property {number} Economic_situation_index
 * @property {number} Tax_rate - As decimal
 *
 * // Personelle
 * @property {number} Turnover
 * @property {number} Production_minimum_wage
 * @property {number} Marketing_minimum_wage
 *
 * // Production
 * @property {Machine | undefined} Machine_on_offer
 * @property {number} External_storage_price - Per item
 * @property {number} Energy_price - Per unit of energy
 * @property {number} Material_price - Per unit of material
 *
 * @property {number} Machine_depreciation_rate - In decimal
 */

export {};
