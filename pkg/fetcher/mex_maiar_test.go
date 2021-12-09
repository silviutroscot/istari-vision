package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/silviutroscot/istari-vision/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MexEconomicsFetcherMaiar_GetMexEconomics(t *testing.T) {
	t.Parallel()
	log.SetLevel(log.DebugLevel)

	t.Run("live_testnet", func(t *testing.T) {
		fetcher := MexEconomicsFetcherMaiar{
			ApiEndpoint: "https://testnet-exchange-graph.elrond.com/graphql",
		}

		economics, err := fetcher.FetchMexEconomics()
		require.NoError(t, err)

		// todo: add some asserts on price etc

		t.Logf("lockedRewardsAPR=%s unlockedRewardsAPR=%s price=%s",
			economics.LockedRewardsAPR.String(), economics.UnlockedRewardsAPR.String(), economics.Price.String())
	})

	t.Run("offline", func(t *testing.T) {
		mockResponse := `{
    "data": {
        "farms": [
            {
                "unlockedRewardsAPR": "89.718773317113517867",
				"lockedRewardsAPR": "389.541783567255098943",
                "farmingToken": {
                    "identifier": "EGLDMEX-1331c2",
                    "name": "EGLDMEXLP"
                },
                "farmToken": {
                    "name": "EGLDMEXLPStaked"
                },
                "farmedTokenPriceUSD": "0.00019940206605890213786",
                "farmedToken": {
                    "identifier": "MEX-45ebaa",
                    "name": "MEX"
                }
            },
            {
                "unlockedRewardsAPR": "153.471132951271740273",
                "lockedRewardsAPR": "358.060229713905366985",
                "farmingToken": {
                    "identifier": "WEGLDWUSDC-5b2b41",
                    "name": "WEGLDWUSDCLPToken"
                },
                "farmToken": {
                    "name": "EGLDUSDCLPStaked"
                },
                "farmedTokenPriceUSD": "0.00019940206605890213786",
                "farmedToken": {
                    "identifier": "MEX-45ebaa",
                    "name": "MEX"
                }
            },
            {
                "unlockedRewardsAPR": "8.577608884741437292",
				"lockedRewardsAPR": "102.931306616897247551",
                "farmingToken": {
                    "identifier": "MEX-45ebaa",
                    "name": "MEX"
                },
                "farmToken": {
                    "name": "MEXStaked"
                },
                "farmedTokenPriceUSD": "0.00019940206605890213786",
                "farmedToken": {
                    "identifier": "MEX-45ebaa",
                    "name": "MEX"
                }
            }
        ]
    }
}`
		mockHandlerLogic := &mockHandler{
			responseFunc: func(r *http.Request) ([]byte, int, error) {
				return []byte(mockResponse), http.StatusOK, nil
			},
		}
		handler := http.NewServeMux()
		handler.Handle("/graphql", mockHandlerLogic)

		t.Run("ok", func(t *testing.T) {
			server := httptest.NewUnstartedServer(handler)
			server.Start()
			defer server.Close()

			fetcher := MexEconomicsFetcherMaiar{
				ApiEndpoint: server.URL + "/graphql",
			}
			economics, err := fetcher.FetchMexEconomics()
			require.NoError(t, err)

			assert.Equal(t, "8.577608885", economics.UnlockedRewardsAPR.String())
			assert.Equal(t, "102.9313066", economics.LockedRewardsAPR.String())
			assert.Equal(t, "0.0001994020661", economics.Price.String())
		})

		t.Run("err_response", func(t *testing.T) {
			server := httptest.NewUnstartedServer(handler)
			server.Start()
			defer server.Close()

			fetcher := MexEconomicsFetcherMaiar{
				ApiEndpoint: server.URL + "/graphql?errorCode=400",
			}

			_, err := fetcher.FetchMexEconomics()
			require.Error(t, err)
			assert.Containsf(t, err.Error(), "response status code 400", "expected status code 400 error")
		})

		t.Run("err_conn", func(t *testing.T) {
			server := httptest.NewUnstartedServer(handler)
			server.Start()

			fetcher := MexEconomicsFetcherMaiar{
				ApiEndpoint: server.URL + "/graphql",
			}

			server.Close()
			_, err := fetcher.FetchMexEconomics()

			require.Error(t, err)
			assert.Containsf(t, err.Error(), "connect: connection refused", "expected connection error")
		})
	})
}
