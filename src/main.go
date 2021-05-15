package main

//main package
//for the webserver
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		fmt.Println(err)
	} else {
		message = string(p)
	}
	fmt.Println("Received message: " + message + " from client")
	parsedRequest := ObjectStructures.RequestObject{}
	err = json.Unmarshal([]byte(message), &parsedRequest)
	if err != nil || parsedRequest.Purpose != "Auth" {
		ErrorHelper.InvalidRequestError(w, r)
	}
	//if validuser continue connection, else close socket
	if userHelper.ValidateUser(parsedRequest.Code) {
		fmt.Println("uID has been validated. Progressing")
		PoolHelper.InitInputHandler(socketConn, roomList)
	} else {
		socketConn.Close()
	}

}

func setupRoutes() {
	http.HandleFunc("/", ErrorHelper.InvalidRouteError)
	http.HandleFunc("/ws", handleWebsocketInput)
}

func main() {
	roomList = make(map[string]PoolHelper.Pool)
	fmt.Println("Server init started")
	setupRoutes()
	log.Println(http.ListenAndServe(":8080", nil))
}
