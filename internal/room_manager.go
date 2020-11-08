package internal

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type RoomManager struct {
	rooms map[string]room
}

func NewRoomManager() *RoomManager {
	roomManager := &RoomManager{rooms: map[string]room{}}
	go roomManager.garbageCollect()
	return roomManager
}

func (r *RoomManager) Join(roomName, userName string, ws *websocket.Conn) {
	room, ok := r.rooms[roomName]
	if !ok {
		r.createRoom(roomName)
		room = r.rooms[roomName]
	}
	room.Join(ws, userName)
	log.Println(userName + " joined to " + roomName)
}

func (r *RoomManager) Exit(roomName, userName string, ws *websocket.Conn) {
	room := r.rooms[roomName]
	room.Exit(ws, userName)
	log.Println(userName + " exited from " + roomName)
}

func (r *RoomManager) Speak(roomName string, message Message) {
	if room, ok := r.rooms[roomName]; ok {
		room.broadcast <- message
	}
}

func (r *RoomManager) createRoom(roomName string) {
	r.rooms[roomName] = NewRoom(roomName)
}

func (r *RoomManager) garbageCollect() {
	for {
		for roomName, room := range r.rooms {
			if len(room.clients) == 0 {
				close(r.rooms[roomName].closeChannel)
				delete(r.rooms, roomName)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
