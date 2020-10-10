package main

import (
	"github.com/gorilla/websocket"
	"log"
	"multi-channel-chat/internal"
	"net/http"
)

var roomManager = internal.NewRoomManager()
var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/chat", handleChat)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade get request to ws", err)
	}
	go roomManager.Join(getRoomName(r), getUserName(r), ws)

	defer ws.Close()
	for {
		var message internal.Message
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("failed to read message", err)
			roomManager.Exit(getRoomName(r), getUserName(r), ws)
			break
		}
		message.UserName = getUserName(r)
		roomManager.Speak(getRoomName(r), message)
	}
}

func getRoomName(r *http.Request) string {
	return r.URL.Query().Get("room")
}

func getUserName(r *http.Request) string {
	return r.URL.Query().Get("user")
}
