package ObjectStructures

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	PlayerName string
	ID         string
	Conn       *websocket.Conn
	Pool       *Pool
}

type ClientStct struct {
	Mu      sync.Mutex
	Clients map[*Client]bool
}

type HighScore struct {
	Mu         sync.Mutex
	Highscores map[int]HighScoreStruct
}

type UserStates struct {
	Mu         sync.RWMutex
	Userstates map[int]PlayerPosition
}

type Pool struct {
	UserJoin      chan *Client
	UserLeave     chan *Client
	Clients       ClientStct
	Broadcast     chan ReturnMessage
	TimeList      sync.Map
	TimeListSet   chan HighScoreStruct
	UserStateList sync.Map
	UserStateSet  chan PlayerPosition
	KillPool      bool
	LobbyTime     uint64
	LobbyData     LobbyData
}
type PlayerPosition struct {
	Name      string `json:"PlayerName"`
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
	Name  string `json:"name"`
	Token string `json:"token"`
	Skin  string `json:"skinID"`
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
