package LobbyHelper

import (
	"math/rand"
	"time"
)

//all valid letters for a lobby ID
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

//return lobby ID
func GenerateRoomID() string {
	rand.Seed(time.Now().UnixNano())
	id := make([]rune, 10)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}

//placeholder for feature development
func AlterMap() string {
	return ""
}
