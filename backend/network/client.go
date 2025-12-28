package network

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Socket   *websocket.Conn
	Receive  chan []byte
	Room     *Room
	Username string
}

type WSMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type OutgoingMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

func (c *Client) Read() {
	defer func() {
		c.Room.Leave <- c
		close(c.Receive)
		c.Socket.Close()
	}()
	for {
		_, raw, err := c.Socket.ReadMessage() //Pause execution until client sends something
		if err != nil {
			return
		}

		var msg WSMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println(err)
			continue
		}

		// We do not read "From" from frontend
		// We trust the authenticated socket identity
		from := c.Username

		//Get recipient's room
		recipientRoomID := UserRoom(msg.To)
		recipientRoom := GetRoom(recipientRoomID)

		//Send to recipient's inbox
		recipientRoom.Broadcast <- OutgoingMessage{
			From:    from,
			To:      msg.To,
			Content: msg.Content,
		}
	}
}

func (c *Client) Write() {
	for msg := range c.Receive {
		if err := c.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
