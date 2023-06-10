package websocket

import (
	"encoding/json"
	"log"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"

type MsgBody struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	UserID      string `json:"user_id"`
	RoomID      string `json:"room_id"`
}

type Message struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
