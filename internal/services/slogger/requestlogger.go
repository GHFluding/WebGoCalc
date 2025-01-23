package sl

import (
	"log/slog"
	"net/http"
	"test/internal/server/http/middleware"

	"github.com/gin-gonic/gin"
)

// LogRequestInfo logs the details of incoming requests or errors.
// use with all gin.HandlerFunc
func LogRequestInfo(log *slog.Logger, level string, c *gin.Context, message string, err error, extraFields map[string]interface{}) {
	// Default fields to log for all requests
	fields := []slog.Attr{
		slog.String("requestId", middleware.RequestIdFromContext(c)),
		slog.String("url", c.Request.URL.Path),
		slog.String("method", c.Request.Method),
	}

	// Add new fields from extraFields parameter
	for key, value := range extraFields {
		switch v := value.(type) {
		case string:
			fields = append(fields, slog.String(key, v))
		case int:
			fields = append(fields, slog.Int(key, v))
		case error:
			fields = append(fields, slog.String(key, error.Error(v))) //slog.Error(key, v) have problems
		}
	}

	//
	// Convert fields from []slog.Attr to []any (required for log.Error and log.Info)
	var fieldsAny []any
	for _, f := range fields {
		fieldsAny = append(fieldsAny, f)
	}

	// Log at the correct level
	if level == "error" && err != nil {
		fieldsAny = append(fieldsAny, slog.String("error", error.Error(err))) // Log the error field (slog.Error(key, v) have problems)
		log.Error(message, fieldsAny...)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
	} else if level == "info" {
		log.Info(message, fieldsAny...)
	}
}
