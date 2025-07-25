package main

import (
	"WiSim/simulation"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"golang.org/x/net/websocket"
)

const (
	method_len      = 5
	min_message_len = 8 // eg.: gDecs {}
)

type Server struct {
	conns       map[*websocket.Conn]Player
	conns_mutex sync.Mutex

	methods map[string]func(*Server, *websocket.Conn, []byte)
}

func NewServer() *Server {
	return &Server{
		conns:   make(map[*websocket.Conn]Player),
		methods: make(map[string]func(*Server, *websocket.Conn, []byte)),
	}
}

func (s *Server) addMethod(method_name string, method_func func(*Server, *websocket.Conn, []byte)) {
	s.methods[method_name] = method_func
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming conn: ", ws.RemoteAddr())

	s.conns_mutex.Lock()
	s.conns[ws] = Player{true, false, len(s.conns)}
	s.conns_mutex.Unlock()

	s.readLoop(ws)
}

func (s *Server) broadcast(message []byte) {
	for client := range s.conns {
		client.Write(message)
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		length_read, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			println("Read Error: ", err)
			continue
		}
		message := buf[:length_read]

		if len(message) < min_message_len {
			println("Invalid method: message too short")
		}
		fmt.Printf("Recieving message: %s\n", string(message))

		method := message[:method_len]

		method_func, method_exists := s.methods[string(method)]
		if !method_exists {
			fmt.Fprintf(ws, "RE %s ER \"invalid method\"", method)
			fmt.Printf("RE %s ER \"invalid method\"", method)
			continue
		}

		method_func(s, ws, message[method_len:])
	}
}

type Player struct {
	Active  bool
	Ready   bool
	Company int
}

var (
	PORT                int
	PLAYER_COUNT        int
	ACTIVE_PLAYER_COUNT int
)

var (
	sim_config simulation.Sim_config
	gamestate  simulation.Game_state
)

const USAGE_MESSAGE = "USAGE: ./server <Desired Port> <Num Players> <New/Load> <Game_filepath/name>"

const (
	error_game_full         = "error_game_full"
	error_not_authorised    = "error_not_authorised"
	error_unexpected        = "error_unexpected"
	error_invalid_company   = "error_invalid_company"
	error_no_company        = "error_no_company"
	error_invalid_decisions = "error_invalid_decisions"
	error_json              = "error_json"
)

func getDecisions(s *Server, ws *websocket.Conn, _ []byte) {
	reply := bytes.NewBuffer(make([]byte, 0))
	defer func() { ws.Write(reply.Bytes()) }()

	fmt.Fprint(reply, "RE gDecs ")

	if len(gamestate.Companies) <= s.conns[ws].Company {
		fmt.Fprint(reply, "ER ", error_invalid_company)
		return
	}

	json, err := json.Marshal(
		(gamestate.Companies[s.conns[ws].Company].Get_decisions()))
	if err != nil {
		fmt.Fprint(reply, "ER ", error_json)
		return
	}

	reply.Write(json)
}

func getCompany(s *Server, ws *websocket.Conn, _ []byte) {
	reply := bytes.NewBuffer(make([]byte, 0))
	defer func() { ws.Write(reply.Bytes()) }()

	fmt.Fprint(reply, "RE gComp ")

	if len(gamestate.Companies) <= s.conns[ws].Company {
		fmt.Fprint(reply, "ER ", error_invalid_company)
		return
	}

	json, err := json.Marshal(
		gamestate.Companies[s.conns[ws].Company])
	if err != nil {
		fmt.Fprint(reply, "ER ", error_json)
		return
	}

	fmt.Fprint(reply, "OK ")
	reply.Write(json)
}

func getExternalFactors(s *Server, ws *websocket.Conn, _ []byte) {
	reply := bytes.NewBuffer(make([]byte, 0))
	defer func() { ws.Write(reply.Bytes()) }()

	fmt.Fprint(reply, "RE gExFa ")

	if len(gamestate.Companies) <= s.conns[ws].Company {
		fmt.Fprint(reply, "ER ", error_invalid_company)
		return
	}

	json, err := json.Marshal(
		(gamestate.External_factors))
	if err != nil {
		fmt.Fprint(reply, "ER ", error_json)
		return
	}

	fmt.Fprint(reply, "OK ")
	reply.Write(json)
}

func setDecisions(s *Server, ws *websocket.Conn, message []byte) {
	reply := bytes.NewBuffer(make([]byte, 0))
	defer func() { ws.Write(reply.Bytes()) }()

	fmt.Fprint(reply, "RE sDecs ")

	decisions := simulation.Decisions{}
	err := json.Unmarshal(message[7:], &decisions)
	if err != nil {
		fmt.Fprint(reply, "ER ", error_json)
		return
	}

	if len(gamestate.Companies) <= s.conns[ws].Company {
		fmt.Fprint(reply, "ER ", error_invalid_company)
		return
	}
	gamestate.Current_decisions[s.conns[ws].Company] = decisions

	fmt.Fprint(reply, "OK {}")
}

