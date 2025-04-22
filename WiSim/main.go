package main

import (
	"WiSim/simulation"
	"embed"
	"encoding/json"
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
} = struct {
	is_loaded bool
	state     simulation.Game_state
}{false, simulation.Game_state{}}

var old_game_state struct {
	is_loaded bool
	state     simulation.Game_state
} = struct {
	is_loaded bool
	state     simulation.Game_state
}{false, simulation.Game_state{}}

func main() {
	sim_config_file, err := os.ReadFile("simulation/Config/sim_config.json")
	if err != nil {
		println("Error loading sim_config.json")
		os.Exit(1)
	}

	var sim_config simulation.Sim_config
	err = json.Unmarshal(sim_config_file, &sim_config)
	if err != nil {
		println("Error in sim_config.json")
		os.Exit(2)
	}

	game_state.state = simulation.New_game(sim_config, 4, "Test_ui")
	game_state.is_loaded = true

	game_state.state.Current_decisions, err = simulation.Get_decisions("simulation/Saves/Test_game-0/Decisions", 4)
	if err != nil {
		log.Fatal(err.Error())
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
