package internal

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

type room struct {
	clients      map[*websocket.Conn]User
	broadcast    chan Message
	closeChannel chan struct{}
	roomName     string
}

func NewRoom(roomName string) room {
	room := room{clients: map[*websocket.Conn]User{}, broadcast: make(chan Message), closeChannel: make(chan struct{}), roomName: roomName}
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
		case <-r.closeChannel:
			return
		case message := <-r.broadcast:
			go r.broadcastMessages(message)
			r.transferMessages(message)
		default:
		}
	}
}

func (r *room) broadcastMessages(message Message) {
	for client := range r.clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("error occurred while writing message to client:", err)
			client.Close()
			delete(r.clients, client)
		}
	}
}

type server struct {
	name     string
	endpoint string
}

func (r *room) transferMessages(message Message) {
	if message.Transferred {
		return
	}
	servers := []server{
		{"chat1", "http://chat_1:8080"},
		{"chat2", "http://chat_2:8080"},
	}
	for _, s := range servers {
		go func(server server) {
			if server.name == os.Getenv("name") {
				return
			}
			url := server.endpoint + "/transfer?room=" + r.roomName
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(message.ToTransferredJson()))
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			client.Do(req)
		}(s)
	}
}
