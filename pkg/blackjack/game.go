package blackjack

// GameStatus represents the current state of the game.
type GameStatus uint8

const (
	UnknownState GameStatus = iota
	WaitingToStart
	WaitingForPlayer
	Done
)

// Game represents the state of a blackjack game.
type Game struct {
	// ID is the unique identifier for the game.
	ID string
	// Status is the current status of the game.
	Status GameStatus

	// Players contains all players in the game.
	Players map[string]*Player
	// PlayerOrder is the order of player turns in the game.
	playerOrder []string
	// NextPlayer is the player whose turn it is.
	NextPlayer string
	// Winner is the winner of the game.
	// If the game is still in-progress, this field is empty.
	Winner string
}
