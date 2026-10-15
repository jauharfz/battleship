package lobby

import (
	"context"
	"log/slog"
	"sync"
	"time"
	"uuid"

	"github.com/jauharfz/battleship/internal/auth"
)

// handler check if exist room in matchmaker
// -> if exist just issue a ticket&delete map + check exp before that
// -> if doesn't store in the map with room key and player_id value
type Matchmaker struct {
	mu            sync.Mutex // needs for every manipulate CustomRooms
	Pool          chan string
	Waiting       *QueuedPlayer
	CustomRooms   map[string]*QueuedPlayer // K:roomID,
	CustomPlayers map[string]*QueuedPlayer // K:playedID
	ticketStore   *auth.TicketStore
	sweepInterval time.Duration
}

// NOTE: Custom player exclusively pairing w/ Custom room only, random matchmaker just use pool -> waiting
func NewMatchmaker(ticketStore *auth.TicketStore, sweepInterval time.Duration) *Matchmaker {
	return &Matchmaker{
		Pool:          make(chan string),
		ticketStore:   ticketStore,
		CustomRooms:   make(map[string]*QueuedPlayer),
		CustomPlayers: make(map[string]*QueuedPlayer),
		sweepInterval: sweepInterval,
	}
}

// NOTE: last checked just checking random matchmaker only for now, implementation below the ticker
type QueuedPlayer struct {
	PlayerID    string
	LastChecked time.Time // just assign time.Now(), check with matchmaker ticker below
}

func NewQueuedPlayer(playerID string) *QueuedPlayer {
	return &QueuedPlayer{
		PlayerID:    playerID,
		LastChecked: time.Now(),
	}
}

// from handler = queue <- playerID
func (m *Matchmaker) Matching(wg *sync.WaitGroup, ctx context.Context) {
	slog.Info("matchmaking started")
	ticker := time.NewTicker(m.sweepInterval)
	for {
		select {
		case <-ctx.Done():
			slog.Info("matchmaking shutdown gracefully")
			return
		case id := <-m.Pool:
			if m.Waiting == nil {
				m.Waiting = NewQueuedPlayer(id)
				// m.CustomPlayers[id] = m.Waiting
				continue
			}
			if m.Waiting.PlayerID == id {
				continue
			}
			roomID := uuid.New().String()
			m.ticketStore.Issue(m.Waiting.PlayerID, roomID)
			m.ticketStore.Issue(id, roomID)
			m.Waiting = nil
		case <-ticker.C:
			m.mu.Lock()
			// check for random waiter
			if m.Waiting != nil {
				if time.Since(m.Waiting.LastChecked) >= m.sweepInterval+2 {
					m.Waiting = nil
				}
			}

			// check for waiter which assign the custom room before
			for room, queuedPlayer := range m.CustomRooms {
				if time.Since(queuedPlayer.LastChecked) >= m.sweepInterval+2 {
					delete(m.CustomRooms, room)
				}
			}
			m.mu.Unlock()

		}
	}
}

func (m *Matchmaker) DeleteCustomMatch(roomID, playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.CustomRooms, roomID)
	delete(m.CustomPlayers, playerID)
}

func (m *Matchmaker) AssignCustomMatch(roomID, playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	queuedPlayer := NewQueuedPlayer(playerID)
	m.CustomRooms[roomID] = queuedPlayer
	m.CustomPlayers[playerID] = queuedPlayer
}

func (m *Matchmaker) UpdateLastChecked(playerID string, isCustom bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if isCustom {
		m.CustomPlayers[playerID].LastChecked = time.Now()
	} else {
		m.Waiting.LastChecked = time.Now()
	}
}

// TODO: handle the condition if the client is going away while matchmaking
