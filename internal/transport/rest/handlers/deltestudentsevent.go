package handler

import (
	"context"
	"log/slog"
	"strconv"
	"test/internal/database/postgres"
	"test/internal/server/http/middleware"
	sl "test/internal/services/slogger"
	"time"

	"github.com/gin-gonic/gin"
)

// DeleteEventByIdHandler удаляет студента по ID.
// @Summary      Удалить студента
// @Description  Удаляет студента из базы данных по его ID.
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        id   path      int64    "ID события" format(id)
// @Success      204
// @Failure      400  {object}  map[string]interface{} "неверные данные"
// @Failure      404  {object}  map[string]interface{} "нет такого id"
// @Router       /api/students/{id} [delete]

func DeleteEventByIdHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventIDParam := c.Param("id")
		if eventIDParam == "" {
			// Log error if event ID is missing
			sl.LogRequestInfo(log, "error", c, "Missing event ID in request", nil, nil)
			return
		}
		// Convert the event ID from string to int64
		eventID, err := strconv.ParseInt(eventIDParam, 10, 64)
		if err != nil {
			// Log error if event ID format is invalid
			sl.LogRequestInfo(log, "error", c, "Invalid event ID format", err, map[string]interface{}{
				"event ID": eventIDParam,
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
			"eventID":   eventID,
		}

		sl.LogRequestInfo(log, "info", c, "Handling DeleteEventById request", nil, extraFields)
		// Attempt to delete the event from the database
		err = db.DeleteStudentById(context.Background(), eventID)
		if err != nil {
			// Log the error if deleting the event fails
			extraFields["error"] = err.Error()
			sl.LogRequestInfo(log, "error", c, "Failed to delete event", err, extraFields)
			return
		}

		// Log success after deletion
		extraFields["duration"] = time.Since(startTime).String()
		sl.LogRequestInfo(log, "info", c, "Successfully deleted event", nil, extraFields)
	}
}
