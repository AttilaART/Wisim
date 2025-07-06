package main

import (
	"WiSim/simulation"
	"errors"
)

func check_request(company int, step int) error {
	if !game_state.is_loaded {
		return errors.New("game hasn't loaded yet")
	}

	if company > len(game_state.state.Companies)-1 || company < 0 {
		return errors.New("invalid company")
	}
	if step > game_state.state.Step {
		return errors.New("this step hasn't been simulated yet")
	}

	if step < 0 {
		return errors.New("step cannot be less than 0")
	}

	return nil
}

func get_assets_sheet_data(company int, step int) ([]simulation.FinanceReportEntry, error) {
	err := check_request(company, step)
	if err != nil {
		return nil, err
	}

	return game_state.state.Companies[company].Reports[step].Balance_sheet.Assets, nil
}

func get_liabilities_data(company int, step int) ([]simulation.FinanceReportEntry, error) {
	err := check_request(company, step)
	if err != nil {
		return nil, err
	}

	return game_state.state.Companies[company].Reports[step].Balance_sheet.Liabilities, nil
}

func get_income_statement_data(company int, step int) ([]simulation.FinanceReportEntry, error) {
	err := check_request(company, step)
	if err != nil {
		return nil, err
	}

	return game_state.state.Companies[company].Reports[step].Balance_sheet.Invoice_log, nil
}
