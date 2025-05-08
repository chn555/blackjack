/*
Package blackjack contains the implementation of the BlackjackServiceServer interface.

This package is intended to be a translation layer
between the gRPC API and the internal business logic
implemented in the blackjack package.
*/
package blackjack

import (
	"context"

	blackjackPb "github.com/chn555/schemas/proto/blackjack/v1"
)

// ServiceServer is an implementation of the BlackjackServiceServer interface.
type ServiceServer struct {
	// UnimplementedBlackjackServiceServer is embedded to ensure forward compatibility,
	// and to provide a default implementation of the methods.
	blackjackPb.UnimplementedBlackjackServiceServer
}

// NewGame creates a new persistent game.
// It creates a new deck and a new hand for each of the players.
// All player cards are hidden in the response, since we don't know which player sent the request.
func (s *ServiceServer) NewGame(ctx context.Context, request *blackjackPb.NewGameRequest) (*blackjackPb.Game, error) {
	return s.UnimplementedBlackjackServiceServer.NewGame(ctx, request)
}

// PlayTurn plays a turn in the game.
func (s *ServiceServer) PlayTurn(ctx context.Context, req *blackjackPb.Turn) (*blackjackPb.Game, error) {
	return s.UnimplementedBlackjackServiceServer.PlayTurn(ctx, req)
}

// GetGame returns the state of the game for the given game id.
// It hides cards for all players except the one who sent the request.
func (s *ServiceServer) GetGame(ctx context.Context, req *blackjackPb.GetGameRequest) (*blackjackPb.Game, error) {
	return s.UnimplementedBlackjackServiceServer.GetGame(ctx, req)
}
