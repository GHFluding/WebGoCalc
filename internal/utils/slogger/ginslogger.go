package sl

import (
	"io"
	"log/slog"
	"test/internal/transport/rest/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// logger for gin-gonic
func LoggingMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Add "component" information to all logs in this middleware
		logger := log.With("component", "middleware/logger")
		logger.Info("Logger middleware enabled")

		// Capture the start time of the request
		startTime := time.Now()

		// Continue processing the request
		c.Next()

		// Calculate the duration it took to process the request
		duration := time.Since(startTime)

		// Log the request details
		logger.Info("Request processed",
			"requestId", middleware.RequestIdFromContext(c), // Assuming you have a middleware to extract requestId
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", duration.String(),
			"clientIP", c.ClientIP(),
		)

		// Optionally, you can log the request body or headers if needed
		// For example, if you want to log the request body:
		body, err := io.ReadAll(c.Request.Body)
		if err == nil {
			logger.Info("Request body", "body", string(body))
		}
	}
}
