package handler

import (
	"context"
	"log/slog"
	"strconv"
	"test/internal/database/postgres"
	sl "test/internal/services/slogger"
	"test/internal/transport/rest/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// DeleteStudentByIdHandler удаляет студента по ID.
// @Summary      Удалить студента
// @Description  Удаляет студента из базы данных по его ID.
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        id   path      int64  "ID студента" format(id)
// @Success      204
// @Failure      400  {object}  map[string]interface{} "неверные данные"
// @Failure      404  {object}  map[string]interface{} "нет такого id"
// @Router       /api/students/{id} [delete]

func DeleteStudentByIdHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve student ID from request parameters
		studentIDParam := c.Param("id")
		if studentIDParam == "" {
			// Log error if student ID is missing
			sl.LogRequestInfo(log, "error", c, "Missing student ID in request", nil, nil)
			return
		}

		// Convert the student ID from string to int64
		studentID, err := strconv.ParseInt(studentIDParam, 10, 64)
		if err != nil {
			// Log error if student ID format is invalid
			sl.LogRequestInfo(log, "error", c, "Invalid student ID format", err, map[string]interface{}{
				"studentIDParam": studentIDParam,
			})
			return
		}

		// Create context and log the start of the request
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		extraFields := map[string]interface{}{
			"requestId": requestID,
			"url":       c.Request.URL.Path,
			"method":    c.Request.Method,
			"studentID": studentID,
		}
		sl.LogRequestInfo(log, "info", c, "Handling DeleteStudentById request", nil, extraFields)

		// Attempt to delete the student from the database
		err = db.DeleteStudentById(context.Background(), studentID)
		if err != nil {
			// Log the error if deleting the student fails
			extraFields["error"] = err.Error()
			sl.LogRequestInfo(log, "error", c, "Failed to delete student", err, extraFields)
			return
		}

		// Log success after deletion
		extraFields["duration"] = time.Since(startTime).String()
		sl.LogRequestInfo(log, "info", c, "Successfully deleted student", nil, extraFields)
	}
}
