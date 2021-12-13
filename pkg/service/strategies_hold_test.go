package service

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_HoldStrategy_TargetPriceEqualsCurrentPrice(t *testing.T) {
	service := Service{}

	// test for both Egld and MEX that there are no gains and the USD balance uses the target prices
	t.Run("target price equals initial price", func(t *testing.T) {
		// arrange test
		input := StrategiesInput{
			MexTargetPrice: big.NewFloat(0.000777909073989448632753),
			EgldTargetPrice: big.NewFloat(285.433),
			EgldTokensInvested: big.NewFloat(0.45),
			MexTokensInvested: big.NewFloat(94869182.3086),
		}

		expectedMexUSDValue := big.Float{}
		expectedMexUSDValue.Mul(input.MexTargetPrice, input.MexTokensInvested)

		expectedEgldUSDValue := big.Float{}
		expectedEgldUSDValue.Mul(input.EgldTargetPrice, input.EgldTokensInvested)
		// act
		mexHoldStrategyResult, err := service.HoldStrategy(TokenTypeMex, &input)

		// assert
		assert.Nil(t, err, "expected no error from MEX HOLD strategy, got %s", err)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *mexHoldStrategyResult.ProfitInEgld),
			"the EGLD profit %v is different from the expected EGLD profit %v", mexHoldStrategyResult.ProfitInEgld, BigFloatZero)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *mexHoldStrategyResult.ProfitInMex),
			"the MEX profit %v is different from the expected MEX profit %v", &mexHoldStrategyResult.ProfitInMex, BigFloatZero)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *mexHoldStrategyResult.ProfitInUSD),
			"the profit USD value %v is different from the expected profit value %v", &mexHoldStrategyResult.ProfitInUSD, BigFloatZero)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *mexHoldStrategyResult.TotalBalanceInEgld),
			"the EGLD balance %d is different from the expected EGLD balance %d", &mexHoldStrategyResult.TotalBalanceInEgld, BigFloatZero)
		assert.True(t, BigFloatsAreEqual(*input.MexTokensInvested, *mexHoldStrategyResult.TotalBalanceInMex),
			"the MEX balance %d is different from the expected MEX balance %d", &mexHoldStrategyResult.TotalBalanceInMex, input.MexTokensInvested)
		assert.True(t, BigFloatsAreEqual(expectedMexUSDValue, *mexHoldStrategyResult.TotalBalanceInUsd),
			"the USD value of Mex %d is different from the expected USD balance %d", &mexHoldStrategyResult.TotalBalanceInUsd, expectedMexUSDValue)

		// act
		egldHoldStrategyResult, err := service.HoldStrategy(TokenTypeEgld, &input)
		// assert
		assert.Nil(t, err, "expected no error from EGLD HOLD strategy, got %s", err)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *egldHoldStrategyResult.ProfitInEgld),
			"the EGLD profit %v is different from the expected EGLD profit %v", &egldHoldStrategyResult.ProfitInEgld, BigFloatZero)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *egldHoldStrategyResult.ProfitInMex),
			"the MEX profit %v is different from the expected MEX profit %v", &egldHoldStrategyResult.ProfitInMex, BigFloatZero)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *egldHoldStrategyResult.ProfitInUSD),
			"the profit USD value %v is different from the expected profit value %v", &egldHoldStrategyResult.ProfitInUSD, BigFloatZero)
		assert.True(t, BigFloatsAreEqual(*input.EgldTokensInvested, *egldHoldStrategyResult.TotalBalanceInEgld),
			"the EGLD balance %d is different from the expected EGLD balance %d", &egldHoldStrategyResult.TotalBalanceInEgld, input.EgldTokensInvested)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *egldHoldStrategyResult.TotalBalanceInMex),
			"the MEX balance %d is different from the expected MEX balance %d", &egldHoldStrategyResult.TotalBalanceInMex, BigFloatZero)
		assert.True(t, BigFloatsAreEqual(expectedEgldUSDValue, *egldHoldStrategyResult.TotalBalanceInUsd),
			"the USD value of Mex %d is different from the expected USD balance %d", egldHoldStrategyResult.TotalBalanceInUsd, expectedEgldUSDValue)
	})
}

