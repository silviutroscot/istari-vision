package webservice

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/silviutroscot/istari-vision/pkg/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.DebugMode)
}

type API struct {
	Address string

	engine  *gin.Engine
	service *service.Service
}

// NewAPI creates a new instance of a WebServer, which encapsulates the router and the dependencies of the WebService
func NewAPI(service *service.Service) *API {
	// create a new HTTP router engine
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/health.txt", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return &API{
		engine:  engine,
		service: service,
	}
}

// handleRateLimiting is a middleware that handles the rate limiting of the requests by only allowing 'limit'
// requests to happen, for the given 'duration' starting from the first request (identified by 'ipHeader' http header)
func handleRateLimiting(keyBase, ipHeader string, limit int64, duration time.Duration, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cc := context.WithTimeout(context.Background(), time.Second*5)
		defer cc()

		ip := c.GetHeader(ipHeader)

		result, err := cache.Incr(ctx, keyBase+":"+ip).Result()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if result == 1 {
			_ = cache.Expire(ctx, keyBase+":"+ip, duration).Err()
		} else if result > limit {
			_ = c.AbortWithError(http.StatusTooManyRequests, fmt.Errorf("too many requests"))
			return
		}
	}
}

func (api *API) Setup() error {
	// Enable CORS; CORS allows browser to accept and 'authorize' requests from the right (expected) 'site' and 'cross site' (as defined in RFCxxxx).
	// todo: update the variable `CORS_ORIGINS` in the .env file when it will be in production to use the right domains only
	corsOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
	if len(corsOrigins) == 0 || corsOrigins[0] == "" {
		corsOrigins = []string{"*"}
	}

	api.engine.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Minute, // todo: increase time to 12h; this represents for how long it will be cached in browser
	}))

	apiGroup := api.engine.Group("/api", handleRateLimiting("rate_limit", "X-Client-IP", 200, time.Minute, api.service.Cache))
	{
		apiGroup.GET("/egld_staking_providers", api.HandleGetEgldStakingProviders)
		apiGroup.GET("/prices", api.HandleGetPrices)
		apiGroup.POST("/calculate_profit", api.HandlePostCalculateProfit)
	}

	return nil
}

func (api *API) Run() error {
	return api.engine.Run(api.Address)
}
