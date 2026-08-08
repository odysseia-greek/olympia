package mouseion

import (
	"testing"

	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
)

func TestInMemoryStoreCopiesAddedEvents(t *testing.T) {
	store := NewInMemoryStore()
	event := &pb.RequestEvent{Path: "/original", SessionId: "session-1"}

	if err := store.Add([]*pb.RequestEvent{event}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	event.Path = "/mutated"

	events, err := store.GetBySession("session-1")
	if err != nil {
		t.Fatalf("GetBySession() error = %v", err)
	}
	if got := events[0].Path; got != "/original" {
		t.Fatalf("stored path = %q, want %q", got, "/original")
	}
}

func TestInMemoryStoreCopiesReturnedEvents(t *testing.T) {
	store := NewInMemoryStore()
	if err := store.Add([]*pb.RequestEvent{{Path: "/original"}}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	events, err := store.GetByPath("/original")
	if err != nil {
		t.Fatalf("GetByPath() error = %v", err)
	}
	events[0].Path = "/mutated"

	events, err = store.GetByPath("/original")
	if err != nil {
		t.Fatalf("GetByPath() error = %v", err)
	}
	if got := events[0].Path; got != "/original" {
		t.Fatalf("stored path = %q, want %q", got, "/original")
	}
}

func TestInMemoryStoreSkipsNilEvents(t *testing.T) {
	store := NewInMemoryStore()
	if err := store.Add([]*pb.RequestEvent{nil}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	events, err := store.GetRecent(0)
	if err != nil {
		t.Fatalf("GetRecent() error = %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("len(events) = %d, want 0", len(events))
	}
}

func TestInMemoryStoreEvictsOldestEvents(t *testing.T) {
	store := NewInMemoryStoreWithLimit(2)
	if err := store.Add([]*pb.RequestEvent{
		{Path: "/first", SessionId: "old-session"},
		{Path: "/second", SessionId: "current-session"},
		{Path: "/third", SessionId: "current-session"},
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	recent, err := store.GetRecent(0)
	if err != nil {
		t.Fatalf("GetRecent() error = %v", err)
	}
	if len(recent) != 2 || recent[0].Path != "/third" || recent[1].Path != "/second" {
		t.Fatalf("recent events = %v, want third and second", recent)
	}
	old, err := store.GetBySession("old-session")
	if err != nil {
		t.Fatalf("GetBySession() error = %v", err)
	}
	if len(old) != 0 {
		t.Fatalf("old session retained %d events, want 0", len(old))
	}
}
