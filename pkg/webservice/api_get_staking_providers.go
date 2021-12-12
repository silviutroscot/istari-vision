package webservice

import (
	"errors"
	"net/http"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
)

func (api *API) HandleGetEgldStakingProviders(c *gin.Context) {
	stakingProviders, err := api.service.GetStakingProviders()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "no staking providers found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"staking_providers": stakingProviders,
	})
}