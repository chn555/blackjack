package blackjack

import (
	"context"
	"fmt"
)

const blackjack = 21

// TurnAction represents the action taken by the player during their turn
type TurnAction uint8

const (
	Unknown TurnAction = iota
	Hit
	Stand
)

// Turn represents a player's turn in the game
type Turn struct {
	Action     TurnAction
	PlayerName string
}

// NewTurn creates a new turn with the given action
func NewTurn(action TurnAction, playerName string) *Turn {
	return &Turn{
		Action:     action,
		PlayerName: playerName,
	}
}

// PlayTurn processes the player's turn
func (g *Game) PlayTurn(ctx context.Context, turn *Turn) error {
	playerName := turn.PlayerName

	if g.playerIsOutOfTurn(playerName) {
		return fmt.Errorf("not your turn")
	}

	player, ok := g.Players[playerName]
	if !ok {
		return fmt.Errorf("player %s not found", playerName)
	}

	err := g.takeAction(ctx, turn, player)
	if err != nil {
		return fmt.Errorf("failed to take action: %w", err)
	}

	player.calculateScore()
	g.Players[playerName] = player

	if playerHitBlackjack(player) {
		g.Status = Done
		g.Winner = playerName
		return nil
	}

	if playerWentOver(player) {
		player.Bust = true
		g.removePlayer(playerName)
		g.Players[playerName] = player

		if g.onlyOnePlayerLeft() {
			g.Status = Done
			g.Winner = g.playerOrder[0]
			return nil
		}
	}

	g.NextPlayer = g.getNextPlayer(playerName)
	return nil
}

func (g *Game) onlyOnePlayerLeft() bool {
	return len(g.playerOrder) == 1
}

func playerWentOver(player *Player) bool {
	return player.Score > blackjack
}

func playerHitBlackjack(player *Player) bool {
	return player.Score == blackjack
}

func (g *Game) playerIsOutOfTurn(playerName string) bool {
	return g.Status != WaitingForPlayer || g.NextPlayer != playerName
}

func (g *Game) takeAction(ctx context.Context, turn *Turn, player *Player) error {
	switch turn.Action {
	case Hit:
		err := player.Hand.PullCard(ctx)
		if err != nil {
			return fmt.Errorf("failed to draw card: %w", err)
		}
	case Stand:
	default:
		return fmt.Errorf("unknown turn: %v", turn)
	}
	return nil
}

// remove player
func (g *Game) removePlayer(playerName string) {
	for i, player := range g.playerOrder {
		if player == playerName {
			g.playerOrder = append(g.playerOrder[:i], g.playerOrder[i+1:]...)
			break
		}
	}
}

func (g *Game) getNextPlayer(currentPlayer string) (playerName string) {
	for i, player := range g.playerOrder {
		if player == currentPlayer {
			if i+1 < len(g.playerOrder) {
				return g.playerOrder[i+1]
			}
			return g.playerOrder[0]
		}
	}
	return ""
}

func (g *Game) playDealerHand(ctx context.Context) error {
	return g.PlayTurn(ctx,
		NewTurn(
			Stand, DealerName,
		),
	)
}
