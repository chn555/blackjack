package blackjack

type Player struct {
	Hand  *Hand
	Score int
	Bust  bool
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
