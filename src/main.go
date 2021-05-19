package main

//main package
//for the webserver
import (
	"encoding/json"
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

var roomList map[string]PoolHelper.Pool

func handleWebsocketInput(w http.ResponseWriter, r *http.Request) {
	var socketConn, err = SocketHelper.WsEndpoint(w, r)
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
	if userHelper.ValidateUser(parsedRequest.Name) {
		ErrorHelper.OutputToConsole("Update", " valid Player with name "+parsedRequest.Name+"has connected ")
		if parsedRequest.Name == "" {
			return
		}
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
	router := mux.NewRouter()

	roomList = make(map[string]PoolHelper.Pool)
	newRoom := PoolHelper.NewPool()
	go newRoom.Start()
	roomList["devTest"] = *newRoom
	ErrorHelper.OutputToConsole("Update", "initializing server...")
	setupRoutes(router)
	ErrorHelper.OutputToConsole("Update", "Server online")
	corsObj := handlers.AllowedOrigins([]string{"*"})
	log.Println(http.ListenAndServeTLS(":8080", "./cert/certificate.pem", "./cert/key.pem", handlers.CORS(corsObj)(router)))
}

/*
auth package: name, skin

type: 1 - join
      2 - just new Highscore
      3 - just new pos
      4 - new chat
      5 - all


Highscore: highscore: integer
	==> name will be attached by server based on name attached to connection

NewPlayerPos: {
		name: null (will be included by server)
		x: int
		y: int
		xVel: int
		yVel: int
		isDashing: bool
	} => server appends the name itself based on name attached to connection


ChatMessage : string => if starts with / and user is dev ==> command



so: ClientMessage = {
	type: int
	highscore: int/null
	PlayerPosUpdate: object NewPlayerPos/null
	Chatmessage: string/null
	}
*/
