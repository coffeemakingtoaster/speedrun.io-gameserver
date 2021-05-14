package PoolHelper

import {
	"fmt"
	"log"
	objectStructures "gameserver.speedrun.io/Helper/Objecthelper"
}

func NewPool() *objectStructures.Pool {
	return &Pool{
		Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan Message),
	}
}
