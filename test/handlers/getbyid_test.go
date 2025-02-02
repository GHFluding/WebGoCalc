package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"test/internal/database/postgres"
	handler "test/internal/transport/rest/handlers"
	"test/test/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func String(t pgtype.Time) string {
	return time.UnixMicro(t.Microseconds).UTC().Format("15:04:05")
}

func TestGetStudentByIdHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("Успешное получение студента", func(t *testing.T) {
		mockDB := &mocks.MockDB{
			QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
				slog.Info("Mock DB QueryRow called", "sql", sql, "args", args)
				orderTime := pgtype.Time{
					Microseconds: time.Date(0, 1, 1, 15, 30, 0, 0, time.UTC).UnixMicro(),
					Valid:        true,
				}

				return &mocks.MockRow{
					Data: []interface{}{
						int64(1),
						"John Doe",
						"10A",
						"Springfield HS",
						int16(5),
						orderTime,
						int16(2500),
					},
				}
			},
		}

		queries := postgres.New(mockDB)
		router := gin.New()
		router.GET("/api/students/:id", handler.GetStudentByIdHandler(*queries, testLogger))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/students/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Student struct {
				ID        int64       `json:"ID"`
				Name      string      `json:"Name"`
				SClass    string      `json:"SClass"`
				School    string      `json:"School"`
				OrderDay  int16       `json:"OrderDay"`
				OrderTime pgtype.Time `json:"OrderTime"`
				OrderCost int16       `json:"OrderCost"`
			} `json:"student"`
		}
		t.Log("Response Body:", w.Body.String())

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), response.Student.ID)
		assert.Equal(t, "John Doe", response.Student.Name)
		assert.Equal(t, "10A", response.Student.SClass)
		assert.Equal(t, "Springfield HS", response.Student.School)
		assert.Equal(t, int16(5), response.Student.OrderDay)
		assert.Equal(t, int16(2500), response.Student.OrderCost)

		formattedTime := time.UnixMicro(response.Student.OrderTime.Microseconds).UTC().Format("15:04:05")
		assert.Equal(t, "15:30:00", formattedTime)

	})
}
