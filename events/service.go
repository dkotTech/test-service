package events

import (
	"context"
	"github.com/google/uuid"
	"maps"
	"slices"
	"sync"
	"test-service/helpers"
)

type (
	EventKind string

	Event struct {
		Kind EventKind   `json:"kind"`
		Data interface{} `json:"data"`
	}

	Service interface {
		// BroadcastEvent send a new event to all current listeners
		BroadcastEvent(ctx context.Context, eventKind EventKind, data interface{})
		// RegisterConnection create a EventClient for handle connection outside of service
		RegisterConnection(ctx context.Context, events []EventKind) EventClient
	}

	// EventClient interface for methods with event client must implement
	EventClient interface {
		// Shutdown remove a listener from query
		Shutdown() error
		// Subscribe get chanel with events
		Subscribe() chan Event
		// EventsSubscribed return a list of EventKind
		EventsSubscribed() []EventKind
	}
)

const (
	LeaderboardChanges EventKind = "leaderboard_changes"
	Deposit            EventKind = "deposit"
	Withdraw           EventKind = "withdraw"
)

var defaultEvents = []EventKind{LeaderboardChanges, Deposit, Withdraw}

type connection struct {
	id           uuid.UUID
	events       map[EventKind]struct{}
	eventsChanel chan Event

	shutdown func(id uuid.UUID) error
}

func (c *connection) Shutdown() error {
	return c.shutdown(c.id)
}

func (c *connection) Subscribe() chan Event {
	return c.eventsChanel
}

func (c *connection) EventsSubscribed() []EventKind {
	return slices.Collect(maps.Keys(c.events))
}

type eventsService struct {
	connections map[uuid.UUID]connection

	accessMutex sync.Mutex
}

func NewEventsService() Service {
	return &eventsService{
		connections: map[uuid.UUID]connection{},
	}
}

func (s *eventsService) RegisterConnection(ctx context.Context, events []EventKind) EventClient {
	id := uuid.New()
	// if no events subscribed, subscribe to default list
	if len(events) == 0 {
		events = defaultEvents
	}

	con := connection{
		id:           id,
		events:       helpers.SliceToUniqMap(events),
		eventsChanel: make(chan Event, 10),
		shutdown: func(id uuid.UUID) error {
			s.accessMutex.Lock()
			defer s.accessMutex.Unlock()
			delete(s.connections, id)
			return nil
		},
	}

	s.accessMutex.Lock()
	defer s.accessMutex.Unlock()
	s.connections[id] = con

	return &con
}

func (s *eventsService) BroadcastEvent(ctx context.Context, eventKind EventKind, data interface{}) {
	s.accessMutex.Lock()
	defer s.accessMutex.Unlock()

	// send event over all connections
	for _, conn := range s.connections {
		if _, found := conn.events[eventKind]; !found {
			continue
		}

		conn.eventsChanel <- Event{
			Kind: eventKind,
			Data: data,
		}
	}
}
