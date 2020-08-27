// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

// Hub maintains the set of active admins and broadcasts messages to the
// admins.
type Hub struct {
	// Registered clients.
	admins map[*Client]bool
	users  map[*Client]bool

	status map[string]*UserStatus

	// Inbound messages from the clients.
	broadcastAdmin chan interface{}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type UserStatus struct {
	Online bool
}

func newHub() *Hub {
	return &Hub{
		broadcastAdmin: make(chan interface{}),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		admins:         make(map[*Client]bool),
		users:          make(map[*Client]bool),
		status:         make(map[string]*UserStatus),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if client.isAdmin {
				h.admins[client] = true
			} else {
				h.users[client] = true
			}
			h.status[client.username] = &UserStatus{Online: true}
			go broadcastUserStatus()
		case client := <-h.unregister:
			if _, ok := h.admins[client]; ok {
				delete(h.admins, client)
				close(client.send)
			} else if _, ok := h.admins[client]; ok {
				delete(h.users, client)
				close(client.send)
			}
			if _, ok := h.status[client.username]; ok {
				delete(h.status, client.username)
				go broadcastUserStatus()
			}
		case message := <-h.broadcastAdmin:
			for client := range h.admins {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.admins, client)
				}
			}
		}
	}
}
