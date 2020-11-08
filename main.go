package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"multi-channel-chat/internal"
	"net/http"
	"os"
)

var roomManager = internal.NewRoomManager()
var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/chat", handleChat)
	http.HandleFunc("/transfer", handleTransfer)
	log.Fatal(http.ListenAndServe(":"+getPort(), nil))
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade get request to ws", err)
		return
	}
	roomManager.Join(getRoomName(r), getUserName(r), ws)

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

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	var message internal.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request")
		return
	}
	go roomManager.Speak(room, message)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func getPort() string {
	chatPort := os.Getenv("port")
	if len(chatPort) == 0 {
		chatPort = "8080"
	}
	return chatPort
}

func getRoomName(r *http.Request) string {
	return r.URL.Query().Get("room")
}

func getUserName(r *http.Request) string {
	return r.URL.Query().Get("user")
}
