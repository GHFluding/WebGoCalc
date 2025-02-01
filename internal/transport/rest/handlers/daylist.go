package handler

import (
	"context"
	"log/slog"
	"net/http"
	"test/internal/database/postgres"
	"test/internal/transport/rest/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	sl "test/internal/services/slogger"
)

// DayListHandler возвращает список событий на день.
// @Summary      Получить список событий
// @Description  Возвращает события календаря на указанный день.
// @Tags         calendar
// @Accept       json
// @Produce      json
// @Param        date  query  string  true  "Дата в формате YYYY-MM-DD"
// @Success      200  {array}  nocsqlcpg.StudentEventSwagger
// @Failure      400  {object}  map[string]interface{} "неверные данные"
// @Failure 500 {object} map[string]interface{} "Ошибка сервера"
// @Router       /api/calendar [get]
func DayListHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		startTime := time.Now()

		// Retrieve request ID from middleware context
		requestID := middleware.RequestIdFromContext(c)

		// Log the start of the request with relevant details (using the new logger)
		extraFields := map[string]interface{}{
			"requestId": requestID,
			"url":       c.Request.URL.Path,
			"method":    c.Request.Method,
		}
		sl.LogRequestInfo(log, "info", c, "Handling DayList request", nil, extraFields)

		// Create the 'today' value with the current date
		today := pgtype.Date{Time: startTime, Valid: true}

		// Fetch the events for today from the database

		students, err := db.GetEventsByDate(context.Background(), today)
		if err != nil {
			// Log the error if retrieving students fails
			extraFields["error"] = err.Error()
			sl.LogRequestInfo(log, "error", c, "Failed to retrieve students for this date", err, extraFields)
			return
		}

		// Log the successful retrieval of data
		extraFields["count"] = len(students)
		extraFields["duration"] = time.Since(startTime).String()
		sl.LogRequestInfo(log, "info", c, "Successfully retrieved students", nil, extraFields)

		// Return the list of students
		c.JSON(http.StatusOK, gin.H{
			"students": students,
		})
	}
}
