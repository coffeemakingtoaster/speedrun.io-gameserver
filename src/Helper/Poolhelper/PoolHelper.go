package PoolHelper

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	ApiHelper "gameserver.speedrun.io/Helper/APIHelper"
	ErrorHelper "gameserver.speedrun.io/Helper/Errorhelper"
	LobbyHelper "gameserver.speedrun.io/Helper/Lobbyhelper"
	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
	SocketHelper "gameserver.speedrun.io/Helper/Sockethelper"
	"github.com/gorilla/websocket"
)

type MapPool struct {
	Mu   sync.RWMutex
	Maps map[string]ObjectStructures.Pool
}

//creates new pool
func NewPool() *ObjectStructures.Pool {
	return &ObjectStructures.Pool{
		UserJoin:     make(chan *ObjectStructures.Client),
		UserLeave:    make(chan *ObjectStructures.Client),
		Clients:      ObjectStructures.ClientStct{Clients: make(map[*ObjectStructures.Client]bool)},
		Broadcast:    make(chan ObjectStructures.ReturnMessage),
		TimeListSet:  make(chan ObjectStructures.HighScoreStruct),
		UserStateSet: make(chan ObjectStructures.PlayerPosition),
	}
}

func PoolUpdate(pool *ObjectStructures.Pool, isPermanent bool) {
	steps := 0
	for true {
		//if pool is closed end goroutine
		if pool.KillPool {
			break
		}
		//check every minute if Lobby is empty and delete if true
		if steps >= 300 {
			if len(pool.Clients.Clients) <= 0 && !isPermanent {
				pool.KillPool = true
				break
			}
			steps = 0
		}

		//if lobby has run out -> mapchange
		/*
			if uint64(time.Now().Second()) >= pool.LobbyTime {
				pool.LobbyData.MapCode = LobbyHelper.AlterMap()
				pool.Broadcast <- ObjectStructures.ReturnMessage{Type: 1, LobbyData: pool.LobbyData, Highscore: ([]ObjectStructures.HighScoreStruct{}), PlayerPos: ([]ObjectStructures.PlayerPosition{}), ChatMessage: ""}
				ApiHelper.ReportLobbyChange(pool.LobbyData)
			}
		*/
		pool.Broadcast <- createPositionList(pool)
		time.Sleep(200 * time.Millisecond)
		steps += 1
	}

}

func createPositionList(pool *ObjectStructures.Pool) ObjectStructures.ReturnMessage {
	var currentPlayers []ObjectStructures.PlayerPosition
	pool.UserStateList.Range(func(key, value interface{}) bool {
		currentPlayers = append(currentPlayers, value.(ObjectStructures.PlayerPosition))
		return true
	})
	return ObjectStructures.ReturnMessage{Type: 3, PlayerPos: currentPlayers, ChatMessage: ""}
}

