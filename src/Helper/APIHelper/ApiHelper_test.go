package ApiHelper

/*
func TestLobbyLifecycle(t *testing.T) {
	ReportLobby(ObjectStructures.LobbyData{
		ID:        "unittest",
		MapCode:   "nothing",
		LobbyName: "mock value",
	})
	resp, err := http.Get("https://api.speedrun.io/v1/lobbies/unittest")
	if err != nil {
		t.Error("Cannt connect to api")
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 || resp.StatusCode == 409 {
		t.Error("Lobby either not deleted properly last time or not found")
	}
	CloseLobby(ObjectStructures.LobbyData{
		ID:        "unittest",
		MapCode:   "nothing",
		LobbyName: "mock value",
	})
	resp, err = http.Get("https://api.speedrun.io/v1/lobbies/unittest")
	if err != nil {
		t.Error("Cannt connect to api")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Error("Lobby delete not working properly")
	}
}
*/
