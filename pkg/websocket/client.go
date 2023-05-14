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

// constants for better communication management
const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

type Client struct {
	ID    string
	Token string
	Conn  *websocket.Conn
	Pool  *Pool
	Send  chan []SocketMessage
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

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The pool closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			var msg []byte
			msg, err = json.Marshal(message)
			if err != nil {
				return
			}
			w.Write(msg)

			// Attach queued chat messages to the current websocket message.
			var queuedMsg []byte
			n := len(c.Send)
			for i := 0; i < n; i++ {
				queuedMsg, err = json.Marshal(<-c.Send)
				if err != nil {
					return
				}
				w.Write(queuedMsg)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Declare a new MsgBody struct.
	var body MsgBody

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		// Parse the message into our MsgBody struct.
		// TODO: Handle errors correctly.
		err = json.Unmarshal(p, &body)

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
		"contentType": body.ContentType,
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
