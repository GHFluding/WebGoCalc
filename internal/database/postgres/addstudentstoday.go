package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func MoveStudentToday(db Queries, log *slog.Logger) error {
	now := time.Now()
	weekdayInt := int16(now.Weekday())
	weekday := pgtype.Int2{
		Int16: weekdayInt,
		Valid: true,
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
