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
	    EBIT: number;
	    Taxes: number;
	    Profit: number;
	
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
	        this.EBIT = source["EBIT"];
	        this.Taxes = source["Taxes"];
	        this.Profit = source["Profit"];
	    }
	}
	export class Personelle_report {
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
	    Marketing_Number_of_employees: number;
	    Marketing_Number_of_hires: number;
	    Marketing_Avg_pay: number;
	    Marketing_Minimum_pay: number;
	    Marketing_Maximum_pay: number;
	    Marketing_Standard_dev_pay: number;
	    Marketing_Minimum_skill: number;
	    Marketing_Maximum_skill: number;
	    Marketing_Avg_skill: number;
	    Marketing_Standard_dev_skill: number;
	    Production_Number_of_employees: number;
	    Production_Number_of_hires: number;
	    Production_Avg_pay: number;
	    Production_Minimum_pay: number;
	    Production_Maximum_pay: number;
	    Production_Standard_dev_pay: number;
	    Production_Minimum_skill: number;
	    Production_Maximum_skill: number;
	    Production_Avg_skill: number;
	    Production_Standard_dev_skill: number;
	    Production_avg_productivity: number;
	    Production_minimum_productivity: number;
	    Prouduction_maximum_productivity: number;
	
	    static createFrom(source: any = {}) {
	        return new Personelle_report(source);
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
	        this.Marketing_Number_of_employees = source["Marketing_Number_of_employees"];
	        this.Marketing_Number_of_hires = source["Marketing_Number_of_hires"];
	        this.Marketing_Avg_pay = source["Marketing_Avg_pay"];
	        this.Marketing_Minimum_pay = source["Marketing_Minimum_pay"];
	        this.Marketing_Maximum_pay = source["Marketing_Maximum_pay"];
	        this.Marketing_Standard_dev_pay = source["Marketing_Standard_dev_pay"];
	        this.Marketing_Minimum_skill = source["Marketing_Minimum_skill"];
	        this.Marketing_Maximum_skill = source["Marketing_Maximum_skill"];
	        this.Marketing_Avg_skill = source["Marketing_Avg_skill"];
	        this.Marketing_Standard_dev_skill = source["Marketing_Standard_dev_skill"];
	        this.Production_Number_of_employees = source["Production_Number_of_employees"];
	        this.Production_Number_of_hires = source["Production_Number_of_hires"];
	        this.Production_Avg_pay = source["Production_Avg_pay"];
	        this.Production_Minimum_pay = source["Production_Minimum_pay"];
	        this.Production_Maximum_pay = source["Production_Maximum_pay"];
	        this.Production_Standard_dev_pay = source["Production_Standard_dev_pay"];
	        this.Production_Minimum_skill = source["Production_Minimum_skill"];
	        this.Production_Maximum_skill = source["Production_Maximum_skill"];
	        this.Production_Avg_skill = source["Production_Avg_skill"];
	        this.Production_Standard_dev_skill = source["Production_Standard_dev_skill"];
	        this.Production_avg_productivity = source["Production_avg_productivity"];
	        this.Production_minimum_productivity = source["Production_minimum_productivity"];
	        this.Prouduction_maximum_productivity = source["Prouduction_maximum_productivity"];
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

