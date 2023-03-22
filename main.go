package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fivekage/stay.chatting/pkg/websocket"
)

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Host)

	// upgrade this connection to a WebSocket connection
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Println(w, "%+V\n", err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	go websocket.Writer(ws)
	websocket.Reader(ws)
}

func setupRoutes() {
	// map our "/ws" endpoint to the "serveWs" function
	http.HandleFunc("/ws", serveWs)
}

func main() {
	log.Println("Starting server...")
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
