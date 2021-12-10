package webservice

import (
	"errors"
	"net/http"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
)

// HandleGetPrices returns a JSON containing the live price for EGLD and MEX
func (api *API) HandleGetPrices(c *gin.Context) {
	economics, err := api.service.GetEconomics()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no prices found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prices": economics.Prices,
	})