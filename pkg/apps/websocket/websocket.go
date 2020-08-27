package websocket

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/apps/authentication"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler/manager"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	baseURL = "/ws"
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
		for _, task := range manager.GetAllTaskStatus() {
			hub.broadcastAdmin <- broadcast{
				Type:  "progress",
				Value: task,
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func serveWS(w http.ResponseWriter, r *http.Request, hub *Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	username := authentication.GetUsernameFromSession(r)

	client := &Client{username: username, conn: conn, send: make(chan interface{}, 1), isAdmin: func() bool { return authentication.IsAdminFromSession(r) }()}
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
