package SocketHelper

import (
	"net/http"
	"testing"

	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
)

var (
	thttp *testing.T
)

type echoServer struct {
	upgrader websocket.Upgrader
	Done     chan struct{}
	Ready    chan struct{}
}

func TestWsEndpoint(t *testing.T) {

	thttp = t

	var s = &echoServer{}

	//instantiate mockup server
	var d = wstest.NewDialer(s)

	//Connect to mockup server
	DummyConn, _, err := d.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	//close Connection to the server
	DummyConn.Close()

	<-s.Done
}

func TestSender(t *testing.T) {

	thttp = t

	var s = &echoServer{}

	//instantiate mockup server
	var d = wstest.NewDialer(s)

	//Connect to mockup server
	DummyConn, _, err := d.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	var returnValue = ObjectStructures.ReturnMessage{}
	err = DummyConn.ReadJSON(&returnValue)
	if err != nil {
		t.Fatal(err)
	}

	if returnValue.Type != 0 {
		t.Errorf("Invalid type was returned")
	}

	//close Connection to the server
	DummyConn.Close()

	<-s.Done
}

//Mock a httpserver => replace main.go and bypass user authentification
func (s *echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)

	s.Done = make(chan struct{})
	defer close(s.Done)

	conn, err := WsEndpoint(w, r)
	if err != nil {
		thttp.Errorf("Error while upgrading")
	}

	Sender(conn, ObjectStructures.ReturnMessage{Type: 0, LobbyData: (ObjectStructures.LobbyData{}), Highscore: ([]ObjectStructures.HighScoreStruct{}), PlayerPos: ([]ObjectStructures.PlayerPosition{}), ChatMessage: ""})

}
