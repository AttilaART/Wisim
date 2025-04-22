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

}

