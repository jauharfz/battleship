package lobby

import (
	"sync"

	"github.com/jauharfz/battleship/internal/room"
)

type Matchmaker struct {
	mu      sync.Mutex
	Waiting *room.Client
}

func (m *Matchmaker) Enqueue(client *room.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Waiting == nil || m.Waiting == client {
		m.Waiting = client
		return
	}

	room := room.NewRoom()
	go room.Run()

	room.Register <- m.Waiting
	room.Register <- client
	m.Waiting = nil
}

func (m *Matchmaker) Dequeue() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Waiting = nil
}
