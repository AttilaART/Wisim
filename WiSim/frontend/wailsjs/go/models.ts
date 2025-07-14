export namespace int {
	
	export enum int {
	    production = 0,
	    marketing = 1,
	    prodcution_personelle = 2,
	    marketing_personelle = 3,
	    other_personelle = 4,
	    facilities = 5,
	    logistics = 6,
	    materials = 7,
	    energy = 8,
	    product_development = 9,
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
	    other = 20,
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
	export class Balance_sheet {
	    Bank_balance: number;
	    Invoice_log: FinanceReportEntry[];
	    Assets: FinanceReportEntry[];
	    Liabilities: FinanceReportEntry[];
	
	    static createFrom(source: any = {}) {
	        return new Balance_sheet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Bank_balance = source["Bank_balance"];
	        this.Invoice_log = this.convertValues(source["Invoice_log"], FinanceReportEntry);
	        this.Assets = this.convertValues(source["Assets"], FinanceReportEntry);
	        this.Liabilities = this.convertValues(source["Liabilities"], FinanceReportEntry);
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
	export class Employee {
	    Id: number;
	    Name: string;
	    Employer: number;
	    Employee_type: number;
	    Motivation: number;
	    Skill: number;
	    Extra_training: number;
	    Pay: number;
	    Bonus: number;
	    Working_hours: number;
	
	    static createFrom(source: any = {}) {
	        return new Employee(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Employer = source["Employer"];
	        this.Employee_type = source["Employee_type"];
	        this.Motivation = source["Motivation"];
	        this.Skill = source["Skill"];
	        this.Extra_training = source["Extra_training"];
	        this.Pay = source["Pay"];
	        this.Bonus = source["Bonus"];
	        this.Working_hours = source["Working_hours"];
	    }
	}
	export class Machine {
	    Id: number;
	    Production_capacity: number;
	    Required_workers: number;
	    Minimum_workers: number;
	    Assigned_workers_ptr: Employee[];
	    Energy_use: number;
	    Value: number;
	
	    static createFrom(source: any = {}) {
	        return new Machine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Production_capacity = source["Production_capacity"];
	        this.Required_workers = source["Required_workers"];
	        this.Minimum_workers = source["Minimum_workers"];
	        this.Assigned_workers_ptr = this.convertValues(source["Assigned_workers_ptr"], Employee);
	        this.Energy_use = source["Energy_use"];
	        this.Value = source["Value"];
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
	export class Warehouse {
	    Id: number;
	    Capacity: number;
	    Operating_costs: number;
	    Value: number;
	
	    static createFrom(source: any = {}) {
	        return new Warehouse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Capacity = source["Capacity"];
	        this.Operating_costs = source["Operating_costs"];
	        this.Value = source["Value"];
	    }
	}
	export class Product {
	    Id: number;
	    Name: string;
	    Weight: number;
	    Material_use: number;
	    Production_cost: number;
	    Base_material_use: number;
	    Base_production_cost: number;
	    Base_quality: number;
	    Base_ecology: number;
	    Base_durability: number;
	    Ethics_factor: number;
	    Quality_factor: number;
	    Ecology_factor: number;
	    Durabilty: number;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Weight = source["Weight"];
	        this.Material_use = source["Material_use"];
	        this.Production_cost = source["Production_cost"];
	        this.Base_material_use = source["Base_material_use"];
	        this.Base_production_cost = source["Base_production_cost"];
	        this.Base_quality = source["Base_quality"];
	        this.Base_ecology = source["Base_ecology"];
	        this.Base_durability = source["Base_durability"];
	        this.Ethics_factor = source["Ethics_factor"];
	        this.Quality_factor = source["Quality_factor"];
	        this.Ecology_factor = source["Ecology_factor"];
	        this.Durabilty = source["Durabilty"];
	    }
	}
	export class Offer {
	    Product: Product;
	    Price: number;
	    Promotion_quality: number;
	    // Go type: struct { Quantity float64; Style_quality float32; Style_ecology float32; Style_ethics float32; Style_durability float32 }
	    Promotion_goal: any;
	
	    static createFrom(source: any = {}) {
	        return new Offer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Product = this.convertValues(source["Product"], Product);
	        this.Price = source["Price"];
	        this.Promotion_quality = source["Promotion_quality"];
	        this.Promotion_goal = this.convertValues(source["Promotion_goal"], Object);
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
	export class Marketing_statistics {
	    Product: Product_statistics;
	    Price: number;
	    Bang_for_buck: number;
	    Promotion_quantity: number;
	    Promotion_quality: number;
	
	    static createFrom(source: any = {}) {
	        return new Marketing_statistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Product = this.convertValues(source["Product"], Product_statistics);
	        this.Price = source["Price"];
	        this.Bang_for_buck = source["Bang_for_buck"];
	        this.Promotion_quantity = source["Promotion_quantity"];
	        this.Promotion_quality = source["Promotion_quality"];
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
	    Avr_ethics_factor: number;
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
	        this.Avr_ethics_factor = source["Avr_ethics_factor"];
	        this.Avr_bang_for_buck_factor = source["Avr_bang_for_buck_factor"];
	    }
	}
	export class Product_statistics {
	    Quality: number;
	    Durabilty: number;
	    Ethics: number;
	    Ecology: number;
	
	    static createFrom(source: any = {}) {
	        return new Product_statistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Quality = source["Quality"];
	        this.Durabilty = source["Durabilty"];
	        this.Ethics = source["Ethics"];
	        this.Ecology = source["Ecology"];
	    }
	}
	export class Sales_report {
	    Product_statistics: Product_statistics;
	    Company_sales_statistics: Sales_statistics;
	    Marketing_statistics: Marketing_statistics;
	
	    static createFrom(source: any = {}) {
	        return new Sales_report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Product_statistics = this.convertValues(source["Product_statistics"], Product_statistics);
	        this.Company_sales_statistics = this.convertValues(source["Company_sales_statistics"], Sales_statistics);
	        this.Marketing_statistics = this.convertValues(source["Marketing_statistics"], Marketing_statistics);
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
	export class Personelle_report {
	    General: Personelle_sub_report;
	    Marketing: Personelle_sub_report;
	    Production: Personelle_sub_report;
	
	    static createFrom(source: any = {}) {
	        return new Personelle_report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.General = this.convertValues(source["General"], Personelle_sub_report);
	        this.Marketing = this.convertValues(source["Marketing"], Personelle_sub_report);
	        this.Production = this.convertValues(source["Production"], Personelle_sub_report);
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
	export class Financial_Report {
	    // Go type: struct { Gross_sales float64; Other_income float64; Cost_of_sales float64; Gross_profit float64 }
	    Income: any;
	    // Go type: struct { Advertising float64; Facilities_and_logistics float64; Research_and_development float64; Total_operating_expenses float64 }
	    Operating_expenses: any;
	    // Go type: struct { Write_offs float64; Loan_interest float64; Loan_repayment float64; Bridge_loan_intrest float64; Bridge_loan_repayment float64; Other float64; Total_non_operating_expenses float64; Income_before_tax float64; Taxes float64; Net_income float64; Cashflow float64 }
	    Non_operating_expenses: any;
	
	    static createFrom(source: any = {}) {
	        return new Financial_Report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Income = this.convertValues(source["Income"], Object);
	        this.Operating_expenses = this.convertValues(source["Operating_expenses"], Object);
	        this.Non_operating_expenses = this.convertValues(source["Non_operating_expenses"], Object);
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
	export class Report {
	    Month: number;
	    Financial_Report: Financial_Report;
	    Balance_sheet: Balance_sheet;
	    Personelle_report: Personelle_report;
	    Production_report: Production_report;
	    Sales_report: Sales_report;
	
	    static createFrom(source: any = {}) {
	        return new Report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Month = source["Month"];
	        this.Financial_Report = this.convertValues(source["Financial_Report"], Financial_Report);
	        this.Balance_sheet = this.convertValues(source["Balance_sheet"], Balance_sheet);
	        this.Personelle_report = this.convertValues(source["Personelle_report"], Personelle_report);
	        this.Production_report = this.convertValues(source["Production_report"], Production_report);
	        this.Sales_report = this.convertValues(source["Sales_report"], Sales_report);
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
	export class Decisions_research {
	    Quality: number;
	    Durability: number;
	    Ecology: number;
	    Promotion: number;
	    Production_cost: number;
	
	    static createFrom(source: any = {}) {
	        return new Decisions_research(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Quality = source["Quality"];
	        this.Durability = source["Durability"];
	        this.Ecology = source["Ecology"];
	        this.Promotion = source["Promotion"];
	        this.Production_cost = source["Production_cost"];
	    }
	}
	export class Decisions_product {
	    // Go type: struct { Quality float32; Ecology float32; Ethical_sourcing float32 }
	    Materials: any;
	    // Go type: struct { Quality float32; Ecological_energy float32; Material_efficiency float32; Durability float32; Max_durability int }
	    Manufacturing: any;
	
	    static createFrom(source: any = {}) {
	        return new Decisions_product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Materials = this.convertValues(source["Materials"], Object);
	        this.Manufacturing = this.convertValues(source["Manufacturing"], Object);
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
	export class Decisions_marketing {
	    Price: number;
	    Product: Decisions_product;
	    // Go type: struct { Quantity float64; Style_quality float32; Style_ecology float32; Style_ethics float32; Style_durability float32 }
	    Promotion: any;
	
	    static createFrom(source: any = {}) {
	        return new Decisions_marketing(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Price = source["Price"];
	        this.Product = this.convertValues(source["Product"], Decisions_product);
	        this.Promotion = this.convertValues(source["Promotion"], Object);
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
	export class Decisions {
	    // Go type: struct { Sales_prediction int }
	    Predictions: any;
	    // Go type: struct { Set_bank_loan float64 }
	    Finances: any;
	    Marketing: Decisions_marketing;
	    // Go type: struct { Production_deltas []simulation
	    Employees: any;
	    // Go type: struct { Production_goal int; Machines []simulation
	    Production: any;
	    Research: Decisions_research;
	
	    static createFrom(source: any = {}) {
	        return new Decisions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Predictions = this.convertValues(source["Predictions"], Object);
	        this.Finances = this.convertValues(source["Finances"], Object);
	        this.Marketing = this.convertValues(source["Marketing"], Decisions_marketing);
	        this.Employees = this.convertValues(source["Employees"], Object);
	        this.Production = this.convertValues(source["Production"], Object);
	        this.Research = this.convertValues(source["Research"], Decisions_research);
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
	export class Company {
	    Id: number;
	    Name: string;
	    Balance: number;
	    Loans: number;
	    Bridge_loans: number;
	    Decision_history: Decisions[];
	    Reports: Report[];
	    Global_quality_factor: number;
	    Base_marketing_strength: number;
	    Offer: Offer;
	    Orders: number;
	    Warehouses: Warehouse[];
	    Items_in_storage: number;
	    Machines: Machine[];
	
	    static createFrom(source: any = {}) {
	        return new Company(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Balance = source["Balance"];
	        this.Loans = source["Loans"];
	        this.Bridge_loans = source["Bridge_loans"];
	        this.Decision_history = this.convertValues(source["Decision_history"], Decisions);
	        this.Reports = this.convertValues(source["Reports"], Report);
	        this.Global_quality_factor = source["Global_quality_factor"];
	        this.Base_marketing_strength = source["Base_marketing_strength"];
	        this.Offer = this.convertValues(source["Offer"], Offer);
	        this.Orders = source["Orders"];
	        this.Warehouses = this.convertValues(source["Warehouses"], Warehouse);
	        this.Items_in_storage = source["Items_in_storage"];
	        this.Machines = this.convertValues(source["Machines"], Machine);
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
	
	
	
	
	export class Delta_WiSim_simulation_Employee_ {
	    Change: number;
	    Item: Employee;
	
	    static createFrom(source: any = {}) {
	        return new Delta_WiSim_simulation_Employee_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Change = source["Change"];
	        this.Item = this.convertValues(source["Item"], Employee);
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
	export class Delta_WiSim_simulation_Machine_ {
	    Change: number;
	    Item: Machine;
	
	    static createFrom(source: any = {}) {
	        return new Delta_WiSim_simulation_Machine_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Change = source["Change"];
	        this.Item = this.convertValues(source["Item"], Machine);
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
	export class Delta_WiSim_simulation_Warehouse_ {
	    Change: number;
	    Item: Warehouse;
	
	    static createFrom(source: any = {}) {
	        return new Delta_WiSim_simulation_Warehouse_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Change = source["Change"];
	        this.Item = this.convertValues(source["Item"], Warehouse);
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
	
	export class External_factors {
	    Month: number;
	    Inflation: number;
	    Intrest_rate: number;
	    Bridge_loans_intrest_rate: number;
	    Economic_situation_index: number;
	    Tax_rate: number;
	    Turnover: number;
	    Production_minimum_wage: number;
	    Marketing_minimum_wage: number;
	    Machine_on_offer: Machine;
	    External_storage_price: number;
	    Energy_price: number;
	    Material_price: number;
	    Machine_depreciation_rate: number;
	
	    static createFrom(source: any = {}) {
	        return new External_factors(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Month = source["Month"];
	        this.Inflation = source["Inflation"];
	        this.Intrest_rate = source["Intrest_rate"];
	        this.Bridge_loans_intrest_rate = source["Bridge_loans_intrest_rate"];
	        this.Economic_situation_index = source["Economic_situation_index"];
	        this.Tax_rate = source["Tax_rate"];
	        this.Turnover = source["Turnover"];
	        this.Production_minimum_wage = source["Production_minimum_wage"];
	        this.Marketing_minimum_wage = source["Marketing_minimum_wage"];
	        this.Machine_on_offer = this.convertValues(source["Machine_on_offer"], Machine);
	        this.External_storage_price = source["External_storage_price"];
	        this.Energy_price = source["Energy_price"];
	        this.Material_price = source["Material_price"];
	        this.Machine_depreciation_rate = source["Machine_depreciation_rate"];
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
	
	
	
	
	
	
	
	
	
	
	
	
	

}

