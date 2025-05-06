package blackjack

import (
	"context"
	"fmt"
	"sync"
)

// GameStore is an interface for storing and retrieving game state
type GameStore interface {
	Get(ctx context.Context, id string) (*Game, error)
	Put(ctx context.Context, id string, game *Game) error
}

// InMemoryGameStore is an in-memory implementation of GameStore
type InMemoryGameStore struct {
	state   map[string]*Game
	rwMutex *sync.RWMutex
}

// NewInMemoryGameStore creates a new InMemoryGameStore
func NewInMemoryGameStore() *InMemoryGameStore {
	return &InMemoryGameStore{
		state:   make(map[string]*Game),
		rwMutex: &sync.RWMutex{},
	}
}

// Get retrieves the game state by ID
func (s *InMemoryGameStore) Get(_ context.Context, id string) (*Game, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	game, ok := s.state[id]
	if !ok {
		return nil, fmt.Errorf("game %s not found", id)
	}

	return game, nil
}

// Put stores the game state by ID
func (s *InMemoryGameStore) Put(_ context.Context, id string, game *Game) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.state[id] = game
	return nil
}
