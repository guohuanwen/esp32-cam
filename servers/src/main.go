package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"text/template"
	"github.com/gorilla/websocket"
)

const (
	PART_BOUNDARY      = "123456789000000000000987654321"
	STREAM_CONTENT_TYPE = "multipart/x-mixed-replace;boundary=" + PART_BOUNDARY
	STREAM_BOUNDARY    = "\r\n--" + PART_BOUNDARY + "\r\n"
	STREAM_PART        = "Content-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10240,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var frameChan chan[] byte
func init() {
	frameChan = make(chan []byte, 16)
}

func watchHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", STREAM_CONTENT_TYPE)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, _ := w.(http.Flusher)

	for i := 0; ; i++ {
		var data []byte
		timeout := time.NewTimer(time.Millisecond * 100)
		select {
		case data = <-frameChan:

		case <-timeout.C:
			continue
		}
		fmt.Fprintf(w, "%v", STREAM_BOUNDARY)
		fmt.Fprintf(w, STREAM_PART, i)
		_, err = w.Write(data)
		flusher.Flush()
		if err != nil {
			break
		}
	}
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var ws *websocket.Conn

	log.Printf("video\n")
	ws, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}

	defer ws.Close()

	for {
		t, p, e := ws.ReadMessage()
		log.Printf("type %v data len %v err %v", t, len(p), e)
		if e != nil {
			break
		}

		select {
		case frameChan <- p:
			log.Printf("put frame")
		default:
			log.Printf("drop frame")
		}
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var ws *websocket.Conn

	ws, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}

	defer ws.Close()

	for {
		var data []byte
		select {
		case data = <-frameChan:
		default:
			continue
		}

		err = ws.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			break
		}
	}
}

func canvasHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./canvas.html"))
	var data = map[string]interface{}{}
	tmpl.Execute(w, data)
}

func main() {
	defer log.Printf("Program Exited...")
	log.Printf("start");
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/watch", watchHandler)
	serverMux.HandleFunc("/video", videoHandler)
	serverMux.HandleFunc("/canvas", canvasHandler)
    serverMux.HandleFunc("/stream", streamHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	server.ListenAndServe()
}