package main

import (
	// "fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		log.Info("Serve ws : ", p)

		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request static ", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request index")
	http.ServeFile(w, r, "views/index.html")
}

func startServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}

func main() {
	// go h.run()
	http.HandleFunc("/ws", serveWs)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/public/", staticHandler)
	startServer()
}
