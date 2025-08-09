package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nishujangra/coinmatch/lib/controllers"
	"github.com/nishujangra/coinmatch/lib/middlewares"
)

func APIRoutes(r *gin.Engine, db *sql.DB) {
	apiController := controllers.NewAPIController(db)

	// Only for ADMIN Route
	r.POST("/api/pairs", middlewares.AuthenticateAdmin(), apiController.AddCurrency)

	// placing limiting orders
	r.POST("/api/orders", apiController.AddOrder)

	// View order book
	r.GET("/api/orderbook", apiController.ViewOrderbook)

	// Get User Orders
	r.GET("/api/orders", apiController.GetUserOrder)

	// (Optional) DELETE order
	r.DELETE("/api/orders/:id", apiController.CancelOrder)
}
