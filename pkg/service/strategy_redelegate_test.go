package service

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// todo: test that redelegation with a redelegation interval > investment interval is equal to the result from staking on that investment interval
// todo: add more tests
func TestService_RedelegateStrategy(t *testing.T) {
	// allow the tests to run in parallel
	t.Parallel()

	// test setup
	service := Service{}

	t.Run("redelegate weekly for one year", func(t *testing.T) {
		// todo: the apy here is calculated using an online calculator; it does not have enough accuracy
		// arrange the test (inputs and expected values)
		input := StrategiesInput{
			EgldTokensInvested: big.NewFloat(4.8),
			EgldTargetPrice: big.NewFloat(380.54),
			EgldAPR: big.NewFloat(13.4),
			MexTargetPrice: big.NewFloat(0.00057854),
			MexTokensInvested: big.NewFloat(498345.7098),
			MexAPRUnlocked: big.NewFloat(50.0),
			MexAPRLocked: big.NewFloat(1254.43),
			InvestmentDurationInDays: 365,
			RedelegationIntervalInDays: 7,
			RewardsInLockedMEX: false,
		}

		initialMexPrice := big.NewFloat(0.00194567073989448632753)
		initialEgldPrice := big.NewFloat(240.50)

		// set the expected values
		expectedEgldAPY := big.NewFloat(14.319654)
		expectedMexAPY := big.NewFloat(64.48)

		expectedMexEarned := &big.Float{}
		expectedMexEarned.Mul(input.MexTokensInvested, expectedMexAPY)
		expectedMexEarnedInUSD := &big.Float{}
		expectedMexEarnedInUSD.Mul(expectedMexEarned, initialMexPrice)

		expectedEgldEarned := &big.Float{}
		expectedEgldEarned.Mul(input.EgldTokensInvested, expectedEgldAPY)
		expectedEgldEarned.Quo(expectedEgldEarned, BigFloatOneHundred)
		expectedEgldEarnedInUSD := &big.Float{}
		expectedEgldEarnedInUSD.Mul(expectedEgldEarned, initialEgldPrice)

		expectedTotalEgld := &big.Float{}
		expectedTotalEgld.Add(expectedEgldEarned, input.EgldTokensInvested)
		expectedTotalEgldInUSD := &big.Float{}
		expectedTotalEgldInUSD.Mul(expectedTotalEgld, input.EgldTargetPrice)

		// act
		result, err := service.RedelegateStrategy(TokenTypeEgld, &input, initialEgldPrice)

		// assert
		// todo: use a more reliable comparison way than converting it to string and taking the first 5 decimals
		assert.Nil(t, err, "expected no error from EGLD REDELEGATE strategy, got %s", err)
		assert.Equal(t, expectedEgldEarned.Text('f', 5), result.ProfitInEgld.Text('f', 5),
			"the amount of EGLD earned %v is different from the expected amount of EGLD earned %v", result.ProfitInEgld, expectedEgldEarned)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *result.ProfitInMex),
			"the amount of EGLD earned %v is different frm the expected amount of EGLD earned %v", result.ProfitInMex, expectedEgldEarned)
		assert.Equal(t, expectedEgldEarnedInUSD.Text('f', 3), result.ProfitInUSD.Text('f', 3),
			"the USD value of the earned EGLD %v is different from the expected EGLD value of the earned MEX %v",
			result.ProfitInUSD, &expectedEgldEarnedInUSD)
		assert.Equal(t, expectedTotalEgld.Text('f', 5), result.TotalBalanceInEgld.Text('f', 5),
			"the EGLD balance %v is different from the expected EGLD balance %v", result.TotalBalanceInEgld, expectedTotalEgld)
		assert.True(t, BigFloatsAreEqual(*BigFloatZero, *result.TotalBalanceInMex),
			"the MEX balance %v is different from the expected MEX balance %v", result.TotalBalanceInMex, BigFloatZero)
		assert.Equal(t, expectedTotalEgldInUSD.Text('f', 2), result.TotalBalanceInUsd.Text('f', 2),
			"the USD value of Mex %v is different from the expected balance %v", result.TotalBalanceInUsd, expectedTotalEgldInUSD)
	})

	t.Run("redelegation period larger than investment period", func(t *testing.T) {
		input := StrategiesInput{
			EgldTokensInvested: big.NewFloat(10),
			EgldTargetPrice: big.NewFloat(390.4),
			EgldAPR: big.NewFloat(15),
			MexTargetPrice: big.NewFloat(0.00057854),
			MexTokensInvested: big.NewFloat(498345.7098),
			MexAPRUnlocked: big.NewFloat(50.0),
			MexAPRLocked: big.NewFloat(1254.43),
			InvestmentDurationInDays: 20,
			RedelegationIntervalInDays: 30,
			RewardsInLockedMEX: false,
		}

		egldInitialPrice := big.NewFloat(280)
		mexInitialPrice := big.NewFloat(0.004)

		// for this test case, the result should be equal to just staking the tokens
		egldStakeResult, err := service.StakeStrategy(TokenTypeEgld, &input, egldInitialPrice)
		assert.Nil(t, err, "expected no error from EGLD Staking strategy, got %s", err)

		egldRedelegateResult, err := service.RedelegateStrategy(TokenTypeEgld, &input, egldInitialPrice)
		assert.Nil(t, err, "expected no error from EGLD Redelegate strategy, got %s", err)

		// assert that the results are equal
		assert.True(t, egldStakeResult.Equals(egldRedelegateResult), "expected strategies results for stake " +
			"and stake+redelegate to be equal")

		mexStateResult, err := service.StakeStrategy(TokenTypeMex, &input, mexInitialPrice)
		assert.Nil(t, err, "expected no error from MEX Staking strategy, got %s", err)

		mexRedelegateResult, err := service.RedelegateStrategy(TokenTypeMex, &input, mexInitialPrice)
		assert.Nil(t, err, "expected no error from MEX Redelegate strategy, got %s", err)

		// assert that the results are equal
		assert.True(t, mexStateResult.Equals(mexRedelegateResult), "expected strategies results for stake " +
			"and stake+redelegate to be equal")
	})
}
