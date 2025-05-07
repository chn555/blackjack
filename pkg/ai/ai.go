package ai

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	blackjackPb "github.com/chn555/schemas/proto/blackjack/v1"
)

type Strategy interface {
	DecideAction(hand *blackjackPb.Hand) (blackjackPb.Turn_TURN_ACTION, error)
}

type AI struct {
	gameClient blackjackPb.BlackjackServiceClient
	games      map[string]*Game
}

type Game struct {
	gameID     string
	playerName string

	strategy Strategy
}

func NewAI(gameClient blackjackPb.BlackjackServiceClient) *AI {
	return &AI{
		gameClient: gameClient,
		games:      make(map[string]*Game),
	}
}

func (a *AI) Start() {
	go a.startGameLoop()
}

func (a *AI) startGameLoop() {
	for range time.Tick(5 * time.Second) {
		for _, game := range a.games {
			if game.strategy == nil {
				delete(a.games, game.gameID)
				slog.Info("no strategy, removed from loop",
					"gameID", game.gameID,
					"playerName", game.playerName)
				continue
			}

			gameState, err := a.gameClient.GetGame(context.Background(), &blackjackPb.GetGameRequest{
				GameId:     game.gameID,
				PlayerName: game.playerName,
			},
			)
			if err != nil {
				delete(a.games, game.gameID)
				slog.Error("failed to get game state",
					"error", err, "gameID", game.gameID,
					"playerName", game.playerName)
				continue
			}

			if gameState.GetStatus() == blackjackPb.Game_GAME_STATUS_DONE {
				delete(a.games, game.gameID)
				slog.Info("game is done, removed from loop",
					"gameID", game.gameID,
					"playerName", game.playerName)
				continue
			}

			if gameState.GetNextPlayer() != game.playerName {
				slog.Info("not my turn",
					"gameID", game.gameID,
					"playerName", game.playerName)
				continue
			}

			playerHand, ok := gameState.PlayerHands[game.playerName]
			if !ok {
				delete(a.games, game.gameID)
				slog.Error("failed to get player hand",
					"error", err, "gameID", game.gameID,
					"playerName", game.playerName)
				continue
			}

			action, err := game.strategy.DecideAction(playerHand)
			if err != nil {
				slog.Error("failed to decide action",
					"error", err, "gameID", game.gameID,
					"playerName", game.playerName)
			}
			resp, err := a.gameClient.PlayTurn(context.Background(), &blackjackPb.Turn{
				GameId:     game.gameID,
				PlayerName: game.playerName,
				Action:     action,
			})
			if err != nil {
				slog.Error("failed to play turn",
					"error", err, "gameID", game.gameID,
					"playerName", game.playerName)
			} else {
				slog.Info("played turn",
					"gameID", game.gameID,
					"playerName", game.playerName,
					"response", resp)
			}
		}
	}
}

func (a *AI) AddGame(gameID string, playerName string, strat Strategy) error {
	_, exists := a.games[gameID]
	if exists {
		return fmt.Errorf("game %s already exists", gameID)
	}

	game := &Game{
		gameID:     gameID,
		playerName: playerName,
		strategy:   strat,
	}
	a.games[gameID] = game
	return nil
}
