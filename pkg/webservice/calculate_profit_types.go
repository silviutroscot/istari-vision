package webservice

import (
	"fmt"
	"math/big"
	
	"github.com/silviutroscot/istari-vision/pkg/service"
)

// todo: move the comments for fields here
type CalculateStrategiesRequestPayload struct {
	// If MEXTokensInvested is provided, the PercentageOfPortfolioInEGLD and PercentageOfPortfolioInMEX is not used as we don't do any swaps
	// If MexTokensInvested is not provided and PercentageOfPortfolioInMEX > 0, convert that percentage of the EGLD invested to MEX as we get the request, using the live prices
	EGLDTokensInvested          string `json:"egld-tokens-invested"`
	MEXTokensInvested           string `json:"mex-tokens-invested"`
	PercentageOfPortfolioInEGLD string `json:"egld-pct"`
	PercentageOfPortfolioInMEX  string `json:"mex-pct"`

	// RewardsInLockedMEX is true if the user wants their MEX rewards to be in LockedMEX and false otherwise
	RewardsInLockedMEX bool   `json:"mex-rewards-locked"`
	EgldTargetPrice    string `json:"egld-price-target"`

	//MexTargetPrice can be 0 if the user wants to invest EGLD only
	MexTargetPrice           string `json:"mex-price-target"`
	InvestmentDurationInDays int    `json:"target-date-days"`
	RedelegationPeriodInDays int    `json:"redelegation-interval"`
	StakingProvider          string `json:"egld-staking-provider"`
}

// ToStrategiesInput returns an instance of service.StrategiesInput representing the parsed inputs and a list of
// errors for invalid fields
// todo: add unit tests
// todo: refactor use of repetitive parsing into generic function
func (payload *CalculateStrategiesRequestPayload) ToStrategiesInput() (*service.StrategiesInput, []error) {
	var err error
	var errs []error

	strategiesInput := &service.StrategiesInput{
		EgldTokensInvested:          &big.Float{},
		MexTokensInvested:           &big.Float{},
		PercentageOfPortfolioInEgld: &big.Float{},
		PercentageOfPortfolioInMex:  &big.Float{},
		RewardsInLockedMEX:          payload.RewardsInLockedMEX,
		EgldTargetPrice:             &big.Float{},
		MexTargetPrice:              &big.Float{},
		EgldAPR:                     &big.Float{},
		MexAPRLocked:                &big.Float{},
		MexAPRUnlocked:              &big.Float{},
		InvestmentDurationInDays:    0,
		RedelegationIntervalInDays:  0,
		StakingProvider:             payload.StakingProvider,
	}

	if payload.EGLDTokensInvested != "" {
		strategiesInput.EgldTokensInvested, err = parseBigFloat(payload.EGLDTokensInvested)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed parsing field 'EGLDTokensInvested': %w", err))
		} else if strategiesInput.EgldTokensInvested.Cmp(service.BigFloatZero) == -1 {
			errs = append(errs, fmt.Errorf("failed validating field 'EGLDTokensInvested' value '%s': %w", payload.EGLDTokensInvested, err))
		}
	}

	if payload.MEXTokensInvested != "" {
		strategiesInput.MexTokensInvested, err = parseBigFloat(payload.MEXTokensInvested)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed parsing field 'MexTokensInvested': %w", err))
		} else if strategiesInput.MexTokensInvested.Cmp(service.BigFloatZero) == -1 {
			errs = append(errs, fmt.Errorf("failed validating field 'MexTokensInvested' value '%s': %w", payload.MEXTokensInvested, err))
		}
	}

	if payload.PercentageOfPortfolioInEGLD != "" {
		strategiesInput.PercentageOfPortfolioInEgld, err = parseBigFloat(payload.PercentageOfPortfolioInEGLD)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed parsing field 'PercentageOfPortfolioInEgld': %w", err))
		} else if strategiesInput.PercentageOfPortfolioInEgld.Cmp(service.BigFloatZero) == -1 || strategiesInput.PercentageOfPortfolioInEgld.Cmp(service.BigFloatOneHundred) == 1 {
			errs = append(errs, fmt.Errorf("failed validating field 'PercentageOfPortfolioInEgld' value '%s': %w", payload.PercentageOfPortfolioInEGLD, err))
		}
	}

	if payload.PercentageOfPortfolioInMEX != "" {
		strategiesInput.PercentageOfPortfolioInMex, err = parseBigFloat(payload.PercentageOfPortfolioInMEX)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed parsing field 'PercentageOfPortfolioInMex': %w", err))
		} else if strategiesInput.PercentageOfPortfolioInMex.Cmp(service.BigFloatZero) == -1 || strategiesInput.PercentageOfPortfolioInMex.Cmp(service.BigFloatOneHundred) == 1 {
			errs = append(errs, fmt.Errorf("failed validating field 'PercentageOfPortfolioInMex' value '%s': %w", payload.PercentageOfPortfolioInMEX, err))
		}
	}

	strategiesInput.MexTargetPrice, err = parseBigFloat(payload.MexTargetPrice)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed parsing field 'MexTargetPrice': %w", err))
	} else if strategiesInput.MexTargetPrice.Cmp(service.BigFloatZero) == -1 {
		errs = append(errs, fmt.Errorf("failed validating field 'MexTargetPrice' value '%s': %w", payload.MexTargetPrice, err))
	}

	strategiesInput.EgldTargetPrice, err = parseBigFloat(payload.EgldTargetPrice)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed parsing field 'EgldTargetPrice': %w", err))
	} else if strategiesInput.EgldTargetPrice.Cmp(service.BigFloatZero) == -1 {
		errs = append(errs, fmt.Errorf("failed validating field 'EgldTargetPrice' value '%s': %w", payload.EgldTargetPrice, err))
	}

	strategiesInput.InvestmentDurationInDays = payload.InvestmentDurationInDays
	strategiesInput.RedelegationIntervalInDays = payload.RedelegationPeriodInDays

	// verify that the sum of percentages in MEX and EGLD is 100
	percentageSum := &big.Float{}
	percentageSum.Add(strategiesInput.PercentageOfPortfolioInEgld, strategiesInput.PercentageOfPortfolioInMex)
	if !service.BigFloatsAreEqual(*percentageSum, *service.BigFloatOneHundred) {
		errs = append(errs, fmt.Errorf("the sum of the pecentages is %s, not 100%", percentageSum.String()))
	}

	if len(errs) != 0 {
		return nil, errs
	}

	return strategiesInput, nil
}

func parseBigFloat(input string) (*big.Float, error) {
	f, _, err := big.ParseFloat(input, 10, 0, big.ToNearestEven)
	return f, err
}
