package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ValidateAuthorizationHeader verifies if the custom authorization header is present on request
func ValidateAuthorizationHeader(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	envAuthorization := os.Getenv("REQ_VALIDATOR")

	if authorizationHeader != envAuthorization {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Request was not possible due to missing or invalid header",
		})
		return
	}
	log.Println("Request authorized")
	c.Next()
}
