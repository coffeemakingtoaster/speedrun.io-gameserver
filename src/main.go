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
	parsedRequest := ObjectStructures.Message{}
	err = json.Unmarshal([]byte(message), &parsedRequest)
	if err != nil || parsedRequest.Type != 0 {
		ErrorHelper.InvalidRequestError(w, r)
	}
	//if validuser continue connection, else close socket
	if userHelper.ValidateUser(parsedRequest.Data[0]) {
		ErrorHelper.OutputToConsole("Update", " valid Player with name "+parsedRequest.Data[0]+"has connected ")
		PoolHelper.InitInputHandler(socketConn, roomList, parsedRequest.Data[0])
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
