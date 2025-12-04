package simulation

import (
	"maps"
	"math"
	"math/rand"
	"slices"

	"github.com/pehringer/simd"
)

type Interval struct {
	Start       int
	Stop_before int
}

func split_load(thread_count int, array_len int) []Interval {
	thread_people_range := make([]Interval, thread_count)

	count_per_thread := array_len / thread_count
	remainder := array_len % thread_count
	offset := 0

	for i := range thread_people_range {
		thread_people_range[i].Start = offset
		thread_people_range[i].Stop_before = offset + count_per_thread
		offset += count_per_thread

		if remainder > 0 {
			thread_people_range[i].Stop_before += 1
			remainder -= 1
			offset += 1
		}
	}
	return thread_people_range
}

func round(num float64, decimal_place int) float64 {
	num = num * math.Pow(10, (float64(decimal_place)))
	num = math.Round(num)
	num = num / math.Pow(10, (float64(decimal_place)))
	return num
}

type Number interface {
	int | float64 | float32
}

func clamp[V Number](num V, max V) V {
	if num > max {
		return max
	}
	return num
}

func rand_income(mean_income int, standard_dev int) int {
	income := -1
	for income < 1000 {
		income = int(rand.NormFloat64()*float64(standard_dev)) + mean_income
	}
	return income
}

type Employee_pool map[int]*Employee

func (c *Company) Get_employees_ids(employee_type Employee_type) []int {
	return c.employeePool.Get_employees_of_company(c.ID, employee_type)
}

func (employee_pool Employee_pool) Get_employees_of_company(company_id int, employee_type Employee_type) (employees_ids_of_company []int) {
	for id := range employee_pool {
		if employee_pool[id].Employer == company_id {
			switch employee_type {
			case Employee_type_all:
				employees_ids_of_company = append(employees_ids_of_company, id)
			case employee_pool[id].EmployeeType:
				employees_ids_of_company = append(employees_ids_of_company, id)
			}
		}
	}

	return employees_ids_of_company
}

func (c *Company) SetEmployeePool(p Employee_pool) {
	c.employeePool = p
}

func (c *Company) SetProductComponents(components *ProductComponents) {
	c.productComponents = components
}

func (employee_pool Employee_pool) Get_avr_skill(companyID int, employeeType Employee_type) (avgSkill float32) {
	employeesIDs := employee_pool.Get_employees_of_company(companyID, employeeType)
	for _, id := range employeesIDs {
		avgSkill += employee_pool[id].Skill
	}

	return avgSkill / float32(len(employeesIDs))
}

func (g *GameState) resetCurrentDecisions() {
	for i := range g.Companies {
		g.CurrentDecisions[i].resetDecisions()
		g.DecisionsSubmitted[i] = false
	}
}

func (d *Decisions) resetDecisions() {
	oldPredictionSales := d.Predictions.ProductSales
	maps.Copy(d.Predictions.ProductSales, oldPredictionSales)

	d.Employees.MarketingDeltas = make([]Delta[Employee], 0)
	d.Employees.ProductionDeltas = make([]Delta[Employee], 0)

	d.Production.Machines = make([]Delta[Machine], 0)
	d.Production.Logistics = make([]Delta[Warehouse], 0)

	oldProductDecisions := d.Products
	maps.Copy(d.Products, oldProductDecisions)
}

func SynchroniseCompanyWithDecisions(company Company, decisions Decisions) Company {
	company.Name = decisions.General.CompanyName
	company.CEO = decisions.General.CEO

	for ID, d := range decisions.Products {
		offer := company.Offers[ID]
		offer.Product = d.Product
		offer.Price = d.Price
		offer.Outdated = d.Outdated
		offer.PromotionQuality = 0
		offer.Promotion = struct {
			Quantity   float32
			Quality    float32
			Price      float32
			Ecology    float32
			Ethics     float32
			Durability float32
		}{
			d.Promotion.Quantity,
			d.Promotion.Quality,
			d.Promotion.Price,
			d.Promotion.Ecology,
			d.Promotion.Ethics,
			d.Promotion.Durability,
		}
		offer.ProductStats, _ = CalculateProductStats(d.Product, *company.productComponents)

		company.Offers[ID] = offer
	}

	company.Machines = slices.Clone(company.Machines)

MachineLoop:
	for _, d := range decisions.Production.Machines {
		if d.Change == Delta_New {

			for _, m := range company.Machines {
				if d.Item.ID == m.ID {
					continue MachineLoop
				}
			}

			company.Machines = append(company.Machines, d.Item)
		}
	}

	// make sure to avoid null / undefined

	assignCompanySlicesAndMaps(&company)

	return company
}

