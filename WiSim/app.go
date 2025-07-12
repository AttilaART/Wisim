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

func (a *App) Initial_app_load() error { //
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

func (a *App) New_simulation(num_companies int) (int, error) {
	game_state.state = simulation.New_game(game_state.config, num_companies, "Test_ui")
	game_state.is_loaded = true
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

func (a *App) Get_simulation_step() int {
	if !game_state.is_loaded {
		return -1
	}

	return game_state.state.Step
}

func (a *App) Get_company(company int) (simulation.Company, error) {
	err := check_request(company, game_state.state.Step)
	if err != nil {
		return simulation.Company{}, err
	}

	return game_state.state.Companies[company], nil
}

func (a *App) Get_external_factors() (simulation.External_factors, error) {
	return game_state.state.External_factors, nil
}

func (a *App) Get_decisions(company, step int) (simulation.Decisions, error) {
	err := check_request(company, step)
	if err != nil {
		if err.Error() == "game hasn't loaded yet" {
			return game_state.state.Companies[company].Get_decisions(), nil
		}
		return simulation.Decisions{}, err
	}

	return game_state.state.Companies[company].Get_decisions(), nil
}

func (a *App) Temp_generate_new_emloyee(company, step int, employee_type simulation.Employee_type) (int, error) {
	err := check_request(company, step)
	if err != nil {
		return -1, err
	}

	return game_state.state.Generate_employee(10000, 8, employee_type, 1).Id, nil
}

func (a *App) Get_unemployed(employee_type simulation.Employee_type) (Unemployed []*simulation.Employee) {
	for i := range game_state.state.Employees {
		if game_state.state.Employees[i].Employer == simulation.Employee_employer_none && game_state.state.Employees[i].Employee_type == employee_type {
			Unemployed = append(Unemployed, &game_state.state.Employees[i])
		}
	}

	return Unemployed
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

	game_state.state.Current_decisions[company] = decisions
	game_state.state.Decisions_submitted[company] = true

	return nil
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
		return game_state.state.Step, err
	}

	return game_state.state.Step, nil
}
