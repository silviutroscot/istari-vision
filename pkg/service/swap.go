package service

import (
	"math/big"
)

const (
	FloatPrecision = 30
)

// SwapTokens returns the amount of destinationToken corresponding to the amount of firstToken provided, based on their USD price
func SwapTokens(amount, firstTokenValueInUSD, destinationTokenValueInUSD *big.Float) *big.Float {
	egldToMexRate := &big.Float{}
	egldToMexRate.Quo(firstTokenValueInUSD, destinationTokenValueInUSD)

	convertedMexAmount := &big.Float{}
	convertedMexAmount.Mul(amount, egldToMexRate)

	return convertedMexAmount
}