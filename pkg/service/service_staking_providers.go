package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/silviutroscot/istari-vision/pkg/fetcher"
)

func (s *Service) GetStakingProviders() ([]fetcher.EgldStakingProvider, error) {
	ctx, cc := context.WithTimeout(context.Background(), time.Second*5)
	defer cc()

	result, err := s.Cache.Get(ctx, "staking_providers_egld").Result()
	if err != nil {
		return nil, err
	}

	var providers []fetcher.EgldStakingProvider

	if err := json.Unmarshal([]byte(result), &providers); err != nil {
		return nil, err
	}

	return providers, nil
}