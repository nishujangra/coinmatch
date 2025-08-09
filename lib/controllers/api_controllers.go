package controllers

import (
	"container/heap"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nishujangra/coinmatch/lib/engine"
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
	var orderReq models.OrderRequest

	if err := c.BindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON body",
		})

		return
	}

	order := orderReq.ToOrder()

	// save to order table in DB
	_, err := apiController.DB.Exec(
		"INSERT INTO orders (pair, side, price, quantity, user_id, status) VALUES ($1, $2, $3, $4, $5, 'open')",
		order.Pair, order.Side, order.Price, order.Quantity, order.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to insert orders into database",
			"details": err.Error(),
		})
		return
	}

	book, exists := engine.Books[order.Pair]
	if !exists {
		book = &engine.OrderBook{}
		heap.Init(&book.BuyPQ)
		heap.Init(&book.SellPQ)
		engine.Books[order.Pair] = book
	}

	go engine.MatchOrder(order, book)

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

func (apiController *APIController) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")

	var status string
	err := apiController.DB.QueryRow("SELECT status FROM orders WHERE id = $1", orderID).Scan(&status)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	// Only cancel if order is open or partial
	if status != "open" && status != "partial" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be canceled"})
		return
	}

	// Update status to canceled
	_, err = apiController.DB.Exec("UPDATE orders SET status = 'canceled' WHERE id = $1", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order canceled successfully",
		"id":      orderID,
	})

}
