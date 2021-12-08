package fetcher

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EgldStakingProviders_FetchStakingProviders(t *testing.T) {
	// t.Parallel() allows Go to run tests in parallel
	t.Parallel()

	// mock a server that returns the same type of answer as the Elrond staking provider API;
	// we need this do avoid relying on the Elrond API for testing as any downtime they have may lead to us having flaky tests
	t.Run("offline", func(t *testing.T) {
		mockHandlerLogic := &mockHandler{
			responseFunc: func(r *http.Request) ([]byte, int, error) {
				query := r.URL.Query()

				providerCount, _ := strconv.Atoi(query.Get("providerCount"))
				providers := make([]EgldStakingProvider, 0, providerCount)

				for idx := 0; idx < providerCount; idx++ {
					providers = append(providers, EgldStakingProvider{
						ServiceFee: rand.Float64(),
						APR:        rand.Float64(),
						Identity:   "random_provider_" + strconv.Itoa(rand.Int()),
					})
				}

				payload, err := json.Marshal(providers)
				if err != nil {
					return nil, http.StatusTeapot, err
				}

				return payload, http.StatusOK, nil
			},
		}
		handler := http.NewServeMux()
		handler.Handle("/offline/staking_providers", mockHandlerLogic)

		t.Run("ok", func(t *testing.T) {
			server := httptest.NewUnstartedServer(handler)
			server.Start() // assigns random listing local address
			defer server.Close()

			// create a fetcher and fetch from the mock server
			egldFetcher := EgldStakingProvidersElrond{}
			egldFetcher.ApiEndpoint = server.URL + "/offline/staking_providers?providerCount=77"
			providers, err := egldFetcher.FetchStakingProviders()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(providers) != 77 {
				t.Fatalf("expected 77 providers, got %d", len(providers))
			}

			for _, provider := range providers {
				assert.True(t, strings.HasPrefix(provider.Identity, "random_provider_"), "missing prefix", provider)
				assert.GreaterOrEqual(t, provider.ServiceFee, 0.0, provider.Identity, provider)
				assert.GreaterOrEqual(t, len(provider.Identity), 1, provider.Identity, provider)
			}
		})

		// test if we got a status code different from HTTPSuccess our fetcher returns an error
		t.Run("err_response", func(t *testing.T) {
			server := httptest.NewUnstartedServer(handler)
			server.Start() // assigns random listing local address
			defer server.Close()

			egldFetcher := EgldStakingProvidersElrond{}
			egldFetcher.ApiEndpoint = server.URL + "/offline/staking_providers?providerCount=77&errorCode=400"
			providers, err := egldFetcher.FetchStakingProviders()

			require.Error(t, err)
			assert.Containsf(t, err.Error(), "Response status code 400", "expected status code 400 error")
			assert.Nil(t, providers)
		})

		// test if we were unable to connect to the server an error is returned
		t.Run("err_conn", func(t *testing.T) {
			server := httptest.NewUnstartedServer(handler)
			server.Start() // assigns random listing local address

			egldFetcher := EgldStakingProvidersElrond{}
			egldFetcher.ApiEndpoint = server.URL + "/offline/staking_providers?providerCount=77"
			server.Close()
			providers, err := egldFetcher.FetchStakingProviders()

			require.Error(t, err)
			assert.Containsf(t, err.Error(), "connect: connection refused", "expected connection error")
			assert.Nil(t, providers)
		})
	})

	// test fetching the staking providers using the live Elrond API
	t.Run("live", func(t *testing.T) {
		if testing.Short() {
			t.SkipNow()
		}

		egldFetcher := EgldStakingProvidersElrond{
			ApiEndpoint: EgldStakingProvidersEndpoint,
		}

		providers, err := egldFetcher.FetchStakingProviders()
		if err != nil {
			t.Fatalf("failed fetcher fetch staking providers: %v", err)
		}

		if len(providers) == 0 {
			t.Errorf("expected providers list to be not empty")
		}

		for idx, provider := range providers {
			assert.GreaterOrEqual(t, provider.APR, 0.0, provider.Identity, idx, provider)
			assert.GreaterOrEqual(t, provider.ServiceFee, 0.0, provider.Identity, idx, provider)
			assert.GreaterOrEqual(t, len(provider.Identity), 1, provider.Identity, idx, provider)
		}

		t.Logf("providers: %+v", providers)
	})
}
