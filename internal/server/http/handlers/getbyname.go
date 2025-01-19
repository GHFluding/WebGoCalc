package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"test/internal/database/postgres"
	"test/internal/server/http/middleware"
	sl "test/internal/services/slogger"
	"time"

	"github.com/gin-gonic/gin"
)

// GetStudentByIdHandler возвращает информацию о студенте по ID.
// @Summary      Получить информацию о студенте
// @Description  Возвращает полную информацию о студенте по его ID.
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID студента"
// @Success      200  {object}  postgres.StudentSwagger
// @Failure      400  {object}  map[string]interface{} "неверные данные"
// @Failure      404  {object}  map[string]interface{} "нет такого id"
// @Router       /api/students/{id} [get]

func GetStudentByIdHandler(db postgres.Queries, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve student ID from the request parameters
		studentIDParam := c.Param("id")
		if studentIDParam == "" {
			// Log the missing student ID error
			sl.LogRequestInfo(log, "error", c, "Missing student ID in request", nil, nil)

			// Return bad request response
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Student ID is required",
			})
			return
		}

		// Convert student ID from string to int64
		studentID, err := strconv.ParseInt(studentIDParam, 10, 64)
		if err != nil {
			// Log invalid student ID format error
			sl.LogRequestInfo(log, "error", c, "Invalid student ID format", err, map[string]interface{}{
				"studentIDParam": studentIDParam,
			})

			// Return bad request response
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid student ID format",
			})
			return
		}

		// Start the request timer and log the incoming request
		startTime := time.Now()
		requestID := middleware.RequestIdFromContext(c)
		extraFields := map[string]interface{}{
			"requestId": requestID,
			"url":       c.Request.URL.Path,
			"method":    c.Request.Method,
			"studentID": studentID,
		}
		sl.LogRequestInfo(log, "info", c, "Handling GetStudentById request", nil, extraFields)

		// Fetch student data from the database
		student, err := db.GetStudentById(context.Background(), studentID)
		if err != nil {
			// Log the error when failing to retrieve the student
			extraFields["error"] = err.Error()
			sl.LogRequestInfo(log, "error", c, "Failed to retrieve student", err, extraFields)

			// If student is not found in the database
			if err.Error() == "no rows in result set" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Student not found",
				})
				return
			}

			// Other errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve student",
			})
			return
		}

		// Log success after retrieving student data
		extraFields["studentId"] = student.ID
		extraFields["duration"] = time.Since(startTime).String()
		sl.LogRequestInfo(log, "info", c, "Successfully retrieved student", nil, extraFields)

		// Return student data in the response
		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}
