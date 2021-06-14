package ApiHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

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
}

func GetRandomMapFromApi() (string, error) {
	_, err := http.Get("https://api.speedrun.io")
	if err != nil {
		return "", err
	}
	return "resp.Request.Body.Read()", nil
}

func ReportLobbyChange(lobby ObjectStructures.LobbyData) {
	return
}

func ReportLobby(lobby ObjectStructures.LobbyData) {
	secret, err := ioutil.ReadFile("/cert/jwtSecret.txt")
	if err != nil {
		fmt.Println(err)
	}
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
		MaxPlayerCount: 32,
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

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+string(secret))
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

func CloseLobby(lobby ObjectStructures.LobbyData) {
	secret, err := ioutil.ReadFile("/cert/jwtSecret.txt")
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("DELETE", apiUrl+"/v1/lobbies/"+lobby.ID, nil)
	if err != nil {
		ErrorHelper.OutputToConsole("Error", err.Error())
	}
	req.Header.Add("Authorization", "Basic "+string(secret))
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
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
