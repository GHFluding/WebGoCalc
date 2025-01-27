package handler

import (
	"log/slog"
	"test/internal/database/postgres"
	nocsqlcpg "test/internal/database/postgres/nosqlcpg"
	"test/internal/server/http/middleware"
	sl "test/internal/services/slogger"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateStudentHandler создает нового студента.
// @Summary      Создать студента
// @Description  Добавляет нового студента в базу данных.
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        student  body  postgres.CreateStudentSwagger  true  "Данные студента"
// @Success      201  {object}  postgres.StudentSwagger
// @Failure      400  {object}  map[string]interface{} "неверные данные"
// @Failure 500 {object} map[string]interface{} "Ошибка сервера"
// @Router       /api/students [post]
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
		s, student, err := nocsqlcpg.CreateStudentData(db, c, log)

		// Handle errors
		if err != nil {
			extraFields["string"] = s
			extraFields["error"] = err.Error()
			sl.LogRequestInfo(log, "error", c, "Failed to create student", err, extraFields)
			return
		}

		// Log successful creation
		extraFields["studentId"] = student.ID
		extraFields["duration"] = time.Since(startTime).String()
		sl.LogRequestInfo(log, "info", c, "Successfully created student", nil, extraFields)
	}
}