func assignCompanySlicesAndMaps(company *Company) {
	if len(company.Machines) == 0 {
		company.Machines = make([]Machine, 0)
	}

	if company.Offers == nil {
		company.Offers = make(map[string]Offer)
	}

	if company.Reports == nil {
		company.Reports = make([]Report, 0)
	}

	for i := range company.Reports {
		if company.Reports[i].BalanceSheet.Assets == nil {
			company.Reports[i].BalanceSheet.Assets = make([]FinanceReportEntry, 0)
		}

		if company.Reports[i].BalanceSheet.Liabilities == nil {
			company.Reports[i].BalanceSheet.Liabilities = make([]FinanceReportEntry, 0)
		}

		if company.Reports[i].BalanceSheet.InvoiceLog == nil {
			company.Reports[i].BalanceSheet.InvoiceLog = make([]FinanceReportEntry, 0)
		}
	}
}

func ValidateDecisions(d Decisions) Decisions {
	for i, p := range d.Products {
		p.Promotion.Quantity = max(p.Promotion.Quantity, 0)
		d.Products[i] = p
	}

	return d
}

func deleteByIndex[V any](s []V, index ...int) []V {
	toBeDeleted := make([]bool, len(s))
	for _, i := range index {
		toBeDeleted[i] = true
	}

	var out []V
	for i, el := range s {
		if !toBeDeleted[i] {
			out = append(out, el)
		}
	}
	return out
}

func deleteByID[V interface{ get_id() int }](s []V, id ...int) []V {
	var IDsToDelete []int
	for i := range s {
		for ii := range id {
			if s[i].get_id() == id[ii] {
				IDsToDelete = append(IDsToDelete, i)
			}
		}
	}

	return deleteByIndex(s, IDsToDelete...)
}

func avr[V Number](values []V) V {
	var total V = 0

	for _, n := range values {
		total += n
	}

	return total / V(len(values))
}

/*
func max[V Number](values ...V) V {
	var max_val V = values[0]

	for _, n := range values {
		if max_val < n {
			max_val = n
		}
	}

	return max_val
}
*/

/*
func min[V Number](values ...V) V {
	var min_val V = values[0]

	for _, n := range values {
		if min_val > n {
			min_val = n
		}
	}

	return min_val
}
*/

func std_dev[V Number](values ...V) V {
	avr := avr(values)

	var Sigma V = 0
	for _, n := range values {
		Sigma += (n - avr) * (n - avr)
	}

	return V(math.Sqrt(float64(Sigma) / float64(len(values))))
}

func exponential(base, x, scale float64) float64 {
	return math.Pow(base, x) * scale
}

func scalar_product32(a, b, result []float32) float32 {
	var product float32

	simd.MulFloat32(a, b, result)
	for _, s := range result {
		product += s
	}
	return product
}

func sumFunc[V any](s []V, v func(V) float64) float64 {
	total := 0.
	for _, e := range s {
		total += v(e)
	}

	return total
}

func CompanyValuation(quartarlyCashflow, totalAssets, interestRate float64) float64 {
	cashflowPerYear := quartarlyCashflow * 4
	const years = 5

	return cashflowPerYear/(1+interestRate)*years + totalAssets
}

func quartarlyCashflow(c Company) float64 {
	reports := c.Reports[max(len(c.Reports)-4, 0):len(c.Reports)]

	totalCashflow := 0.
	for _, r := range reports {
		totalCashflow += r.FinancialReport.Totals.Cashflow
	}

	if len(reports) != 4 {
		totalCashflow = (totalCashflow * 4) / float64(len(reports))
	}

	return totalCashflow
}

func zeroIfNaN[float float64 | float32](v float) float {
	if math.IsNaN(float64(v)) {
		return 0
	}

	return v
}
