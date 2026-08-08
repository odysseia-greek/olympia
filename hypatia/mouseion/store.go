package mouseion

import (
	"sync"
	"time"

	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
	"google.golang.org/protobuf/proto"
)

const DefaultMaxEvents = 10_000

// EventStore is the storage interface. Swap the in-memory implementation for
// a Postgres, SQLite, or Redis backend without touching the handler layer.
type EventStore interface {
	Add(events []*pb.RequestEvent) error
	GetBySession(sessionID string) ([]*pb.RequestEvent, error)
	GetByPath(path string) ([]*pb.RequestEvent, error)
	GetRecent(limit int) ([]*pb.RequestEvent, error)
}

// inMemoryStore is the default, non-persistent implementation.
type inMemoryStore struct {
	mu     sync.RWMutex
	events []*storedEvent
	max    int

	// secondary indexes for fast lookups
	bySession map[string][]*storedEvent
	byPath    map[string][]*storedEvent
}

type storedEvent struct {
	proto     *pb.RequestEvent
	timestamp time.Time
}

func NewInMemoryStore() EventStore {
	return NewInMemoryStoreWithLimit(DefaultMaxEvents)
}

func NewInMemoryStoreWithLimit(maxEvents int) EventStore {
	if maxEvents <= 0 {
		maxEvents = DefaultMaxEvents
	}
	return &inMemoryStore{
		bySession: make(map[string][]*storedEvent),
		byPath:    make(map[string][]*storedEvent),
		max:       maxEvents,
	}
}

func (s *inMemoryStore) Add(events []*pb.RequestEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range events {
		if e == nil {
			continue
		}
		e = proto.Clone(e).(*pb.RequestEvent)
		ts, err := time.Parse(time.RFC3339Nano, e.Timestamp)
		if err != nil {
			ts = time.Now().UTC()
		}
		se := &storedEvent{proto: e, timestamp: ts}
		s.events = append(s.events, se)

		if e.SessionId != "" {
			s.bySession[e.SessionId] = append(s.bySession[e.SessionId], se)
		}
		s.byPath[e.Path] = append(s.byPath[e.Path], se)
	}
	s.trim()
	return nil
}

func (s *inMemoryStore) trim() {
	if len(s.events) <= s.max {
		return
	}

	s.events = s.events[len(s.events)-s.max:]
	s.bySession = make(map[string][]*storedEvent)
	s.byPath = make(map[string][]*storedEvent)
	for _, event := range s.events {
		if event.proto.SessionId != "" {
			s.bySession[event.proto.SessionId] = append(s.bySession[event.proto.SessionId], event)
		}
		s.byPath[event.proto.Path] = append(s.byPath[event.proto.Path], event)
	}
}

func (s *inMemoryStore) GetBySession(sessionID string) ([]*pb.RequestEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := s.bySession[sessionID]
	out := make([]*pb.RequestEvent, 0, len(entries))
	for _, e := range entries {
		out = append(out, proto.Clone(e.proto).(*pb.RequestEvent))
	}
	return out, nil
}

func (s *inMemoryStore) GetByPath(path string) ([]*pb.RequestEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := s.byPath[path]
	out := make([]*pb.RequestEvent, 0, len(entries))
	for _, e := range entries {
		out = append(out, proto.Clone(e.proto).(*pb.RequestEvent))
	}
	return out, nil
}

func (s *inMemoryStore) GetRecent(limit int) ([]*pb.RequestEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.events)
	if limit <= 0 || limit > total {
		limit = total
	}
	// events are appended in order, so the most recent are at the tail
	out := make([]*pb.RequestEvent, 0, limit)
	for i := total - 1; i >= total-limit; i-- {
		out = append(out, proto.Clone(s.events[i].proto).(*pb.RequestEvent))
	}
	return out, nil
}
