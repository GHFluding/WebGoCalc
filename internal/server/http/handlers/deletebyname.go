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

func DeleteStudentByNameHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Getting name from request
		studentName := c.Param("name")
		if studentName == "" {
			log.Error("Missing student name in request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Student name is required",
			})
			return
		}

		// create context and log start of request
		ctx := context.Background()
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		log.Info("Handling DeleteStudentByName request",
			"requestId", requestID,
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
			"studentName", studentName,
		)

		// delete in db
		err := db.DeleteStudentByName(ctx, studentName)
		if err != nil {
			log.Error("Failed to delete student",
				"requestId", requestID,
				"error", err.Error(),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete student",
			})
			return
		}

		// log request complete
		log.Info("Successfully deleted student",
			"requestId", requestID,
			"studentName", studentName,
			"duration", time.Since(startTime).String(),
		)

		// return request complete
		c.JSON(http.StatusOK, gin.H{
			"message": "Student successfully deleted",
		})
	}
}
