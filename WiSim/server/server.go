package main

import (
	"WiSim/simulation"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type Server struct {
	conns      map[*websocket.Conn]Player
	connsMutex sync.Mutex

	methods map[string]func(*Server, *websocket.Conn, Message[any])
}

func NewServer() *Server {
	return &Server{
		conns:   make(map[*websocket.Conn]Player),
		methods: make(map[string]func(*Server, *websocket.Conn, Message[any])),
	}
}

func (s *Server) addMethod(methodName string, methodFunc func(*Server, *websocket.Conn, Message[any])) {
	s.methods[methodName] = methodFunc
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming conn: ", ws.RemoteAddr())

	s.connsMutex.Lock()
	s.conns[ws] = Player{true, false, -1}
	s.connsMutex.Unlock()

	s.readLoop(ws)
}

func broadcast[V any](s *Server, message Message[V]) error {
	for client := range s.conns {
		err := client.WriteJSON(message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) readLoop(ws *websocket.Conn) {
	defer func() {
		println("Websocket client", ws.RemoteAddr().String(), "disconnected")
		s.connsMutex.Lock()
		delete(s.conns, ws)
		s.connsMutex.Unlock()
	}()

	defer func() {
		r := recover()
		fmt.Printf("Panic: %#v \n", r)
	}()

	for {
		message := Message[any]{}
		err := ws.ReadJSON(&message)

		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
			break
		} else if err != nil {
			println(err.Error())
			err := ws.WriteJSON(Message[any]{Method: "", IsResponse: true, Error: err.Error()})
			if err != nil {
				println(fmt.Errorf("unexpected Error durng JSON encoding: %w", err).Error())
			}
			continue
		}

		fmt.Printf("Recieving message: %+v\n", message)

		methodFunc, methodExists := s.methods[message.Method]
		if !methodExists {
			err := ws.WriteJSON(Message[any]{Method: message.Method, IsResponse: true, Error: "given method doesn't exist"})
			if err != nil {
				println("Unexpected Error durng JSON encoding: ", err.Error())
			}
			continue
		}

		methodFunc(s, ws, message)
	}
}

type Player struct {
	Active  bool
	Ready   bool
	Company int
}

type Message[V any] struct {
	Method     string
	IsResponse bool
	Error      string
	Data       *V
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const usageMessage = "USAGE: ./server <Desired Port> <Num Players> <New/Load> <Game_filepath/name>"

const (
	errorGameFull         = "error_game_full"
	errorNotAuthorised    = "error_not_authorised"
	errorUnexpected       = "error_unexpected"
	errorInvalidCompany   = "error_invalid_company"
	errorNoCompany        = "error_no_company"
	errorInvalidDecisions = "error_invalid_decisions"
	errorJSON             = "error_json"
)

func getDecisions(s *Server, ws *websocket.Conn, message Message[any]) {
	reply := Message[simulation.Decisions]{Method: message.Method, IsResponse: true}
	defer func() {
		err := ws.WriteJSON(reply)
		if err != nil {
			println("Error writing JSON to websocket: ", err.Error())
		}
	}()

	if len(gamestate.Companies) <= s.conns[ws].Company {
		reply.Error = errorInvalidCompany
		return
	}

	decisions := gamestate.Companies[s.conns[ws].Company].Get_decisions()
	reply.Data = &decisions
}

func getCompany(s *Server, ws *websocket.Conn, message Message[any]) {
	reply := Message[simulation.Company]{Method: message.Method, IsResponse: true}
	defer func() {
		err := ws.WriteJSON(reply)
		if err != nil {
			println("Error writing JSON to websocket: ", err.Error())
		}
	}()

	if len(gamestate.Companies) <= s.conns[ws].Company {
		reply.Error = errorInvalidCompany
		return
	}

	company := gamestate.Companies[s.conns[ws].Company]
	reply.Data = &company
}

func getExternalFactors(s *Server, ws *websocket.Conn, message Message[any]) {
	reply := Message[simulation.External_factors]{Method: message.Method, IsResponse: true}
	defer func() {
		err := ws.WriteJSON(reply)
		if err != nil {
			println("Error writing JSON to websocket: ", err.Error())
		}
	}()

	externalFactors := gamestate.External_factors
	reply.Data = &externalFactors
}

func setCompany(s *Server, ws *websocket.Conn, message Message[any]) {
	reply := Message[bool]{Method: message.Method, IsResponse: true}

	defer func() {
		err := ws.WriteJSON(reply)
		if err != nil {
			println("Error writing JSON to websocket: ", err.Error())
		}
	}()

	var requestedCompanyID int
	tempStruct := struct{ ID int }{}
	err := mapstructure.Decode(*message.Data, &tempStruct)
	if err != nil {
		reply.Error = err.Error()
		return
	}

	requestedCompanyID = tempStruct.ID

	playerGotRequestedID := false
	reply.Data = &playerGotRequestedID

	if requestedCompanyID >= PLAYER_COUNT {
		reply.Error = "ID too high"
		return
	}

	s.connsMutex.Lock()
	defer s.connsMutex.Unlock()

	for conn := range s.conns {
		player := s.conns[conn]
		if player.Company == requestedCompanyID {
			reply.Error = "ID taken"
			return
		}
	}

	playerGotRequestedID = true

	player := s.conns[ws]
	player.Company = requestedCompanyID
	s.conns[ws] = player
}

func setDecisions(s *Server, ws *websocket.Conn, message Message[any]) {
	reply := Message[any]{Method: message.Method, IsResponse: true}
	defer func() {
		err := ws.WriteJSON(reply)
		if err != nil {
			println("Error writing JSON to websocket: ", err.Error())
		}
	}()

	decisions := simulation.Decisions{}
	err := mapstructure.Decode(message.Data, &decisions)
	if err != nil {
		reply.Error = err.Error()
		return
	}

	if len(gamestate.Companies) <= s.conns[ws].Company {
		reply.Error = errorInvalidCompany
		return
	}

	gamestate.Current_decisions[s.conns[ws].Company] = decisions
}

func setReady(s *Server, ws *websocket.Conn, message Message[any]) {
	reply := Message[any]{Method: message.Method, IsResponse: true}
	defer func() {
		err := ws.WriteJSON(reply)
		if err != nil {
			println("Error writing JSON to websocket: ", err.Error())
		}
	}()

	s.connsMutex.Lock()
	defer s.connsMutex.Lock()
	player, exists := s.conns[ws]
	if !exists {
		reply.Error = errorUnexpected
		return
	}

	player.Ready = true
	s.conns[ws] = player

	println("Company ", player.Company, "ready!")

	for client := range s.conns {
		if !s.conns[client].Ready {
			return
		}
	}
	println("Every player is ready: running simulation")

	runSimulation(s)
}

func setUnReady(s *Server, ws *websocket.Conn, message Message[any]) {
	reply := Message[any]{Method: message.Method, IsResponse: true}
	defer func() {
		err := ws.WriteJSON(reply)
		if err != nil {
			println("Error writing JSON to websocket: ", err.Error())
		}
	}()

	s.connsMutex.Lock()
	defer s.connsMutex.Unlock()

	player, exists := s.conns[ws]
	if !exists {
		reply.Error = errorUnexpected
		return
	}

	player.Ready = false
	s.conns[ws] = player
	println("Company ", player.Company, "unready!")
}

func runSimulation(s *Server) {
	broadcast(s, Message[any]{Method: "bSim_starting", IsResponse: false})

	simDoneMessage := Message[any]{Method: "bSim_done", IsResponse: false}
	defer func() {
		err := broadcast(s, simDoneMessage)
		if err != nil {
			println("Error broadcasting JSON to websockets: ", err.Error())
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			simDoneMessage.Error = fmt.Sprint("Critical Simulation Error: ", r)
		}
	}()

	err := gamestate.Simulate_step()
	if err != nil {
		simDoneMessage.Error = err.Error()
	}
}

func calculateProductStats(s *Server, ws *websocket.Conn, message Message[any]) {
	reply := Message[simulation.Product]{Method: message.Method, IsResponse: true}
	defer func() {
		err := ws.WriteJSON(reply)
		if err != nil {
			println("Error writing JSON to websocket: ", err.Error())
		}
	}()

	decisions := struct {
		Product  simulation.Decisions_product
		Research simulation.Decisions_research
	}{}
	err := mapstructure.Decode(message.Data, &decisions)
	if err != nil {
		reply.Error = err.Error()
		return
	}

	if len(gamestate.Companies) <= s.conns[ws].Company {
		reply.Error = errorInvalidCompany
		return
	}

	product := gamestate.Companies[s.conns[ws].Company].Calculate_product(decisions.Product, decisions.Research)
	reply.Data = &product
}

func sendChat(s *Server, ws *websocket.Conn, message Message[any]) {
	reply := Message[string]{Method: message.Method, IsResponse: true}
	defer func() {
		err := ws.WriteJSON(reply)
		if err != nil {
			println("Error writing JSON to websocket: ", err.Error())
		}
	}()

	if message.Data != nil {
		chat := fmt.Sprint(*message.Data)
		reply.Data = &chat
		broadcast(s, reply)
		fmt.Printf("CHAT (%d): %s\n", s.conns[ws].Company, *message.Data)
	}
}

func main() {
	// Check args, format "[EXECUTABLE NAME] <PORT> <Number of Players>"
	if len(os.Args) < 4 {
		log.Fatal(usageMessage)
	}

	var err error
	PORT, err = strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %s \n%s\n", err.Error(), usageMessage)
	}

	PLAYER_COUNT, err = strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Error: %s \n%s\n", err.Error(), usageMessage)
	}

	// load sim_config & start game
	sim_config, err = simulation.Load_sim_config(filepath.Dir(os.Args[0]) + "/sim_config.json")
	if err != nil {
		log.Fatalf("Error: %s \n", err.Error())
	}

	gamestate = simulation.New_game(sim_config, PLAYER_COUNT, "TempGameName")

	// Prepare server

	server := NewServer()

	server.addMethod("gDecisions", getDecisions)
	server.addMethod("gCompany", getCompany)
	server.addMethod("gExternal_factors", getExternalFactors)

	server.addMethod("sCompany", setCompany)
	server.addMethod("sDecisions", setDecisions)
	server.addMethod("sReady", setReady)
	server.addMethod("sUnready", setUnReady)

	server.addMethod("fProduct_stats", calculateProductStats)

	server.addMethod("bChat", sendChat)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			println("Upgrade failed", err.Error())
			return
		}
		server.handleWS(conn)
	})

	println("Serving on port: ", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal(err)
	}
}
