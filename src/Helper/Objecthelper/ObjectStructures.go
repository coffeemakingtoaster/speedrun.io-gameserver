package ObjectStructures

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	PlayerName string
	ID         string
	Conn       *websocket.Conn
	Pool       *Pool
}

type Pool struct {
	UserJoin      chan *Client
	UserLeave     chan *Client
	Clients       map[*Client]bool
	Broadcast     chan string
	TimeList      map[int]HighScoreStruct
	TimeListSet   chan HighScoreStruct
	UserStateList map[int]PlayerPosition
	UserStateSet  chan PlayerPosition
}

type PlayerPosition struct {
	Name      string `json:"name"`
	PosX      int    `json:"y"`
	PosY      int    `json:"x"`
	VelX      int    `json:"yVel"`
	VelY      int    `json:"xVel"`
	IsDashing bool   `json:"isDashing"`
}

type LobbyData struct {
	ID        string `json:"LobbyID"`
	MapCode   string `json:"MapID"`
	LobbyName string `json:"Name"`
}

type ClientMessage struct {
	Type        int            `json:"type"`
	LobbyData   LobbyData      `json:"LobbyData"`
	Highscore   int64          `json:"highscore"`
	PlayerPos   PlayerPosition `json:"playerpos"`
	ChatMessage string         `json:"chat"`
}

type AuthMessage struct {
	Name string `json:"name"`
	Skin string `json:"skinID"`
}

type HighScoreStruct struct {
	PlayerName string
	Time       int64
}

type PlayerStats struct {
	PlayerName string `json:"Playername"`
	PositionX  int32  `json:"x"`
	PositionY  int32  `json:"y"`
	VelocityX  int32  `json:"xVel"`
	VelocityY  int32  `json:"yVel"`
	IsDashing  bool   `json:"isDashing"`
}

/*
type: 1 lobbydata
type: 2 return Highscores + pos
type: 3 return highscores
type: 4 return pos
type: 5 return chatmessage
*/
type ReturnMessage struct {
	Type        int               `json:"type"`
	LobbyData   LobbyData         `json:"LobbyData"`
	Highscore   []HighScoreStruct `json:"highscore"`
	PlayerPos   []PlayerPosition  `json:"playerpos"`
	ChatMessage string            `json:"chat"`
}
