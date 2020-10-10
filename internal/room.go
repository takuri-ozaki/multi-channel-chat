package internal

import (
	"github.com/gorilla/websocket"
	"log"
)

type room struct {
	clients      map[*websocket.Conn]User
	broadcast    chan Message
	closeChannel chan struct{}
}

func NewRoom() room {
	room := room{clients: map[*websocket.Conn]User{}, broadcast: make(chan Message), closeChannel: make(chan struct{})}
	go room.handleMessages()
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

func (r *room) handleMessages() {
	for {
		select {
		case <- r.closeChannel:
			return
		case message := <-r.broadcast:
			r.broadcaseMessages(message)
		default:
		}
	}
}

func (r *room) broadcaseMessages(message Message) {
	for client := range r.clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("error occurred while writing message to client:", err)
			client.Close()
			delete(r.clients, client)
		}
	}
}
