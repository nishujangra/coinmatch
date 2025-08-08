package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nishujangra/coinmatch/lib/models"
)

type APIController struct {
	DB *sql.DB
}

func NewAPIController(db *sql.DB) APIController {
	return APIController{
		DB: db,
	}
}

func (apiController *APIController) AddCurrency(c *gin.Context) {
	var currencyPair models.CurrencyPairRequest

	if err := c.BindJSON(&currencyPair); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON body",
		})

		return
	}

	// TODO : add to currency_pair table in DB
	_, err := apiController.DB.Exec(
		"INSERT INTO currency_pairs (base_currency, quote_currency) VALUES ($1, $2)",
		currencyPair.Base, currencyPair.Quote,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to insert currency pair into database",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "added successfully",
	})
}

func (apiController *APIController) AddOrder(c *gin.Context) {
	var order models.OrderRequest

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON body",
		})

		return
	}

	// TODO: save to order table in DB

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully added order to the table",
	})
}
