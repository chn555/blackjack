package ai

import (
	blackjackPb "github.com/chn555/schemas/proto/blackjack/v1"
)

type DealerStrategy struct {
}

func (g *DealerStrategy) DecideAction(hand *blackjackPb.Hand) (blackjackPb.Turn_TURN_ACTION, error) {
	// Implement the greedy strategy logic here
	// For example, always hit if the player's hand value is less than 17
	if hand.Score < 17 {
		return blackjackPb.Turn_TURN_ACTION_HIT, nil
	}
	return blackjackPb.Turn_TURN_ACTION_STAND, nil
}
