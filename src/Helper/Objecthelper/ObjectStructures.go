package ObjectStructures

import (
	"github.com/gorilla/websocket"
)

type RequestObject struct {
	Purpose string `json:"purpose"`
	Code    string `json:"code"`
}

type client struct {
	PlayerName string
	ID         string
	Conn       *websocket.Conn
	Pool       *Pool
}

type Pool struct {
	UserJoin  chan *client
	UserLeave chan *client
	Clients   map[*client]bool
	Broadcast chan Message
}

type Message struct {
	Type int      `json:"type"`
	Data []string `json:"data"`
}

type HighScoreStruct struct {
	PlayerName string
	Time       int64
}

type PlayerStats struct {
	PlayerName string `json:"Playername"`
	PositionX  int    `json:"x"`
	PositionY  int    `json:"y"`
	VelocityX  int    `json:"xVel"`
	VelocityY  int    `json:"yVel"`
	IsDashing  bool   `json:"isDashing"`
}
