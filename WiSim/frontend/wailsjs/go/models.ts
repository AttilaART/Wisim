export namespace int {
	
	export enum int {
	    production = 0,
	    marketing = 1,
	    personelle = 2,
	    logistics = 3,
	    materials = 4,
	    energy = 5,
	    product_development = 6,
	    employee_training = 7,
	    loans = 8,
	    loan_intrest = 9,
	    bridge_loans = 10,
	    bridge_loan_intrest = 11,
	    taxes = 12,
	    sales = 13,
	    severance = 14,
	    predictions = 15,
	    other = 16,
	}

}

export namespace simulation {
	
	export class Decisions {
	    Sales_projection: number;
	    Selling_price: number;
	    Marketing: number;
	    Quality_development_investment: number;
	    Ecological_production_investment: number;
	    Durability_development_investment: number;
	    Production_target: number;
	    Purchase_of_machines: number;
	    Selling_of_machines: number;
	    Material_quality: number;
	    Percentage_of_ecological_energy: number;
	    Purchase_of_warehouses: number;
	    New_hires_in_production: number;
	    New_hires_in_marketing: number;
	    Base_pay_for_production: number;
	    Base_pay_for_marketing: number;
	    Raise_for_production_personelle: number;
	    Raise_for_marketing_personelle: number;
	    Severance_for_production_personelle: number;
	    Severance_for_marketing_personelle: number;
	    Working_hours_for_production: number;
	    Working_hours_for_marketing: number;
	    Investment_in_production_training: number;
	    Investment_in_marketing_training: number;
	    Increase_of_loans: number;
	    Dividends: number;
	
	    static createFrom(source: any = {}) {
	        return new Decisions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Sales_projection = source["Sales_projection"];
	        this.Selling_price = source["Selling_price"];
	        this.Marketing = source["Marketing"];
	        this.Quality_development_investment = source["Quality_development_investment"];
	        this.Ecological_production_investment = source["Ecological_production_investment"];
	        this.Durability_development_investment = source["Durability_development_investment"];
	        this.Production_target = source["Production_target"];
	        this.Purchase_of_machines = source["Purchase_of_machines"];
	        this.Selling_of_machines = source["Selling_of_machines"];
	        this.Material_quality = source["Material_quality"];
	        this.Percentage_of_ecological_energy = source["Percentage_of_ecological_energy"];
	        this.Purchase_of_warehouses = source["Purchase_of_warehouses"];
	        this.New_hires_in_production = source["New_hires_in_production"];
	        this.New_hires_in_marketing = source["New_hires_in_marketing"];
	        this.Base_pay_for_production = source["Base_pay_for_production"];
	        this.Base_pay_for_marketing = source["Base_pay_for_marketing"];
	        this.Raise_for_production_personelle = source["Raise_for_production_personelle"];
	        this.Raise_for_marketing_personelle = source["Raise_for_marketing_personelle"];
	        this.Severance_for_production_personelle = source["Severance_for_production_personelle"];
	        this.Severance_for_marketing_personelle = source["Severance_for_marketing_personelle"];
	        this.Working_hours_for_production = source["Working_hours_for_production"];
	        this.Working_hours_for_marketing = source["Working_hours_for_marketing"];
	        this.Investment_in_production_training = source["Investment_in_production_training"];
	        this.Investment_in_marketing_training = source["Investment_in_marketing_training"];
	        this.Increase_of_loans = source["Increase_of_loans"];
	        this.Dividends = source["Dividends"];
	    }
	}
	export class FinanceReportEntry {
	    Name: string;
	    Group: number;
	    Info: string;
	    Cash_cost: boolean;
	    Value: number;
	
	    static createFrom(source: any = {}) {
	        return new FinanceReportEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Group = source["Group"];
	        this.Info = source["Info"];
	        this.Cash_cost = source["Cash_cost"];
	        this.Value = source["Value"];
	    }
	}
	export class Financial_Report {
	    Total_income: number;
	    Loan_repayments: number;
	    Bridge_loan_repayments: number;
	    New_bridge_loan: number;
	    Non_cash_costs: number;
	    Cash_costs: number;
	    Total_expenses_before_tax: number;
	    Total_expenses_after_tax: number;
	    Operating_profit: number;
	    Taxes: number;
	    Net_Profit: number;
	
