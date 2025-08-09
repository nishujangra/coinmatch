package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

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

	// add to currency_pair table in DB
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

	// save to order table in DB
	_, err := apiController.DB.Exec(
		"INSERT INTO orders (pair, side, price, quantity, user_id, status) INTO VALUES ($1, $2, $3, $4, $5, 'open')",
		order.Pair, order.Side, order.Price, order.Quantity, order.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to insert orders into database",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully added order to the table",
	})
}

func (apiController *APIController) ViewOrderbook(c *gin.Context) {
	pair := c.Query("pair")
	depthStr := c.DefaultQuery("depth", "10")

	if pair == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "pair parameter is required",
		})
		return
	}

	depth, err := strconv.Atoi(depthStr)
	if err != nil || depth < 0 {
		depth = 0
	}

	// Get order book snapshot from matching engine
	// Return top 10 BUY and SELL orders (sorted)
	// snapshot, err := apiController.MatchingEngine.GetOrderBookSnapshot(pair, depth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get order book snapshot: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pair":  pair,
		"depth": depth,
	})
}

func (apiController *APIController) GetUserOrder(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id parameter is required",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id parameter must be valid integer",
		})
		return
	}

	// Query user orders from db
	rows, err := apiController.DB.Query(
		"SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query user orders: " + err.Error(),
		})
		return
	}

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan order data: " + err.Error(),
			})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"orders":  orders,
		"count":   len(orders),
	})
}

func (apiController *APIController) DeleteOrder(c *gin.Context) {
	order_id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"order_id": order_id,
	})
}
