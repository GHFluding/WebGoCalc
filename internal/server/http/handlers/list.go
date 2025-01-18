package handler

import (
	"log/slog"
	"net/http"
	"test/internal/database/postgres"
	"test/internal/server/http/middleware"
	"time"

	"github.com/gin-gonic/gin"

	sl "test/internal/services/slogger"
)

func ListStudentsHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start the request timer
		startTime := time.Now()

		// Retrieve request ID for correlation
		requestID := middleware.RequestIdFromContext(c)

		// Prepare fields for logging
		extraFields := map[string]interface{}{
			"requestId": requestID,
			"url":       c.Request.URL.Path,
			"method":    c.Request.Method,
		}

		// Log the incoming request
		sl.LogRequestInfo(log, "info", c, "Handling ListStudents request", nil, extraFields)

		// Fetch the list of students from the database
		students, err := db.ListStudents(c.Request.Context())
		if err != nil {
			// Log the error if fetching students fails
			extraFields["error"] = err.Error()
			sl.LogRequestInfo(log, "error", c, "Failed to retrieve students", err, extraFields)

			// Return error response to the client
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve students",
			})
			return
		}

		// Log success after retrieving students
		extraFields["count"] = len(students)
		extraFields["duration"] = time.Since(startTime).String()
		sl.LogRequestInfo(log, "info", c, "Successfully retrieved students", nil, extraFields)

		// Return the list of students as a response
		c.JSON(http.StatusOK, gin.H{
			"students": students,
		})
	}
}
