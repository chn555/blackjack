package blackjack

import (
	"fmt"

	"github.com/google/uuid"
)

// GameStatus represents the current state of the game.
type GameStatus uint8

const (
	UnknownState GameStatus = iota
	WaitingToStart
	WaitingForPlayer
	Done
)

// Game represents the state of a blackjack game.
type Game struct {
	// ID is the unique identifier for the game.
	ID string
	// Status is the current status of the game.
	Status GameStatus

	// Players contains all players in the game.
	Players map[string]*Player
	// PlayerOrder is the order of player turns in the game.
	playerOrder []string
	// NextPlayer is the player whose turn it is.
	NextPlayer string
	// Winner is the winner of the game.
	// If the game is still in-progress, this field is empty.
	Winner string
}

// DealerName is the name of the dealer player in the game.
const DealerName = "Dealer"

// NewGame seeds the player hands and starts the game.
// The dealer is always the first player.
// To progress the game call Game.PlayTurn.
func NewGame(players ...string) (*Game, error) {
	g := &Game{
		ID:          uuid.New().String(),
		Status:      WaitingToStart,
		Players:     map[string]*Player{},
		playerOrder: []string{},
	}

	players = append([]string{DealerName}, players...)

	err := g.createPlayers(players)
	if err != nil {
		return nil, fmt.Errorf("failed to create players: %w", err)
	}

	g.Status = WaitingForPlayer
	g.NextPlayer = g.playerOrder[0]

	return g, nil
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

// createPlayer creates a player with a hand of cards,
// and adds them to the game.
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
