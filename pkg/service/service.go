package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/silviutroscot/istari-vision/pkg/fetcher"
)

// Service encapsulates the business logic and computation
type Service struct {
	Cache *redis.Client

	// note: can be switched to list of fetchers, or map of fetchers
	EgldPriceFetcher            fetcher.EgldPriceFetcher
	MexEconomicsFetcher         fetcher.MexEconomicsFetcher
	EgldStakingProvidersFetcher fetcher.EgldStakingProvidersFetcher
}
