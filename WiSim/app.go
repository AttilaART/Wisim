package main

import (
	"WiSim/simulation"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Get_simulation_step() int {
	if !game_state.is_loaded {
		return -1
	}

	return game_state.state.Step
}

func (a *App) Get_Decisions(company, step int) (simulation.Decisions, error) {
	err := check_request(company, step)
	if err != nil {
		if err.Error() == "game hasn't loaded yet" {
			return game_state.state.Companies[company].Get_decisions(), nil
		}
		return simulation.Decisions{}, err
	}
	return game_state.state.Companies[company].Get_decisions(), nil
}

func (a *App) Get_External_Factors() simulation.External_factors {
	return game_state.state.External_factors
}

func (a *App) Get_bank_balance(company int) (float64, error) {
	if !game_state.is_loaded {
		return 0, errors.New("game hasn't loaded yet")
	}

	if company > len(game_state.state.Companies) {
		return 0, errors.New("invalid company")
	}

	return game_state.state.Companies[company].Balance, nil
}

func (a *App) Get_current_stock(company int) (int, error) {
	if !game_state.is_loaded {
		return 0, errors.New("game hasn't loaded yet")
	}

	if company > len(game_state.state.Companies) {
		return 0, errors.New("invalid company")
	}

	return game_state.state.Companies[company].Items_in_storage, nil
}

func (a *App) Get_accounting_data(company int, step int, data string) ([]simulation.FinanceReportEntry, error) {
	var financial_data []simulation.FinanceReportEntry
	var err error
	switch data {
	case "Income_statement":
		financial_data, err = get_income_statement_data(company, step)
		if err != nil {
			return nil, err
		}
	case "Liabilities":
		financial_data, err = get_liabilities_data(company, step)
		if err != nil {
			return nil, err
		}
	case "Assets":
		financial_data, err = get_assets_sheet_data(company, step)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("given data-category invalid")
	}

	return financial_data, nil
}

func (a *App) Get_reports(company, step int) (simulation.Report, error) {
	err := check_request(company, step)
	if err != nil {
		return simulation.Report{}, nil
	}

	return game_state.state.Companies[company].Reports[step], nil
}

//func (a *App) Get_availible_decisions(step int) (simulation.Decisions, error) {
//err := check_request(0, step)
//if err != nil {
//	if err.Error() != "this step hasn't been simulated yet" {
//		return simulation.Decisions{}, err
//	}
//}
//
//if step > (game_state.state.Step + 1) {
//	return simulation.Decisions{}, errors.New("this step hasn't been simulated yet")
//}
//
//return simulation.Decisions{}, nil
//}

func (a *App) Temp_generate_new_emloyee(company, step int, employee_type int) (int, error) {
	err := check_request(company, step)
	if err != nil {
		return -1, err
	}

	return game_state.state.Generate_employee(10000, 8, employee_type, 1).Id, nil
}

func (a *App) Get_past_decisions(company int, step int) (simulation.Decisions, error) {
	err := check_request(company, step)
	if err != nil {
		return simulation.Decisions{}, err
	}

	return game_state.state.Companies[company].Decision_history[step], nil
}

func (a *App) Submit_decisions(company int, decisions simulation.Decisions) error {
	err := check_request(company, 0)
	if err != nil {
		if err.Error() != "this step hasn't been simulated yet" {
			return err
		}
	}

	// Link employee id to corresponding pointer
	for i := range decisions.Employees.Marketing_actions {
		decisions.Employees.Marketing_actions[i], err = game_state.state.Link_employees_to_action(decisions.Employees.Marketing_actions[i])
		if err != nil {
			return err
		}
	}

	// TODO: Link employees assigned to machines to their pointers

	game_state.state.Current_decisions[company] = decisions
	game_state.state.Decisions_submitted[company] = true

	return nil
}

func (a *App) Get_action_employee(action simulation.Employee_action) (simulation.Employee, error) {
	e, err := game_state.state.Employees.Find_employee_by_id(action.Employee_id)

	return *e, err
}

func (a *App) Trigger_simulation(force bool) (new_step int, err error) {
	if !game_state.is_loaded {
		return 0, errors.New("game hasn't loaded yet")
	}

	for i := range game_state.state.Decisions_submitted {
		if force {
			game_state.state.Decisions_submitted[i] = true
		}

		if !game_state.state.Decisions_submitted[i] {
			return 0, errors.New("not all companies' decisions have been submitted")
		}

		fmt.Printf("decisions_loaded: %t\n", game_state.state.Decisions_submitted[i])
	}

	old_game_state.state = game_state.state
	old_game_state.is_loaded = true

	err = game_state.state.Simulate_step()
	if err != nil {
		(game_state.state) = old_game_state.state
		return 0, err
	}

	return game_state.state.Step, nil
}

func (a *App) Revert_simulation() (int, error) {
	if !game_state.is_loaded {
		return 0, errors.New("game hasn't loaded yet")
	}
	if game_state.state.Step == old_game_state.state.Step {
		return 0, errors.New("cannot go back to previous step")
	}
	if !old_game_state.is_loaded {
		return 0, errors.New("cannot go back to previous step")
	}

	game_state.state = old_game_state.state

	return game_state.state.Step, nil
}

func (a *App) New_simulation(num_companies int) (int, error) {
	game_state.state = simulation.New_game(game_state.config, num_companies, "Test_ui")

	game_state.is_loaded = true

	//game_state.state.Current_decisions, err = simulation.Get_decisions("simulation/Saves/Test_game-0/Decisions", 4)
	//if err != nil {
	//	return 0, err
	//}

	return game_state.state.Step, nil
}

func (a *App) Initial_app_load() error {
	sim_config_file, err := os.ReadFile(program_info.Data_dir + "config/sim_config.json")
	if err != nil {
		return errors.New("error loading sim_config.json at '" + program_info.Data_dir + "config/sim_config.json'")
	}

	err = json.Unmarshal(sim_config_file, &game_state.config)
	if err != nil {
		return errors.New("error in sim_config.json")
	}

	println("app loaded successfully")

	return nil
}

//  let option = {
//    title: {
//      text: "",
//    },
//    tooltip: {},
//    legend: {
//      data: ["sales"],
//    },
//    xAxis: {
//      data: ["Shirts", "Cardigans", "Chiffons", "Pants", "Heels", "Socks"],
//    },
//    yAxis: {},
//    series: [
//      {
//        name: "sales",
//        type: "bar",
//        data: [5, 20, 36, 10, 10, 20],
//      },
//    ],
//  };

type Graph_data struct {
	xAxis  []string
	yAxis  []string
	series []Series
}

type Series struct {
	name string
	data []float64
}

// Request format: company
//func (a *App) Get_graph_data(request string, start_month int, end_month int) (Graph_data error){
//  tokens := strings.Split(request, " ")
//}
