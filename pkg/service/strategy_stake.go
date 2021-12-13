package service

import (
	"fmt"
	"math/big"

	"github.com/silviutroscot/istari-vision/pkg/log"
)

// StakeStrategy returns a StrategyResult representing the result of staking the token but not reinvesting the returns
func (s *Service) StakeStrategy(tokenType TokenType, input *StrategiesInput, tokenInitialPrice *big.Float) (*StrategyResult, error) {
	tokenBalance := input.EgldTokensInvested
	tokenAPR := input.EgldAPR
	targetPrice := input.EgldTargetPrice

	if tokenType == TokenTypeMex {
		tokenBalance = input.MexTokensInvested
		targetPrice = input.MexTargetPrice
		if input.RewardsInLockedMEX {
			tokenAPR = input.MexAPRLocked
		} else {
			tokenAPR = input.MexAPRUnlocked
		}
	}
	log.Info("tokenAPR is %s", tokenAPR.String())

	currentUSDValueFloat := big.Float{}
	currentUSDValueFloat.Mul(tokenInitialPrice, tokenBalance)

	// divide the number of days by the number of days of a year to get the percentage of the APR we get
	daysInYear := big.NewFloat(365.0)
	investmentDurationInDays := big.NewFloat(float64(input.InvestmentDurationInDays))
	percentageOfTheYearReceivingAPR := &big.Float{}
	percentageOfTheYearReceivingAPR.Quo(investmentDurationInDays, daysInYear)
	log.Info("investmentDurationInDays is %s", investmentDurationInDays.String())
	log.Info("percentageOfTheYearReceivingAPR is %s", percentageOfTheYearReceivingAPR.String())

	APRToBeReceived := &big.Float{}
	APRToBeReceived.Mul(tokenAPR, percentageOfTheYearReceivingAPR)
	APRToBeReceived.Quo(APRToBeReceived, big.NewFloat(100.0))
	log.Info("apr to be received is %s", APRToBeReceived)
	log.Info("tokens balance is %s", tokenBalance)

	tokensReceivedFromStaking := &big.Float{}
	tokensReceivedFromStaking.Mul(APRToBeReceived, tokenBalance)

	// calculate the ROI as the earned tokens / initial tokens balance
	roi := &big.Float{}
	roi.Quo(tokensReceivedFromStaking, tokenBalance)
	roi.Mul(roi, BigFloatOneHundred)
	log.Info("tokens received from staking are %s", tokensReceivedFromStaking)

	// USDValueOfEarnedTokens represents the USD value of the earned tokens using the current price of the token
	USDValueOfEarnedTokens := &big.Float{}
	USDValueOfEarnedTokens.Mul(tokensReceivedFromStaking, tokenInitialPrice)

	totalTokensBalance := &big.Float{}
	totalTokensBalance.Add(tokensReceivedFromStaking, tokenBalance)

	totalUSDValue := &big.Float{}
	totalUSDValue.Mul(totalTokensBalance, targetPrice)

	result := NewStrategyResult()
	result.TotalBalanceInUsd = totalUSDValue
	result.ProfitInUSD = USDValueOfEarnedTokens
	result.ROI = roi

	switch tokenType {
	case TokenTypeEgld:
		result.ProfitInEgld = tokensReceivedFromStaking
		result.TotalBalanceInEgld = totalTokensBalance
	case TokenTypeMex:
		result.ProfitInMex = tokensReceivedFromStaking
		result.TotalBalanceInMex = totalTokensBalance
	default:
		log.Error("unknown token type")
		return nil, fmt.Errorf("unknown token type '%d'", tokenType)
	}

	return result, nil
}
