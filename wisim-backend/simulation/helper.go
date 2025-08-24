package simulation

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
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

type Employee_pool []Employee

func (employee_pool Employee_pool) Find_employee_by_id(id int) *Employee {
	for i := range employee_pool {
		if employee_pool[i].Id == id {
			return &employee_pool[i]
		}
	}
	return nil
}

func (employee_pool *Employee_pool) Get_employees_of_company(company_id int, employee_type Employee_type) (employees_of_company []*Employee) {
	for i := range *employee_pool {
		if (*employee_pool)[i].Employer == company_id {
			if employee_type == Employee_type_all {
				employees_of_company = append(employees_of_company, &(*employee_pool)[i])
			} else if employee_type == (*employee_pool)[i].Employee_type {
				employees_of_company = append(employees_of_company, &(*employee_pool)[i])
			}
		}
	}

	return employees_of_company
}

func (employee_pool *Employee_pool) Get_avr_skill(company_id int, employee_type Employee_type) (avrg_skill float32) {
	employees := employee_pool.Get_employees_of_company(company_id, employee_type)
	for _, e := range employees {
		avrg_skill += e.Skill
	}

	return avrg_skill / float32(len(employees))
}

func (c *Company) Get_decisions() Decisions {
	var decisions Decisions
	if len(c.Decision_history) >= 1 {
		decisions = c.Decision_history[len(c.Decision_history)-1]
	} else {
		fmt.Println("No decision history!")
	}

	//  make sure these are more than 0.1 (otherwise simulation breaks)
	more_than_0 := []float32{
		decisions.Marketing.Product.Materials.Quality,
		decisions.Marketing.Product.Materials.Ecology,
		decisions.Marketing.Product.Materials.Ethical_sourcing,

		decisions.Marketing.Product.Manufacturing.Quality,
		decisions.Marketing.Product.Manufacturing.Durability,
		decisions.Marketing.Product.Manufacturing.Ecological_energy,
		decisions.Marketing.Product.Manufacturing.Material_efficiency,
	}

	for i := range more_than_0 {
		if more_than_0[i] < 0.1 {
			more_than_0[i] = 0.1
		}
	}

	decisions.Marketing.Product.Materials.Quality = more_than_0[0]
	decisions.Marketing.Product.Materials.Ecology = more_than_0[1]
	decisions.Marketing.Product.Materials.Ethical_sourcing = more_than_0[2]

	decisions.Marketing.Product.Manufacturing.Quality = more_than_0[3]
	decisions.Marketing.Product.Manufacturing.Durability = more_than_0[4]
	decisions.Marketing.Product.Manufacturing.Ecological_energy = more_than_0[5]
	decisions.Marketing.Product.Manufacturing.Material_efficiency = more_than_0[6]

	if decisions.Marketing.Product.Manufacturing.Max_durability < 1 {
		decisions.Marketing.Product.Manufacturing.Max_durability = 1
	}

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

func check_product(p Product) error {
	errorString := ""
	if p.Base_quality <= 0 {
		errorString += "Base_quality <= 0;"
	}
	if p.Base_durability <= 0 {
		errorString += "Base_durability <= 0;"
	}
	if p.Base_production_cost <= 0 {
		errorString += "Invalid_product: Base_production_cost <= 0;"
	}
	if p.Base_ecology <= 0 {
		errorString += "Base_ecology <= 0;"
	}
	if p.Base_material_use <= 0 {
		errorString += "Base_material_use <= 0;"
	}
	if p.Production_cost <= 0 {
		errorString += "Production_cost <= 0;"
	}
	if p.Weight <= 0 {
		errorString += "Weight <= 0;"
	}
	if p.Material_use <= 0 {
		errorString += "Material_use <= 0;"
	}
	if p.Durabilty < 0 {
		errorString += "Durabilty < 0;"
	}

	if math.IsInf(float64(p.Production_cost), 1) {
		errorString += "Production_cost == +Inf;"
	}

	if errorString != "" {
		return errors.New("Invalid_product: " + errorString)
	}
	return nil
}

func avr[V Number](values []V) V {
	var total V = 0

	for _, n := range values {
		total += n
	}

	return total / V(len(values))
}

func max[V Number](values []V) V {
	var max_val V = 0

	for _, n := range values {
		if max_val < n {
			max_val = n
		}
	}

	return max_val
}

func min[V Number](values []V) V {
	var min_val V = 0

	for _, n := range values {
		if min_val > n {
			min_val = n
		}
	}

	return min_val
}

func std_dev[V Number](values []V) V {
	avr := avr(values)

	var Sigma V = 0
	for _, n := range values {
		Sigma += (n - avr) * (n - avr)
	}

	return V(math.Sqrt(float64(Sigma) / float64(len(values))))
}
