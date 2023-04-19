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
			informClients(message, pool)
		case client := <-pool.Unregister:
			var message = client.ID + " left"
			delete(pool.Clients, client)
			informClients(message, pool)
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Println("Error :", err)
					return
				}
			}
		}
	}
}

// informClients sends a message to all clients in
// the pool when a new client joins or leaves
func informClients(message string, pool *Pool) {
	log.Println("Size of Connection Pool: ", len(pool.Clients))
	//log.Println(message)
	for client := range pool.Clients {
		client.Conn.WriteJSON(SocketMessage{
			Type: 1, Body: MsgBody{
				Content:     message,
				ContentType: "text",
				UserID:      "system"}})
	}
}
