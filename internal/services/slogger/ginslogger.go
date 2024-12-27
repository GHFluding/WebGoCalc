package sl

import (
	"log/slog"
	"test/internal/server/http/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// logger for gin-gonic
func LoggingMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		log.With("component", "middleware/logger")
		log.Info("logger middleware enable")

		startTime := time.Now()

		c.Next()

		// log
		duration := time.Since(startTime)
		log.Info(
			"Request processed",
			"requestId", middleware.RequestIdFromContext(c),
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", duration.String(),
			"clientIP", c.ClientIP(),
		)
	}
}
