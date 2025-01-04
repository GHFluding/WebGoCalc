package postgres

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type request struct {
	ID        int64  `json:"id" `
	Name      string `json:"name" binding:"required"`
	Class     string `json:"class"`
	School    string `json:"school"`
	OrderDay  int16  `json:"order_day"`
	OrderTime string `json:"order_time"` // Expected string format Time
	OrderCost int16  `json:"order_cost"`
}

// check empty fields in struct
func setDefaultValuesManually(req *request, student *Student, arg *UpdateStudentByNameParams) {
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
	log.Debug("request data: ", req)
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
func CreateStudentData(db Queries, c *gin.Context, log *slog.Logger) (string, Student, error) {
	req, err := parseRequest(c, log)
	if err != nil {
		return "Invalid request payload", Student{}, err
	}
	pgTime, err := parseTime(c, &req, log)
	if err != nil {
		return "Invalid time format", Student{}, err
	}
	arg := CreateStudentParams{
		Name:      req.Name,
		SClass:    pgtype.Text{String: req.Class, Valid: true},
		School:    pgtype.Text{String: req.School, Valid: true},
		OrderDay:  pgtype.Int2{Int16: req.OrderDay, Valid: true},
		OrderTime: pgTime,
		OrderCost: pgtype.Int2{Int16: req.OrderCost, Valid: true},
	}
	ctx := context.Background()
	student, err := db.CreateStudent(ctx, arg)
	if err != nil {
		return "Failed to create student", Student{}, err
	}
	return "Successfully created student", student, err
}

// update students
func UpdateStudentData(db Queries, c *gin.Context, log *slog.Logger) (string, error) {
	req, err := parseRequest(c, log)
	if err != nil {
		return "Invalid request payload", err
	}
	ctx := context.Background()
	student, ok := db.GetStudentByName(ctx, req.Name)
	pgTime, err := parseTime(c, &req, log)
	if err != nil {
		return "Invalid time format", err
	}
	arg := UpdateStudentByNameParams{
		Name:      req.Name,
		SClass:    pgtype.Text{String: req.Class, Valid: true},
		School:    pgtype.Text{String: req.School, Valid: true},
		OrderDay:  pgtype.Int2{Int16: req.OrderDay, Valid: true},
		OrderTime: pgTime,
		OrderCost: pgtype.Int2{Int16: req.OrderCost, Valid: true},
	}
	if ok == nil {
		setDefaultValuesManually(&req, &student, &arg)
	}
	ctx = context.Background()
	err = db.UpdateStudentByName(ctx, arg)
	if err != nil {
		return "Failed to update student", err
	}
	return "Successfully update student", err
}
