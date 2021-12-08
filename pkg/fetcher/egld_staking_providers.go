package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/silviutroscot/istari-vision/pkg/log"
)

type EgldStakingProvidersElrond struct {
	ApiEndpoint string
}

func (sp *EgldStakingProvidersElrond) FetchStakingProviders() ([]EgldStakingProvider, error) {
	// make HTTP call to retrieve the staking providers
	res, err := httpClient.Get(sp.ApiEndpoint)
	if err != nil {
		log.Error("Error retrieving the EGLD staking providers from endpoint %s: %s", sp.ApiEndpoint, err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("error retrieving the EGLD staking providers from endpoint %s: Response status code %d",
			sp.ApiEndpoint, res.StatusCode)
		log.Error("%s", err)
		return nil, err
	}

	var providers []EgldStakingProvider

	body, _ := io.ReadAll(res.Body)
	if err = json.Unmarshal(body, &providers); err != nil {
		log.Error("Error unmarshalling the response from GET %s: %s", sp.ApiEndpoint, err)
		return nil, err
	}

	// give "Unknown<ID>" names to the staking providers that are unnamed so we can offer our users more staking options
	unknownCount := 0
	for idx := range providers {
		if providers[idx].Identity == "" {
			providers[idx].Identity = "unknown_" + strconv.Itoa(unknownCount)
			unknownCount++
		}
	}

	return providers, nil
}
