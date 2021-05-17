package PoolHelper

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	ErrorHelper "gameserver.speedrun.io/Helper/Errorhelper"
	LobbyHelper "gameserver.speedrun.io/Helper/Lobbyhelper"
	objectStructures "gameserver.speedrun.io/Helper/Objecthelper"
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
	Broadcast     chan objectStructures.Message
	TimeList      map[int]objectStructures.HighScoreStruct
	TimeListSet   chan objectStructures.HighScoreStruct
	UserStateList map[int]objectStructures.PlayerStats
	UserStateSet  chan objectStructures.PlayerStats
}

//creates new pool
func NewPool() *Pool {
	return &Pool{
		UserJoin:      make(chan *Client),
		UserLeave:     make(chan *Client),
		Clients:       make(map[*Client]bool),
		Broadcast:     make(chan objectStructures.Message),
		TimeList:      make(map[int]objectStructures.HighScoreStruct),
		TimeListSet:   make(chan objectStructures.HighScoreStruct),
		UserStateList: make(map[int]objectStructures.PlayerStats),
		UserStateSet:  make(chan objectStructures.PlayerStats),
	}
}

//handles interaction with the pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.UserJoin:
			pool.Clients[client] = true
			// TODO: There is redundant code here that needs to be removed when refactoring
			var currentHighscores []string
			for _, element := range pool.TimeList {
				currentHighscores = append(currentHighscores, element.PlayerName)
				currentHighscores = append(currentHighscores, strconv.FormatInt(element.Time, 10))
			}
			var returnMessage []string
			for _, element := range pool.UserStateList {
				isDash := "false"
				if element.IsDashing {
					isDash = "true"
				}
				returnMessage = append(returnMessage, element.PlayerName)
				returnMessage = append(returnMessage, strconv.Itoa(element.PositionX))
				returnMessage = append(returnMessage, strconv.Itoa(element.PositionY))
				returnMessage = append(returnMessage, strconv.Itoa(element.VelocityX))
				returnMessage = append(returnMessage, strconv.Itoa(element.VelocityY))
				returnMessage = append(returnMessage, isDash)
			}
			ErrorHelper.OutputToConsole("Update", "User "+client.PlayerName+" joined")
			SocketHelper.Sender(client.Conn, objectStructures.Message{Type: 3, Data: returnMessage})
			SocketHelper.Sender(client.Conn, objectStructures.Message{Type: 2, Data: currentHighscores})
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, objectStructures.Message{Type: 1, Data: []string{"User joined..."}})
			}
			break
		case client := <-pool.UserLeave:
			delete(pool.Clients, client)
			for index, element := range pool.UserStateList {
				if element.PlayerName == client.PlayerName {
					//delete player from list
					delete(pool.UserStateList, index)
					break
				}
			}
			ErrorHelper.OutputToConsole("Update", "User "+client.PlayerName+" left")
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, objectStructures.Message{Type: 1, Data: []string{"User " + client.PlayerName + " disconnected..."}})
			}
			break
		case message := <-pool.Broadcast:
			ErrorHelper.OutputToConsole("Log", "Broadcasting")
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, message)
			}
		case userToUpdate := <-pool.TimeListSet:
			foundUser := false
			/*
				for index, element := range pool.TimeList {
					pool.TimeList[index] = element
				}
			*/
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
			var returnMessage []string
			for _, element := range pool.TimeList {
				returnMessage = append(returnMessage, element.PlayerName)
				returnMessage = append(returnMessage, strconv.FormatInt(element.Time, 10))
			}
			fmt.Println("Sending new Highscore list")
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, objectStructures.Message{Type: 2, Data: returnMessage})
			}
			fmt.Println(returnMessage)
			break

		case userToUpdate := <-pool.UserStateSet:
			foundUser := false
			/*
				for index, element := range pool.UserStateList {
					pool.TimeList[index] = element
				}
			*/
			for index, element := range pool.UserStateList {
				if element.PlayerName == userToUpdate.PlayerName {
					pool.UserStateList[index] = userToUpdate
					foundUser = true
					break
				}
			}
			if !foundUser {
				ErrorHelper.OutputToConsole("Warning", "User not found in Positionlist. Adding user to List...")
				pool.UserStateList[len(pool.UserStateList)+1] = userToUpdate
			}
			var returnMessage []string
			// TODO: We have to get on and continue so this has to stay for now. BUT IT HAS TO BE REWORKED
			for _, element := range pool.UserStateList {
				isDash := "false"
				if element.IsDashing {
					isDash = "true"
				}
				returnMessage = append(returnMessage, element.PlayerName)
				returnMessage = append(returnMessage, strconv.Itoa(element.PositionX))
				returnMessage = append(returnMessage, strconv.Itoa(element.PositionY))
				returnMessage = append(returnMessage, strconv.Itoa(element.VelocityX))
				returnMessage = append(returnMessage, strconv.Itoa(element.VelocityY))
				returnMessage = append(returnMessage, isDash)
			}
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, objectStructures.Message{Type: 3, Data: returnMessage})
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
		decodedPayload := objectStructures.Message{}
		json.Unmarshal(p, &decodedPayload)
		//if player is not in room
		if c.Pool == nil {
			//create room if no roomID was passed
			if decodedPayload.Data == nil {
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
				if roomPool, b := poolList[decodedPayload.Data[0]]; b {
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
	decodedPayload := objectStructures.Message{}
	json.Unmarshal(payload, &decodedPayload)
	if decodedPayload.Type == 1 {
		time, err := strconv.Atoi(decodedPayload.Data[0])
		if err != nil {
			return
		}
		data := objectStructures.HighScoreStruct{
			PlayerName: c.PlayerName,
			Time:       int64(time),
		}
		ErrorHelper.OutputToConsole("Update", "User "+c.PlayerName+" send a new personal highscore!")
		c.Pool.TimeListSet <- data
		return
	} else if decodedPayload.Type == 2 {
		x, _ := strconv.Atoi(decodedPayload.Data[1])
		y, _ := strconv.Atoi(decodedPayload.Data[2])
		Velx, _ := strconv.Atoi(decodedPayload.Data[3])
		Vely, _ := strconv.Atoi(decodedPayload.Data[4])
		isDashing := false
		if decodedPayload.Data[4] == "true" {
			isDashing = true
		}
		decodedPlayerState := objectStructures.PlayerStats{
			PlayerName: decodedPayload.Data[0],
			PositionX:  x,
			PositionY:  y,
			VelocityX:  Velx,
			VelocityY:  Vely,
			IsDashing:  isDashing,
		}
		json.Unmarshal([]byte(decodedPayload.Data[0]), &decodedPlayerState)
		c.Pool.UserStateSet <- decodedPlayerState
	}
}

func InitInputHandler(conn *websocket.Conn, m map[string]Pool, username string) {
	c := &Client{
		PlayerName: username,
		ID:         "1234",
		Conn:       conn,
		Pool:       nil,
	}
	SocketHelper.Sender(c.Conn, objectStructures.Message{Type: 1, Data: []string{"credentials confirmed"}})
	c.HandleInput(m)
}
