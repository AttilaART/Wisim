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
	err := check_request(company, step)
	if err != nil {
		return simulation.Financial_Report{}, err
	}

	return game_state.state.Companies[company].Reports[step].Financial_Report, nil
}

func (a *App) Get_sales_statistics(company int, step int) (simulation.Sales_statistics, error) {
	err := check_request(company, step)
	if err != nil {
		if err.Error() != "invalid company" {
			return simulation.Sales_statistics{}, err
		}
	}

	if company == -1 {
		return game_state.state.Market_sales_statistics[step], nil
	} else if company < -1 || company > len(game_state.state.Companies)-1 {
		return simulation.Sales_statistics{}, errors.New("invalid company")
	}

	return game_state.state.Companies[company].Reports[step].Sales_report.Company_sales_statistics, nil
}

func (a *App) Get_marketing_statistics(company int, step int) (simulation.Marketing_statistics, error) {
	err := check_request(company, step)
	if err != nil {
		return simulation.Marketing_statistics{}, err
	}

	return game_state.state.Companies[company].Reports[step].Sales_report.Marketing_statistics, nil
}

func (a *App) Get_product_statistics(company int, step int) (simulation.Product_statistics, error) {
	err := check_request(company, step)
	if err != nil {
		return simulation.Product_statistics{}, err
	}

	return game_state.state.Companies[company].Reports[step].Sales_report.Product_statistics, nil
}

func (a *App) Get_personelle_report(company int, step int, employee_type string) (simulation.Personelle_sub_report, error) {
	err := check_request(company, step)
	if err != nil {
		return simulation.Personelle_sub_report{}, err
	}

	switch employee_type {
	case "General":
		return game_state.state.Companies[company].Reports[step].Personelle_report.General, nil
	case "Production":
		return game_state.state.Companies[company].Reports[step].Personelle_report.Production, nil
	case "Marketing":
		return game_state.state.Companies[company].Reports[step].Personelle_report.Marketing, nil
	}

	return simulation.Personelle_sub_report{}, errors.New("invalid employee type")
}

func (a *App) Get_production_report(company int, step int) (simulation.Production_report, error) {
	err := check_request(company, step)
	if err != nil {
		return simulation.Production_report{}, err
	}

	return game_state.state.Companies[company].Reports[step].Production_report, nil
}

func (a *App) Get_availible_decisions(step int) (simulation.Decisions, error) {
	err := check_request(0, step)
	if err != nil {
		if err.Error() != "this step hasn't been simulated yet" {
			return simulation.Decisions{}, err
		}
	}

	if step > (game_state.state.Step + 1) {
		return simulation.Decisions{}, errors.New("this step hasn't been simulated yet")
	}

	return simulation.Decisions{}, nil
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

func (a *App) Trigger_simulation(force bool) (int, error) {
	if !game_state.is_loaded {
		return 0, errors.New("game hasn't loaded yet")
	}

	for _, d := range game_state.state.Decisions_submitted {
		if !d {
			return 0, errors.New("not all companies' decisions have been submitted")
		}

		fmt.Printf("decisions_loaded: %t", d)
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

func (a *App) New_simulation() (int, error) {
	game_state.state = simulation.New_game(game_state.config, 1, "Test_ui")
	game_state.is_loaded = true

	//game_state.state.Current_decisions, err = simulation.Get_decisions("simulation/Saves/Test_game-0/Decisions", 4)
	//if err != nil {
	//	return 0, err
	//}

	return game_state.state.Step, nil
}

func (a *App) Initial_app_load() error {
	sim_config_file, err := os.ReadFile("simulation/Config/sim_config.json")
	if err != nil {
		return errors.New("error loading sim_config.json")
	}

	err = json.Unmarshal(sim_config_file, &game_state.config)
	if err != nil {
		return errors.New("error in sim_config.json")
	}

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
