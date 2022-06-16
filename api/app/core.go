package app

import (
	"crypto/rand"
	"math"
	"math/big"
)

// GenerateRandomNumber generates random integer of n digits.
func GenerateRandomNumber(numberOfDigits int) (uint64, error) {
	maxLimit := uint64(int(math.Pow10(numberOfDigits)) - 1)
	lowLimit := uint64(math.Pow10(numberOfDigits - 1))

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(maxLimit)))
	if err != nil {
		return 0, err
	}
	randomNumberInt := randomNumber.Uint64()

	// Handling integers between 0, 10^(n-1) .. for n=4, handling cases between (0, 999)
	if randomNumberInt <= lowLimit {
		randomNumberInt += lowLimit
	}

	// Never likely to occur, kust for safe side.
	if randomNumberInt > maxLimit {
		randomNumberInt = maxLimit
	}
	return randomNumberInt, nil
}
