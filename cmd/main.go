package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nishujangra/coinmatch/lib/config"
	"github.com/nishujangra/coinmatch/lib/routes"
)

func main() {
	dbConfig, err := config.BuildDataBaseConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Add api routes
	routes.APIRoutes(router)

	router.Run() // listen and serve on 0.0.0.0:8080
}
