package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Recovered from panic", "error", err)
				c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
			}
		}()
		c.Next()
	}
}
