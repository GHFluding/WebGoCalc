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

// Handler для получения информации о студенте по имени
func GetStudentByNameHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получение имени студента из параметров запроса
		studentName := c.Param("name")
		if studentName == "" {
			log.Error("Missing student name in request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Student name is required",
			})
			return
		}

		// Создание контекста и логирование начала запроса
		ctx := context.Background()
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		log.Info("Handling GetStudentByName request",
			"requestId", requestID,
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
			"studentName", studentName,
		)

		// Получение данных студента из базы данных
		student, err := db.GetStudentByName(ctx, studentName)
		if err != nil {
			log.Error("Failed to retrieve student",
				"requestId", requestID,
				"error", err.Error(),
			)

			// Если студент не найден, возвращаем 404
			if err.Error() == "no rows in result set" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Student not found",
				})
				return
			}

			// Другие ошибки
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve student",
			})
			return
		}

		// Логируем успешное получение данных
		log.Info("Successfully retrieved student",
			"requestId", requestID,
			"studentId", student.ID,
			"duration", time.Since(startTime).String(),
		)

		// Возвращаем данные о студенте
		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}
