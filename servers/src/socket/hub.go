// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

type Broadcast struct {
	data   []byte
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
			if (client.sender != UID_UPLOADER) {//客户端连上了
				h.broadcast <- &Broadcast{data: []byte(ACTION_OPEN_CAMERA), sender: client.sender}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			if len(h.clients) == 1 {//只剩下uploader
				h.broadcast <- &Broadcast{data: []byte(ACTION_CLOSE_CAMERA), sender: client.sender}
			}
		case broadcast := <-h.broadcast:
			//uploader的视频
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
