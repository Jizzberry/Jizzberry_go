package websocket

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/gorilla/websocket"
)

type Client struct {
	isAdmin bool

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}

	username string
}

func (c *Client) writeData() {
	for {
		select {
		case value, ok := <-c.send:
			if !ok {
				return
			}
			err := c.conn.WriteJSON(value)
			if err != nil {
				helpers.LogError(err.Error())
				return
			}
		}
	}
}

func (c *Client) readData() {
	defer func() {
		hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			helpers.LogError(err.Error())
		}
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				helpers.LogInfo(err.Error())
			}
			break
		}
	}
}
