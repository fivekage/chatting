package websocket

import (
	"encoding/json"
	"fmt"
)

const welcomeMessage = "%s joined"
const goodbyeMessage = "%s left"

type Room struct {
	Name       string
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *MsgBody
}

// NewRoom creates a new Room
func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *MsgBody),
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {
	for {
		select {

		case client := <-room.Register:
			room.registerClientInRoom(client)

		case client := <-room.Unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.Broadcast:
			room.broadcastToClientsInRoom(message)
		}

	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.notifyClientJoined(client)
	room.Clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.Clients[client]; ok {
		room.notifyClientLeft(client)
		delete(room.Clients, client)
	}
}

func (room *Room) broadcastToClientsInRoom(message *MsgBody) {
	messageString, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}

	for client := range room.Clients {
		client.Send <- messageString
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	message := &MsgBody{
		Content:     fmt.Sprintf(welcomeMessage, client.GetID()),
		ContentType: "text",
		UserID:      "system"}

	room.broadcastToClientsInRoom(message)
}

func (room *Room) notifyClientLeft(client *Client) {
	message := &MsgBody{
		Content:     fmt.Sprintf(goodbyeMessage, client.GetID()),
		ContentType: "text",
		UserID:      "system"}

	room.broadcastToClientsInRoom(message)
}

func (room *Room) GetName() string {
	return room.Name
}
