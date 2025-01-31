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

// UpdateStudentByIdHandler обновляет информацию о студенте по ID.
// @Summary      Обновить информацию о студенте
// @Description  Позволяет обновить определенные данные студента по его ID.
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        id      path      int64                          true  "ID студента" format(id)
// @Param        student body      nocsqlcpg.UpdateStudentSwagger true  "Данные для обновления"
// @Success      200     {object}  nocsqlcpg.StudentSwagger
// @Failure      400  {object}  map[string]interface{} "неверные данные"
// @Failure      404  {object}  map[string]interface{} "нет такого id"
// @Router       /api/students/{id} [patch]

func UpdateStudentByIdHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request JSON into a struct (this assumes your request body is properly structured)
		var req struct {
			Name string `json:"name"`
		}

		// Start the timer for request duration
		startTime := time.Now()

		// Extract the request ID for correlation
		requestID := middleware.RequestIdFromContext(c)

		// Log the incoming request with essential details
		extraFields := map[string]interface{}{
			"requestId": requestID,
			"url":       c.Request.URL.Path,
			"method":    c.Request.Method,
		}
		sl.LogRequestInfo(log, "info", c, "Handling UpdateStudent request", nil, extraFields)

		// Bind JSON from the request body to the struct
		if err := c.ShouldBindJSON(&req); err != nil {
			sl.LogRequestInfo(log, "error", c, "Failed to bind request data", err, extraFields)
			return
		}

		// Call the database function to update the student data
		s, err := nocsqlcpg.UpdateStudentData(db, c, log)
		if err != nil {
			// Log the error and return a failure response
			extraFields["string"] = s //more massage with error
			extraFields["error"] = err.Error()
			sl.LogRequestInfo(log, "error", c, "Failed to update student", err, extraFields)
			return
		}

		// Log successful student update, include duration
		extraFields["name"] = req.Name
		extraFields["duration"] = time.Since(startTime).String()
		sl.LogRequestInfo(log, "info", c, "Successfully updated student", nil, extraFields)
	}
}
