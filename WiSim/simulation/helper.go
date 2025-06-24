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

func (employees Employee_pool) find_employee_by_id(id int) (*Employee, error) {
	for i := range employees {
		if employees[i].Id == id {
			return &employees[i], nil
		}
	}
	return nil, errors.New(fmt.Sprint("could not find employee with Id ", id))
}

func (s *Game_state) Link_employees_to_action(a Employee_action) (Employee_action, error) {
	var err error
	a.employee, err = s.Employees.find_employee_by_id(a.Employee_id)
	return a, err
}

func (c *Company) Get_decisions() Decisions {
	var decisions Decisions
	if len(c.Decision_history) >= 1 {
		decisions = c.Decision_history[len(c.Decision_history)-1]
		decisions.Production.Machines = c.Machines
		decisions.Production.Logistics = c.Warehouses

		for _, e := range c.Production_personelle {
			last_month_action := c.get_previous_employee_action_by_id(e.Id)

			action := Employee_action{
				Employee_id:    e.Id,
				employee:       e,
				Extra_training: last_month_action.Extra_training,
				Pay:            e.Pay,
				Bonus:          e.Bonus,

				Working_hours: e.Working_hours,

				Status: Existing,
			}

			decisions.Employees.Production_actions = append(decisions.Employees.Production_actions, action)
		}

		for _, e := range c.Marketing_personelle {
			last_month_action := c.get_previous_employee_action_by_id(e.Id)

			action := Employee_action{
				Employee_id:    e.Id,
				employee:       e,
				Extra_training: last_month_action.Extra_training,
				Pay:            e.Pay,
				Bonus:          e.Bonus,

				Working_hours: e.Working_hours,

				Status: Existing,
			}

			decisions.Employees.Marketing_actions = append(decisions.Employees.Marketing_actions, action)
		}
	} else {
		fmt.Println("No decision history!")
	}

	return decisions
}

func (c *Company) get_previous_employee_action_by_id(id int) Employee_action {
	e, err := c.employee_pool.find_employee_by_id(id)
	if err != nil {
		panic(err)
	}
	var search_through *[]Employee_action
	if e.Employee_type == Marketing_employee {
		search_through = &c.Decision_history[len(c.Decision_history)-1].Employees.Marketing_actions
	} else if e.Employee_type == Production_employee {
		search_through = &c.Decision_history[len(c.Decision_history)-1].Employees.Production_actions
	} else {
		panic(fmt.Sprintf("unknown employee type: %d (employee %d %s", e.Employee_type, e.Id, e.Name))
	}

	for _, a := range *search_through {
		if a.Employee_id == e.Id {
			return a
		}
	}

	fmt.Println("action not found")
	return Employee_action{
		Employee_id: e.Id,
		employee:    e,

		Pay:   e.Pay,
		Bonus: e.Bonus,

		Working_hours: e.Working_hours,

		Status: Existing,
	}
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
