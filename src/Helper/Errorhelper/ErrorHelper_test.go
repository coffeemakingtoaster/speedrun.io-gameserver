package ErrorHelper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidRouteError(t *testing.T) {
	//Request not going to /ws should not be accepted
	expected := "Invalid Route. If you are trying to reach the game API please interact with api.speedrun.io"
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	InvalidRouteError(res, req)

	// is response code correct?
	if res.Code != http.StatusBadRequest {
		t.Log(res.Code)
		t.Errorf("Http connections don´t get the appropriate response Code")
	}

	//parse Body
	msgbdy := string(res.Body.Bytes())
	msgbdy = strings.ReplaceAll(msgbdy, "Bad Request", "")
	msgbdy = strings.TrimSpace(msgbdy)

	//is response message correct?
	if msgbdy != expected {
		t.Log(msgbdy)
		t.Errorf("Http connections don´t get the appropriate response")
	}
}

func TestConnectionNotWebsocketError(t *testing.T) {
	//Http request should not be accepted
	//http request
	expected := "Error: Connection to the /ws part of the gameserver should only be via websockets"
	req, err := http.NewRequest("GET", "/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	ConnectionNotWebsocketError(res, req)

	// is response code correct?
	if res.Code != http.StatusBadRequest {
		t.Errorf("Http connections don´t get the appropriate response Code")
	}

	//parse Body
	msgbdy := string(res.Body.Bytes())
	msgbdy = strings.ReplaceAll(msgbdy, "Bad Request", "")
	msgbdy = strings.TrimSpace(msgbdy)
	//is response message correct?
	if msgbdy != expected {
		t.Errorf("Http connections don´t get the appropriate response")
	}
}

func TestInvalidRequestError(t *testing.T) {
	//Http request should not be accepted
	//http request
	expected := "Error: The request is invalid in this context"
	req, err := http.NewRequest("GET", "/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	InvalidRequestError(res, req)

	// is response code correct?
	if res.Code != http.StatusBadRequest {
		t.Errorf("Http connections don´t get the appropriate response Code")
	}

	//parse Body
	msgbdy := string(res.Body.Bytes())
	msgbdy = strings.ReplaceAll(msgbdy, "Bad Request", "")
	msgbdy = strings.TrimSpace(msgbdy)
	//is response message correct?
	if msgbdy != expected {
		t.Errorf("Http connections don´t get the appropriate response")
	}
}

//Testing OutputToConsole is pointless...
