package userHelper

import (
	"fmt"

	"github.com/gorilla/websocket"

	SocketHelper "gameserver.speedrun.io/Helper/Sockethelper"
)

func ValidateUser(uID string) bool {
	fmt.Println("Validated User: " + uID)
	return true
}

func InputHandler(conn *websocket.Conn) {
	message := SocketHelper.Reader(conn)
	fmt.Println(message)
	conn.Close()
	fmt.Println("discornnected from client")
}
