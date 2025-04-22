package main

import (
	"WiSim/simulation"
	"context"
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

func (a *App) Get_accounting_data(company int, year int, data string) ([]simulation.FinanceReportEntry, error) {
	var financial_data []simulation.FinanceReportEntry
	var err error
	switch data {
	case "Income_statement":
		financial_data, err = get_income_statement_data(company, year)
		if err != nil {
			return nil, err
		}
	case "Liabilities":
		financial_data, err = get_liabilities_data(company, year)
		if err != nil {
			return nil, err
		}
	case "Assets":
		financial_data, err = get_assets_sheet_data(company, year)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("given data-category invalid")
	}

	return financial_data, nil
}

func (a *App) Get_financial_report(company int, year int) (simulation.Financial_Report, error) {
	if game_state_pointer == nil {
		return simulation.Financial_Report{}, errors.New("game hasn't loaded yet")
	}

	if company > len((*game_state_pointer).Companies)-1 || company < 0 {
		return simulation.Financial_Report{}, errors.New("invalid company")
	}

	return (*game_state_pointer).Companies[company].Reports[year].Financial_Report, nil
}
