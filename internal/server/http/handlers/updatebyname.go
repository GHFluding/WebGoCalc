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

func UpdateStudentByNameHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := `json:"name"`
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		log.Info("Handling UpdateStudents request",
			"requestId", requestID,
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
		)
		//request for db
		s, err := postgres.UpdateStudentData(db, c, log)
		//errors
		if err != nil {
			log.Error(s, "error", sl.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//log request
		log.Info("Successfully updated student",
			"requestId", requestID,
			"name", name,
			"duration", time.Since(startTime).String(),
		)

		// return
		c.JSON(http.StatusOK, gin.H{
			"message": "Student updated successfully",
		})
	}
}