func setReady(s *Server, ws *websocket.Conn, message []byte) {
	reply := bytes.NewBuffer(make([]byte, 0))
	defer func() { ws.Write(reply.Bytes()) }()

	fmt.Fprint(reply, "RE sRedy ")

	s.conns_mutex.Lock()
	player, exists := s.conns[ws]
	if !exists {
		s.conns_mutex.Unlock()
		fmt.Fprint(reply, "ER ", error_unexpected)
		return
	}

	player.Ready = true
	s.conns[ws] = player
	s.conns_mutex.Unlock()

	fmt.Fprint(reply, "OK {}")

	s.conns_mutex.Lock()
	for client := range s.conns {
		if !s.conns[client].Ready {
			s.conns_mutex.Unlock()
			return
		}
	}

	runSimulation(s)
}

func setUnReady(s *Server, ws *websocket.Conn, message []byte) {
	reply := bytes.NewBuffer(make([]byte, 0))
	defer func() { ws.Write(reply.Bytes()) }()

	fmt.Fprint(reply, "RE sURdy ")

	s.conns_mutex.Lock()
	player, exists := s.conns[ws]
	if !exists {
		s.conns_mutex.Unlock()
		fmt.Fprint(reply, "ER ", error_unexpected)
		return
	}

	player.Ready = false
	s.conns[ws] = player
	s.conns_mutex.Unlock()

	fmt.Fprint(reply, "OK {}")
}

func runSimulation(s *Server) {
	s.broadcast(([]byte)("bSimS {}"))

	result := struct {
		Success bool
		Message string
	}{true, ""}

	defer func() {
		if r := recover(); r != nil {
			result.Message = fmt.Sprint("Unxpected Error: ", r)
			result.Success = false
		}

		result_json, err := json.Marshal(result)
		if err != nil {
			panic(fmt.Errorf("error packing json of simulation: %w", err))
		}

		s.broadcast(append(([]byte)("bSimD "), result_json...))
	}()

	err := gamestate.Simulate_step()
	if err != nil {
		result.Message = err.Error()
		result.Success = false
	}
}

func calculateProductStats(s *Server, ws *websocket.Conn, message []byte) {
	reply := bytes.NewBuffer(make([]byte, 0))

	defer func() { ws.Write(reply.Bytes()) }()

	fmt.Fprint(reply, "RE fProd")

	decisions := struct {
		Product  simulation.Decisions_product
		Research simulation.Decisions_research
	}{}
	err := json.Unmarshal(message[7:], &decisions)
	if err != nil {
		fmt.Fprint(reply, "ER ", error_json)
		return
	}

	if len(gamestate.Companies) <= s.conns[ws].Company {
		fmt.Fprint(reply, "ER ", error_invalid_company)
		return
	}
	product := gamestate.Companies[s.conns[ws].Company].Calculate_product(decisions.Product, decisions.Research)

	product_json, err := json.Marshal(product)
	if err != nil {
		fmt.Fprint(reply, "ER ", error_json)
	}

	fmt.Fprint(reply, "OK ")
	reply.Write(product_json)
}

func sendChat(s *Server, ws *websocket.Conn, message []byte) {
	reply := bytes.NewBuffer(make([]byte, 0))

	defer func() { ws.Write(reply.Bytes()) }()

	fmt.Fprint(reply, "RE bChat ")

	message_struct := struct{ Message []byte }{}
	err := json.Unmarshal(message[7:], &message_struct)
	if err != nil {
		fmt.Fprint(reply, " ER ", err)
	}

	fmt.Fprint(reply, "OK {}")
	s.broadcast(message_struct.Message)
}

func main() {
	// Check args, format "[EXECUTABLE NAME] <PORT> <Number of Players>"
	if len(os.Args) < 4 {
		log.Fatal(USAGE_MESSAGE)
	}

	var err error
	PORT, err = strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %s \n%s\n", err.Error(), USAGE_MESSAGE)
	}

	PLAYER_COUNT, err = strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Error: %s \n%s\n", err.Error(), USAGE_MESSAGE)
	}

	// load sim_config & start game
	sim_config, err = simulation.Load_sim_config(filepath.Dir(os.Args[0]) + "/sim_config.json")
	if err != nil {
		log.Fatalf("Error: %s \n", err.Error())
	}

	gamestate = simulation.New_game(sim_config, PLAYER_COUNT, "TempGameName")

	// Prepare server

	server := NewServer()

	server.addMethod("gDecs", getDecisions)
	server.addMethod("gComp", getCompany)
	server.addMethod("gExFa", getExternalFactors)

	server.addMethod("sDecs", setDecisions)
	server.addMethod("sRedy", setReady)
	server.addMethod("sURdy", setUnReady)

	server.addMethod("fProd", calculateProductStats)

	server.addMethod("bChat", sendChat)

	http.Handle("/connect", websocket.Handler(server.handleWS))

	println("Serving on port: ", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal(err)
	}
}
