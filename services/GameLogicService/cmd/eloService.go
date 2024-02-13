package cmd

import (
	"math"
)

func CalculateEloRating(player1Rating, player2Rating, result float64) (float64, float64) {
	// You can adjust the k-factor based on the game's characteristics
	kFactor := 32.0

	expectedResultPlayer1 := 1 / (1 + math.Pow(10, (player2Rating-player1Rating)/400))
	expectedResultPlayer2 := 1 / (1 + math.Pow(10, (player1Rating-player2Rating)/400))

	newRatingPlayer1 := player1Rating + kFactor*(result-expectedResultPlayer1)
	newRatingPlayer2 := player2Rating + kFactor*((1-result)-expectedResultPlayer2)

	return math.Round(newRatingPlayer1), math.Round(newRatingPlayer2)
}