package blackjack

// Player represents a player in the game, consisting of the player's hand, their score, and whether they have gone bust.
type Player struct {
	// Hand represents the player's hand in the game.
	Hand *Hand
	// Score represents the player's score in the game.
	Score int
	// Bust indicates whether the player has gone bust (over 21).
	Bust bool
}
