package ApiHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	ErrorHelper "gameserver.speedrun.io/Helper/Errorhelper"
	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
)

var apiUrl = "https://api.speedrun.io"

type LobbyReport struct {
	ID             string `json:"lobbyCode"`
	MapCode        string `json:"mapSlug"`
	LobbyName      string `json:"Name"`
	IP             string `json:"ip"`
	Region         string `json:"region"`
	MaxPlayerCount int    `json:"maxPlayerCount"`
	PlayerCount    int    `json:"playerCount"`
}

func GetRandomMapFromApi() (string, error) {
	_, err := http.Get("https://api.speedrun.io")
	if err != nil {
		return "", err
	}
	return "resp.Request.Body.Read()", nil
}

func ReportClientChange(playerCount int, lobby ObjectStructures.LobbyData) {
	ReportLobbyChange(LobbyReport{PlayerCount: playerCount, ID: lobby.ID, MapCode: lobby.MapCode, IP: getIP(), MaxPlayerCount: 69})
}

func ReportMapChange(lobby ObjectStructures.LobbyData) {
	ReportLobbyChange(LobbyReport{MapCode: lobby.MapCode, ID: lobby.ID, IP: getIP(), MaxPlayerCount: 69})
}

func ReportLobbyChange(data LobbyReport) {
	ip := getIP()
	if ip == "" {
		ErrorHelper.OutputToConsole("Error", "No valid local IP found")
	}

	requestData, err := json.Marshal(data)
	fmt.Println(string(requestData))
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}

	req, err := http.NewRequest("PATCH", apiUrl+"/v1/lobbies/"+data.ID, bytes.NewBuffer(requestData))
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}

	doRequest(req)
}

func ReportLobby(lobby ObjectStructures.LobbyData) {

	ip := getIP()
	if ip == "" {
		ErrorHelper.OutputToConsole("Error", "No valid local IP found")
	}
	data := LobbyReport{
		ID:             lobby.ID,
		MapCode:        lobby.MapCode,
		LobbyName:      lobby.LobbyName,
		IP:             ip,
		Region:         "EUW",
		MaxPlayerCount: 69,
		PlayerCount:    0,
	}

	requestData, err := json.Marshal(data)
	fmt.Println(string(requestData))
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}

	req, err := http.NewRequest("POST", apiUrl+"/v1/lobbies", bytes.NewBuffer(requestData))
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}

	doRequest(req)

}

func CloseLobby(lobby ObjectStructures.LobbyData) {
	secret, err := ioutil.ReadFile("./cert/apiSecret.txt")
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("DELETE", apiUrl+"/v1/lobbies/"+lobby.ID, nil)
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}
	req.Header.Add("Authorization", "Basic "+strings.TrimSuffix(string(secret), "\n"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

//sends request to API
func doRequest(req *http.Request) {
	secret, err := ioutil.ReadFile("./cert/apiSecret.txt")
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+strings.TrimSuffix(string(secret), "\n"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func getIP() string {
	return "gameserver.speedrun.io"
	/*
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return ""
		}
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	*/
}
