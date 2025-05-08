package blackjack

import (
	"context"
	"fmt"

	deckPb "github.com/chn555/schemas/proto/deck/v1"
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

	// DeckClient is the client used to interact with the deck service.
	deckClient deckPb.DeckServiceClient
	// DeckID is the ID of the deck used in the game.
	DeckID string
}

// DealerName is the name of the dealer player in the game.
const DealerName = "Dealer"

// NewGame creates a new deck, seeds the player hands and starts the game.
// The dealer is always the first player.
// To progress the game call Game.PlayTurn.
func NewGame(ctx context.Context, deckClient deckPb.DeckServiceClient, players ...string) (*Game, error) {
	if deckClient == nil {
		return nil, fmt.Errorf("deck client is nil")
	}

	deckID, err := generateDeck(ctx, deckClient)
	if err != nil {
		return nil, fmt.Errorf("failed to generate deck: %w", err)
	}

	g := &Game{
		ID:          uuid.New().String(),
		deckClient:  deckClient,
		DeckID:      deckID,
		Status:      WaitingToStart,
		Players:     map[string]*Player{},
		playerOrder: []string{},
	}

	players = append([]string{DealerName}, players...)

	err = g.createPlayers(ctx, players)
	if err != nil {
		return nil, fmt.Errorf("failed to create players: %w", err)
	}

	g.Status = WaitingForPlayer
	g.NextPlayer = g.playerOrder[0]

	return g, nil
}

// generateDeck creates a new deck using the DeckServiceClient and returns its ID.
func generateDeck(ctx context.Context, deckClient deckPb.DeckServiceClient) (string, error) {
	req := &deckPb.CreateDeckRequest{
		Shuffle: true,
	}

	resp, err := deckClient.Create(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create deck: %w", err)
	}

	return resp.GetDeckId(), nil
}

func (g *Game) createPlayers(ctx context.Context, players []string) error {
	for _, player := range players {
		err := g.createPlayer(ctx, player)
		if err != nil {
			return err
		}
	}
	return nil
}

// createPlayer creates a player with a hand of cards,
// and adds them to the game.
func (g *Game) createPlayer(ctx context.Context, player string) error {
	p := &Player{}
	var err error

	p.Hand, err = g.newHandWithCards(ctx)
	if err != nil {
		return fmt.Errorf("failed to create player hand: %w", err)
	}

	p.calculateScore()
	g.Players[player] = p
	g.playerOrder = append(g.playerOrder, player)
	return nil
}

func (g *Game) newHandWithCards(ctx context.Context) (*Hand, error) {
	hand, err := NewHand(g.DeckID, g.deckClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create hand: %w", err)
	}

	for range 2 {
		if err := hand.PullCard(ctx); err != nil {
			return nil, fmt.Errorf("failed to pull card: %w", err)
		}
	}

	return hand, nil
}
