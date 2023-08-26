package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v, Payload: %v", err, msg)
			delete(clients, ws)
			break
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)

	go handleMessages()

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
