package PoolHelper

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	ErrorHelper "gameserver.speedrun.io/Helper/Errorhelper"
	LobbyHelper "gameserver.speedrun.io/Helper/Lobbyhelper"
	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
	SocketHelper "gameserver.speedrun.io/Helper/Sockethelper"
	"github.com/gorilla/websocket"
)

type MapPool struct {
	Mu   sync.Mutex
	Maps map[string]ObjectStructures.Pool
}

//creates new pool
func NewPool() *ObjectStructures.Pool {
	return &ObjectStructures.Pool{
		UserJoin:      make(chan *ObjectStructures.Client),
		UserLeave:     make(chan *ObjectStructures.Client),
		Clients:       ObjectStructures.ClientStct{Clients: make(map[*ObjectStructures.Client]bool)},
		Broadcast:     make(chan ObjectStructures.ReturnMessage),
		TimeList:      ObjectStructures.HighScore{Highscores: make(map[int]ObjectStructures.HighScoreStruct)},
		TimeListSet:   make(chan ObjectStructures.HighScoreStruct),
		UserStateList: ObjectStructures.UserStates{Userstates: make(map[int]ObjectStructures.PlayerPosition)},
		UserStateSet:  make(chan ObjectStructures.PlayerPosition),
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
		pool.UserStateList.Mu.Lock()
		var currentPlayers []ObjectStructures.PlayerPosition
		var UserListCopy = pool.UserStateList.Userstates
		pool.UserStateList.Mu.Unlock()
		for _, element := range UserListCopy {
			currentPlayers = append(currentPlayers, element)
		}
		pool.Broadcast <- ObjectStructures.ReturnMessage{Type: 3, LobbyData: (ObjectStructures.LobbyData{}), Highscore: ([]ObjectStructures.HighScoreStruct{}), PlayerPos: currentPlayers, ChatMessage: ""}
		time.Sleep(200 * time.Millisecond)
		steps += 1
	}

}

//handles interaction with the pool
func Start(isPermanent bool, pool *ObjectStructures.Pool) {
	fmt.Println("started Lobby")
	go PoolUpdate(pool, isPermanent)

	for {
		//if loby is empty and not meant to be permanent => close it
		if pool.KillPool {
			break
		}
		select {
		case client := <-pool.UserJoin:
			pool.Clients.Clients[client] = true
			// TODO: There is redundant code here that needs to be removed when refactoring
			var currentHighscores []ObjectStructures.HighScoreStruct
			for _, element := range pool.TimeList.Highscores {
				currentHighscores = append(currentHighscores, element)
			}
			var currentPlayers []ObjectStructures.PlayerPosition
			for _, element := range pool.UserStateList.Userstates {
				currentPlayers = append(currentPlayers, element)
			}
			ErrorHelper.OutputToConsole("Update", "User "+client.PlayerName+" joined")
			SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 4, LobbyData: (ObjectStructures.LobbyData{}), Highscore: currentHighscores, PlayerPos: currentPlayers, ChatMessage: ""})
			fmt.Println("Send data to user")
			for client, _ := range pool.Clients.Clients {
				SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 5, LobbyData: (ObjectStructures.LobbyData{}), Highscore: []ObjectStructures.HighScoreStruct{}, PlayerPos: []ObjectStructures.PlayerPosition{}, ChatMessage: "User joined " + client.PlayerName + "!"})
			}
			break
		case client := <-pool.UserLeave:
			delete(pool.Clients.Clients, client)
			for index, element := range pool.UserStateList.Userstates {
				if element.Name == client.PlayerName {
					//delete player from list
					delete(pool.UserStateList.Userstates, index)
					break
				}
			}
			ErrorHelper.OutputToConsole("Update", "User "+client.PlayerName+" left")
			for c, _ := range pool.Clients.Clients {
				SocketHelper.Sender(c.Conn, ObjectStructures.ReturnMessage{Type: 5, LobbyData: (ObjectStructures.LobbyData{}), Highscore: []ObjectStructures.HighScoreStruct{}, PlayerPos: []ObjectStructures.PlayerPosition{}, ChatMessage: "User left " + client.PlayerName + "!"})
			}
			break
		case message := <-pool.Broadcast:
			for client, _ := range pool.Clients.Clients {
				SocketHelper.Sender(client.Conn, message)
			}
		case userToUpdate := <-pool.TimeListSet:
			foundUser := false
			for index, element := range pool.TimeList.Highscores {
				if element.PlayerName == userToUpdate.PlayerName {
					pool.TimeList.Highscores[index] = userToUpdate
					foundUser = true
					break
				}
			}
			if !foundUser {
				ErrorHelper.OutputToConsole("Warning", "User not found in Highscorelist. Adding user to List...")
				pool.TimeList.Highscores[len(pool.TimeList.Highscores)+1] = userToUpdate
			}
			var currentHighscores []ObjectStructures.HighScoreStruct
			for _, element := range pool.TimeList.Highscores {
				currentHighscores = append(currentHighscores, element)
			}
			fmt.Println("Sending new Highscore list", currentHighscores)
			for client, _ := range pool.Clients.Clients {
				SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 2, LobbyData: (ObjectStructures.LobbyData{}), Highscore: currentHighscores, PlayerPos: []ObjectStructures.PlayerPosition{}, ChatMessage: ""})
			}
			break

		case userToUpdate := <-pool.UserStateSet:
			foundUser := false
			for index, element := range pool.UserStateList.Userstates {
				if element.Name == userToUpdate.Name {
					pool.UserStateList.Userstates[index] = userToUpdate
					foundUser = true
					break
				}
			}
			if !foundUser {
				ErrorHelper.OutputToConsole("Warning", "User not found in Positionlist. Adding user to List...")
				pool.UserStateList.Userstates[len(pool.UserStateList.Userstates)+1] = userToUpdate
			}
			/*
				var returnMessage []ObjectStructures.PlayerPosition
				// TODO: We have to get on and continue so this has to stay for now. BUT IT HAS TO BE REWORKED
				for _, element := range pool.UserStateList {
					returnMessage = append(returnMessage, element)
				}
				for client, _ := range pool.Clients {
					SocketHelper.Sender(client.Conn, ObjectStructures.ReturnMessage{Type: 3, LobbyData: (ObjectStructures.LobbyData{}), Highscore: []ObjectStructures.HighScoreStruct{}, PlayerPos: returnMessage, ChatMessage: ""})
				}
			*/
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
			if decodedPayload.LobbyData.ID == "" {
				CreateRoom(c, poolList)

			} else {
				if roomPool, b := poolList.Maps[decodedPayload.LobbyData.ID]; b {
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
	fmt.Println(decodedPayload.PlayerPos)
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
