package fetcher

import "math/big"

// EgldPriceFetcher fetch the EGLD price in USD, source agnostic
// note: having it as an interface makes it much easier to write tests for it and makes it more future proof
type EgldPriceFetcher interface {
	FetchEgldPrice() (*big.Float, error)
}

// EgldStakingProvidersFetcher fetch the list of Egld staking providers;
// todo: verify which staking providers have enough capacity left for the amount of tokens to be invested
// todo: consider the fees of the staking providers  
// note: having it as an interface makes it much easier to write tests for it and makes it more future proof
type EgldStakingProvidersFetcher interface {
	FetchStakingProviders() ([]EgldStakingProvider, error)
}

// MexEconomicsFetcher retrieves the price for MEX and the APR for loecked and unlocked MEX as rewards from staking; 
type MexEconomicsFetcher interface {
	FetchMexEconomics() (economics MexEconomics, err error)
}
