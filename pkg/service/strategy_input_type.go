package service

import (
	"math/big"
)

// StrategiesInput represents a parsed and preprocessed request from an user to calculate their estimated gains
type StrategiesInput struct {
	EgldTokensInvested          *big.Float
	MexTokensInvested           *big.Float
	PercentageOfPortfolioInEgld *big.Float
	PercentageOfPortfolioInMex  *big.Float
	RewardsInLockedMEX          bool
	EgldTargetPrice             *big.Float
	MexTargetPrice              *big.Float
	EgldAPR                     *big.Float
	MexAPRLocked                *big.Float
	MexAPRUnlocked              *big.Float
	InvestmentDurationInDays    int
	RedelegationIntervalInDays  int
	StakingProvider             string
}
