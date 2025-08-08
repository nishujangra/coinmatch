package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nishujangra/coinmatch/lib/controllers"
	"github.com/nishujangra/coinmatch/lib/middlewares"
)

func APIRoutes(r *gin.Engine) {
	r.POST("/api/pairs", middlewares.AuthenticateAdmin(), controllers.AddCurrency)
}
