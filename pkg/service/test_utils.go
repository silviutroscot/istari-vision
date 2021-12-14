package service

import (
	"math"
	"math/big"
)

var (
	EPSILON = big.NewFloat(1/math.Pow(10, 10))
)

func BigFloatsAreEqual(x, y big.Float) bool {
	copyX := &big.Float{}
	copyY := &big.Float{}
	x.Copy(copyX)
	y.Copy(copyY)
	return copyX.Sub(copyX, copyY).Abs(copyX).Cmp(EPSILON) < 1
}
