
package service

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/silviutroscot/istari-vision/pkg/fetcher"
)

type Prices struct {
	EGLD string `json:"egld"`
	MEX  string `json:"mex"`
}

type EgldPrice struct {
	Price *big.Float
}

// Economics encapsulates the USD prices for MEX and EGLD, alongside the APR and APRMultiplier for MEX farm
type Economics struct {
	Prices       Prices
	mexEconomics fetcher.MexEconomics
}

func (s *Service) GetEconomics() (Economics, error) {
	var economics Economics

	ctx, cc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cc()

	// mex parsing
	{
		result, err := s.Cache.Get(ctx, "mex_economics").Result()
		if err != nil {
			return economics, err
		}

		if err := json.Unmarshal([]byte(result), &economics.mexEconomics); err != nil {
			return economics, err
		}

		economics.Prices.MEX = economics.mexEconomics.Price.String()
	}

	// egld parsing
	{
		result, err := s.Cache.Get(ctx, "egld_price").Result()
		if err != nil {
			return economics, err
		}

		if err := json.Unmarshal([]byte(result), &economics.Prices.EGLD); err != nil {
			return economics, err
		}
	}

	return economics, nil
}