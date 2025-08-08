package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nishujangra/coinmatch/lib/models"
)

func AddCurrency(c *gin.Context) {
	var currencyPair models.CurrencyPairRequest

	if err := c.BindJSON(&currencyPair); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON body",
		})

		return
	}

	// TODO : add to currency_pair table in DB

	c.JSON(http.StatusCreated, gin.H{
		"message": "added successfully",
	})
}
