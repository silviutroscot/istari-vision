package service

import (
"math/big"

"github.com/silviutroscot/istari-vision/pkg/log"
)

var (
	BigFloatZero       = big.NewFloat(0)
	BigFloatDaysInYear = big.NewFloat(365)
	BigFloatOneHundred = big.NewFloat(100)
)

// SwapTokensToMatchDistribution returns the amount of EGLD and MEX needed to match the required distribution, based on the USD value
func (s *Service) SwapTokensToMatchDistribution(EGLDTokensBalance, MEXTokensBalance, EGLDPercentage, EGLDPrice, MEXPrice *big.Float) (*big.Float, *big.Float) {
	// compute the total EGLD to know the amount of EGLD needed to match the distribution
	mexInlEgldBalance := SwapTokens(MEXTokensBalance, MEXPrice, EGLDPrice)

	totalEgldBalance := &big.Float{}
	totalEgldBalance.Add(mexInlEgldBalance, EGLDTokensBalance)

	// targetEGLDBalance represents the EGLD we should have to match the distribution requirement
	targetEGLDBalance := &big.Float{}
	targetEGLDBalance.Quo(totalEgldBalance, BigFloatOneHundred)
	targetEGLDBalance.Mul(targetEGLDBalance, EGLDPercentage)

	// the remaining amount of EGLD can be converted to MEX and that will represent that percentage of the wallet value that should be in MEX
	var targetMEXBalanceInEGLD big.Float
	targetMEXBalanceInEGLD.Sub(totalEgldBalance, targetEGLDBalance)
	log.Debug("totalEgldBalance is %v", &totalEgldBalance)
	log.Debug("targetMEXBalanceInEGLD is %v", &targetMEXBalanceInEGLD)

	mexTargetBalance := SwapTokens(&targetMEXBalanceInEGLD, EGLDPrice, MEXPrice)

	return targetEGLDBalance, mexTargetBalance
}

