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

func UpdateStudentByNameHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parsing input data
		var req struct {
			Name      string `json:"name" binding:"required"`
			Clas      string `json:"clas" binding:"required"`
			Scool     string `json:"scool" binding:"required"`
			OrderDay  int16  `json:"order_day" binding:"required"`
			OrderTime string `json:"order_time" binding:"required"` // Expectet string format Time
			OrderCost int16  `json:"order_cost" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Invalid request payload", "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		// Parsing time
		orderTime, err := time.Parse("15:04", req.OrderTime)
		if err != nil {
			log.Error("Invalid time format", "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid time format (expected HH:MM)",
			})
			return
		}

		// Time change
		pgOrderTime := pgtype.Time{
			Microseconds: int64(orderTime.Hour()*3600+orderTime.Minute()*60+orderTime.Second()) * 1_000_000,
			Valid:        true,
		}

		// SQL-request args
		arg := postgres.UpdateStudentByNameParams{
			Name:      req.Name,
			Clas:      pgtype.Text{String: req.Clas, Valid: true},
			Scool:     pgtype.Text{String: req.Scool, Valid: true},
			OrderDay:  pgtype.Int2{Int16: req.OrderDay, Valid: true},
			OrderTime: pgOrderTime,
			OrderCost: pgtype.Int2{Int16: req.OrderCost, Valid: true},
		}

		// DB update
		ctx := context.Background()
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		log.Info("Handling UpdateStudentByName request",
			"requestId", requestID,
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
		)

		err = db.UpdateStudentByName(ctx, arg)
		if err != nil {
			log.Error("Failed to update student",
				"requestId", requestID,
				"error", err.Error(),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update student",
			})
			return
		}

		// log
		log.Info("Successfully updated student",
			"requestId", requestID,
			"name", req.Name,
			"duration", time.Since(startTime).String(),
		)

		// return
		c.JSON(http.StatusOK, gin.H{
			"message": "Student updated successfully",
		})
	}
}
