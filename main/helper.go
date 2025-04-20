package main

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
