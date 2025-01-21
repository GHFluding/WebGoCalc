package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func MoveStudentToday(db Queries, log *slog.Logger) error {
	now := time.Now()
	weekdayString := now.Weekday().String()
	var weekday pgtype.Int2
	switch weekdayString {
	case "Monday":
		weekday = pgtype.Int2{
			Int16: 1,
			Valid: true,
		}
	case "Tuesday":
		weekday = pgtype.Int2{
			Int16: 2,
			Valid: true,
		}
	case "Wednesday":
		weekday = pgtype.Int2{
			Int16: 3,
			Valid: true,
		}
	case "Thursday":
		weekday = pgtype.Int2{
			Int16: 4,
			Valid: true,
		}
	case "Friday":
		weekday = pgtype.Int2{
			Int16: 5,
			Valid: true,
		}
	case "Saturday":
		weekday = pgtype.Int2{
			Int16: 6,
			Valid: true,
		}
	case "Sunday":
		weekday = pgtype.Int2{
			Int16: 7,
			Valid: true,
		}
	}
	students, err := db.GetStudentsByDay(context.Background(), weekday)
	if err != nil {
		return err
	}

	pgnow := pgtype.Date{
		Time:             now,
		InfinityModifier: 0,
		Valid:            true,
	}

	for i := range students {
		addParasm := AddEventsForDayParams{
			StudentID: students[i].ID,
			Column2:   pgnow,
			OrderTime: students[i].OrderTime,
			OrderCost: students[i].OrderCost.Int16,
		}
		db.AddEventsForDay(context.Background(), addParasm)
	}
	return nil
}
