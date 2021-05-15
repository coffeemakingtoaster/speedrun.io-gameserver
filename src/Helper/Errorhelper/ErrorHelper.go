package ErrorHelper

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	ObjectStructures "gameserver.speedrun.io/Helper/Objecthelper"
	SocketHelper "gameserver.speedrun.io/Helper/Sockethelper"
)

func InvalidRouteError(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Invalid Route. If you are trying to reach the game API please interact with api.speedrun.io")
}

func ConnectionNotWebsocketError(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Error: Connection to the /ws part of the gameserver should only be via websockets")
}

func InvalidRequestError(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Error: The request is invalid in this context")
}

func InvalidRoomIDError(conn *websocket.Conn) {
	SocketHelper.Sender(conn, ObjectStructures.Message{Type: 0, Data: []string{"Error! Invalid room code"}})
}
