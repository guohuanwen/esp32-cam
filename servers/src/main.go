package main

import (
	"camera/src/socket"
	"log"
	"net/http"
)

const (
	PART_BOUNDARY       = "123456789000000000000987654321"
	STREAM_CONTENT_TYPE = "multipart/x-mixed-replace;boundary=" + PART_BOUNDARY
	STREAM_BOUNDARY     = "\r\n--" + PART_BOUNDARY + "\r\n"
	STREAM_PART         = "Content-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n"
)

func init() {
	log.Print("init")
}

func main() {
	defer log.Printf("serves exited...")
	log.Printf("serves start")
	serverMux := http.NewServeMux()
	var hub *socket.Hub
	hub = socket.SetupHub()
	socket.SetupUploader(serverMux, hub)
	socket.SetupClient(serverMux, hub)
	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}
	server.ListenAndServe()
}