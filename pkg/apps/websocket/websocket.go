package websocket

import (
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

type progress struct {
	Key   string `json:"uid"`
	Value int    `json:"value"`
}

type Client struct {
	isAdmin bool

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}
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
			hub.broadcastAdmin <- progress{
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

	client := &Client{conn: conn, send: make(chan interface{}, 1), isAdmin: true}
	hub.register <- client

	go client.writeData()
	go updateProgress()
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
