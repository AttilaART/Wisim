package main

import (
	"WiSim/simulation"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var game_state struct {
	is_loaded bool
	state     simulation.Game_state
	config    simulation.Sim_config
} = struct {
	is_loaded bool
	state     simulation.Game_state
	config    simulation.Sim_config
}{false, simulation.Game_state{}, simulation.Sim_config{}}

var old_game_state struct {
	is_loaded bool
	state     simulation.Game_state
	config    simulation.Sim_config
} = struct {
	is_loaded bool
	state     simulation.Game_state
	config    simulation.Sim_config
}{false, simulation.Game_state{}, simulation.Sim_config{}}

func main() {
	sim_config_file, err := os.ReadFile("simulation/Config/sim_config.json")
	if err != nil {
		println("Error loading sim_config.json")
		os.Exit(1)
	}

	err = json.Unmarshal(sim_config_file, &game_state.config)
	if err != nil {
		println("Error in sim_config.json")
		os.Exit(2)
	}

	game_state.state.Current_decisions, err = simulation.Get_decisions(
		fmt.Sprintf("Saves/%s-%d/Decisions", game_state.state.Game_name, game_state.state.Step),
		len(game_state.state.Companies),
	)
	if err != nil {
		println(err.Error())
		log.Fatal("Failed to get current decisions")
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "WiSim",
		Width:     1440,
		Height:    900,
		MinWidth:  1000,
		MinHeight: 400,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		EnumBind: []interface{}{simulation.AllGroups},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
