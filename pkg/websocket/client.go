package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type MsgBody struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	UserID      string `json:"user_id"`
}

type SocketMessage struct {
	Type int     `json:"type"`
	Body MsgBody `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	// Declare a new MsgBody struct.
	var body MsgBody

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Parse the message into our MsgBody struct.
		// TODO: Handle errors correctly.
		err = json.Unmarshal(p, &body)

		message := SocketMessage{Type: messageType, Body: body}
		c.Pool.Broadcast <- message
		log.Printf("Message Received: %+v\n", message)
	}
}
