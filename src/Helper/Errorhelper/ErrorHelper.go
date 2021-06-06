package ErrorHelper

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
)

func InvalidRouteError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Invalid Route. If you are trying to reach the game API please interact with api.speedrun.io")
}

func ConnectionNotWebsocketError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Error: Connection to the /ws part of the gameserver should only be via websockets")
}

func InvalidRequestError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Error: The request is invalid in this context")
}

/*
func InvalidRoomIDError(conn *websocket.Conn) {
		err := SocketHelper.Sender(conn, ObjectStructures.Message{Type: 0, Data: []string{"Error! Invalid room code"}})
		if err != nil {
			OutputToConsole("Error", "Message could not be send")
		}
	OutputToConsole("Error", "Invalid room ID provided")
}
*/

func OutputToConsole(mode string, message string) {
	c := color.New(color.FgWhite)
	if mode == "Error" {
		c = color.New(color.FgRed)
	} else if mode == "Warning" {
		c = color.New(color.FgYellow)
	} else if mode == "Update" {
		c = color.New(color.FgGreen)
	}
	c.Println("[ " + mode + " ]: " + message)
}
