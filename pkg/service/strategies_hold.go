package service

import (
	"fmt"
	"math/big"

	"github.com/silviutroscot/istari-vision/pkg/log"
)

type TokenType uint8

const (
	TokenTypeUndefined TokenType = iota
	TokenTypeEgld
	TokenTypeMex

	// FloatingPointAccuracy sets the accuracy when converting to string
	FloatingPointAccuracy = 10
)

// StrategyResultJSON represents a StrategyResult but with all fields formatted to have 10^-10 accuracy
type StrategyResultJSON struct {
	ProfitInEgld       string
	ProfitInMex        string
	ProfitInUSD        string
	TotalBalanceInEgld string
	TotalBalanceInMex  string
	TotalBalanceInUsd  string
	ROI                string
}





// NewStrategyResult returns a 0 value StrategyResult
func NewStrategyResult() *StrategyResult {
	return &StrategyResult{
		ProfitInEgld:       &big.Float{},
		ProfitInMex:        &big.Float{},
		ProfitInUSD:        &big.Float{},
		TotalBalanceInEgld: &big.Float{},
		TotalBalanceInMex:  &big.Float{},
		TotalBalanceInUsd:  &big.Float{},
		ROI:                &big.Float{},
	}
}

// HoldStrategy returns a StrategyResult which represent what happens if we only hold the token (either MEX or EGLD)
func (s *Service) HoldStrategy(tokenType TokenType, input *StrategiesInput) (*StrategyResult, error) {
	targetPrice := input.EgldTargetPrice
	tokenBalance := input.EgldTokensInvested
	if tokenType == TokenTypeMex {
		targetPrice = input.MexTargetPrice
		tokenBalance = input.MexTokensInvested
	}

	result := NewStrategyResult()
	result.TotalBalanceInUsd.Mul(targetPrice, tokenBalance)

	log.Info("token balance is %s", tokenBalance)
	switch tokenType {
	case TokenTypeEgld:
		result.TotalBalanceInEgld.Copy(tokenBalance)
	case TokenTypeMex:
		result.TotalBalanceInMex.Copy(tokenBalance)
	default:
		log.Error("unknown token type")
		return nil, fmt.Errorf("unknown token type '%d'", tokenType)
	}

	return result, nil
}