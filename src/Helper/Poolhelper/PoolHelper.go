package PooleHelper

import (
	"encoding/json"
	"fmt"
	"log"

	ErrorHelper "gameserver.speedrun.io/Helper/Errorhelper"
	LobbyHelper "gameserver.speedrun.io/Helper/Lobbyhelper"
	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
	SocketHelper "gameserver.speedrun.io/Helper/Sockethelper"
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
	TimeList      map[int]ObjectStructures.HighScoreStruct
	TimeListSet   chan ObjectStructures.HighScoreStruct
	UserStateList map[int]ObjectStructures.PlayerPosition
	UserStateSet  chan ObjectStructures.PlayerPosition
}

//creates new pool
func NewPool() *Pool {
	return &Pool{
		UserJoin:      make(chan *Client),
		UserLeave:     make(chan *Client),
		Clients:       make(map[*Client]bool),
		Broadcast:     make(chan string),
		TimeList:      make(map[int]ObjectStructures.HighScoreStruct),
		TimeListSet:   make(chan ObjectStructures.HighScoreStruct),
		UserStateList: make(map[int]ObjectStructures.PlayerPosition),
		UserStateSet:  make(chan ObjectStructures.PlayerPosition),
	}
}

//handles interaction with the pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.UserJoin:
			pool.Clients[client] = true
			// TODO: There is redundant code here that needs to be removed when refactoring
			var currentHighscores []ObjectStructures.HighScoreStruct
			for _, element := range pool.TimeList {
				currentHighscores = append(currentHighscores, element)
			}
			var currentPlayers []ObjectStructures.PlayerPosition
			for _, element := range pool.UserStateList {
				currentPlayers = append(currentPlayers, element)
			}
			ErrorHelper.OutputToConsole("Update", "User "+client.PlayerName+" joined")
			SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 2, LobbyData: (ObjectStructures.LobbyData{}), Highscore: currentHighscores, PlayerPos: currentPlayers, ChatMessage: ""})
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 5, LobbyData: (ObjectStructures.LobbyData{}), Highscore: []ObjectStructures.HighScoreStruct{}, PlayerPos: []ObjectStructures.PlayerPosition{}, ChatMessage: "User joined " + client.PlayerName + "!"})
			}
			break
		case client := <-pool.UserLeave:
			delete(pool.Clients, client)
			for index, element := range pool.UserStateList {
				if element.Name == client.PlayerName {
					//delete player from list
					delete(pool.UserStateList, index)
					break
				}
			}
			ErrorHelper.OutputToConsole("Update", "User "+client.PlayerName+" left")
			for c, _ := range pool.Clients {
				SocketHelper.Sender(c.Conn, ObjectStructures.ReturnMessage{Type: 5, LobbyData: (ObjectStructures.LobbyData{}), Highscore: []ObjectStructures.HighScoreStruct{}, PlayerPos: []ObjectStructures.PlayerPosition{}, ChatMessage: "User left " + client.PlayerName + "!"})
			}
			break
		case message := <-pool.Broadcast:
			ErrorHelper.OutputToConsole("Log", "Broadcasting")
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 5, LobbyData: (ObjectStructures.LobbyData{}), Highscore: []ObjectStructures.HighScoreStruct{}, PlayerPos: []ObjectStructures.PlayerPosition{}, ChatMessage: message})
			}
		case userToUpdate := <-pool.TimeListSet:
			foundUser := false
			for index, element := range pool.TimeList {
				if element.PlayerName == userToUpdate.PlayerName {
					pool.TimeList[index] = userToUpdate
					foundUser = true
					break
				}
			}
			if !foundUser {
				ErrorHelper.OutputToConsole("Warning", "User not found in Highscorelist. Adding user to List...")
				pool.TimeList[len(pool.TimeList)+1] = userToUpdate
			}
			var currentHighscores []ObjectStructures.HighScoreStruct
			for _, element := range pool.TimeList {
				currentHighscores = append(currentHighscores, element)
			}
			fmt.Println("Sending new Highscore list")
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 2, LobbyData: (ObjectStructures.LobbyData{}), Highscore: currentHighscores, PlayerPos: []ObjectStructures.PlayerPosition{}, ChatMessage: ""})
			}
			break

		case userToUpdate := <-pool.UserStateSet:
			foundUser := false
			for index, element := range pool.UserStateList {
				if element.Name == userToUpdate.Name {
					pool.UserStateList[index] = userToUpdate
					foundUser = true
					break
				}
			}
			if !foundUser {
				ErrorHelper.OutputToConsole("Warning", "User not found in Positionlist. Adding user to List...")
				pool.UserStateList[len(pool.UserStateList)+1] = userToUpdate
			}
			var returnMessage []ObjectStructures.PlayerPosition
			// TODO: We have to get on and continue so this has to stay for now. BUT IT HAS TO BE REWORKED
			for _, element := range pool.UserStateList {
				returnMessage = append(returnMessage, element)
			}
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 3, LobbyData: (ObjectStructures.LobbyData{}), Highscore: []ObjectStructures.HighScoreStruct{}, PlayerPos: returnMessage, ChatMessage: ""})
			}
			break
		}
	}
}

