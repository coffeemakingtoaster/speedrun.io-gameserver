package userHelper

import (
	"fmt"

	"github.com/gorilla/websocket"

	PoolHelper "gameserver.speedrun.io/Helper/Poolhelper"
)

func ValidateUser(uID string) bool {
	fmt.Println("Validated User: " + uID)
	return true
}

func InitInputHandler(conn *websocket.Conn) {
	c := &PoolHelper.Client{
		Conn: conn,
		Pool: nil,
	}
	c.HandleInput()
}
