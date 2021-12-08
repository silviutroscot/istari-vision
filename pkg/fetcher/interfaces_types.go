package fetcher

import "math/big"

type EgldStakingProvider struct {
	ServiceFee float64 `json:"serviceFee"`
	APR        float64 `json:"apr"`
	Identity   string  `json:"identity"`
}

type MexEconomics struct {
	LockedRewardsAPR   *big.Float
	UnlockedRewardsAPR *big.Float
	Price              *big.Float
}
