// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import "log"

type Broadcast struct {
	data []byte
	sender int
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Broadcast

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Broadcast),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case broadcast := <-h.broadcast:
			//uploader的视频
			log.Printf("uid %d - %s", broadcast.sender, broadcast.data)
			for client := range h.clients {
				if broadcast.sender == UID_UPLOADER {
					if client.sender != UID_UPLOADER {
						select {
						case client.send <- broadcast.data:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				} else {
					if client.sender == UID_UPLOADER {
						select {
						case client.send <- broadcast.data:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		}
	}
}
