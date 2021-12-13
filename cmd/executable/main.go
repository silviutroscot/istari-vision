package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"

	"github.com/silviutroscot/istari-vision/pkg/fetcher"
	"github.com/silviutroscot/istari-vision/pkg/log"
	"github.com/silviutroscot/istari-vision/pkg/service"
	"github.com/silviutroscot/istari-vision/pkg/webservice"
)

func main() {
	fmt.Print("Hello, world")
	if err := run(); err != nil {
		log.Error("runtime error: %s", err.Error())
	}
}

func getEnv(key, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return value
}

func run() error {
	cache := redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
		DB:   3, // todo: remove this
	})

	s := service.Service{
		Cache:                       cache,
		EgldPriceFetcher:            &fetcher.EgldPriceFetcherCoingecko{ApiEndpoint: getEnv("FETCHER_ENDPOINT_EGLD_PRICE_CG", fetcher.EgldPriceFetcherCoingekoEndpoint)},
		MexEconomicsFetcher:         &fetcher.MexEconomicsFetcherMaiar{ApiEndpoint: getEnv("FETCHER_ENDPOINT_MEXECO_MAIAR", fetcher.MexMaiarFetcherEndpoint)},
		EgldStakingProvidersFetcher: &fetcher.EgldStakingProvidersElrond{ApiEndpoint: getEnv("FETCHER_ENDPOINT_EGLD_STAKING", fetcher.EgldStakingProvidersEndpoint)},
	}

	if getEnv("CACHE_WARMUP", "") == "1" {
		if err := s.CacheWarmup(); err != nil {
			return err
		}
	}

	api := webservice.NewAPI(&s)
	api.Address = os.Getenv("API_ADDRESS")

	if err := api.Setup(); err != nil {
		return fmt.Errorf("failed api setup: %w", err)
	}

	go func() {
		s.CacheCron(context.Background())
		panic("cron stopped")
	}()

	return api.Run()
}