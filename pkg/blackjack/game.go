package blackjack

import (
	"context"
	"fmt"

	aiPb "github.com/chn555/schemas/proto/ai/v1"
	deckPb "github.com/chn555/schemas/proto/deck/v1"
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

	deckClient deckPb.DeckServiceClient
	aiClient   aiPb.AiServiceClient
	DeckID     string
}

const DealerName = "Dealer"

// NewGame creates a game, call Game.Start to start the game
func NewGame(ctx context.Context, deckClient deckPb.DeckServiceClient, aiClient aiPb.AiServiceClient) (*Game, error) {
	if deckClient == nil {
		return nil, fmt.Errorf("deck client is nil")
	}

	if aiClient == nil {
		return nil, fmt.Errorf("ai client is nil")
	}

	deckID, err := generateDeck(ctx, deckClient)
	if err != nil {
		return nil, fmt.Errorf("failed to generate deck: %w", err)
	}

	return &Game{
		ID:          uuid.New().String(),
		deckClient:  deckClient,
		aiClient:    aiClient,
		DeckID:      deckID,
		Status:      WaitingToStart,
		Players:     map[string]*Player{},
		playerOrder: []string{},
	}, nil
}

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

// Start seeds the player hands and starts the game
// The dealer is always the first player
// To progress the game call Game.PlayTurn
func (g *Game) Start(ctx context.Context, players ...string) error {
	players = append([]string{DealerName}, players...)

	err := g.createPlayers(ctx, players)
	if err != nil {
		return fmt.Errorf("failed to create players: %w", err)
	}

	_, err = g.aiClient.PlayGame(ctx, &aiPb.PlayGameRequest{
		GameId:     g.ID,
		PlayerName: DealerName,
		Strategy:   "dealer",
	})

	if err != nil {
		return fmt.Errorf("failed to start dealer game: %w", err)
	}

	g.Status = WaitingForPlayer
	g.NextPlayer = g.playerOrder[0]

	return nil
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
