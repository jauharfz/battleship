package room

import "sync"

type RoomManager struct {
	mu    sync.Mutex
	Rooms map[string]*Room
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) RegisterClient(roomID string, client *Client) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exist := rm.Rooms[roomID]
	if !exist {
		room = NewRoom(roomID)
		rm.Rooms[roomID] = room
		go room.Run()
	}
	room.Register <- client
}
