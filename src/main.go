package main

//main package
//for the webserver
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	ErrorHelper "gameserver.speedrun.io/Helper/Errorhelper"
	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
	PoolHelper "gameserver.speedrun.io/Helper/Poolhelper"
	SocketHelper "gameserver.speedrun.io/Helper/Sockethelper"
	userHelper "gameserver.speedrun.io/Helper/Userhelper"
)

var roomList PoolHelper.MapPool
var tokenSecret []byte

//handles incoming http requests over the websocket protocoll
func handleWebsocketInput(w http.ResponseWriter, r *http.Request) {
	var socketConn, err = SocketHelper.WsEndpoint(w, r)
	//if connection was not http => give error
	if err != nil {
		ErrorHelper.ConnectionNotWebsocketError(w, r)
		return
	}
	_, p, err := socketConn.ReadMessage()
	var message string
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	} else {
		message = string(p)
	}
	parsedRequest := ObjectStructures.AuthMessage{}
	err = json.Unmarshal([]byte(message), &parsedRequest)
	//if validuser continue connection, else close socket
	if ok, err := userHelper.ValidateJWSToken(parsedRequest.Token, tokenSecret, parsedRequest.Name); err == nil && ok {
		ErrorHelper.OutputToConsole("Update", " valid Player with name "+parsedRequest.Name+"has connected ")
		SocketHelper.Sender(socketConn, ObjectStructures.ReturnMessage{Type: 42, LobbyData: (ObjectStructures.LobbyData{}), Highscore: ([]ObjectStructures.HighScoreStruct{}), PlayerPos: ([]ObjectStructures.PlayerPosition{}), ChatMessage: ""})
		PoolHelper.InitInputHandler(socketConn, roomList, parsedRequest.Name)
	} else {
		socketConn.Close()
	}

}

func setupRoutes(router *mux.Router) {
	router.HandleFunc("/", ErrorHelper.InvalidRouteError)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebsocketInput(w, r)
	})
}

func main() {
	//read very secret token
	cont, err := ioutil.ReadFile("./cert/jwtSecret.txt")
	if err != nil {
		fmt.Println(err)
	}
	tokenSecret = []byte(cont)
	router := mux.NewRouter()

	//Create List of all active lobbies
	roomList = PoolHelper.MapPool{Maps: make(map[string]ObjectStructures.Pool)}
	newRoom := PoolHelper.NewPool()
	go PoolHelper.Start(true, newRoom)
	newRoom.LobbyData = ObjectStructures.LobbyData{
		ID:        "devTest",
		MapCode:   "dome",
		LobbyName: "not considered",
	}
	//open permanent room devtest
	roomList.Maps["devTest"] = *newRoom
	ErrorHelper.OutputToConsole("Update", "initializing server...")
	setupRoutes(router)
	ErrorHelper.OutputToConsole("Update", "Server online")
	corsObj := handlers.AllowedOrigins([]string{"*"})
	//log.Println(http.ListenAndServe(":8080", handlers.CORS(corsObj)(router)))
	log.Println(http.ListenAndServeTLS(":443", "./cert/certificate.pem", "./cert/key.pem", handlers.CORS(corsObj)(router)))
}
