package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

// define a reader which will listen for new messages
// being sent to our WebSocket endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// log that message for dev purposes
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("aa ", err)
			return
		}
	}
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is online")
	})
	// map our "/ws" endpoint to the "serveWs" function
	http.HandleFunc("/ws", serveWs)
}

func main() {
	fmt.Println("Starting server...")
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
		log.Printf("Defaulting to port %s", port)
	}
	setupRoutes()
	log.Println("Chatting backend server is listening on port " + port + ".")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
