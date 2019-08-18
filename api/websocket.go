package api

import (
	"log"

	"github.com/rgamba/evtwebsocket"
)

var conn = evtwebsocket.Conn {
	// When connection is established
	OnConnected: func(w *evtwebsocket.Conn) {
		log.Println("Connected")
	},
	// When a message arrives
	OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
		log.Printf("OnMessage: %s\n", msg)
	},
    // When the client disconnects for any reason
	OnError: func(err error) {
		log.Printf("** ERROR **\n%s\n", err.Error())
	},
	// Auto reconnect on error
	Reconnect: true,
}

func Connect() {
	err := conn.Dial("ws://predb.ovh/api/v1/ws", "")
	if err != nil {
		log.Fatal(err)
	}
}