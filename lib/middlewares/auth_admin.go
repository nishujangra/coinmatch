package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthenticateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Load .env
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, continuing...")
		}

		admin_token := os.Getenv("ADMIN_TOKEN")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication Header missing",
			})
			return
		}

		token_parts := strings.Split(authHeader, " ")

		if len(token_parts) != 2 || strings.ToLower(token_parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
			})
			return
		}

		token := token_parts[1]

		if token != admin_token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Auth Token",
			})

			return
		}

		c.Next()
	}
}
