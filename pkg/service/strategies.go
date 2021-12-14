package service

import (
	"fmt"
	"math/big"

	"github.com/silviutroscot/istari-vision/pkg/fetcher"
	"github.com/silviutroscot/istari-vision/pkg/log"
)

func (s *Service) CalculateStrategies(input *StrategiesInput, egldStakingProviders []fetcher.EgldStakingProvider, economics Economics) (map[string]StrategyResultJSON, error) {
	// wrapper over all the strategies results
	result := make(map[string]StrategyResultJSON)

	egldInitialPrice, _, err := big.ParseFloat(economics.Prices.EGLD, 10, 0, big.ToNearestEven)
	if err != nil {
		log.Error("error converting EGLD price string to float: %s", err)
		return result, err
	}

	mexInitialPrice := economics.mexEconomics.Price
	input.MexAPRLocked.Copy(economics.mexEconomics.LockedRewardsAPR)
	input.MexAPRUnlocked.Copy(economics.mexEconomics.UnlockedRewardsAPR)

	// if the portfolio percentage distribution is provided, simulate the swap to match the distribution
	// todo: if the user provides percentage of portofolio in egld do this, else don't match the distribution as it is not provided
	egldToBeInvested, mexToBeInvested :=
		s.SwapTokensToMatchDistribution(input.EgldTokensInvested, input.MexTokensInvested, input.PercentageOfPortfolioInEgld, egldInitialPrice, mexInitialPrice)
	input.EgldTokensInvested = egldToBeInvested
	input.MexTokensInvested = mexToBeInvested

	// find the staking provider in the list of staking providers and retrieve its APR
	egldStakingProviderWasFound := false
	for _, egldStakingProvider := range egldStakingProviders {
		if egldStakingProvider.Identity == input.StakingProvider {
			input.EgldAPR.SetFloat64(egldStakingProvider.APR)
			egldStakingProviderWasFound = true
			break
		}
	}

	if !egldStakingProviderWasFound {
		message := fmt.Sprintf("error finding the staking provider %s", input.StakingProvider)
		log.Error(message)
		return result, fmt.Errorf(message)
	}

	// if there is egld invested, compute the results for HOLD, STAKE and STAKE + REDELEGATE
	if egldToBeInvested.Cmp(EPSILON) == 1 {
		// HOLD
		egldHoldResult, err := s.HoldStrategy(TokenTypeEgld, input)
		if err != nil {
			log.Error("error calculating HOLD strategy for EGLD: %s", err)
			return result, err
		}
		log.Info("HOLD result token balance is %s", egldHoldResult.TotalBalanceInEgld.String())
		result["egld_hold"] = egldHoldResult.MarshallToJSON()

		// Stake
		egldStakeResult, err := s.StakeStrategy(TokenTypeEgld, input, egldInitialPrice)
		if err != nil {
			log.Error("error calculating STAKE strategy for EGLD: %s", err)
			return result, err
		}
		log.Info("egld stake result ROI is %v", egldStakeResult)
		result["egld_stake"] = egldStakeResult.MarshallToJSON()

		// Redelegate
		egldRedelegateResult, err := s.RedelegateStrategy(TokenTypeEgld, input, egldInitialPrice)
		if err != nil {
			log.Error("error calculating REDELEGATE strategy for EGLD: %s", err)
			return result, err
		}
		result["egld_redelegate"] = egldRedelegateResult.MarshallToJSON()
	}

	// if there is MEX invested, compute the results for STAKE and STAKE + REDELEGATE
	if mexToBeInvested.Cmp(EPSILON) == 1 {
		// Stake
		mexStakeResult, err := s.StakeStrategy(TokenTypeMex, input, mexInitialPrice)
		if err != nil {
			log.Error("error calculating STAKE strategy for MEX: %s", err)
			return result, err
		}
		result["mex_stake"] = mexStakeResult.MarshallToJSON()

		// Redelegate
		mexRedelegateResult, err := s.RedelegateStrategy(TokenTypeMex, input, mexInitialPrice)
		if err != nil {
			log.Error("error calculating REDELEGATE strategy for MEX: %s", err)
			return result, err
		}
		result["mex_redelegate"] = mexRedelegateResult.MarshallToJSON()
	}

	return result, nil
}
