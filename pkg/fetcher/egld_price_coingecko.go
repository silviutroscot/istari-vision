package fetcher

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/silviutroscot/istari-vision/pkg/log"
)

// EgldPriceFetcherCoingecko contains the URL from where we fetch the EGLD price and a function to get the price
type EgldPriceFetcherCoingecko struct {
	ApiEndpoint string
	ElrondId    string
	Currency    string
}

func (e *EgldPriceFetcherCoingecko) FetchEgldPrice() (*big.Float, error) {
	elrondId, currency := e.ElrondId, e.Currency
	if elrondId == "" {
		elrondId = "elrond-erd-2"
	}
	if currency == "" {
		currency = "USD"
	}

	req, err := http.NewRequest(http.MethodGet, e.ApiEndpoint, nil)
	if err != nil {
		log.Error("error creating the EGLD price request for endpoint %s: %s", e.ApiEndpoint, err)
		return nil, err
	}

	// set query parameters for the request; this makes a copy of the URL and we set it to the request after setting all the parameters
	query := req.URL.Query()
	query.Set("ids", elrondId)
	query.Set("vs_currencies", currency)
	req.URL.RawQuery = query.Encode()

	res, err := httpClient.Do(req)
	if err != nil {
		log.Error("error retrieving the EGLD price from endpoint %s: %s", e.ApiEndpoint, err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("error retrieving the EGLD price from endpoint %s: Response status code %d",
			e.ApiEndpoint, res.StatusCode)
		log.Error("%s", err)
		return nil, err
	}

	var response struct {
		Token map[string]float64 `json:"elrond-erd-2"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Error("Error unmarshalling the response from GET %s: %s", e.ApiEndpoint, err)
		return nil, err
	}

	return big.NewFloat(response.Token[strings.ToLower(currency)]), nil
}
