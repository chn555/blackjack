/*
Package blackjack provides a simple blackjack game implementation.

It comprises the following base components:

A Hand represents a player's hand in the game, consisting of the players cards.

A Player represents a player in the game, consisting of the player's hand, their score, and whether they have gone bust.

A Game represents the state of the game, including the players, their order, and the current status of the game.

When a game is created, the dealer is always the first player.
The game starts with the dealer and players being dealt two cards each.

To progress the game each player in turn can choose to hit or stand.
To play a turn, call the PlayTurn method on the game.
After a turn has been played, the game state is updated and the Game.NextPlayer is updated.
The dealer always stands.
*/
package blackjack
