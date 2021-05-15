package SocketHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Something went wrong")
	}
	// helpful log statement to show connections
	log.Println("Client Connected")
	return ws, nil

}

/*
func Reader(conn *websocket.Conn) string {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return ""
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return ""
		}

		return string(p)
	}
}
*/

func Sender(conn *websocket.Conn, payload ObjectStructures.Message) {
	m, err := json.Marshal(payload)
	fmt.Println("Sending out " + string(m))
	err = conn.WriteJSON(payload)
	if err != nil {
		fmt.Println("Error: Cannot send to client")
	}
}
