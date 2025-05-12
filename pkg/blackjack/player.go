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

func (p *Player) calculateScore() {
	points := 0

	for _, card := range p.Hand.Cards {
		value := int(card.GetValue())
		if value == 1 && points+11 <= 21 {
			value = 11 // Ace can be 1 or 11, prefer 11 if it doesn't bust
		}
		if value > 10 {
			value = 10 // Face cards are worth 10
		}
		points += value
	}

	p.Score = points
}