func (c *Client) HandleInput(poolList map[string]Pool) {

	defer func() {
		c.Pool.UserLeave <- c
		c.Conn.Close()
	}()

	ErrorHelper.OutputToConsole("Update", "Input from client "+c.PlayerName+" is now being handled")
	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		decodedPayload := ObjectStructures.ClientMessage{}
		json.Unmarshal(p, &decodedPayload)
		//if player is not in room
		if c.Pool == nil {
			//create room if no roomID was passed
			if decodedPayload.LobbyData.ID == "" {
				newRoom := NewPool()
				go newRoom.Start()
				c.Pool = *&newRoom
				Id := LobbyHelper.GenerateRoomID()
				for {
					if _, ok := poolList[Id]; ok {
						Id = LobbyHelper.GenerateRoomID()
					} else {
						break
					}
				}
				fmt.Println("RoomID: " + Id)
				poolList[Id] = *newRoom
				ErrorHelper.OutputToConsole("Update", "Created new Room(ID:"+Id+")")
				c.Pool.UserJoin <- c
			} else {
				if roomPool, b := poolList[decodedPayload.LobbyData.ID]; b {
					c.Pool = &roomPool
					c.Pool.UserJoin <- c
				} else {
					ErrorHelper.InvalidRoomIDError(c.Conn)
				}
			}
		} else {
			GenerateMessage(p, c)
		}
	}
}

func GenerateMessage(payload []byte, c *Client) {
	decodedPayload := ObjectStructures.ClientMessage{}
	json.Unmarshal(payload, &decodedPayload)
	fmt.Println(decodedPayload.PlayerPos)
	// 1 is reserver for join
	// 2 => new Highscore
	if decodedPayload.Type == 2 {
		time := decodedPayload.Highscore
		data := ObjectStructures.HighScoreStruct{
			PlayerName: c.PlayerName,
			Time:       time,
		}
		ErrorHelper.OutputToConsole("Update", "User "+c.PlayerName+" send a new personal highscore!")
		c.Pool.TimeListSet <- data
		return
		//new Player Position
	} else if decodedPayload.Type == 3 {
		newPos := decodedPayload.PlayerPos
		newPos.Name = c.PlayerName
		c.Pool.UserStateSet <- newPos
	}
}

func InitInputHandler(conn *websocket.Conn, m map[string]Pool, username string) {
	c := &Client{
		PlayerName: username,
		ID:         "1234",
		Conn:       conn,
		Pool:       nil,
	}
	//SocketHelper.Sender(c.Conn, ObjectStructures.Message{Type: 1, Data: []string{"credentials confirmed"}})
	c.HandleInput(m)
}
