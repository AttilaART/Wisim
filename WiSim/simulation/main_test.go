package simulation

import (
	"math/rand"
	"slices"
	"testing"
)

func Test_round(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		num           float64
		decimal_place int
		want          float64
	}{
		// TODO: Add test cases.
		struct {
			name          string
			num           float64
			decimal_place int
			want          float64
		}{name: "Check 0", num: 0, decimal_place: 0, want: 0},
		struct {
			name          string
			num           float64
			decimal_place int
			want          float64
		}{name: "Check multiple positive decimal places", num: 123.4567, decimal_place: 2, want: 123.46},
		struct {
			name          string
			num           float64
			decimal_place int
			want          float64
		}{name: "Check multiple negative decimal places", num: 78910.1112, decimal_place: -2, want: 78900},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := round(tt.num, tt.decimal_place)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("round() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assign_workers(t *testing.T) {
	gen_test := func(name string) struct {
		name     string
		machines []Machine
		workers  []Employee
		want     []Machine
		want2    int
	} {
		test := struct {
			name     string
			machines []Machine
			workers  []Employee
			want     []Machine
			want2    int
		}{}
		test.name = name
		gen_machines := func(order []int) []Machine {
			worst_machine := Machine{Production_capacity: 10000, Required_workers: 10}

			machines := make([]Machine, len(order))
			for i, factor := range order {
				current_machine := worst_machine
				current_machine.Required_workers -= factor
				current_machine.Production_capacity += factor * 1000
				machines[i] = current_machine
			}
			return machines
		}
		test.machines = gen_machines([]int{1, 3, 4, 2, 3})
		test.workers = func(number_of_workers int) []Employee {
			workers := make([]Employee, number_of_workers)
			for i := range workers {
				workers[i] = generate_employee(0, 8, production_employee, rand.Float32())
			}
			return workers
		}(20)
		assign_w := func(machines []Machine, workers []Employee) ([]Machine, int) {
			machines = slices.SortedFunc(func(yield func(Machine) bool) {
				for _, m := range machines {
					if !yield(m) {
						return
					}
				}
			}, func(a, b Machine) int {
				if a.Production_capacity < b.Production_capacity {
					return 1
				} else if a.Production_capacity == b.Production_capacity {
					return 0
				}
				return -1
			})

			ii := 0
			missing_workers := 0
			for i, m := range machines {
				for range m.Required_workers {
					if ii > len(workers)-1 {
						break
					}
					machines[i].Assigned_workers = append(machines[i].Assigned_workers, workers[ii])
					ii++
				}
				missing_workers += len(machines[i].Assigned_workers) - m.Required_workers
			}

			unasigned_workers := len(workers) - ii + 1

			if unasigned_workers > missing_workers {
				return machines, unasigned_workers
			}
			return machines, -missing_workers
		}

		test.want, test.want2 = assign_w(test.machines, test.workers)
		return test
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		machines []Machine
		workers  []Employee
		want     []Machine
		want2    int
	}{
		gen_test("test1"),
		gen_test("test2"),
		gen_test("test3"),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := assign_workers(tt.machines, tt.workers)
			// TODO: update the condition below to compare got with tt.want.
			test_passed := func(got, want []Machine) bool {
				if len(got) != len(want) {
					return false
				}

				for i := range want {
					if slices.Equal(got[i].Assigned_workers, want[i].Assigned_workers) {
						return false
					}
					if got[i].Production_capacity != want[i].Production_capacity {
						return false
					}
				}

				return true
			}
			if !test_passed(got, tt.want) {
				t.Errorf("assign_workers() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("assign_workers() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
