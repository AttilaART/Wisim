package main

import (
	"WiSim/simulation"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (a *App) Get_financial_report(company int, step int) (simulation.Financial_Report, error) {
	if !game_state.is_loaded {
		return simulation.Financial_Report{}, errors.New("game hasn't loaded yet")
	}

	if company > len((game_state.state).Companies)-1 || company < 0 {
		return simulation.Financial_Report{}, errors.New("invalid company")
	}

	if step > game_state.state.Step {
		return simulation.Financial_Report{}, errors.New("this step hasn't been simulated yet")
	}

	if step < 0 {
		return simulation.Financial_Report{}, errors.New("step cannot be less than 0")
	}

	s, err := json.MarshalIndent(game_state.state.Companies[company].Reports[step].Financial_Report, "-->", "    ")
	if err != nil {
		return simulation.Financial_Report{}, err
	}
	println(string(s))

	return game_state.state.Companies[company].Reports[step].Financial_Report, nil
}

func (a *App) Trigger_simulation() (int, error) {
	if !game_state.is_loaded {
		return 0, errors.New("game hasn't loaded yet")
	}

	old_game_state.state = (game_state.state)
	old_game_state.is_loaded = true

	err := (game_state.state).Simulate_step()
	if err != nil {
		(game_state.state) = old_game_state.state
		return 0, err
	}

	return game_state.state.Step, nil
}

func (a *App) Revert_simulation() error {
	if !game_state.is_loaded {
		return errors.New("game hasn't loaded yet")
	}
	if game_state.state.Step == old_game_state.state.Step {
		return errors.New("cannot go back to previous step")
	}
	if !old_game_state.is_loaded {
		return errors.New("cannot go back to previous step")
	}

	game_state.state = old_game_state.state

	return nil
}
