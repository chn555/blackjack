package blackjack

import (
	"context"
	"fmt"

	"github.com/chn555/schemas/proto/deck/v1"
)

// Hand represents a player's hand in the game, consisting of the player's cards.
type Hand struct {
	// Cards represents the cards in the player's hand.
	Cards []*v1.Card
	// DeckID is the ID of the deck used in the hand.
	DeckID string
	// deckClient is the client used to interact with the deck service.
	deckClient v1.DeckServiceClient
}

// NewHand creates a new empty hand.
func NewHand(deckID string, deckClient v1.DeckServiceClient) (*Hand, error) {
	if deckID == "" {
		return nil, fmt.Errorf("DeckID is empty")
	}
	if deckClient == nil {
		return nil, fmt.Errorf("deck client is nil")
	}

	return &Hand{DeckID: deckID, deckClient: deckClient}, nil
}

// PullCard pulls a card from the deck and adds it to the hand
func (h *Hand) PullCard(ctx context.Context) error {
	req := &v1.FetchCardRequest{
		Deck: &v1.Deck{DeckId: h.DeckID},
	}

	resp, err := h.deckClient.FetchCard(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to fetch card: %w", err)
	}

	h.Cards = append(h.Cards, resp)
	return nil
}
