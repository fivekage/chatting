package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID    string
	Token string
	Conn  *websocket.Conn
	Pool  *Pool
}

type MsgBody struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	UserID      string `json:"user_id"`
	RoomID      string `json:"room_id"`
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
		// TOFIX: Handle errors correctly.
		err = json.Unmarshal(p, &body)
		if err != nil {
			log.Println(err)
			return
		}

		message := SocketMessage{Type: messageType, Body: body}
		c.Pool.Broadcast <- message
		log.Printf("Message Received: %+v\n", message)

		historizeMessage(body, c)
	}
}

func historizeMessage(body MsgBody, c *Client) {
	// Build the request body as a JSON object
	bodyData := map[string]interface{}{
		"message":     body.Content,
		"writedAt":    time.Now().Format("2006-01-02T15:04:05.000Z"),
		"writedBy":    body.UserID,
		"chatRoomUid": body.RoomID,
	}
	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		// Handle Error
		log.Fatalf("An error occurred during building request %v", err)
	}

	// Build the request
	client := http.Client{}
	var apiBaseUrl = os.Getenv("API_BASE_URL")
	var apiUrl = apiBaseUrl + "message/chatroom"
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		// Handle Error
		log.Fatalf("An error occurred during building request %v", err)
	}

	// Set the request headers
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + c.Token},
	}

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		// Handle Error
		log.Fatalf("An error occurred during sending request %v", err)
	}
	defer res.Body.Close()
}
