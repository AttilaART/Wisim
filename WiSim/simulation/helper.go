package simulation

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

func round(num float64, decimal_place int) float64 {
	num = num * math.Pow(10, (float64(decimal_place)))
	num = math.Round(num)
	num = num / math.Pow(10, (float64(decimal_place)))
	return num
}

func rand_income(mean_income int, standard_dev int) int {
	income := -1
	for income < 1000 {
		income = int(rand.NormFloat64()*float64(standard_dev)) + mean_income
	}
	return income
}

func find_employee_by_id(id int, employees []Employee) (Employee, error) {
	for i := range employees {
		if employees[i].Id == id {
			return employees[i], nil
		}
	}
	return Employee{}, errors.New(fmt.Sprint("could not find employee with Id ", id))
}

func delete_by_index[V Employee | any](s []V, index ...int) []V {
	var s_not_deleted []V
	is_to_be_deleted := func(i int) bool {
		for _, ii := range index {
			if i == ii {
				return true
			}
		}
		return false
	}
	for i, el := range s {
		if !is_to_be_deleted(i) {
			s_not_deleted = append(s_not_deleted, el)
		}
	}
	return s_not_deleted
}
