package socket

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

const (
	UID_UPLOADER = 1
)
var UID_CLIENT = UID_UPLOADER

func uploaderHandler(w http.ResponseWriter, r *http.Request, hub *Hub) {
	tmpl := template.Must(template.ParseFiles("./html/uploader.html"))
	var data = map[string]interface{}{}
	tmpl.Execute(w, data)
}

func uploaderWs(w http.ResponseWriter, r *http.Request, hub *Hub) {
	conn, err := bwUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	uploader := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), sender: UID_UPLOADER}
	uploader.hub.register <- uploader
	go uploader.writeUploaderPump()
	go uploader.readUploaderPump()
}

func clientHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./html/client.html"))
	var data = map[string]interface{}{}
	tmpl.Execute(w, data)
}

func clientWs(w http.ResponseWriter, r *http.Request, hub *Hub,) {
	conn, err := bwUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), sender: GetClientUid()}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}

var bwUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SetupHub() *Hub {
	hub := newHub()
	go hub.run()
	return hub
}

func SetupClient(serverMux *http.ServeMux, hub *Hub) {
	serverMux.HandleFunc("/camera/client", clientHandler)
	serverMux.HandleFunc("/camera/client/ws", func(writer http.ResponseWriter, request *http.Request) {
		clientWs(writer, request, hub)
	})
}

func SetupUploader(serverMux *http.ServeMux, hub *Hub)  {
	serverMux.HandleFunc("/camera/uploader", func(writer http.ResponseWriter, request *http.Request) {
		uploaderHandler(writer, request, hub);
	})
	serverMux.HandleFunc("/camera/uploader/ws", func(writer http.ResponseWriter, request *http.Request) {
		uploaderWs(writer, request, hub);
	})
}

func GetClientUid() int {
	UID_CLIENT = UID_CLIENT + 1
	return UID_CLIENT
}
