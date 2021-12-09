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
