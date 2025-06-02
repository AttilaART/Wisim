package main

import (
	"WiSim/simulation"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	// Load settings
	sim_config_file, err := os.ReadFile("Config/sim_config.json")
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

	// Setup game
	game_state := simulation.New_game(sim_config, 4, "Test_game")

	// Load game
	//game_state, err := load_game("Saves/Test_game-10.json")
	//if err != nil {
	//	println(err.Error())
	//	log.Fatal("Failed to load save")
	//}

	game_state.Current_decisions, err = simulation.Get_decisions(
		fmt.Sprintf("Saves/%s-0/Decisions", game_state.Game_name),
		len(game_state.Companies),
	)
	// game_state.Current_decisions, err = get_decisions(
	//
	//	fmt.Sprintf("Saves/%s-%d/Decisions", game_state.Game_name, game_state.Step),
	//	len(game_state.Companies),
	//
	// )
	if err != nil {
		println(err.Error())
		log.Fatal("Failed to get current decisions")
	}

	// fmt.Printf("%+#v\n", game_state.Current_decisions[0])

	for range 100 {
		err = game_state.Simulate_step()
		if err != nil {
			println(err.Error())
			log.Fatal("An error has occured")
		}
	}
	//for _, c := range game_state.Companies {
	//	company_as_json, err := json.MarshalIndent(c, "", "  ")
	//	if err != nil {
	//		log.Fatal("Error occored printing companies")
	//	}
	//
	//	 println(string(company_as_json))
	//}

	// println("Saving file")
	// err = game_state.Save_game("Saves", false)
	//
	//	if err != nil {
	//		log.Fatal(err.Error())
	//	}
}
