package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	errorHelper "gameserver.speedrun.io/Helper/Errorhelper"
	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
	SocketHelper "gameserver.speedrun.io/Helper/Sockethelper"
	userHelper "gameserver.speedrun.io/Helper/Userhelper"
)

func handleWebsocketInput(w http.ResponseWriter, r *http.Request) {
	var socketConn, err = SocketHelper.WsEndpoint(w, r)
	if err != nil {
		errorHelper.ConnectionNotWebsocketError(w, r)
		return
	}
	message := SocketHelper.Reader(socketConn)
	fmt.Println("Received message: " + message + " from client")
	parsedRequest := ObjectStructures.RequestObject{}
	err = json.Unmarshal([]byte(message), &parsedRequest)
	if err != nil || parsedRequest.Purpose != "Auth" {
		errorHelper.InvalidRequestError(w, r)
	}
	if userHelper.ValidateUser(parsedRequest.Code) {
		fmt.Println("uID has been validated. Progressing")
		SocketHelper.Sender(socketConn, "Credentials have been confirmed")
		userHelper.InputHandler(socketConn)
	}

}

func setupRoutes() {
	http.HandleFunc("/", errorHelper.InvalidRouteError)
	http.HandleFunc("/ws", handleWebsocketInput)
}

func main() {
	fmt.Println("Server init started")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
