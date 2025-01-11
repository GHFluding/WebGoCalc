package handler

import (
	"context"
	"log/slog"
	"net/http"
	"test/internal/database/postgres"
	"test/internal/server/http/middleware"
	"time"

	"github.com/gin-gonic/gin"
)


func GetStudentByNameHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get name from request
		studentName := c.Param("name")
		if studentName == "" {
			log.Error("Missing student name in request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Student name is required",
			})
			return
		}

		// Create context and logs
		ctx := context.Background()
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		log.Info("Handling GetStudentByName request",
			"requestId", requestID,
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
			"studentName", studentName,
		)

		// Get data 
		student, err := db.GetStudentByName(ctx, studentName)
		if err != nil {
			log.Error("Failed to retrieve student",
				"requestId", requestID,
				"error", err.Error(),
			)

			// If Student no in db 404
			if err.Error() == "no rows in result set" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Student not found",
				})
				return
			}

			// Errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve student",
			})
			return
		}

		// Log
		log.Info("Successfully retrieved student",
			"requestId", requestID,
			"studentId", student.ID,
			"duration", time.Since(startTime).String(),
		)

		// Return
		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}
