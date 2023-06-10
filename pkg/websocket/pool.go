package websocket

import "log"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan MsgBody
	Rooms      map[*Room]bool
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan MsgBody),
		Rooms:      make(map[*Room]bool),
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
	for client := range pool.Clients {
		client.Conn.WriteJSON(MsgBody{
			Content:     message,
			ContentType: "text",
			UserID:      "system"})
	}
}

func (server *Pool) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range server.Rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *Pool) createRoom(name string) *Room {
	room := NewRoom(name)
	go room.RunRoom()
	server.Rooms[room] = true

	return room
}
