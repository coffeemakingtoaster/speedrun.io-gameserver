package SocketHelper

import (
	"errors"
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
		return nil, errors.New("Something went wrong")
	}
	// helpful log statement to show connections
	return ws, nil

}

func Sender(conn *websocket.Conn, message ObjectStructures.ReturnMessage) error {
	//m, err := json.Marshal(payload)
	err := conn.WriteJSON(message)
	if err != nil {
		return err
	}
	return nil
}
