package internal

import (
	"github.com/gorilla/websocket"
	"log"
)

type room struct {
	clients   map[*websocket.Conn]User
	broadcast chan Message
}

func NewRoom() room {
	room := room{clients: map[*websocket.Conn]User{}, broadcast: make(chan Message)}
	go room.broadcastMessages()
	return room
}

func (r *room) Join(ws *websocket.Conn, userName string) {
	r.clients[ws] = NewUser(userName)
	r.broadcast <- Message{Message: "joined!", UserName: userName, System: true}
}

func (r *room) Exit(ws *websocket.Conn, userName string) {
	delete(r.clients, ws)
	r.broadcast <- Message{Message: "exited", UserName: userName, System: true}
}

func (r *room) broadcastMessages() {
	for {
		message := <-r.broadcast
		for client := range r.clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error occurred while writing message to client: %v", err)
				client.Close()
				delete(r.clients, client)
			}
		}
	}
}
