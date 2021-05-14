package LobbyHelper

import (
	"encoding/json"

	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
)

func CreateLobby() {

}

func GenerateMessage(payload []byte) {
	decodedPayload := ObjectStructures.Message{}
	json.Unmarshal(payload, decodedPayload)
}
