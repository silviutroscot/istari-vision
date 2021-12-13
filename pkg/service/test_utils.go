package service

import (
	"math"
	"math/big"
)

var (
	EPSILON = big.NewFloat(1/math.Pow(10, 10))
)

func BigFloatsAreEqual(x, y big.Float) bool {

	return x.Sub(&x, &y).Abs(&x).Cmp(EPSILON) < 1
}
