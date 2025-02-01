package pgmodels

import (
	"time"
)

// Datatype for swagger docs
type DateSwagger struct {
	Date string
}

type StudentSwagger struct {
	ID        int64
	Name      string
	SClass    string
	School    string
	OrderDay  int16
	OrderTime OrderTime
	OrderCost int16
}

type OrderTime struct {
	Microseconds int64
	Valid        bool
}

type StudentEventSwagger struct {
	ID         int64
	StudentID  int64
	EventDate  EventDate
	OrderTime  OrderTime
	OrderCost  int16
	OrderCheck OrderCheck
}
type EventDate struct {
	Time             time.Time
	InfinityModifier int8
	Valid            bool
}
type OrderCheck struct {
	Bool  bool
	Valid bool
}
type CreateStudentSwagger struct {
	Name      string
	SClass    string
	School    string
	OrderDay  int16
	OrderTime string
	OrderCost int16
}
type UpdateStudentSwagger struct {
	ID        int64
	Name      string
	SClass    string
	School    string
	OrderDay  int16
	OrderTime OrderTime
	OrderCost int16
}
