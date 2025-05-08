/*
Package blackjack contains the implementation of the BlackjackServiceServer interface.

This package is intended to be a translation layer
between the gRPC API and the internal business logic
implemented in the blackjack package.
*/
package blackjack

import (
	"context"
	"fmt"

	"github.com/chn555/blackjack/pkg/blackjack"
	blackjackPb "github.com/chn555/schemas/proto/blackjack/v1"
	deckPb "github.com/chn555/schemas/proto/deck/v1"
)

// ServiceServer is an implementation of the BlackjackServiceServer interface.
type ServiceServer struct {
	// UnimplementedBlackjackServiceServer is embedded to ensure forward compatibility,
	// and to provide a default implementation of the methods.
	blackjackPb.UnimplementedBlackjackServiceServer
	deckClient deckPb.DeckServiceClient
	store      blackjack.GameStore
}

// NewServiceServer creates a new BlackjackServiceServer.
func NewServiceServer(store blackjack.GameStore, deckClient deckPb.DeckServiceClient) (*ServiceServer, error) {
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	if deckClient == nil {
		return nil, fmt.Errorf("deck client is nil")
	}

	return &ServiceServer{store: store, deckClient: deckClient}, nil
}

// NewGame creates a new persistent game.
// It creates a new deck and a new hand for each of the players.
// All player cards are hidden in the response, since we don't know which player sent the request.
func (s *ServiceServer) NewGame(ctx context.Context, request *blackjackPb.NewGameRequest) (*blackjackPb.Game, error) {
	if len(request.GetPlayerNames()) == 0 {
		return nil, fmt.Errorf("no player names provided")
	}

	game, err := blackjack.NewGame(ctx, s.deckClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	err = s.store.Put(ctx, game.ID, game)
	if err != nil {
		return nil, fmt.Errorf("failed to store game: %w", err)
	}

	protoGame := gameToProto(game)

	// Since we do not know which player sent the request,
	// we hide all player hands in the response.
	hideHandForPlayer(protoGame, "")
	return protoGame, nil
}

func gameToProto(game *blackjack.Game) *blackjackPb.Game {
	return &blackjackPb.Game{
		GameId:      game.ID,
		NextPlayer:  game.NextPlayer,
		Winner:      game.Winner,
		Status:      mapGameStatusToProto(game.Status),
		PlayerHands: mapPlayerHandsToProto(game.Players),
	}
}

func mapPlayerHandsToProto(players map[string]*blackjack.Player) map[string]*blackjackPb.Hand {
	handMap := make(map[string]*blackjackPb.Hand)
	for name, player := range players {
		handMap[name] = &blackjackPb.Hand{
			Cards:  player.Hand.Cards,
			Score:  int32(player.Score),
			IsBust: player.Bust,
		}
	}
	return handMap
}

func mapGameStatusToProto(status blackjack.GameStatus) blackjackPb.Game_GAME_STATUS {
	switch status {
	case blackjack.WaitingToStart:
		return blackjackPb.Game_GAME_STATUS_WAITING_TO_START
	case blackjack.WaitingForPlayer:
		return blackjackPb.Game_GAME_STATUS_WAITING_FOR_PLAYER
	case blackjack.Done:
		return blackjackPb.Game_GAME_STATUS_DONE
	default:
		return blackjackPb.Game_GAME_STATUS_UNSPECIFIED
	}
}

// hideHandForPlayer hides the hands of all players except the one specified.
func hideHandForPlayer(game *blackjackPb.Game, playerName string) {
	for name, player := range game.PlayerHands {
		if name == blackjack.DealerName {
			player.Cards = player.Cards[:1]
			player.Score = 0
			continue
		}
		if name != playerName {
			player.Cards = nil
			player.Score = 0
		}
	}
}

// PlayTurn plays a turn in the game.
func (s *ServiceServer) PlayTurn(ctx context.Context, req *blackjackPb.Turn) (*blackjackPb.Game, error) {
	game, err := s.store.Get(ctx, req.GetGameId())
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	turn := blackjack.NewTurn(
		mapProtoTurnToTurn(req.GetAction()),
		req.GetPlayerName(),
	)
	err = game.PlayTurn(ctx, turn)
	if err != nil {
		return nil, fmt.Errorf("failed to play out of turn: %w", err)
	}

	err = s.store.Put(ctx, game.ID, game)
	if err != nil {
		return nil, fmt.Errorf("failed to store game: %w", err)
	}

	protoGame := gameToProto(game)
	hideHandForPlayer(protoGame, req.GetPlayerName())
	return protoGame, nil
}

func mapProtoTurnToTurn(action blackjackPb.Turn_TURN_ACTION) blackjack.TurnAction {
	switch action {
	case blackjackPb.Turn_TURN_ACTION_HIT:
		return blackjack.Hit
	case blackjackPb.Turn_TURN_ACTION_STAND:
		return blackjack.Stand
	default:
		return blackjack.Unknown
	}
}

// GetGame returns the state of the game for the given game id.
// It hides cards for all players except the one who sent the request.
func (s *ServiceServer) GetGame(ctx context.Context, req *blackjackPb.GetGameRequest) (*blackjackPb.Game, error) {
	game, err := s.store.Get(ctx, req.GetGameId())
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	protoGame := gameToProto(game)
	hideHandForPlayer(protoGame, req.GetPlayerName())
	return protoGame, nil
}
