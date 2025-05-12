package blackjack

import (
	"github.com/chn555/schemas/proto/deck/v1"
)

// Hand represents a player's hand in the game, consisting of the player's cards.
type Hand struct {
	// Cards represents the cards in the player's hand.
	Cards []*v1.Card
}

// NewHand creates a new empty hand.
func NewHand() (*Hand, error) {
	return &Hand{}, nil
}
