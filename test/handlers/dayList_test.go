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

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

const EventID = int64(1)

func TestDayListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	OrderCheck := pgtype.Bool{
		Bool:  false,
		Valid: true,
	}
	today := time.Date(2023, time.June, 15, 0, 0, 0, 0, time.UTC)
	EventDate := pgtype.Date{
		Time:             today,
		InfinityModifier: 0,
		Valid:            true,
	}
	OrderTime := pgtype.Time{
		Microseconds: time.Date(0, 1, 1, 15, 30, 0, 0, time.UTC).UnixMicro(),
		Valid:        true,
	}
	t.Run("Успешное получение events", func(t *testing.T) {
		mockDB := &mocks.MockDB{
			QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
				slog.Info("Mock DB Query called", "sql", sql, "args", args)

				data := [][]interface{}{
					{
						EventID,     // StudentEventsID
						StudentID,   // StudentID
						StudentName, // StudentName
						EventDate,   // EventDate
						OrderTime,   // OrderTime
						OrderCost,   // OrderCost
						OrderCheck,
					},
				}
				return mocks.NewMockRows(data), nil
			},
		}

		queries := postgres.New(mockDB)
		router := gin.New()
		//maybe not safe
		patch := monkey.Patch(time.Now, func() time.Time {
			return today
		})
		defer patch.Unpatch()

		router.GET("/api/calendar/", handler.DayListHandler(*queries, testLogger))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/calendar/", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Events []struct {
				ID          int64       `json:"StudentEventsID"`
				StudentID   int64       `json:"StudentID"` // Измените на int64, если это число
				StudentName string      `json:"StudentName"`
				EventDate   pgtype.Date `json:"EventDate"` // Преобразуем в строку
				OrderTime   pgtype.Time `json:"OrderTime"` // Преобразуем время в строку
				OrderCost   int16       `json:"OrderCost"`
				OrderCheck  pgtype.Bool `json:"OrderCheck"`
			} `json:"students"`
		}
		t.Log("Response Body:", w.Body.String())

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, EventID, response.Events[0].ID)
		assert.Equal(t, StudentID, response.Events[0].StudentID)
		assert.Equal(t, StudentName, response.Events[0].StudentName)
		assert.Equal(t, EventDate, response.Events[0].EventDate)
		assert.Equal(t, OrderTime, response.Events[0].OrderTime)
		assert.Equal(t, OrderCost, response.Events[0].OrderCost)
		assert.Equal(t, OrderCheck, response.Events[0].OrderCheck)

	})
}
