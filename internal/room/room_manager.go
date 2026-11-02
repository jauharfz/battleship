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

// TODO: unlock mu before sending to channel
func (rm *RoomManager) RegisterClient(roomID string, client *Client) {
	rm.mu.Lock()

	room, exist := rm.Rooms[roomID]
	if !exist {
		room = NewRoom(roomID, rm)
		rm.Rooms[roomID] = room
		go room.Run()
	}
	rm.mu.Unlock()
	room.Register <- client
}

// TODO: CREATE UnregisterClient
// and fix possibility concurent bug at room when delete id
func (rm *RoomManager) UnregisterClient(id string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	delete(rm.Rooms, id)
}
