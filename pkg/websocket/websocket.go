package websocket

import (
	"io"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

// We'll need to define an Upgrader.
// This will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection.
	// This will allow us to make requests from our
	// frontend app to this server.
	// For now, we won't check anything and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}


func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}

// define a reader which will listen for new messages
// being sent to our WebSocket endpoint
func Reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// log that message for dev purposes
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("error reading: ", err)
			return
		}
	}
}

// define a writer which will send out
// received messages to all other clients
func Writer(conn *websocket.Conn) {
	for {
		log.Println("Sending")
		messageType, r, err := conn.NextReader()
		if err != nil {
			log.Println(err)
			return
		}
		w, err := conn.NextWriter(messageType)
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := io.Copy(w, r); err != nil {
			log.Println(err)
			return
		}
		if err := w.Close(); err != nil {
			log.Println(err)
			return
		}
	}
}
