package ApiHelper

import (
	"testing"

	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
)

// Do is the mock client's `Do` func

func TestReportLobby(t *testing.T) {
	setApiURL("")
	data := ReportLobby(ObjectStructures.LobbyData{
		ID:        "unittest",
		MapCode:   "nothing",
		LobbyName: "mock value",
	})
	if data.ID != "unittest" {
		t.Errorf("ID does not match")
	}
	if data.MapCode != "nothing" {
		t.Errorf("Mapcode does not match")
	}
	if data.Region != "EUW" {
		t.Errorf("Region does not match")
	}
	if data.LobbyName != "mock value" {
		t.Errorf("Lobbyname does not match")
	}
	if data.PlayerCount != 0 {
		t.Errorf("Playercount has invalid number at this state")
	}
}

func TestCloseLobby(t *testing.T) {
	setApiURL("")
	closeUrl := CloseLobby(ObjectStructures.LobbyData{ID: "mock"})
	if closeUrl != "/v1/lobbies/mock" {
		t.Errorf("Delete sends to wrong url")
	}
}

func TestReportClientChange(t *testing.T) {
	setApiURL("")
	data := ReportClientChange(5, ObjectStructures.LobbyData{
		ID:        "unittest",
		MapCode:   "nothing",
		LobbyName: "mock value",
	})
	if data.ID != "unittest" {
		t.Errorf("ID does not match")
	}
	if data.MapCode != "nothing" {
		t.Errorf("Mapcode does not match")
	}
	if data.PlayerCount != 5 {
		t.Errorf("Playercount has invalid number at this state")
	}
}

func TestReportMapChange(t *testing.T) {
	setApiURL("")
	data := ReportMapChange(ObjectStructures.LobbyData{
		ID:        "unittest",
		MapCode:   "also nothing",
		LobbyName: "mock value",
	})
	if data.ID != "unittest" {
		t.Errorf("ID does not match")
	}
	if data.MapCode != "also nothing" {
		t.Errorf("Mapcode does not match")
	}
}
