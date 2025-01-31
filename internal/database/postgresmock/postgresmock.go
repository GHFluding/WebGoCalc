// Файл: internal/database/postgres/mocks.go
package postgres

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockQueries реализует интерфейс Queries для тестов
type MockQueries struct {
	mock.Mock
}

// CreateStudent мок для метода CreateStudent
func (m *MockQueries) CreateStudent(ctx context.Context, params CreateStudentParams) (Student, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(Student), args.Error(1)
}

// GetStudentByID мок для метода GetStudentByID
func (m *MockQueries) GetStudentByID(ctx context.Context, id int32) (Student, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Student), args.Error(1)
}

// UpdateStudent мок для метода UpdateStudent
func (m *MockQueries) UpdateStudent(ctx context.Context, params UpdateStudentParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

// DeleteStudent мок для метода DeleteStudent
func (m *MockQueries) DeleteStudent(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ListStudents мок для метода ListStudents
func (m *MockQueries) ListStudents(ctx context.Context) ([]Student, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Student), args.Error(1)
}

// GetEventsByDate мок для метода GetEventsByDate
func (m *MockQueries) GetEventsByDate(ctx context.Context, date Date) ([]StudentEvent, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]StudentEvent), args.Error(1)
}

// ... Добавьте остальные методы интерфейса Queries
