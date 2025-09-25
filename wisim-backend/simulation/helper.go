package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"

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
			if employee_type == Employee_type_all {
				employees_ids_of_company = append(employees_ids_of_company, id)
			} else if employee_type == employee_pool[id].EmployeeType {
				employees_ids_of_company = append(employees_ids_of_company, id)
			}
		}
	}

	return employees_ids_of_company
}

func (employee_pool Employee_pool) Get_avr_skill(company_id int, employee_type Employee_type) (avrg_skill float32) {
	employees_ids := employee_pool.Get_employees_of_company(company_id, employee_type)
	for _, id := range employees_ids {
		avrg_skill += employee_pool[id].Skill
	}

	return avrg_skill / float32(len(employees_ids))
}

func (c *Company) Get_decisions() Decisions {
	var decisions Decisions = Decisions{}
	if len(c.DecisionHistory) >= 1 {
		decisions = c.DecisionHistory[len(c.DecisionHistory)-1]
	} else {
		decisions = Decisions{
			Products: make(map[string]Decisions_product),
			Research: Decisions_research{
				Quality:         1000,
				Durability:      1000,
				Ecology:         1000,
				Promotion:       1000,
				Production_cost: 1000,
			},
			Production: struct {
				Machines  []Delta[Machine]
				Logistics []Delta[Warehouse]
			}{
				make([]Delta[Machine], 0),
				make([]Delta[Warehouse], 0),
			},
			Employees: struct {
				Production_deltas []Delta[Employee]
				Marketing_deltas  []Delta[Employee]
				Severance_pay     float32
			}{
				make([]Delta[Employee], 0),
				make([]Delta[Employee], 0),
				10000,
			},
		}

		for productId := range c.Offers {
			decisions.Products[productId] = Decisions_product{
				Price: 350,
				Materials: struct {
					Quality         float32
					Ecology         float32
					EthicalSourcing float32
				}{
					50,
					50,
					50,
				},
				Manufacturing: struct {
					Quality            float32
					EcologicalEnergy   float32
					MaterialEfficiency float32
					Durability         float32
					MaxDurability      int
				}{
					50,
					50,
					50,
					50,
					5,
				},
			}
		}

		fmt.Println("No decision history!")
	}

	//  make sure these are more than 0.1 (otherwise simulation breaks)
	for id := range decisions.Products {
		ProductDecisions := decisions.Products[id]
		manufacturingReflect := reflect.ValueOf(&ProductDecisions.Manufacturing)
		for i := range manufacturingReflect.NumField() {
			field := manufacturingReflect.Field(i)
			if field.CanFloat() {
				field.SetFloat(max(field.Float(), 0.1))
			}
		}

		materialsReflect := reflect.ValueOf(&ProductDecisions.Materials)
		for i := range materialsReflect.NumField() {
			field := materialsReflect.Field(i)
			if field.CanFloat() {
				field.SetFloat(max(field.Float(), 0.1))
			}
		}
	}

	// initialise slices
	decisions.Employees.Marketing_deltas = make([]Delta[Employee], 0)
	decisions.Employees.Production_deltas = make([]Delta[Employee], 0)

	decisions.Production.Logistics = make([]Delta[Warehouse], 0)
	decisions.Production.Machines = make([]Delta[Machine], 0)

	return decisions
}

func delete_by_index[V any](s []V, index ...int) []V {
	to_be_deleted := make([]bool, len(s))
	for _, i := range index {
		to_be_deleted[i] = true
	}

	var out []V
	for i, el := range s {
		if !to_be_deleted[i] {
			out = append(out, el)
		}
	}
	return out
}

func delete_by_id[V interface{ get_id() int }](s []V, id ...int) []V {
	var indexes_to_delete []int
	for i := range s {
		for ii := range id {
			if s[i].get_id() == id[ii] {
				indexes_to_delete = append(indexes_to_delete, i)
			}
		}
	}

	return delete_by_index(s, indexes_to_delete...)
}

func check_product(p Product) string {
	errorString := ""
	if p.ProductionCost <= 0 {
		errorString += "Production_cost <= 0;"
		p.ProductionCost = 0.01
	}
	if p.Weight <= 0 {
		errorString += "Weight <= 0;"
		p.Weight = 0.01
	}
	if p.MaterialUse <= 0 {
		errorString += "Material_use <= 0;"
		p.MaterialUse = 0
	}
	if p.Durabilty < 0 {
		errorString += "Durabilty < 0;"
	}

	if math.IsInf(float64(p.ProductionCost), 1) {
		errorString += "Production_cost == +Inf;"
	}

	if errorString != "" {
		return "Invalid_product: " + errorString
	}
	return ""
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
