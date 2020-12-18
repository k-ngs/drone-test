package main

import (
	"log"
	"net/http"
	"sync"
)


type InMemoryPlayerStore struct {
	mu sync.RWMutex
	store map[string]int
}

func NewImMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		store: map[string]int{},
		mu: sync.RWMutex{},
	}
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.store[name]
}

func (i *InMemoryPlayerStore) GetLeague() (league []Player) {
	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}
	return
}

func main() {
	server := NewPlayerServer(NewImMemoryPlayerStore())

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen %v", err)
	}
}
