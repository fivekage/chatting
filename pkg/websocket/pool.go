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
			pool.Clients[client] = true
			log.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				log.Println(client)
				client.Conn.WriteJSON(SocketMessage{
					Type: 1, Body: MsgBody{
						Content:     "New user joined...",
						ContentType: "text",
						UserID:      "system"}})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(SocketMessage{
					Type: 1, Body: MsgBody{
						Content:     "User disconnected...",
						ContentType: "text",
						UserID:      "system"}})
			}
		case message := <-pool.Broadcast:
			log.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
