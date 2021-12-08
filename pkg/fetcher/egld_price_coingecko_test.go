package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EgldPriceFetcherCoingecko_GetEgldPrice(t *testing.T) {
	t.Parallel()

	t.Run("offline", func(t *testing.T) {
		mockHandlerLogic := &mockHandler{
			responseFunc: func(r *http.Request) ([]byte, int, error) {
				mockResponse := `{"elrond-erd-2": {"usd": 31.41500}}`
				return []byte(mockResponse), http.StatusOK, nil
			},
		}
		handler := http.NewServeMux()
		handler.Handle("/api/v3/simple/price", mockHandlerLogic)

		t.Run("ok", func(t *testing.T) {
			server := httptest.NewUnstartedServer(handler)
			server.Start()
			defer server.Close()

			priceFetcher := EgldPriceFetcherCoingecko{
				ApiEndpoint: server.URL + "/api/v3/simple/price",
			}

			price, err := priceFetcher.FetchEgldPrice()
			require.NoError(t, err)

			priceFloat, _ := price.Float64()
			assert.Equal(t, 31.415, priceFloat)
		})

		t.Run("err_response", func(t *testing.T) {
			server := httptest.NewUnstartedServer(handler)
			server.Start()
			defer server.Close()

			priceFetcher := EgldPriceFetcherCoingecko{
				ApiEndpoint: server.URL + "/api/v3/simple/price?errorCode=400",
			}

			_, err := priceFetcher.FetchEgldPrice()
			require.Error(t, err)
			assert.Containsf(t, err.Error(), "Response status code 400", "expected status code 400 error")
		})

		t.Run("err_conn", func(t *testing.T) {
			server := httptest.NewUnstartedServer(handler)
			server.Start() // assigns random listing local address

			priceFetcher := EgldPriceFetcherCoingecko{
				ApiEndpoint: server.URL + "/api/v3/simple/price",
			}
			server.Close()
			_, err := priceFetcher.FetchEgldPrice()

			require.Error(t, err)
			assert.Containsf(t, err.Error(), "connect: connection refused", "expected connection error")
		})
	})

	t.Run("live", func(t *testing.T) {
		if testing.Short() {
			t.SkipNow()
		}

		priceFetcher := EgldPriceFetcherCoingecko{
			ApiEndpoint: EgldPriceFetcherCoingekoEndpoint,
		}

		price, err := priceFetcher.FetchEgldPrice()
		assert.NoError(t, err)

		priceFloat, _ := price.Float64()
		assert.GreaterOrEqual(t, priceFloat, 0.0)

		t.Logf("egld price usd: %f", priceFloat)
	})
}
