package websocket

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/apps/authentication"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/manager"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	baseURL   = "/ws"
	component = "websocket"
)

var hub *Hub

type Websocket struct {
}

type broadcast struct {
	Type  string      `json:"type"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Client struct {
	isAdmin bool

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}

	username string
}

var upgrader = websocket.Upgrader{}

func (Websocket) Register(r *mux.Router) {
	socketRouter := r.PathPrefix(baseURL).Subrouter()
	hub = newHub()
	go hub.run()

	socketRouter.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		serveWS(w, r, hub)
	})
}

func updateProgress() {
	for {
		for key, value := range manager.GetAllProgress() {
			hub.broadcastAdmin <- broadcast{
				Type:  "progress",
				Key:   key,
				Value: value,
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func serveWS(w http.ResponseWriter, r *http.Request, hub *Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	username := authentication.GetUsernameFromSession(r)

	client := &Client{username: username, conn: conn, send: make(chan interface{}, 1), isAdmin: func() bool { return authentication.IsAdmin(username) }()}
	hub.register <- client

	go client.writeData()
	go client.readData()
	go updateProgress()

}

func broadcastUserStatus() {
	hub.broadcastAdmin <- broadcast{
		Type:  "status",
		Value: hub.status,
	}
	fmt.Println(hub.broadcastAdmin)
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
				helpers.LogError(err.Error(), component)
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
			helpers.LogError(err.Error(), component)
		}
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				helpers.LogInfo(err.Error(), component)
			}
			break
		}
	}
}
