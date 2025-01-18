package handler

import (
	"log/slog"
	"net/http"
	"test/internal/database/postgres"
	"test/internal/server/http/middleware"
	sl "test/internal/services/slogger"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler for creating a new student in the DB
func CreateStudentHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start the request timer
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)

		// Log the start of the request
		extraFields := map[string]interface{}{
			"requestId": requestID,
			"url":       c.Request.URL.Path,
			"method":    c.Request.Method,
		}
		sl.LogRequestInfo(log, "info", c, "Handling CreateStudent request", nil, extraFields)

		// Request for db
		s, student, err := postgres.CreateStudentData(db, c, log)

		// Handle errors
		if err != nil {
			extraFields["string"] = s
			extraFields["error"] = err.Error()
			sl.LogRequestInfo(log, "error", c, "Failed to create student", err, extraFields)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Log successful creation
		extraFields["studentId"] = student.ID
		extraFields["duration"] = time.Since(startTime).String()
		sl.LogRequestInfo(log, "info", c, "Successfully created student", nil, extraFields)

		// Return result to user
		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}
