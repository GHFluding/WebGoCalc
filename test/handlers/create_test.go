package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"test/internal/database/postgres"
	nocsqlcpg "test/internal/database/postgres/nosqlcpg"
	handler "test/internal/server/http/handlers"
	sl "test/internal/services/slogger"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockQueries реализует интерфейс postgres.Queries
type MockQueries struct {
	mock.Mock
}

func (m *MockQueries) CreateStudent(ctx *gin.Context, params postgres.CreateStudentParams) (postgres.Student, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(postgres.Student), args.Error(1)
}

// Остальные методы интерфейса...

func TestCreateStudentHandler_Success(t *testing.T) {
	// Инициализация
	gin.SetMode(gin.TestMode)
	log := sl.NewMockLogger()

	mockDB := new(MockQueries)
	handlerFunc := handler.CreateStudentHandler(mockDB, log)

	// Мок данных студента
	studentInput := nocsqlcpg.CreateStudentSwagger{
		Name:      "Mihael Dse",
		SClass:    "11B",
		School:    "Schosl Nsme",
		OrderDay:  6,
		OrderTime: "18:30",
		OrderCost: 5000,
	}

	expectedStudent := postgres.Student{
		ID:       1,
		Name:     "Mihael Dse",
		SClass:   "11B",
		School:   "Schosl Nsme",
		OrderDay: 6,
		OrderTime: pgtype.Time{
			Microseconds: 66600000000,
			Valid:        true,
		},
		OrderCost: 5000,
	}

	// Настройка моков
	mockDB.On("CreateStudent", mock.Anything, mock.Anything).
		Return(expectedStudent, nil)

	// Создание тестового запроса
	body, _ := json.Marshal(studentInput)
	req, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handlerFunc(ctx)

	// Проверки
	require.Equal(t, http.StatusCreated, w.Code)

	var response postgres.Student
	json.Unmarshal(w.Body.Bytes(), &response)
	require.Equal(t, expectedStudent.ID, response.ID)
	require.Equal(t, expectedStudent.Name, response.Name)

	mockDB.AssertExpectations(t)
}

func TestCreateStudentHandler_InvalidData(t *testing.T) {
	// Initialization
	gin.SetMode(gin.TestMode)
	log := sl.SetupMockLogger()

	mockDB := new(MockQueries)
	handlerFunc := handler.CreateStudentHandler(mockDB, log)

	// Invalid data
	invalidInput := nocsqlcpg.CreateStudentSwagger{
		OrderCost: 1500,
	}

	// Make Request
	body, _ := json.Marshal(invalidInput)
	req, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Run test
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handlerFunc(ctx)

	// check data
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "неверные данные")
}

// Вспомогательные функции
func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
