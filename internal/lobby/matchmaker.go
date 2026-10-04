package lobby

import (
	"context"
	"log/slog"
	"sync"
	"uuid"

	"github.com/jauharfz/battleship/internal/auth"
)

type Matchmaker struct {
	Pool        chan string
	Waiting     string
	ticketStore *auth.TicketStore
}

func NewMatchmaker(ticketStore *auth.TicketStore) *Matchmaker {
	return &Matchmaker{
		Pool:        make(chan string),
		Waiting:     "",
		ticketStore: ticketStore,
	}
}

// from handler = queue <- playerID
func (m *Matchmaker) Matching(wg *sync.WaitGroup, ctx context.Context) {
	slog.Info("matchmaking started")
	for {
		select {
		case <-ctx.Done():
			slog.Info("matchmaking shutdown gracefully")
			return
		case id := <-m.Pool:
			if m.Waiting == "" {
				m.Waiting = id
				continue
			}
			roomID := uuid.New().String()
			m.ticketStore.Issue(m.Waiting, roomID)
			m.ticketStore.Issue(id, roomID)
			m.Waiting = ""
		}
	}
}

// TODO: handle the condition if the client is going away while matchmaking
