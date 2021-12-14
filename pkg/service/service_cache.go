package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/silviutroscot/istari-vision/pkg/log"
)

func (s *Service) CacheCron(ctx context.Context) {
	t := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-t.C:
			s.updateCache()
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}

func (s *Service) CacheWarmup() error {
	errs := s.updateCache()
	if len(errs) > 0 {
		return fmt.Errorf("errors during cache warmup: %v", errs)
	}
	return nil
}

// todo: re-think caching to include both "computational" caching and "display purpose" caching
func (s *Service) updateCache() []error {
	anyError := make([]error, 0, 3)

	cacheFuncs := []func() error{
		s.updateCacheStakingProviders,
		s.updateCacheMexEconomics,
		s.updateCacheEgldPrice,
	}

	var wg sync.WaitGroup
	wg.Add(len(cacheFuncs))

	for _, f := range cacheFuncs {
		go func(cacheFunc func() error) {
			defer wg.Done()
			err := cacheFunc()
			if err != nil {
				anyError = append(anyError, err)
			}
		}(f)
	}

	if len(anyError) > 0 {
		return anyError
	}

	return nil
}

func (s *Service) updateCacheStakingProviders() error {
	providers, err := s.EgldStakingProvidersFetcher.FetchStakingProviders()
	if err != nil {
		log.Error("error fetching the EGLD staking providers from Elrond API: %s", err)
		return err
	}

	ctx, cc := context.WithTimeout(context.Background(), time.Second*5)
	defer cc()

	data, err := json.Marshal(&providers)
	if err != nil {
		log.Error("error marshalling the EGLD staking providers structure to JSON: %s",
			err)
		return err
	}

	_, err = s.Cache.Set(ctx, "staking_providers_egld", data, 0).Result()
	if err != nil {
		log.Error("error storing the EGLD staking providers in cache: %s", err)
		return err
	}

	return nil
}

func (s *Service) updateCacheMexEconomics() error {
	mexEconomics, err := s.MexEconomicsFetcher.FetchMexEconomics()
	if err != nil {
		log.Error("error fetching the MEX economics from Maiar API: %s", err)
		return err
	}

	ctx, cc := context.WithTimeout(context.Background(), time.Second*5)
	defer cc()

	data, err := json.Marshal(&mexEconomics)
	if err != nil {
		log.Error("error marshalling the MEX economics structure to JSON: %s", err)
		return err
	}

	_, err = s.Cache.Set(ctx, "mex_economics", data, 0).Result()
	if err != nil {
		log.Error("error storing the MEX economics in cache: %s", err)
		return err
	}

	return nil
}

func (s *Service) updateCacheEgldPrice() error {
	egldPrice, err := s.EgldPriceFetcher.FetchEgldPrice()
	if err != nil {
		log.Error("error fetching the EGLD price in USD: %s", err)
		return err
	}

	ctx, cc := context.WithTimeout(context.Background(), time.Second*5)
	defer cc()

	data, err := json.Marshal(&egldPrice)
	if err != nil {
		log.Error("error marshalling the EGLD price in USD to JSON: %s", err)
		return err
	}

	_, err = s.Cache.Set(ctx, "egld_price", data, 0).Result()
	if err != nil {
		log.Error("error storing the EGLD price in USD in cache: %s", err)
		return err
	}

	return nil
}