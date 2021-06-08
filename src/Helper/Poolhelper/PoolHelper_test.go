package PoolHelper

import (
	"fmt"
	"net/http"
	"testing"

	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
)

var (
	wsUpgrader   = websocket.Upgrader{}
	DummyConn    *websocket.Conn
	mockLobbyMap = MapPool{Maps: make(map[string]ObjectStructures.Pool)}
	DummyLobby   ObjectStructures.Pool
)

type echoServer struct {
	upgrader websocket.Upgrader
	Done     chan struct{}
	Ready    chan struct{}
}

type goRoutine struct {
	Done chan struct{}
}

func TestPoolJoin(t *testing.T) {

	var s = &echoServer{}

	//instantiate mockup server
	var d = wstest.NewDialer(s)

	//Connect to mockup server
	DummyConn, _, err := d.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	//dummy join request
	payload := ObjectStructures.ClientMessage{Type: 1, LobbyData: ObjectStructures.LobbyData{ID: "mock", MapCode: "", LobbyName: ""}, Highscore: 0, PlayerPos: ObjectStructures.PlayerPosition{}, ChatMessage: ""}
	DummyConn.WriteJSON(payload)

	var returnValue = ObjectStructures.ReturnMessage{}
	err = DummyConn.ReadJSON(&returnValue)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	//close Connection to the server
	DummyConn.Close()

	//check if response has correct type
	if returnValue.Type != 4 {
		t.Errorf("Join room failed")
	}

	<-s.Done
	DummyLobby.KillPool = true

}

func TestPoolCreate(t *testing.T) {

	var s = &echoServer{}

	//instantiate mockup server
	var d = wstest.NewDialer(s)

	//Connect to mockup server
	DummyConn, _, err := d.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	//dummy join request
	payload := ObjectStructures.ClientMessage{Type: 1, LobbyData: ObjectStructures.LobbyData{ID: "", MapCode: "", LobbyName: ""}, Highscore: 0, PlayerPos: ObjectStructures.PlayerPosition{}, ChatMessage: ""}
	DummyConn.WriteJSON(payload)

	fmt.Println("send message")

	//two blocks because client sends 2 packages when a user joins
	var returnValue = ObjectStructures.ReturnMessage{}
	err = DummyConn.ReadJSON(&returnValue)
	fmt.Println("received message")
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	//close Connection to the server
	DummyConn.Close()

	if returnValue.Type != 4 {
		t.Errorf("Incorrect response code")
	}

	<-s.Done
	DummyLobby.KillPool = true

	fmt.Println("second test success")
}

func TestClientInput(t *testing.T) {

	var s = &echoServer{}

	//instantiate mockup server
	var d = wstest.NewDialer(s)

	//Connect to mockup server
	DummyConn, _, err := d.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	//dummy join request
	payload := ObjectStructures.ClientMessage{Type: 1, LobbyData: ObjectStructures.LobbyData{ID: "", MapCode: "", LobbyName: ""}, Highscore: 0, PlayerPos: ObjectStructures.PlayerPosition{}, ChatMessage: ""}
	DummyConn.WriteJSON(payload)

	//2 blocks because the server sends 2 packages on client join
	var returnValue = ObjectStructures.ReturnMessage{}
	err = DummyConn.ReadJSON(&returnValue)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	returnValue = ObjectStructures.ReturnMessage{}
	err = DummyConn.ReadJSON(&returnValue)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	//send new Highscore
	expectedTime := 666
	payload = ObjectStructures.ClientMessage{Type: 2, LobbyData: ObjectStructures.LobbyData{ID: "", MapCode: "", LobbyName: ""}, Highscore: int64(expectedTime), PlayerPos: ObjectStructures.PlayerPosition{}, ChatMessage: ""}
	DummyConn.WriteJSON(payload)

	//catch update blocks before the highscore changed
	for len(returnValue.Highscore) == 0 {
		err = DummyConn.ReadJSON(&returnValue)
	}
	fmt.Println("time", returnValue)
	if int(returnValue.Highscore[0].Time) != expectedTime {
		t.Errorf("received time differs from expected value")
	}

	fmt.Println("Highscore matched")

	//send new Position
	expectedPos := ObjectStructures.PlayerPosition{Name: "mock", PosX: 1, PosY: 1, VelX: 1, VelY: 1, IsDashing: false}
	payload = ObjectStructures.ClientMessage{Type: 3, LobbyData: ObjectStructures.LobbyData{ID: "", MapCode: "", LobbyName: ""}, Highscore: 0, PlayerPos: expectedPos, ChatMessage: ""}
	DummyConn.WriteJSON(payload)

	//catch update blocks before the PlayerPos changed
	for len(returnValue.PlayerPos) == 0 {
		err = DummyConn.ReadJSON(&returnValue)
	}
	fmt.Println("pos", returnValue)
	if returnValue.PlayerPos[0] != expectedPos {
		t.Errorf("received position differs from expected value")
	}

	fmt.Println("Playerpos matched")

	DummyConn.Close()
	<-s.Done
	DummyLobby.KillPool = true

	fmt.Println("third test success")
}

//Mock a httpserver => replace main.go and bypass user authentification
func (s *echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)

	s.Done = make(chan struct{})
	defer close(s.Done)

	var mockConn *websocket.Conn
	mockConn, err = s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	//this is usually provided by the main class and generated at programmstart
	mockLobbyMap := MapPool{Maps: make(map[string]ObjectStructures.Pool)}
	DummyLobby := NewPool()
	go Start(true, DummyLobby)              //true because otherwise Lobby would be instantly closed due to no connected clients
	mockLobbyMap.Maps["mock"] = *DummyLobby //add Dummylobby

	InitInputHandler(mockConn, mockLobbyMap, "mock")

}