	    static createFrom(source: any = {}) {
	        return new Financial_Report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Total_income = source["Total_income"];
	        this.Loan_repayments = source["Loan_repayments"];
	        this.Bridge_loan_repayments = source["Bridge_loan_repayments"];
	        this.New_bridge_loan = source["New_bridge_loan"];
	        this.Non_cash_costs = source["Non_cash_costs"];
	        this.Cash_costs = source["Cash_costs"];
	        this.Total_expenses_before_tax = source["Total_expenses_before_tax"];
	        this.Total_expenses_after_tax = source["Total_expenses_after_tax"];
	        this.Operating_profit = source["Operating_profit"];
	        this.Taxes = source["Taxes"];
	        this.Net_Profit = source["Net_Profit"];
	    }
	}
	export class Product_statistics {
	    Quality: number;
	    Durabilty: number;
	    Coolness: number;
	    Ecology: number;
	
	    static createFrom(source: any = {}) {
	        return new Product_statistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Quality = source["Quality"];
	        this.Durabilty = source["Durabilty"];
	        this.Coolness = source["Coolness"];
	        this.Ecology = source["Ecology"];
	    }
	}
	export class Marketing_statistics {
	    Product: Product_statistics;
	    Price: number;
	    Bang_for_buck: number;
	    Promotion: number;
	
	    static createFrom(source: any = {}) {
	        return new Marketing_statistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Product = this.convertValues(source["Product"], Product_statistics);
	        this.Price = source["Price"];
	        this.Bang_for_buck = source["Bang_for_buck"];
	        this.Promotion = source["Promotion"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Personelle_sub_report {
	    Number_of_employees: number;
	    Number_of_hires: number;
	    Avg_pay: number;
	    Minimum_pay: number;
	    Maximum_pay: number;
	    Standard_dev_pay: number;
	    Minimum_skill: number;
	    Maximum_skill: number;
	    Avg_skill: number;
	    Standard_dev_skill: number;
	
	    static createFrom(source: any = {}) {
	        return new Personelle_sub_report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Number_of_employees = source["Number_of_employees"];
	        this.Number_of_hires = source["Number_of_hires"];
	        this.Avg_pay = source["Avg_pay"];
	        this.Minimum_pay = source["Minimum_pay"];
	        this.Maximum_pay = source["Maximum_pay"];
	        this.Standard_dev_pay = source["Standard_dev_pay"];
	        this.Minimum_skill = source["Minimum_skill"];
	        this.Maximum_skill = source["Maximum_skill"];
	        this.Avg_skill = source["Avg_skill"];
	        this.Standard_dev_skill = source["Standard_dev_skill"];
	    }
	}
	
	export class Production_report {
	    Machines_purchased: number;
	    Machines_sold: number;
	    Worker_surplus: number;
	    Avg_machine_productivity: number;
	    Products_produced: number;
	    Base_production: number;
	    Bonus_production: number;
	    Material_used: number;
	    Energy_used: number;
	    Warehouses_bought: number;
	
	    static createFrom(source: any = {}) {
	        return new Production_report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Machines_purchased = source["Machines_purchased"];
	        this.Machines_sold = source["Machines_sold"];
	        this.Worker_surplus = source["Worker_surplus"];
	        this.Avg_machine_productivity = source["Avg_machine_productivity"];
	        this.Products_produced = source["Products_produced"];
	        this.Base_production = source["Base_production"];
	        this.Bonus_production = source["Bonus_production"];
	        this.Material_used = source["Material_used"];
	        this.Energy_used = source["Energy_used"];
	        this.Warehouses_bought = source["Warehouses_bought"];
	    }
	}
	export class Sales_statistics {
	    Products_sold: number;
	    Difference_to_previous_month: number;
	    Product_demand: number;
	    Market_share: number;
	    Avr_decision_factor: number;
	    Avr_purchasing_threshold: number;
	    Avr_quality_factor: number;
	    Avr_durability_factor: number;
	    Avr_ecology_factor: number;
	    Avr_price_factor: number;
	    Avr_coolness_factor: number;
	    Avr_bang_for_buck_factor: number;
	
	    static createFrom(source: any = {}) {
	        return new Sales_statistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Products_sold = source["Products_sold"];
	        this.Difference_to_previous_month = source["Difference_to_previous_month"];
	        this.Product_demand = source["Product_demand"];
	        this.Market_share = source["Market_share"];
	        this.Avr_decision_factor = source["Avr_decision_factor"];
	        this.Avr_purchasing_threshold = source["Avr_purchasing_threshold"];
	        this.Avr_quality_factor = source["Avr_quality_factor"];
	        this.Avr_durability_factor = source["Avr_durability_factor"];
	        this.Avr_ecology_factor = source["Avr_ecology_factor"];
	        this.Avr_price_factor = source["Avr_price_factor"];
	        this.Avr_coolness_factor = source["Avr_coolness_factor"];
	        this.Avr_bang_for_buck_factor = source["Avr_bang_for_buck_factor"];
	    }
	}

}

