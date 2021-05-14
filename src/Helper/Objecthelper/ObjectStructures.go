package ObjectStructures

import {
	"github.com/gorilla/websocket"
}

type RequestObject struct {
	Purpose string `json:"purpose"`
	Code    string `json:"code"`
}

type Client struct {
    ID   string
    Conn *websocket.Conn
    Pool *Pool
}

type Pool struct {
	Register chan *Client
	Unregister chan *Client
	Clients map[*Client]bool
	Broadcast chan Message
}
