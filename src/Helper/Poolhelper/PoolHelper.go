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
	UserJoin    chan *Client
	UserLeave   chan *Client
	Clients     map[*Client]bool
	Broadcast   chan objectStructures.Message
	TimeList    map[int]objectStructures.HighScoreStruct
	TimeListSet chan objectStructures.HighScoreStruct
}

//creates new pool
func NewPool() *Pool {
	return &Pool{
		UserJoin:    make(chan *Client),
		UserLeave:   make(chan *Client),
		Clients:     make(map[*Client]bool),
		Broadcast:   make(chan objectStructures.Message),
		TimeList:    make(map[int]objectStructures.HighScoreStruct),
		TimeListSet: make(chan objectStructures.HighScoreStruct),
	}
}

//handles interaction with the pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.UserJoin:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, objectStructures.Message{Type: 1, Data: []string{"User joined..."}})
			}
			break
		case client := <-pool.UserLeave:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, objectStructures.Message{Type: 1, Data: []string{"User Disconnected..."}})
			}
			break
		case message := <-pool.Broadcast:
			for client, _ := range pool.Clients {
				fmt.Println("Sending message in loop")
				SocketHelper.Sender(client.Conn, message)
			}
		case userToUpdate := <-pool.TimeListSet:
			fmt.Println("Updating Pool")
			foundUser := false
			for index, element := range pool.TimeList {
				pool.TimeList[index] = element
			}
			for index, element := range pool.TimeList {
				if element.PlayerName == userToUpdate.PlayerName {
					pool.TimeList[index] = userToUpdate
					foundUser = true
					break
				}
			}
			if !foundUser {
				fmt.Println("User not found in Highscore List. Adding Player to the List")
				pool.TimeList[len(pool.TimeList)+1] = userToUpdate
			}
			fmt.Println(pool.TimeList)
			var returnMessage []string
			for _, element := range pool.TimeList {
				returnMessage = append(returnMessage, element.PlayerName)
				returnMessage = append(returnMessage, strconv.FormatInt(element.Time, 10))
			}
			fmt.Println(returnMessage)
			for client, _ := range pool.Clients {
				SocketHelper.Sender(client.Conn, objectStructures.Message{Type: 1, Data: returnMessage})
			}
			break
		}
	}
}

func (c *Client) HandleInput(poolList map[string]Pool) {

	defer func() {
		log.Println("user left")
		c.Pool.UserLeave <- c
		c.Conn.Close()
	}()

	for {
		fmt.Println("inputhandler running")
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))
		decodedPayload := objectStructures.Message{}
		json.Unmarshal(p, &decodedPayload)
		fmt.Println(decodedPayload)
		//if player is not in room
		if c.Pool == nil {
			fmt.Println("Player is not in room")
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
				fmt.Println("created new room")
				c.Pool.UserJoin <- c
			} else {
				fmt.Println("joining Room")
				if roomPool, b := poolList[decodedPayload.Data[0]]; b {
					c.Pool = &roomPool
					c.Pool.UserJoin <- c
					fmt.Println("joined room")
				} else {
					ErrorHelper.InvalidRoomIDError(c.Conn)
				}
			}
		} else {
			message := objectStructures.Message{Type: messageType, Data: GenerateMessage(p, c)}
			c.Pool.Broadcast <- message
			fmt.Println("Send broadcast")
		}
	}
}

func GenerateMessage(payload []byte, c *Client) []string {
	decodedPayload := objectStructures.Message{}
	json.Unmarshal(payload, &decodedPayload)
	if decodedPayload.Type == 1 {
		fmt.Println("Updating Highscorelist")
		time, err := strconv.Atoi(decodedPayload.Data[0])
		if err != nil {
			return nil
		}
		data := objectStructures.HighScoreStruct{
			PlayerName: "proplayer123",
			Time:       int64(time),
		}
		c.Pool.TimeListSet <- data

		return []string{"we received your highscore and don´t really care"}
	}
	return []string{"we received your message and don´t really care"}
}

func InitInputHandler(conn *websocket.Conn, m map[string]Pool) {
	c := &Client{
		PlayerName: "testuser",
		ID:         "1234",
		Conn:       conn,
		Pool:       nil,
	}
	SocketHelper.Sender(c.Conn, objectStructures.Message{Type: 1, Data: []string{"credentials confirmed"}})
	c.HandleInput(m)
}
