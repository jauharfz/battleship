package auth

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
	"uuid"
)

type Ticket struct {
	PlayerID  string
	RoomID    string
	ExpiresAt time.Time
}

func NewTicket(playerID, roomID string, duration time.Duration) *Ticket {
	return &Ticket{
		PlayerID:  playerID,
		RoomID:    roomID,
		ExpiresAt: time.Now().Add(duration),
	}
}

type TicketStore struct {
	mu            sync.Mutex
	ttl           time.Duration
	sweepInterval time.Duration
	tickets       map[string]*Ticket // ticketID as key
	ticketIDs     map[string]string  // playerID as key
}

func NewTicketStore(ttl, sweepInterval time.Duration) *TicketStore {
	return &TicketStore{
		ttl:           ttl,
		sweepInterval: sweepInterval,
		tickets:       make(map[string]*Ticket),
		ticketIDs:     make(map[string]string),
	}
}

func (ts *TicketStore) Issue(playerID, roomID string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	// create ticket -> store to map -> delete either in consume or sweeper
	ticket := NewTicket(playerID, roomID, ts.ttl)
	ts.ticketIDs[playerID] = uuid.New().String()
	ts.tickets[ts.ticketIDs[playerID]] = ticket
}

func (ts *TicketStore) Get(playerID string) (string, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ticketID, exist := ts.ticketIDs[playerID]
	if !exist {
		return "", errors.New("ticket/token_id doesn't exist")
	}
	_, err := ts.getTicketById(ticketID)
	if err != nil {
		return "", err
	}
	return ticketID, nil
}

// consider use mutex or channel
// refactor the websoceket handler either the ticket itself
func (ts *TicketStore) Consume(ticketID string) (*Ticket, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ticket, err := ts.getTicketById(ticketID)
	if err != nil {
		return nil, err
	}
	delete(ts.tickets, ticketID)
	delete(ts.ticketIDs, ticket.PlayerID)
	return ticket, nil
}

func (ts *TicketStore) getTicketById(ticketID string) (*Ticket, error) {
	ticket, exist := ts.tickets[ticketID]
	if !exist {
		return nil, errors.New("ticket/token doesn't exist")
	}
	if time.Now().After(ticket.ExpiresAt) {
		delete(ts.tickets, ticketID)
		delete(ts.ticketIDs, ticket.PlayerID)
		return nil, errors.New("ticket/token is Expired")
	}
	return ticket, nil
}

// TODO: dont forget to create sweeper
func (ts *TicketStore) TicketSweeper(ctx context.Context) {
	slog.Info("ticket sweeper started")
	ticker := time.NewTicker(ts.sweepInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			slog.Info("ticket sweeper shutdown gracefully")
			return
		case <-ticker.C:
			ts.mu.Lock()
			for id, ticket := range ts.tickets {
				if time.Now().After(ticket.ExpiresAt) {
					delete(ts.tickets, id)
					delete(ts.ticketIDs, ticket.PlayerID)
				}
			}
			ts.mu.Unlock()
		}
	}
}
