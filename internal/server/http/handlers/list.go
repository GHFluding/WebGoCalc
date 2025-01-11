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

func ListStudentsHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		log.Info("Handling ListStudents request",
			"requestId", requestID,
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
		)

		//use List func sql
		students, err := db.ListStudents(ctx)
		//check errors and log
		if err != nil {
			log.Error("Failed to retrieve students",
				"requestId", requestID,
				"error", err.Error(),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve students",
			})
			return
		}
		// return data
		log.Info("Successfully retrieved students",
			"requestId", requestID,
			"count", len(students),
			"duration", time.Since(startTime).String(),
		)
		c.JSON(http.StatusOK, gin.H{
			"students": students,
		})
	}
}
