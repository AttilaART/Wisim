package main

import (
	"WiSim/simulation"
	"errors"
)

func get_assets_sheet_data(company int, step int) ([]simulation.FinanceReportEntry, error) {
	if !game_state.is_loaded {
		return nil, errors.New("game hasn't loaded yet")
	}

	if company > len(game_state.state.Companies)-1 || company < 0 {
		return nil, errors.New("invalid company")
	}
	if step > game_state.state.Step {
		return nil, errors.New("this step hasn't been simulated yet")
	}

	if step < 0 {
		return nil, errors.New("step cannot be less than 0")
	}

	return game_state.state.Companies[company].Reports[step].Balance_sheet.Assets, nil
}

func get_liabilities_data(company int, step int) ([]simulation.FinanceReportEntry, error) {
	if !game_state.is_loaded {
		return nil, errors.New("game hasn't loaded yet")
	}

	if company > len(game_state.state.Companies)-1 || company < 0 {
		return nil, errors.New("invalid company")
	}
	if step > game_state.state.Step {
		return nil, errors.New("this step hasn't been simulated yet")
	}

	if step < 0 {
		return nil, errors.New("step cannot be less than 0")
	}

	return game_state.state.Companies[company].Reports[step].Balance_sheet.Liabilities, nil
}

func get_income_statement_data(company int, step int) ([]simulation.FinanceReportEntry, error) {
	if !game_state.is_loaded {
		return nil, errors.New("game hasn't loaded yet")
	}

	if company > len(game_state.state.Companies)-1 || company < 0 {
		return nil, errors.New("invalid company")
	}

	if step > game_state.state.Step {
		return nil, errors.New("this step hasn't been simulated yet")
	}

	if step < 0 {
		return nil, errors.New("step cannot be less than 0")
	}

	return game_state.state.Companies[company].Reports[step].Balance_sheet.Income_statement, nil
}
