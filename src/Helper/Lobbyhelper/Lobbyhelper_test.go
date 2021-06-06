package LobbyHelper

import (
	"testing"
)

func TestGenerateRoomID(t *testing.T) {
	var newID interface{} = GenerateRoomID()
	var newID2 interface{} = GenerateRoomID()
	_, isString := newID.(string)
	_, isString2 := newID2.(string)
	if !isString || !isString2 {
		t.Errorf("LobbyIDs don´t get generated as strings")
	}

	if len(GenerateRoomID()) != 10 {
		t.Errorf("LobbyIDs have an incorrect length")
	}
}
