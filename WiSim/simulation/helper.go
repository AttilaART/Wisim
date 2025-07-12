package simulation

import (
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
	if float64(num) > float64(max) {
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

func (employee_pool Employee_pool) Get_employees_of_company(company_id int, employee_type Employee_type) (employees_of_company []*Employee) {
	for i := range employee_pool {
		if employee_pool[i].Employer == company_id && employee_pool[i].Employee_type == employee_type {
			employees_of_company = append(employees_of_company, &employee_pool[i])
		}
	}

	return employees_of_company
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
