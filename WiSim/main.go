package main

import (
	"WiSim/simulation"
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed simulation/config simulation/defaults
var default_config embed.FS

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

var program_info struct {
	Working_dir string
	Data_dir    string
}

func main() {
	// get working dir
	exec_name := filepath.Base(os.Args[0])
	working_dir_temp, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	program_info.Working_dir = string(([]byte)(working_dir_temp)[0 : len(working_dir_temp)-len(exec_name)])

	// make config_folder if not present
	_, err = os.ReadDir(program_info.Working_dir + "simulation/Config")
	if !os.IsExist(err) {
		os.CopyFS(program_info.Working_dir, default_config)
	} else if err != nil {
		log.Fatal(err)
	}

	program_info.Data_dir = program_info.Working_dir + "simulation/"

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "WiSim",
		Width:     1440,
		Height:    900,
		MinWidth:  1000,
		MinHeight: 400,
		Mac:       &mac.Options{Preferences: &mac.Preferences{FullscreenEnabled: mac.Enabled}},
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
