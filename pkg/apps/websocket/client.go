package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/apps/jizzberry/stream"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/gorilla/websocket"
	"strconv"
)

type Client struct {
	isAdmin bool

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}

	username string
}

type data struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

const (
	TypeGetStream = "getStreamURL"
)

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
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				helpers.LogInfo(err.Error())
			}
			break
		}
		data := data{}
		err = json.Unmarshal(p, &data)
		if err != nil {
			helpers.LogError(err.Error())
			continue
		}
		c.handleMessage(data)
	}
}

func (c *Client) handleMessage(recieved data) {
	var details map[string]interface{}
	err := json.Unmarshal([]byte(recieved.Data), &details)
	if err != nil {
		helpers.LogError(err.Error())
	}
	switch recieved.Type {
	case TypeGetStream:
		marshal, err := json.Marshal(stream.URLGenerator(func(details map[string]interface{}) (int64, bool, string) {
			if sceneIdVal, ok := details["scene_id"]; ok {
				sceneId, err2 := strconv.ParseInt(fmt.Sprintf("%v", sceneIdVal), 10, 64)
				if err2 != nil {
					helpers.LogError(err2.Error())
					return -1, true, ""
				}
				if playable, ok := details["playable"]; ok {
					if startTime, ok := details["start_time"]; ok {
						return sceneId, playable.(bool), fmt.Sprintf("%f", startTime.(float64))
					}
				}
			}
			return -1, true, ""
		}(details)))

		if err != nil {
			helpers.LogError(err.Error())
			return
		}

		c.send <- data{
			Type: TypeGetStream,
			Data: string(marshal),
		}
	}
}
