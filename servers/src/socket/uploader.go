package socket

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	PART_BOUNDARY       = "123456789000000000000987654321"
	STREAM_CONTENT_TYPE = "multipart/x-mixed-replace;boundary=" + PART_BOUNDARY
	STREAM_BOUNDARY     = "\r\n--" + PART_BOUNDARY + "\r\n"
	STREAM_PART         = "Content-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n"
)

func (up *Client) readUploaderPump() {
	defer func() {
		up.hub.unregister <- up
		up.conn.Close()
	}()
	up.conn.SetReadLimit(maxMessageSize)
	up.conn.SetReadDeadline(time.Now().Add(pongWait))
	up.conn.SetPongHandler(func(string) error { up.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msgType, message, err := up.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		switch msgType {
		case websocket.TextMessage:
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			var text = string(message)
			up.hub.broadcast <- &Broadcast{data: message, sender: up.sender}
			log.Print("rev text " + text)
			if text == "open" {
				log.Print("open")
				log.Print("startPushStream")
			} else if text == "close" {
				log.Print("close")
				log.Print("stopPushStream")
			} else if text == "turnTo0" {
				log.Print("turnTo0")
			} else if text == "turnTo90" {
				log.Print("turnTo90")
			} else if text == "turnTo180" {
				log.Print("turnTo180")
			}
			break
		case websocket.BinaryMessage:
			up.hub.broadcast <- &Broadcast{data: message, sender: up.sender}
			break
		case websocket.CloseMessage:
			log.Print("rev closes")
			break
		}
	}
}

func (c *Client) writeUploaderPump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
