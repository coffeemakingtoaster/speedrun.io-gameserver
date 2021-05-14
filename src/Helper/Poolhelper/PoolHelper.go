package PoolHelper

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	objectStructures "gameserver.speedrun.io/Helper/Objecthelper"
	SocketHelper "gameserver.speedrun.io/Helper/Sockethelper"
	"github.com/gorilla/websocket"
)

type Pool struct {
	UserJoin    chan *Client
	UserLeave   chan *Client
	Clients     map[*Client]bool
	Broadcast   chan objectStructures.Message
	TimeList    map[int]objectStructures.HighScoreStruct
	TimeListSet chan map[int]objectStructures.HighScoreStruct
}

type Client struct {
	PlayerName string
	ID         string
	Conn       *websocket.Conn
	Pool       *Pool
}

//creates new pool
func NewPool() *Pool {
	return &Pool{
		UserJoin:  make(chan *Client),
		UserLeave: make(chan *Client),
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan objectStructures.Message),
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
				m, err := json.Marshal(objectStructures.Message{Type: 1, Data: []string{"User joined..."}})
				if err != nil {
					log.Fatal("Sending response failed due to an ")
				}
				SocketHelper.Sender(client.Conn, string(m))
			}
			break
		case client := <-pool.UserLeave:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				m, err := json.Marshal(objectStructures.Message{Type: 1, Data: []string{"User Disconnected..."}})
				if err != nil {
					log.Fatal("Sending response failed due to an ")
				}
				SocketHelper.Sender(client.Conn, string(m))
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				m, err := json.Marshal(message)
				if err != nil {
					log.Fatal("Sending response failed due to an ")
				}
				SocketHelper.Sender(client.Conn, string(m))
			}
		case timeList := <-pool.TimeListSet:
			for index, element := range timeList {
				pool.TimeList[index] = element
			}
		}
	}
}

func (c *Client) HandleInput() {
	defer func() {
		c.Pool.UserLeave <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := objectStructures.Message{Type: messageType, Data: GenerateMessage(p, c)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}

func GenerateMessage(payload []byte, c *Client) []string {
	decodedPayload := objectStructures.Message{}
	json.Unmarshal(payload, decodedPayload)
	if decodedPayload.Type == 1 {
		currentTimes := c.Pool.TimeList
		for index, element := range currentTimes {
			if element.PlayerName == c.PlayerName {
				time, err := strconv.Atoi(decodedPayload.Data[0])
				if err != nil {
					log.Fatal("Error! Package contained invalid time")
				}
				currentTimes[index] = objectStructures.HighScoreStruct{
					PlayerName: c.PlayerName,
					Time:       int64(time),
				}
				break
			}
		}
		c.Pool.TimeListSet <- currentTimes
		var returnMessage []string
		for _, element := range currentTimes {
			returnMessage = append(returnMessage, element.PlayerName)
			returnMessage = append(returnMessage, strconv.FormatInt(element.Time, 10))
		}
		return returnMessage
	}
	return []string{}
}
