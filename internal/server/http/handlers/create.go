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

// Handler for create new student in db
func CreateStudentHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		log.Info("Handling CreateStudent request",
			"requestId", requestID,
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
		)
		//request for db
		s, student, err := postgres.CreateStudentData(db, c, log)
		//errors
		if err != nil {
			log.Error(s, "error", sl.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//log request
		log.Info(s,
			"requestId", requestID,
			"studentId", student.ID,
			"duration", time.Since(startTime).String(),
		)
		// return result to user
		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}
