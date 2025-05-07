package blackjack

import (
	"context"
	"fmt"

	"github.com/chn555/blackjack/pkg/blackjack"
	aiPb "github.com/chn555/schemas/proto/ai/v1"
	blackjackPb "github.com/chn555/schemas/proto/blackjack/v1"
	deckPb "github.com/chn555/schemas/proto/deck/v1"
)

type ServiceServer struct {
	blackjackPb.UnimplementedBlackjackServiceServer
	deckClient deckPb.DeckServiceClient
	aiClient   aiPb.AiServiceClient
	store      blackjack.GameStore
}

func NewServiceServer(store blackjack.GameStore, deckClient deckPb.DeckServiceClient, aiClient aiPb.AiServiceClient) (*ServiceServer, error) {
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	if deckClient == nil {
		return nil, fmt.Errorf("deck client is nil")
	}
	if aiClient == nil {
		return nil, fmt.Errorf("ai client is nil")
	}
	return &ServiceServer{store: store, deckClient: deckClient, aiClient: aiClient}, nil
}

func (s *ServiceServer) NewGame(ctx context.Context, request *blackjackPb.NewGameRequest) (*blackjackPb.Game, error) {
	if len(request.GetPlayerNames()) == 0 {
		return nil, fmt.Errorf("no player names provided")
	}

	game, err := blackjack.NewGame(ctx, s.deckClient, s.aiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	err = game.Start(ctx, request.GetPlayerNames()...)
	if err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}

	err = s.store.Put(ctx, game.ID, game)
	if err != nil {
		return nil, fmt.Errorf("failed to store game: %w", err)
	}

	protoGame := gameToProto(game)
	hideHandForPlayer(protoGame, "")
	return protoGame, nil
}

// map game to proto
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

func (s *ServiceServer) GetGame(ctx context.Context, req *blackjackPb.GetGameRequest) (*blackjackPb.Game, error) {
	game, err := s.store.Get(ctx, req.GetGameId())
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	protoGame := gameToProto(game)
	hideHandForPlayer(protoGame, req.GetPlayerName())
	return protoGame, nil
}