//handles interaction with the pool
func Start(isPermanent bool, pool *ObjectStructures.Pool) {
	ApiHelper.ReportLobby(pool.LobbyData)
	pool.LobbyTime = uint64(time.Now().Second()) + 600
	go PoolUpdate(pool, isPermanent)
	for {
		//if loby is empty and not meant to be permanent => close it
		if pool.KillPool {
			ApiHelper.CloseLobby(pool.LobbyData)
			break
		}
		select {
		case client := <-pool.UserJoin:
			pool.Clients.Clients[client] = true
			// TODO: There is redundant code here that needs to be removed when refactoring
			var currentHighscores []ObjectStructures.HighScoreStruct
			pool.TimeList.Range(func(key, value interface{}) bool {
				currentHighscores = append(currentHighscores, value.(ObjectStructures.HighScoreStruct))
				return true
			})
			var currentPlayers []ObjectStructures.PlayerPosition
			pool.UserStateList.Range(func(key, value interface{}) bool {
				currentPlayers = append(currentPlayers, value.(ObjectStructures.PlayerPosition))
				return true
			})
			ErrorHelper.OutputToConsole("Update", "User "+client.PlayerName+" joined")
			SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 4, Highscore: currentHighscores, PlayerPos: currentPlayers})
			for client, _ := range pool.Clients.Clients {
				SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 5, ChatMessage: "User joined " + client.PlayerName + "!"})
			}
			break
		case client := <-pool.UserLeave:
			delete(pool.Clients.Clients, client)
			pool.UserStateList.Delete(client.PlayerName)
			ErrorHelper.OutputToConsole("Update", "User "+client.PlayerName+" left")
			for c, _ := range pool.Clients.Clients {
				SocketHelper.Sender(c.Conn, ObjectStructures.ReturnMessage{Type: 5, ChatMessage: "User left " + client.PlayerName + "!"})
			}
			break
		case message := <-pool.Broadcast:
			for client, _ := range pool.Clients.Clients {
				SocketHelper.Sender(client.Conn, message)
			}
		case userToUpdate := <-pool.TimeListSet:
			pool.TimeList.Store(userToUpdate.PlayerName, userToUpdate)
			var currentHighscores []ObjectStructures.HighScoreStruct
			pool.TimeList.Range(func(key, value interface{}) bool {
				currentHighscores = append(currentHighscores, value.(ObjectStructures.HighScoreStruct))
				return true
			})
			for client, _ := range pool.Clients.Clients {
				SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 2, Highscore: currentHighscores})
			}
			break
		case userToUpdate := <-pool.UserStateSet:
			pool.UserStateList.Store(userToUpdate.Name, userToUpdate)
			break
		}
	}
}

func HandleInput(poolList MapPool, c *ObjectStructures.Client) {

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
			fmt.Println(decodedPayload)
			if decodedPayload.LobbyData.ID == "" {
				CreateRoom(c, poolList)

			} else {
				if roomPool, b := poolList.Maps[decodedPayload.LobbyData.ID]; b {
					if roomPool.KillPool == true {
						delete(poolList.Maps, decodedPayload.LobbyData.ID)
						break
					}
					fmt.Println("User joined")
					c.Pool = &roomPool
					c.Pool.UserJoin <- c
				}
			}
		} else {
			GenerateMessage(p, c)
		}
	}
}

func CreateRoom(c *ObjectStructures.Client, poolList MapPool) {
	poolList.Mu.Lock()
	//as this locks the Map I will also check for deleted lobbies here
	for lobby := range poolList.Maps {
		if poolList.Maps[lobby].KillPool {
			delete(poolList.Maps, lobby)
		}
	}
	newRoom := NewPool()
	go Start(false, newRoom)
	c.Pool = *&newRoom
	Id := LobbyHelper.GenerateRoomID()
	for {
		if _, ok := poolList.Maps[Id]; ok {
			Id = LobbyHelper.GenerateRoomID()
		} else {
			break
		}
	}
	fmt.Println("RoomID: " + Id)
	poolList.Maps[Id] = *newRoom
	ErrorHelper.OutputToConsole("Update", "Created new Room(ID:"+Id+")")
	c.Pool.UserJoin <- c
	defer poolList.Mu.Unlock()
}

func GenerateMessage(payload []byte, c *ObjectStructures.Client) {
	decodedPayload := ObjectStructures.ClientMessage{}
	json.Unmarshal(payload, &decodedPayload)
	//fmt.Println(decodedPayload.PlayerPos)
	// 1 is reserver for join
	// 2 => new Highscore
	// 3 => new Player Position
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

func InitInputHandler(conn *websocket.Conn, m MapPool, username string) {

	//TODO set UserID => check if it is ever needed
	c := &ObjectStructures.Client{
		PlayerName: username,
		ID:         "1234",
		Conn:       conn,
		Pool:       nil,
	}
	//SocketHelper.Sender(c.Conn, ObjectStructures.Message{Type: 1, Data: []string{"credentials confirmed"}})
	HandleInput(m, c)
}
