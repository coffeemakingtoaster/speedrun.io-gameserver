package LobbyHelper

import (
	"encoding/json"
	"math/rand"
	"time"

	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateRoomID() string {
	rand.Seed(time.Now().UnixNano())
	id := make([]rune, 10)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)

}

func GenerateMessage(payload []byte) {
	decodedPayload := ObjectStructures.ReturnMessage{}
	json.Unmarshal(payload, decodedPayload)
}
