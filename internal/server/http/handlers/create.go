package handler

import (
	"context"
	"log/slog"
	"net/http"
	"test/internal/database/postgres"
	"test/internal/server/http/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler для добавления студента
func CreateStudentHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Парсинг входных данных
		var req struct {
			Name      string `json:"name" binding:"required"`
			Clas      string `json:"clas" binding:"required"`
			Scool     string `json:"scool" binding:"required"`
			OrderDay  int16  `json:"order_day" binding:"required"`
			OrderTime string `json:"order_time" binding:"required"` // Ожидаем строку в формате времени
			OrderCost int16  `json:"order_cost" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Invalid request payload", "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		// Парсинг времени
		orderTime, err := time.Parse("15:04", req.OrderTime)
		if err != nil {
			log.Error("Invalid time format", "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid time format (expected HH:MM)",
			})
			return
		}

		// Преобразование времени в микросекунды
		pgOrderTime := pgtype.Time{
			Microseconds: int64(orderTime.Hour()*3600+orderTime.Minute()*60+orderTime.Second()) * 1_000_000,
			Valid:        true,
		}

		// Формирование аргументов для SQL-запроса
		arg := postgres.CreateStudentParams{
			Name:      req.Name,
			Clas:      pgtype.Text{String: req.Clas, Valid: true},
			Scool:     pgtype.Text{String: req.Scool, Valid: true},
			OrderDay:  pgtype.Int2{Int16: req.OrderDay, Valid: true},
			OrderTime: pgOrderTime,
			OrderCost: pgtype.Int2{Int16: req.OrderCost, Valid: true},
		}

		// Вставка данных в базу
		ctx := context.Background()
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		log.Info("Handling CreateStudent request",
			"requestId", requestID,
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
		)

		student, err := db.CreateStudent(ctx, arg)
		if err != nil {
			log.Error("Failed to create student",
				"requestId", requestID,
				"error", err.Error(),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create student",
			})
			return
		}

		// Логируем успешную операцию
		log.Info("Successfully created student",
			"requestId", requestID,
			"studentId", student.ID,
			"duration", time.Since(startTime).String(),
		)

		// Возвращаем результат клиенту
		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}
