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
	export class Sales_report {
	    Products_sold: number;
	    Difference_to_previous_month: number;
	
	    static createFrom(source: any = {}) {
	        return new Sales_report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Products_sold = source["Products_sold"];
	        this.Difference_to_previous_month = source["Difference_to_previous_month"];
	    }
	}

}

