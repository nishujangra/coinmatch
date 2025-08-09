package controllers

import (
	"container/heap"
	"database/sql"
	"net/http"
	"sort"
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

	book, ok := engine.Books[pair]
	if !ok {
		// Return empty order book if pair doesn't exist
		c.JSON(http.StatusOK, gin.H{
			"buy":  []models.Order{},
			"sell": []models.Order{},
		})
		return
	}

	var buys, sells []models.Order

	// Convert BuyPQ to slice and sort by price descending
	buys = make([]models.Order, 0, len(book.BuyPQ))
	for _, order := range book.BuyPQ {
		if order.Quantity > 0 {
			buys = append(buys, *order)
		}
	}

	// sort by price descending
	sort.Slice(buys, func(i, j int) bool {
		return buys[i].Price > buys[j].Price
	})

	if depth > 0 && depth < len(buys) {
		buys = buys[:depth]
	}

	// Convert SellPQ to slice and sort by price ascending
	sells = make([]models.Order, 0, len(book.SellPQ))
	for _, order := range book.SellPQ {
		if order.Quantity > 0 {
			sells = append(sells, *order)
		}
	}

	// sort by price ascending
	sort.Slice(sells, func(i, j int) bool {
		return sells[i].Price < sells[j].Price
	})

	if depth > 0 && depth < len(sells) {
		sells = sells[:depth]
	}

	type ViewOrderResponse struct {
		Price    float64
		Quantity float64
	}

	var buyList []ViewOrderResponse

	for _, buy := range buys {
		buyList = append(buyList, ViewOrderResponse{
			Price:    buy.Price,
			Quantity: buy.Quantity,
		})
	}

	var sellList []ViewOrderResponse

	for _, sell := range sells {
		sellList = append(sellList, ViewOrderResponse{
			Price:    sell.Price,
			Quantity: sell.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"buy":  buyList,
		"sell": sellList,
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
		"SELECT id, user_id, pair, side, price, quantity, filled_quantity, status, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query user orders: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.Pair, &order.Side,
			&order.Price, &order.Quantity, &order.FilledQuantity,
			&order.Status, &order.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan order data: " + err.Error(),
			})
			return
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error iterating over orders: " + err.Error(),
		})
		return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be cancelled"})
		return
	}

	// Update status to canceled
	_, err = apiController.DB.Exec("UPDATE orders SET status = 'cancelled' WHERE id = $1", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order canceled successfully",
		"id":      orderID,
	})

}
