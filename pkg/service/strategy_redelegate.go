package service

import (
	"fmt"
	"math/big"

	"github.com/silviutroscot/istari-vision/pkg/log"
)

// RedelegateStrategy returns a StrategyResult representing the result of staking and redelegating the profit each RedelegateIntervalInDays days
func (s *Service) RedelegateStrategy(tokenType TokenType, input *StrategiesInput, tokenInitialPrice *big.Float) (*StrategyResult, error) {
	redelegateIntervalFloat := big.NewFloat(float64(input.RedelegationIntervalInDays))

	// tokenBalance represents the current tokens we have
	tokenBalance := &big.Float{}
	tokenBalance.Copy(input.EgldTokensInvested)
	initialTokenBalance := input.EgldTokensInvested
	tokenAPR := input.EgldAPR
	targetPrice := input.EgldTargetPrice

	if tokenType == TokenTypeMex {
		tokenBalance.Copy(input.MexTokensInvested)
		initialTokenBalance = input.MexTokensInvested
		targetPrice = input.MexTargetPrice
		if input.RewardsInLockedMEX {
			tokenAPR = input.MexAPRLocked
		} else {
			tokenAPR = input.MexAPRUnlocked
		}
	}

	currentUSDValueFloat := big.Float{}
	currentUSDValueFloat.Mul(tokenInitialPrice, tokenBalance)

	// calculate the APR for the redelegation period (i.e. the percentage of initial investment the user gets between 2 redelegations)
	daysInYear := big.NewFloat(365.0)
	aprPerCycle := &big.Float{}
	aprPerCycle.Quo(redelegateIntervalFloat, daysInYear)

	APRToBeReceivedInOneRedelegationCycle := &big.Float{}
	APRToBeReceivedInOneRedelegationCycle.Mul(aprPerCycle, tokenAPR)
	APRToBeReceivedInOneRedelegationCycle.Quo(APRToBeReceivedInOneRedelegationCycle, BigFloatOneHundred)

	log.Info("input.InvestmentDurationInDays is %d", input.InvestmentDurationInDays)
	log.Info("input.RedelegationIntervalInDays is %d", input.RedelegationIntervalInDays)
	// compound the interest for the number of redelegations cycles
	cycleRewardsDays := input.RedelegationIntervalInDays
	for ; cycleRewardsDays <= input.InvestmentDurationInDays; cycleRewardsDays += input.RedelegationIntervalInDays {
		interestReceived := &big.Float{}
		interestReceived.Mul(APRToBeReceivedInOneRedelegationCycle, tokenBalance)
		log.Info("in REDELEGATION the earned interest for cycleDays value %d is %s", cycleRewardsDays, interestReceived.String())
		tokenBalance.Add(tokenBalance, interestReceived)
	}

	cycleRewardsDays = cycleRewardsDays - input.RedelegationIntervalInDays

	log.Info("tokenBalance after cycle is %+v", tokenBalance)
	log.Info("cycleRewardsDays is %+v", cycleRewardsDays)
	log.Info("input.InvestmentDurationInDays is %v", input.InvestmentDurationInDays)
	// compute the rewards for the days left between the last redelegation cycle and the remaining days
	var remainingDays int
	if input.InvestmentDurationInDays >= cycleRewardsDays {
		remainingDays = input.InvestmentDurationInDays - cycleRewardsDays
	} else {
		remainingDays = input.InvestmentDurationInDays
	}
	log.Info("the remaining days are %d", remainingDays)

	remainingDaysFloat := big.NewFloat(float64(remainingDays))

	remainingDaysAsPercentageOfYear := &big.Float{}
	remainingDaysAsPercentageOfYear.Quo(remainingDaysFloat, daysInYear)

	interestReceivedForRemainingDays := &big.Float{}
	interestReceivedForRemainingDays.Mul(remainingDaysAsPercentageOfYear, tokenAPR)
	interestReceivedForRemainingDays.Mul(interestReceivedForRemainingDays, tokenBalance)
	interestReceivedForRemainingDays.Quo(interestReceivedForRemainingDays, BigFloatOneHundred)
	tokenBalance.Add(tokenBalance, interestReceivedForRemainingDays)

	earnedInterestInTokens := &big.Float{}
	earnedInterestInTokens.Sub(tokenBalance, initialTokenBalance)
	// calculate the ROI as (earnedtokens / initialTokens)
	roi := &big.Float{}
	roi.Quo(earnedInterestInTokens, initialTokenBalance)
	roi.Mul(roi, BigFloatOneHundred)

	interestValueInUSD := &big.Float{}
	interestValueInUSD.Mul(earnedInterestInTokens, tokenInitialPrice)

	result := NewStrategyResult()
	result.ProfitInUSD.Copy(interestValueInUSD)
	result.TotalBalanceInUsd.Mul(tokenBalance, targetPrice)
	result.ROI.Copy(roi)

	switch tokenType {
	case TokenTypeEgld:
		result.ProfitInEgld.Copy(earnedInterestInTokens)
		result.TotalBalanceInEgld.Copy(tokenBalance)
	case TokenTypeMex:
		result.ProfitInMex.Copy(earnedInterestInTokens)
		result.TotalBalanceInMex.Copy(tokenBalance)
	default:
		return nil, fmt.Errorf("unknown token type '%d'", tokenType)
	}

	return result, nil
}
