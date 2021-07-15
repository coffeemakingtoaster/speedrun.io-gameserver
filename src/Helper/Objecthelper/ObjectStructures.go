package ObjectStructures

import (
	"sync"

	"github.com/gorilla/websocket"
)

//struct use to represent the client along with its most important data
type Client struct {
	PlayerName string
	ID         string
	Conn       *websocket.Conn
	Pool       *Pool
}

//Mutex construct
type ClientStct struct {
	Mu      sync.Mutex
	Clients map[*Client]bool
}

//Highscore mutex construct
type HighScore struct {
	Mu         sync.Mutex
	Highscores map[int]HighScoreStruct
}

//Userstates mutex construct
type UserStates struct {
	Mu         sync.RWMutex
	Userstates map[int]PlayerPosition
}

//A Lobby with all its channels and attributes.
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

//Playerposition. This includes all data a client needs to render a remote player
type PlayerPosition struct {
	Name      string `json:"PlayerName"`
	PosX      int    `json:"y"`
	PosY      int    `json:"x"`
	VelX      int    `json:"yVel"`
	VelY      int    `json:"xVel"`
	IsDashing bool   `json:"isDashing"`
}

//all metadata of a lobby
type LobbyData struct {
	ID        string `json:"LobbyID"`
	MapCode   string `json:"MapID"`
	LobbyName string `json:"Name"`
}

//Package struct for messages received from client
type ClientMessage struct {
	Type        int            `json:"type"`
	LobbyData   LobbyData      `json:"LobbyData"`
	Highscore   int64          `json:"highscore"`
	PlayerPos   PlayerPosition `json:"playerpos"`
	ChatMessage string         `json:"chat"`
}

//Package struct for authpackage received from client
type AuthMessage struct {
	Name  string `json:"name"`
	Token string `json:"token"`
	Skin  string `json:"skinID"`
}

//Highscores
type HighScoreStruct struct {
	PlayerName string
	Time       int64
}

/*
type: 1 lobbydata
type: 2 return Highscores + pos
type: 3 return highscores
type: 4 return pos
type: 5 return chatmessage
*/
//Package struct for messages send to client
type ReturnMessage struct {
	Type        int               `json:"type"`
	LobbyData   LobbyData         `json:"LobbyData"`
	Highscore   []HighScoreStruct `json:"highscore"`
	PlayerPos   []PlayerPosition  `json:"playerpos"`
	ChatMessage string            `json:"chat"`
}
