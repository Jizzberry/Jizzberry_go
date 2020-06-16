package websocket

import (
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
		return
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
}
