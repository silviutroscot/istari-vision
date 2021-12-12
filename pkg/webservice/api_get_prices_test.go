package webservice

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-redis/redis/v8"
	"github.com/silviutroscot/istari-vision/pkg/service"
	"github.com/stretchr/testify/require"
)

func TestAPI_HandleGetPrices(t *testing.T) {
	cache := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   4,
	})

	cache.Del(context.Background(), "prices")
	cache.Del(context.Background(), "mex_economics")

	t.Cleanup(func() {
		cache.Del(context.Background(), "prices")
		cache.Del(context.Background(), "mex_economics")
		_ = cache.Close()
	})

	s := &service.Service{
		Cache:                       cache,
		EgldPriceFetcher:            nil,
		MexEconomicsFetcher:         nil,
		EgldStakingProvidersFetcher: nil,
	}
	api := NewAPI(s)
	require.NoError(t, api.Setup())

	t.Run("empty cache", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/prices", nil)
		api.engine.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, `{"error":"no prices found"}`, w.Body.String())
	})

	cache.Set(context.Background(), "egld_price", `"3.1415"`, 0)
	cache.Set(context.Background(), "mex_economics", `{"Price":"0.0002184434508"}`, 0)

	t.Run("with cache", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/prices", nil)
		api.engine.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"prices":{"egld":"3.1415","mex":"0.0002184434508"}}`, w.Body.String())
	})
}