package service

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_StakeStrategy(t *testing.T) {
	// allow the tests to run in parallel
	t.Parallel()

	// test setup
	service := Service{}

	// test for MEX that the rewards in MEX, USD and the total portfolio in MEX and USD reflect the earnings
	t.Run("staking rewards for MEX oven one year", func(t *testing.T) {
		// arrange the test (inputs and expected values)
		input := StrategiesInput{
			MexTargetPrice: big.NewFloat(0.003),
			MexTokensInvested: big.NewFloat(498345.7098),
			InvestmentDurationInDays: 365,
			RewardsInLockedMEX: false,
			MexAPRUnlocked: big.NewFloat(50.0),
		}
		initialMexPrice := big.NewFloat(0.00194567073989448632753)

		apr := big.NewFloat(0.5)

		// set the expected values for the MEX owned at the end of staking and its value in USD
		expectedMexEarned := big.Float{}
		expectedMexEarned.Mul(input.MexTokensInvested, apr)

		var expectedMexEarnedInUSD big.Float
		expectedMexEarnedInUSD.Mul(&expectedMexEarned, initialMexPrice)

		var expectedTotalMexBalance big.Float
		expectedTotalMexBalance.Add(input.MexTokensInvested, &expectedMexEarned)

		var expectedTotalBalanceInUSD big.Float
		expectedTotalBalanceInUSD.Mul(&expectedTotalMexBalance, input.MexTargetPrice)

		// act
		result, err := service.StakeStrategy(TokenTypeMex, &input, initialMexPrice)

		// assert
		// todo: investigate why these tests fail even though the actual and expected values are 0
		assert.Nil(t, err, "expected no error from MEX STAKE strategy, got %s", err)
		assert.Equal(t, BigFloatZero.Cmp(result.ProfitInEgld), 0,
			"the amount of EGLD earned %v is different from the expected amount of EGLD earned %v", result.ProfitInEgld, BigFloatZero)
		assert.Equal(t, expectedMexEarned.Cmp(result.ProfitInMex), 0,
			"the amount of MEX earned %v is different frm the expected amount of MEX earned %v", result.ProfitInMex, expectedMexEarned)
		assert.Equal(t, expectedMexEarnedInUSD.Cmp(result.ProfitInUSD), 0,
			"the USD value of the earned MEX %v is different from the expected USD value of the earned MEX %v",
			result.ProfitInUSD, expectedMexEarnedInUSD)
		assert.Equal(t, BigFloatZero.Cmp(result.TotalBalanceInEgld), 0,
			"the EGLD balance %v is different from the expected EGLD balance %v", result.TotalBalanceInEgld, BigFloatZero)
		assert.Equal(t, expectedTotalMexBalance.Cmp(result.TotalBalanceInMex), 0,
			"the MEX balance %v is different from the expected MEX balance %v", result.TotalBalanceInMex, expectedTotalMexBalance)
		assert.Equal(t, expectedTotalBalanceInUSD.Cmp(result.TotalBalanceInUsd), 0,
			"the USD value of Mex %v is different from the expected balance %v", result.TotalBalanceInUsd, &expectedTotalBalanceInUSD)
	})

	// test for MEX that the rewards in MEX, USD and the total portfolio in MEX and USD reflect the earnings
	t.Run("staking rewards for EGLD for a month", func(t *testing.T) {
		// arrange the test (inputs and expected values)
		input := StrategiesInput{
			EgldTokensInvested: big.NewFloat(3.546),
			EgldTargetPrice: big.NewFloat(400.34),
			EgldAPR: big.NewFloat(9),
			PercentageOfPortfolioInEgld: big.NewFloat(9.0),
			InvestmentDurationInDays: 30,
		}

		initialEgldPrice := big.NewFloat(420.567)
		apr := big.NewFloat(0.09)

		// set the expected values for the EGLD owned at the end of staking and its value in USD
		egldUsdValue := big.Float{}
		egldUsdValue.Mul(initialEgldPrice, input.EgldTokensInvested)

		// yearPercentageForRewards stores the division between the staking period and the duration of the year

		redelegationInterval := big.NewFloat(float64(input.InvestmentDurationInDays))
		yearPercentageForRewards := big.Float{}
		yearPercentageForRewards.Quo(redelegationInterval, BigFloatDaysInYear)

		expectedEgldEarned := big.Float{}
		expectedEgldEarned.Mul(input.EgldTokensInvested, apr)
		expectedEgldEarned.Mul(&expectedEgldEarned, &yearPercentageForRewards)

		expectedEgldEarnedInUSD := big.Float{}
		expectedEgldEarnedInUSD.Mul(&expectedEgldEarned, initialEgldPrice)

		expectedTotalEgldBalance := big.Float{}
		expectedTotalEgldBalance.Add(input.EgldTokensInvested, &expectedEgldEarned)

		var expectedTotalBalanceInUSD big.Float
		expectedTotalBalanceInUSD.Mul(&expectedTotalEgldBalance, input.EgldTargetPrice)

		// act
		result, err := service.StakeStrategy(TokenTypeEgld, &input, initialEgldPrice)
		// assert
		assert.Nil(t, err, "expected no error from EGLD STAKE strategy, got %s", err)
		assert.True(t, BigFloatsAreEqual(expectedEgldEarned, *result.ProfitInEgld),
			"the amount of EGLD earned %v is different from the expected amount of EGLD earned %v", result.ProfitInEgld, &expectedEgldEarned)
		assert.Equal(t, BigFloatZero.Cmp(result.ProfitInMex), 0,
			"the amount of EGLD earned %v is different frm the expected amount of EGLD earned %v", result.ProfitInMex, expectedEgldEarned)
		assert.Equal(t, expectedEgldEarnedInUSD.Cmp(result.ProfitInUSD), 0,
			"the USD value of the earned EGLD %v is different from the expected EGLD value of the earned MEX %v",
			&result.ProfitInUSD, &expectedEgldEarnedInUSD)
		assert.Equal(t, expectedTotalEgldBalance.Cmp(result.TotalBalanceInEgld), 0,
			"the EGLD balance %v is different from the expected EGLD balance %v", result.TotalBalanceInEgld, expectedTotalEgldBalance)
		assert.Equal(t, BigFloatZero.Cmp(result.TotalBalanceInMex), 0,
			"the MEX balance %v is different from the expected MEX balance %v", result.TotalBalanceInMex, BigFloatZero)
		assert.Equal(t, expectedTotalBalanceInUSD.Cmp(result.TotalBalanceInUsd), 0,
			"the USD value of Mex %v is different from the expected balance %v", result.TotalBalanceInUsd, expectedTotalBalanceInUSD)
	})

	// todo: add more tests for multiple scenarios (staking periods, prices)
}
