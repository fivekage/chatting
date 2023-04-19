package websocket

import (
	"os"
	"encoding/json"
	"log"
	"time"
	"net/http"
	"net/url"
	"strings"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Token string
	Conn *websocket.Conn
	Pool *Pool
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
		// TODO: Handle errors correctly.
		err = json.Unmarshal(p, &body)

		message := SocketMessage{Type: messageType, Body: body}
		c.Pool.Broadcast <- message
		log.Printf("Message Received: %+v\n", message)

		// historizeMessage(body, c)
		historizeMessage(body, c)
	}
}


func historizeMessage(body MsgBody, c *Client) {

	//Build the request
	data := url.Values{}
    data.Set("message", "body.Content")
	data.Set("writedAt", string(time.Now().UnixMilli()))
	data.Set("writedBy", body.UserID)
	data.Set("chatRoomUid", body.RoomID)

	
	 client := http.Client{}
	 var API_BASE_URL = os.Getenv("API_BASE_URL")
	 var apiUrl = API_BASE_URL + "message/chatroom"
	 req , err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(data.Encode()))
	 if err != nil {
		 //Handle Error
		 log.Fatalf("An Error Occured during building request %v", err)

	 }
	 
	 req.Header = http.Header{
		 "Content-Type": {"application/json"},
		 "Authorization": {"Bearer " + c.Token},
	 }

	 res , err := client.Do(req)
	 if err != nil {
		 //Handle Error
		 log.Fatalf("An Error Occured during request%v", err)
	 }
	 log.Printf(res.Status)
}