package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/philmod/clarifai-go"
	"net/http"
	"os"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	clarifaiMinProb float32 = 0.8
)

var (
	clarifaiId     = os.Getenv("CLARIFAI_ID")
	clarifaiSecret = os.Getenv("CLARIFAI_SECRET")
	clarifaiClient = clarifai.NewClient(clarifaiId, clarifaiSecret)
)

type Message struct {
	Type string   `json:type`
	Pic  string   `json:pic`
	Tags []string `json:tags`
}

func detectTags(tagsToDetect []string, tagsDetected []string, probs []float32) []string {
	var intersections []string
	for i, t := range tagsDetected {
		for _, x := range tagsToDetect {
			if t == x && probs[i] > clarifaiMinProb {
				intersections = append(intersections, t)
			}
		}
	}
	return intersections
}

func tagImage(m *Message, writes chan Message) {
	ind := strings.Index(m.Pic, "base64,") + 7
	tagData, err := clarifaiClient.TagEncoded(clarifai.TagEncodedRequest{EncodedData: m.Pic[ind:]})
	if err != nil {
		log.Error("Error with Clarifai API: ", err)
	} else {
		classes := tagData.Results[0].Result.Tag.Classes
		probs := tagData.Results[0].Result.Tag.Probs
		m.Tags = detectTags(m.Tags, classes, probs)
		log.Info("Tags from Clarifai: ", classes)
		writes <- *m
	}
}

func respondWS(conn *websocket.Conn, writes chan Message) {
	for {
		content, more := <-writes
		conn.WriteJSON(content)
		if !more {
			return
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("New ws connection ")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	writes := make(chan Message)
	go respondWS(conn, writes)

	for {
		m := Message{}
		err := conn.ReadJSON(&m)
		if err != nil {
			log.Error("Error reading json: ", err)
			close(writes)
			return
		}

		log.Info("New image to tag.")

		go tagImage(&m, writes)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request static ", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found.", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	log.Info("Request index")
	http.ServeFile(w, r, "views/index.html")
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/public/", staticHandler)
	startServer()
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
