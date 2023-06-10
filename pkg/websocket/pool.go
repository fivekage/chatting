package websocket

import "log"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan SocketMessage
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan SocketMessage),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			var message = client.ID + " joined"
			pool.Clients[client] = true
			informClients(message, client.RoomID, pool)
		case client := <-pool.Unregister:
			var message = client.ID + " left"
			delete(pool.Clients, client)
			informClients(message, client.RoomID, pool)
		case message := <-pool.Broadcast:
			pool.broadcastMessage(&message)
		}
	}
}

// informClients sends a message to all clients in
// the pool when a new client joins or leaves
func informClients(message string, roomId string, pool *Pool) {
	log.Println("Size of Connection Pool: ", len(pool.Clients))
	var messageObject = SocketMessage{
		Type: 1, Body: MsgBody{
			Content:     message,
			ContentType: "text",
			UserID:      "system",
			RoomID:      roomId}}
	pool.broadcastMessage(&messageObject)
}

// broadcastMessage transmits a given message to
// all clients in the pool and in the same room
func (pool *Pool) broadcastMessage(message *SocketMessage) {
	for client := range pool.Clients {
		// only dispatch message to clients who are the same room as the sender
		if client.RoomID == message.Body.RoomID {
			if err := client.Conn.WriteJSON(message); err != nil {
				log.Println("Error broadcasting to clients:", err)
				return
			}
		}
	}
}
