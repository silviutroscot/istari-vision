package fetcher

import (
	"net/http"
	"time"
)

const (
	// EgldStakingProvidersEndpoint endpoint to fetch the list of staking providers and their details
	EgldStakingProvidersEndpoint = "https://api.elrond.com/providers"

	// EgldPriceFetcherCoingekoEndpoint endpoint to fetch the live price of EGLD
	// todo: provide more endpoints so our system won't have a single point of failure if this provider is down
	EgldPriceFetcherCoingekoEndpoint = "https://api.coingecko.com/api/v3/simple/price"

	// MexMaiarFetcherEndpoint endpoint to fetch the MEX price and the APR for locked and unlocked staking
	MexMaiarFetcherEndpoint = "https://testnet-exchange-graph.elrond.com/graphql"
)

// httpClient will be used as a singleton and can be reused by any request
// note: httpClient acts as a controller for http requests and their respective tcp connections (keeping 'keep-alive' tcp connections pooled for future use)
// note: important to read the response.Body (or io.Discard it) before closing the body, so that httpClient can re-use the TCP connection
var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:    50,
		IdleConnTimeout: 1 * time.Minute,
	},
	Timeout: time.Second * 15,
}
