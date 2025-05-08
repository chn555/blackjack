package blackjack

import (
	"context"
	"fmt"
	"slices"
)

const blackjack = 21

// TurnAction represents the action taken by the player during their turn.
type TurnAction uint8

const (
	Unknown TurnAction = iota
	Hit                // pull another card into the player's hand.
	Stand              // keep the current hand.
)

// Turn represents a player's turn in the game.
type Turn struct {
	// Action is the action taken by the player during their turn.
	Action TurnAction
	// PlayerName is the name of the player taking the turn.
	PlayerName string
}

// NewTurn creates a new turn with the given action.
func NewTurn(action TurnAction, playerName string) *Turn {
	return &Turn{
		Action:     action,
		PlayerName: playerName,
	}
}

// PlayTurn processes the player's turn.
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

	// After the turn action has been executed,
	// we need to check the player's score and modify the game state if needed.
	player.calculateScore()
	g.Players[playerName] = player

	if playerHitBlackjack(player) {
		g.Status = Done
		g.Winner = playerName
		return nil
	}

	if playerWentOver(player) {
		player.Bust = true
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
	playersLeft := 0
	for _, player := range g.Players {
		if !player.Bust {
			playersLeft++
		}
	}
	return playersLeft == 1
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

func (g *Game) removePlayer(playerName string) {
	for i, player := range g.playerOrder {
		if player == playerName {
			g.playerOrder = append(g.playerOrder[:i], g.playerOrder[i+1:]...)
			break
		}
	}
}

func (g *Game) getNextPlayer(currentPlayer string) string {
	curPlayerIndex := slices.Index(g.playerOrder, currentPlayer)
	// if the player is not in the game, return empty string
	if curPlayerIndex == -1 {
		return ""
	}

	for i := range g.playerOrder {
		// We skip the current player and all players that are before them
		if i <= curPlayerIndex {
			continue
		}

		playerName := g.playerOrder[i]
		player := g.Players[playerName]
		if player.Bust {
			continue
		}

		return playerName
	}

	// If we reach here, it means all players after the current player are bust,
	// so we need to loop back to the beginning of the player order
	// if the only player left is the current player, he will be caught in the loop
	for i := range g.playerOrder {
		playerName := g.playerOrder[i]
		player := g.Players[playerName]
		if player.Bust {
			continue
		}

		return playerName
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
