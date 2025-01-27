package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"test/internal/config"
	"test/internal/database/postgres"

	handler "test/internal/server/http/handlers"
	sl "test/internal/services/slogger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// db mock
type MockDBTX struct {
	mock.Mock
}

func (m *MockDBTX) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	argsMock := m.Called(ctx, query, args)
	return argsMock.Get(0).(pgconn.CommandTag), argsMock.Error(1)
}

func (m *MockDBTX) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	argsMock := m.Called(ctx, query, args)
	return argsMock.Get(0).(pgx.Rows), argsMock.Error(1)
}

func (m *MockDBTX) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	argsMock := m.Called(ctx, query, args)
	return argsMock.Get(0).(pgx.Row)
}

// mock rows
type MockRows struct {
	Rows    [][]interface{}
	Columns []string
	Index   int
}

func (m *MockRows) RawValues() [][]byte {
	return nil
}
func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return []pgconn.FieldDescription{}
}
func (m *MockRows) Conn() *pgx.Conn {
	return &pgx.Conn{}
}
func (m *MockRows) Values() ([]any, error) {
	result := make([]any, len(m.Rows))
	for i, backup := range m.Rows {
		result[i] = backup[i]
	}
	return result, nil
}
func (m *MockRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag(`-- name: ListStudents :many
SELECT id, name, s_class, school, order_day, order_time, order_cost FROM students
ORDER BY id
`)
}
func (m *MockRows) Next() bool {
	m.Index++
	return m.Index < len(m.Rows)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	if m.Index >= len(m.Rows) {
		return fmt.Errorf("out of bounds")
	}
	row := m.Rows[m.Index]
	if len(dest) != len(row) {
		return fmt.Errorf("mismatch between destination and source lengths")
	}
	copy(dest, row)
	return nil
}

func (m *MockRows) Close() {}

func (m *MockRows) Err() error {
	return nil
}

// TODO fix
// test func
func TestListStudentsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Создаем мок DBTX
	mockDBTX := new(MockDBTX)

	// Пример данных студентов
	mockRows := &MockRows{
		Rows: [][]interface{}{
			{int64(1), "John Doe", "10A", "High School", int16(1), "63000", int16(100)},
			{int64(2), "Jane Smith", "11B", "Secondary School", int16(2), "52000", int16(150)},
		},
		Columns: []string{"id", "name", "s_class", "school", "order_day", "order_time", "order_cost"},
	}

	// Строка запроса SQL, включая комментарии, как в коде
	listStudentsQuery := `-- name: ListStudents :many
SELECT id, name, s_class, school, order_day, order_time, order_cost FROM students
ORDER BY id
`

	// Настроим мок так, чтобы он ожидал правильный запрос
	mockDBTX.On("Query", mock.Anything, listStudentsQuery, mock.Anything).Return(mockRows, nil)

	// Создаем объект Queries с моком DBTX
	queries := postgres.New(mockDBTX)

	// Настроим конфигурацию для теста
	cfg := config.Config{
		Env: "local", // Указываем окружение
		Storage: config.Storage{
			Path:     "/app/storage",
			Host:     "webgocalc_postgres",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "students",
			AsString: "",
		},
		HTTPServer: config.HTTPServer{
			Address:     "0.0.0.0:8080",
			TimeOut:     5 * time.Second,
			IdleTimeOut: 60 * time.Second,
		},
	}

	// Настроим логгер с учетом окружения
	logger := sl.SetupLogger(cfg.Env)

	// Регистрируем хендлер
	router.GET("/api/students", handler.ListStudentsHandler(*queries, logger))

	// Создаем тестовый запрос
	req, err := http.NewRequest(http.MethodGet, "/api/students", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{
		"students": [
			{"id": 1, "name": "John Doe", "s_class": "10A", "school": "High School", "order_day": 1, "order_time": "10:00:00", "order_cost": 100},
			{"id": 2, "name": "Jane Smith", "s_class": "11B", "school": "Secondary School", "order_day": 2, "order_time": "12:00:00", "order_cost": 150}
		]
	}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	// Проверяем, что мок был вызван
	mockDBTX.AssertExpectations(t)
}
