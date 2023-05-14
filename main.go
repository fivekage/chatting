package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fivekage/stay.chatting/pkg/websocket"
	"github.com/joho/godotenv"
)

// define our WebSocket endpoint
func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket reached on", r.Host)

	// upgrade this connection to a WebSocket connection
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Println(w, "%+V\n", err)
	}

	client := &websocket.Client{
		ID:    r.URL.Query().Get("id"),
		Token: r.URL.Query().Get("token"),
		Conn:  conn,
		Pool:  pool,
		Send:  make(chan []websocket.SocketMessage, 256),
	}

	go client.Write()
	go client.Read()

	pool.Register <- client
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	// map our "/ws" endpoint to the "serveWs" function
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	log.Println("Starting server... version 0.3")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	apiBaseUrl := os.Getenv("API_BASE_URL")
	if port == "" {
		port = "5000"
		log.Printf("Defaulting to port %s", port)
	}
	if apiBaseUrl == "" {
		apiBaseUrl = "http://localhost:5001"
		log.Printf("Defaulting API BASE URL %s", apiBaseUrl)
	}
	setupRoutes()
	log.Println("Chatting backend server is listening on port " + port + ".")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
