package ai

import (
	blackjackPb "github.com/chn555/schemas/proto/blackjack/v1"
)

type GreedyStrategy struct{}

func (g GreedyStrategy) DecideAction(_ *blackjackPb.Hand) (blackjackPb.Turn_TURN_ACTION, error) {
	return blackjackPb.Turn_TURN_ACTION_HIT, nil
}
