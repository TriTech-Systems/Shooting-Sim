package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// client data contains the basic data a client (pi) could need
type clientData struct {
	connection *websocket.Conn
	r, g, b    uint8
	delay      uint
}

// upgrader contains the buffer sizes for the websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// create a new router
	r := mux.NewRouter()

	// create the routes
	r.HandleFunc("/", index)
	r.HandleFunc("/ws", ws)
	r.HandleFunc("/css", css)

	// start the server with the custom router
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
		return
	}
}

// serve the index page
func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// serve the css file
func css(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "style.css")
}

// todo make function that returns a copy of clients for thread safe usage
var clients = make(map[*websocket.Conn]clientData)

// upgrade connection to websocket connection
func ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
	}
	defer c.Close()

	clients[c] = clientData{
		connection: c,
		r:          255,
		g:          0,
		b:          0,
		delay:      300,
	}
	for {
		_, message, err := c.ReadMessage()
		// todo catch close messages that appear as errors
		if err != nil {
			log.Println("message:", err)
			continue
		}
		// this makes it behave like an echo server
		err = c.WriteMessage(1, message)
		if err != nil {
			log.Println("message:", err)
		}
	}
}
