package api_gateway

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// ErrorHandler обрабатывает ошибки, произошедшие в процессе работы.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			log.Println("Error occurred:", err.Error())

			errorResponse := gin.H{
				"error":   "Internal Server Error",
				"details": err.Error(),
			}

			statusCode := http.StatusInternalServerError

			if strings.Contains(err.Error(), "invalid") {
				statusCode = http.StatusBadRequest
				errorResponse["error"] = "Bad Request"
			} else if strings.Contains(err.Error(), "not found") {
				statusCode = http.StatusNotFound
				errorResponse["error"] = "Not Found"
			} else if strings.Contains(err.Error(), "unauthorized") {
				statusCode = http.StatusUnauthorized
				errorResponse["error"] = "Unauthorized"
			} else if strings.Contains(err.Error(), "connection refused") {
				statusCode = http.StatusServiceUnavailable
				errorResponse["error"] = "Service Unavailable"
			}

			c.JSON(statusCode, errorResponse)
		}
	}
}
