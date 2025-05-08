package blackjack

import (
	"fmt"

	"github.com/google/uuid"
)

// GameStatus represents the current state of the game
type GameStatus uint8

const (
	UnknownState GameStatus = iota
	WaitingToStart
	WaitingForPlayer
	Done
)

// Game represents the lifecycle of a blackjack game
// It is not thread-safe and should be used in a single-threaded context
type Game struct {
	ID     string
	Status GameStatus

	Players     map[string]*Player
	playerOrder []string
	NextPlayer  string
	// Winner is empty if game is not Done
	Winner string
}

const DealerName = "Dealer"

// NewGame creates a game, call Game.Start to start the game
func NewGame() *Game {
	return &Game{
		ID:          uuid.New().String(),
		Status:      WaitingToStart,
		Players:     map[string]*Player{},
		playerOrder: []string{},
	}
}

// Start seeds the player hands and starts the game
// The dealer is always the first player
// To progress the game call Game.PlayTurn
func (g *Game) Start(players ...string) error {
	players = append([]string{DealerName}, players...)

	err := g.createPlayers(players)
	if err != nil {
		return fmt.Errorf("failed to create players: %w", err)
	}

	g.Status = WaitingForPlayer
	g.NextPlayer = g.playerOrder[0]

	return nil
}

func (g *Game) createPlayers(players []string) error {
	for _, player := range players {
		err := g.createPlayer(player)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) createPlayer(player string) error {
	p := &Player{}
	var err error

	p.Hand, err = g.newHandWithCards()
	if err != nil {
		return fmt.Errorf("failed to create player hand: %w", err)
	}

	p.calculateScore()
	g.Players[player] = p
	g.playerOrder = append(g.playerOrder, player)
	return nil
}

func (g *Game) newHandWithCards() (*Hand, error) {
	hand, err := NewHand()
	if err != nil {
		return nil, fmt.Errorf("failed to create hand: %w", err)
	}

	return hand, nil
}
