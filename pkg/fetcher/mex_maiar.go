package fetcher

import (
"encoding/json"
"fmt"
"math/big"
"net/http"
"strings"

"github.com/silviutroscot/istari-vision/pkg/log"
)

// get mex price  --------> service ---[if expired]-->  fetcher (Interface)
// get mex apr    __/  |             \___[else]----> cache
// get mex economics __/

type MexEconomicsFetcherMaiar struct {
	ApiEndpoint string
	TokenName   string
}

const mexEconomicsMaiarQuery = `{
  "query": "query {farms {lockedRewardsAPR unlockedRewardsAPR farmingToken{identifier name} farmToken{name} farmedTokenPriceUSD farmedToken {identifier name}}}",
  "variables": {}
}`

func (mf *MexEconomicsFetcherMaiar) FetchMexEconomics() (MexEconomics, error) {
	economics := MexEconomics{
		UnlockedRewardsAPR: new(big.Float),
		LockedRewardsAPR:   new(big.Float),
		Price:              new(big.Float),
	}

	tokenName := mf.TokenName
	if tokenName == "" {
		tokenName = "MEXStaked"
	}

	body := strings.NewReader(mexEconomicsMaiarQuery)

	req, err := http.NewRequest(http.MethodPost, mf.ApiEndpoint, body)
	if err != nil {
		log.Error("error creating the request for MEX economics to endpoint %s: %s", mf.ApiEndpoint, err)
		return economics, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Error("error sending the POST request for MEX economics to endpoint %s: %s", mf.ApiEndpoint, err)
		return economics, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		message := fmt.Sprintf("error in the MEX economics response from endpoint %s: response status code %d",
			mf.ApiEndpoint, res.StatusCode)
		log.Error(message)
		return economics, fmt.Errorf(message)
	}

	var response struct {
		Data struct {
			Farms []struct {
				LockedRewardsAPR   string `json:"lockedRewardsAPR"`
				UnlockedRewardsAPR string `json:"unlockedRewardsAPR"`
				FarmingToken       struct {
					Identifier string `json:"identifier"`
					Name       string `json:"name"`
				} `json:"farmingToken"`
				FarmToken struct {
					Name string `json:"name"`
				} `json:"farmToken"`
				FarmedTokenPriceUSD string `json:"farmedTokenPriceUSD"`
				FarmedToken         struct {
					Identifier string `json:"identifier"`
					Name       string `json:"name"`
				} `json:"farmedToken"`
			} `json:"farms"`
		} `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Error("error decoding the MEX economics JSON response: %s", err)
		return economics, err
	}

	if len(response.Data.Farms) == 0 {
		err = fmt.Errorf("no MEX farms available, unable to retrieve MEX economics")
		log.Info(err.Error())
		return economics, err
	}

	for _, farm := range response.Data.Farms {
		if farm.FarmToken.Name == tokenName {
			_, _, err = economics.Price.Parse(farm.FarmedTokenPriceUSD, 10)
			if err != nil {
				log.Error("error parsing the MEX price in USD from the Maiar API response: %s", err)
				return economics, err
			}

			_, _, err = economics.LockedRewardsAPR.Parse(farm.LockedRewardsAPR, 10)
			if err != nil {
				log.Error("error parsing the MEX LockedRewardsAPR from the Maiar API response: %s", err)
				return economics, err
			}

			_, _, err = economics.UnlockedRewardsAPR.Parse(farm.UnlockedRewardsAPR, 10)
			if err != nil {
				log.Error("error parsing the MEX UnlockedRewardsAPR from the Maiar API response: %s", err)
				return economics, err
			}
		}
	}

	return economics, nil
}
