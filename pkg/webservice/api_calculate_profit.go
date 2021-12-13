package webservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silviutroscot/istari-vision/pkg/log"
)

func (api *API) HandlePostCalculateProfit(c *gin.Context) {
	var requestPayload CalculateStrategiesRequestPayload

	err := c.BindJSON(&requestPayload)
	if err != nil {
		log.Error("error binding the request payload: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	strategiesInput, errs := requestPayload.ToStrategiesInput()
	if errs != nil {
		errsStrings := make([]string, len(errs))
		for i, err := range errs {
			errsStrings[i] = err.Error()
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errsStrings,
		})
		return
	}

	egldStakingProviders, err := api.service.GetStakingProviders()
	if err != nil {
		log.Error("error retrieving the EGLD staking providers: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	economics, err := api.service.GetEconomics()
	if err != nil {
		log.Error("error retrieving the economics: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	results, err := api.service.CalculateStrategies(strategiesInput, egldStakingProviders, economics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"prices":  economics.Prices,
	})
}