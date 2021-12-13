package service

import (
	"math/big"

	"github.com/silviutroscot/istari-vision/pkg/log"
)



type StrategyResult struct {
	// ProfitInEgld the amount of EGLD earned using a strategy
	ProfitInEgld *big.Float
	// ProfitInMex the amount of MEX earned using a strategy
	ProfitInMex *big.Float
	// ProfitInUSD the
	ProfitInUSD *big.Float
	// TotalBalanceInEgld the amount of EGLD owned at the end of the monitored time interval
	TotalBalanceInEgld *big.Float
	// TotalBalanceInMex the amount of MEX owned at the end of the monitored time interval
	TotalBalanceInMex *big.Float
	// TotalBalanceInUsd the value in USD of EGLD + MEX, using their target value
	TotalBalanceInUsd *big.Float
	// ROI the percentage of profit we make in terms of EGLD; i.e. if at the beginning we invested 1 EGLD, and now we have 2 EGLD, ROI=100%
	ROI *big.Float
}

// Equals return true if the other StrategyResult equals the strategy
func (r *StrategyResult) Equals(other *StrategyResult) bool {
	log.Info("StrategyResult is %+v", r)
	log.Info("other is %+v", other)
	return BigFloatsAreEqual(*r.ProfitInEgld, *other.ProfitInEgld) &&
		BigFloatsAreEqual(*r.ProfitInMex, *other.ProfitInMex) &&
		BigFloatsAreEqual(*r.ProfitInUSD, *other.ProfitInUSD) &&
		BigFloatsAreEqual(*r.TotalBalanceInEgld, *other.TotalBalanceInEgld) &&
		BigFloatsAreEqual(*r.TotalBalanceInMex, *other.TotalBalanceInMex) &&
		BigFloatsAreEqual(*r.TotalBalanceInUsd, *other.TotalBalanceInUsd) &&
		BigFloatsAreEqual(*r.ROI, *other.ROI)
}

func (r *StrategyResult) MarshallToJSON() StrategyResultJSON {
	result := StrategyResultJSON{}
	result.ProfitInEgld = r.ProfitInEgld.Text('f', FloatingPointAccuracy)
	result.ProfitInMex = r.ProfitInMex.Text('f', FloatingPointAccuracy)
	result.ProfitInUSD = r.ProfitInUSD.Text('f', FloatingPointAccuracy)
	result.TotalBalanceInEgld = r.TotalBalanceInEgld.Text('f', FloatingPointAccuracy)
	result.TotalBalanceInMex = r.TotalBalanceInMex.Text('f', FloatingPointAccuracy)
	result.TotalBalanceInUsd = r.TotalBalanceInUsd.Text('f', FloatingPointAccuracy)
	// use a more aggressive truncation for ROI
	result.ROI = r.ROI.Text('f', 6)

	return result
}


