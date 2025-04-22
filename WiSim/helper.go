package main

import (
	"WiSim/simulation"
	"errors"
)

func get_assets_sheet_data(company int, year int) ([]simulation.FinanceReportEntry, error) {
	if game_state_pointer == nil {
		return nil, errors.New("game hasn't loaded yet")
	}

	if company > len((*game_state_pointer).Companies)-1 || company < 0 {
		return nil, errors.New("invalid company")
	}

	return (*game_state_pointer).Companies[company].Reports[year].Balance_sheet.Assets, nil
}

func get_liabilities_data(company int, year int) ([]simulation.FinanceReportEntry, error) {
	if game_state_pointer == nil {
		return nil, errors.New("game hasn't loaded yet")
	}

	if company > len((*game_state_pointer).Companies)-1 || company < 0 {
		return nil, errors.New("invalid company")
	}

	return (*game_state_pointer).Companies[company].Reports[year].Balance_sheet.Liabilities, nil
}

func get_income_statement_data(company int, year int) ([]simulation.FinanceReportEntry, error) {
	if game_state_pointer == nil {
		return nil, errors.New("game hasn't loaded yet")
	}

	if company > len((*game_state_pointer).Companies)-1 || company < 0 {
		return nil, errors.New("invalid company")
	}

	return (*game_state_pointer).Companies[company].Reports[year].Balance_sheet.Income_statement, nil
}
