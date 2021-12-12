package service

import (
	"math/big"
	"testing"

	"github.com/silviutroscot/istari-challenge/pkg/log"
	"github.com/stretchr/testify/assert"
)

func Test_SwapTokens(t *testing.T) {
	// allow the tests to run in parallel
	t.Parallel()

	t.Run("equal USD value", func(t *testing.T) {
		// arrange test
		firstTokenValueInUSD := big.NewFloat(0.0009656114709)
		secondTokenValueInUSD := big.NewFloat(0.0009656114709)
		tokenAmount := big.NewFloat(329918.56931317)

		// act
		secondTokenConvertedAmount := SwapTokens(tokenAmount, firstTokenValueInUSD, secondTokenValueInUSD)
		firstTokenConvertedAmount := SwapTokens(tokenAmount, secondTokenValueInUSD, firstTokenValueInUSD)

		// assert
		assert.Equal(t, &tokenAmount, &secondTokenConvertedAmount,
			"the amount of secondToken received after swapping %d firstTokens is %d and the expected amount is %d",
			&tokenAmount, &secondTokenConvertedAmount, &tokenAmount)
		assert.Equal(t, &tokenAmount, &firstTokenConvertedAmount,
			"the amount of firstToken received after swapping %d secondTokens is %d and the expected amount is %d",
			&tokenAmount, &firstTokenConvertedAmount, &tokenAmount)
	})

	t.Run("first token USD value smaller than second token USD value", func(t *testing.T) {
		// arrange test
		firstTokenValueInUSD := big.NewFloat(0.00019940206605890213786)
		secondTokenValueInUSD := big.NewFloat(5456.12)
		tokenAmount := big.NewFloat(1120)
		expectedSecondTokenAmountAfterSwap := big.NewFloat(0.0000409320751717283333)
		expectedFirstTokenAmountAfterSwap := big.NewFloat(30645893098.193332559991599)

		// act
		secondTokenConvertedAmount := SwapTokens(tokenAmount, firstTokenValueInUSD, secondTokenValueInUSD)
		firstTokenConvertedAmount := SwapTokens(tokenAmount, secondTokenValueInUSD, firstTokenValueInUSD)

		// assert and round to precision=30
		assert.Equal(t, expectedSecondTokenAmountAfterSwap.Cmp(secondTokenConvertedAmount.SetPrec(30)), 0
			"the amount of secondToken received after swapping %v firstTokens is %v and the expected amount is %v",
			tokenAmount, secondTokenConvertedAmount, expectedSecondTokenAmountAfterSwap)
		assert.Equal(t, expectedFirstTokenAmountAfterSwap.SetPrec(30), firstTokenConvertedAmount.SetPrec(30),
			"the amount of firstToken received after swapping %v secondTokens is %v and the expected amount is %v",
			tokenAmount, firstTokenConvertedAmount, expectedFirstTokenAmountAfterSwap)
	})

	t.Run("first token USD value larger than second token USD value", func(t *testing.T) {
		// arrange test
		firstTokenValueInUSD := big.NewFloat(48896.02)
		secondTokenValueInUSD := big.NewFloat(0.798416)
		tokenAmount := big.NewFloat(0.3)
		expectedSecondTokenAmountAfterSwap := big.NewFloat(18372.384821947453)
		expectedFirstTokenAmountAfterSwap := big.NewFloat(0.000004898656373259)

		// act
		secondTokenConvertedAmount := SwapTokens(tokenAmount, firstTokenValueInUSD, secondTokenValueInUSD)
		firstTokenConvertedAmount := SwapTokens(tokenAmount, secondTokenValueInUSD, firstTokenValueInUSD)

		// assert and round to precision=30
		assert.Equal(t, expectedSecondTokenAmountAfterSwap.SetPrec(30), secondTokenConvertedAmount.SetPrec(30),
			"the amount of secondToken received after swapping %v firstTokens is %v and the expected amount is %v",
			tokenAmount, secondTokenConvertedAmount, expectedSecondTokenAmountAfterSwap)
		assert.Equal(t, expectedFirstTokenAmountAfterSwap.SetPrec(30), firstTokenConvertedAmount.SetPrec(30),
			"the amount of firstToken received after swapping %v secondTokens is %v and the expected amount is %v",
			tokenAmount, firstTokenConvertedAmount, expectedFirstTokenAmountAfterSwap)
	})
}
