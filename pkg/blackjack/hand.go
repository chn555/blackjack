package blackjack

import (
	"github.com/chn555/schemas/proto/deck/v1"
)

type Hand struct {
	Cards []*v1.Card
}

// NewHand creates a new empty hand
func NewHand() (*Hand, error) {

	return &Hand{}, nil
}
