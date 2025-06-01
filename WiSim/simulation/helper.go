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

func rand_income(mean_income int, standard_dev int) int {
	income := -1
	for income < 1000 {
		income = int(rand.NormFloat64()*float64(standard_dev)) + mean_income
	}
	return income
}

func find_employee_by_id(id int, employees *[]Employee) (Employee, error) {
	for i := range *employees {
		if (*employees)[i].Id == id {
			return (*employees)[i], nil
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
