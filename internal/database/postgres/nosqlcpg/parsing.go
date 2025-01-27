package nocsqlcpg

import (
	"context"
	"log/slog"
	"net/http"
	"test/internal/database/postgres"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// типы для swagger
type StudentSwagger struct {
	ID        int64
	Name      string
	SClass    string
	School    string
	OrderDay  int
	OrderTime string
	OrderCost string
}
type EventSwagger struct {
	CalendarID  int64
	StudentID   int64
	StudentName string
	EventDate   string
	OrderTime   string
	OrderCost   int16
	OrderCheck  bool
}
type CreateStudentSwagger struct {
	Name      string
	SClass    string
	School    string
	OrderDay  int
	OrderTime string
	OrderCost int
}
type UpdateStudentParams struct {
	ID        int64
	Name      string
	SClass    pgtype.Text
	School    pgtype.Text
	OrderDay  pgtype.Int2
	OrderTime pgtype.Time
	OrderCost pgtype.Int2
}

type request struct {
	ID        int64  `json:"id"`
	Name      string `json:"name" `
	Class     string `json:"class"`
	School    string `json:"school"`
	OrderDay  int16  `json:"order_day"`
	OrderTime string `json:"order_time"` // Expected string format Time
	OrderCost int16  `json:"order_cost"`
}

// check empty fields in struct
func setDefaultValuesManually(req *request, student *postgres.Student, arg *postgres.UpdateStudentByIdParams) {

	if req.Name == "" {
		arg.Name = student.Name
	}
	if req.Class == "" {
		arg.SClass = student.SClass
	}
	if req.School == "" {
		arg.School = student.School
	}
	if req.OrderTime == "" {
		arg.OrderTime = student.OrderTime
	}
	if req.OrderDay == 0 {
		arg.OrderDay = student.OrderDay
	}
	if req.OrderCost == 0 {
		arg.OrderCost = student.OrderCost
	}
}

// parse json
func parseRequest(c *gin.Context, log *slog.Logger) (request, error) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid request payload", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return req, err
	}

	return req, nil
}

// parse time
func parseTime(c *gin.Context, r *request, log *slog.Logger) (pgtype.Time, error) {
	var ok error
	orderTime, err := time.Parse("15:04", r.OrderTime)
	if err != nil {
		log.Error("Invalid time format", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid time format (expected HH:MM)",
		})
		ok = err
	}
	return pgtype.Time{
		Microseconds: int64(orderTime.Hour()*3600+orderTime.Minute()*60) * 1_000_000,
		Valid:        true,
	}, ok

}

// create sql request for Create students
func CreateStudentData(db postgres.Queries, c *gin.Context, log *slog.Logger) (string, postgres.Student, error) {
	req, err := parseRequest(c, log)
	if err != nil {
		return "Invalid request payload", postgres.Student{}, err
	}
	pgTime, err := parseTime(c, &req, log)
	if err != nil {
		return "Invalid time format", postgres.Student{}, err
	}
	arg := postgres.CreateStudentParams{
		Name:      req.Name,
		SClass:    req.Class,
		School:    req.School,
		OrderDay:  req.OrderDay,
		OrderTime: pgTime,
		OrderCost: req.OrderCost,
	}
	ctx := context.Background()
	student, err := db.CreateStudent(ctx, arg)
	if err != nil {
		return "Failed to create student", postgres.Student{}, err
	}
	return "Successfully created student", student, err
}

// update students
func UpdateStudentData(db postgres.Queries, c *gin.Context, log *slog.Logger) (string, error) {
	req, err := parseRequest(c, log)
	if err != nil {
		return "Invalid request payload", err
	}
	ctx := context.Background()
	student, ok := db.GetStudentById(ctx, req.ID)
	pgTime, err := parseTime(c, &req, log)
	if err != nil {
		return "Invalid time format", err
	}
	arg := postgres.UpdateStudentByIdParams{
		ID:        req.ID,
		Name:      req.Name,
		SClass:    req.Class,
		School:    req.School,
		OrderDay:  req.OrderDay,
		OrderTime: pgTime,
		OrderCost: req.OrderCost,
	}
	if ok == nil {
		setDefaultValuesManually(&req, &student, &arg)
	}
	ctx = context.Background()
	err = db.UpdateStudentById(ctx, arg)
	if err != nil {
		return "Failed to update student", err
	}
	return "Successfully update student", err
}
