package event_generator

import (
	"context"
	"log/slog"
	"test/internal/database/postgres"
	sl "test/internal/utils/slogger"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Event handlers automatic event creation
type EventGenerator struct {
	queries    *postgres.Queries
	weeksAhead int
}

// New function a new EventGenerator
func New(q *postgres.Queries, weeksAhead int) *EventGenerator {
	return &EventGenerator{
		queries:    q,
		weeksAhead: weeksAhead,
	}
}

// GenerateEvents function create events for all students for configured weeks ahead
func (g *EventGenerator) GenerateEvents(ctx context.Context, log *slog.Logger) error {
	students, err := g.queries.ListStudents(ctx)
	if err != nil {
		log.Debug("failed to fetch students: ", sl.Err(err))
		return err
	}

	// Process each students schedule
	for _, student := range students {
		if err := g.processStudent(ctx, student, log); err != nil {
			log.Debug("Problems: ", slog.Int64("failed to process student, id: ", student.ID), "error: ", sl.Err(err))
			return err
		}
	}
	return nil
}
func (g *EventGenerator) processStudent(ctx context.Context, student postgres.Student, log *slog.Logger) error {
	//Calculate initial event date
	baseDate, err := g.calculateInitDate(student.OrderDay)
	if err != nil {
		return err
	}

	//Generate dates for configuration weeks ahead
	for week := 0; week < g.weeksAhead; week++ {
		eventDate := baseDate.AddDate(0, 0, 7*week)

		//Check if already exists
		exists, err := g.eventExists(ctx, student.ID, eventDate, log)
		if err != nil {
			return err
		}

		//Create event if missing
		if !exists {
			if err := g.createEvent(ctx, student, eventDate, log); err != nil {
				return err
			}
		}
	}
	return nil
}

// calculateInitDate determines the first occurrence date for the students schedule
func (g *EventGenerator) calculateInitDate(orderDay int16) (time.Time, error) {
	now := time.Now()
	currentWeekday := now.Weekday()
	// Convert student's order_day to time.Weekday (0=Sunday)
	targetWeekday := time.Weekday(orderDay % 7)
	// Calculate days until next occurrence
	daysToAdd := (int(targetWeekday) - int(currentWeekday) + 7) % 7
	if daysToAdd == 0 {
		// Schedule for next week if today is the target day
		daysToAdd = 7
	}
	return now.AddDate(0, 0, daysToAdd).Truncate(24 * time.Hour), nil
}

// eventExists checks if an event already exists for given student and date
func (g *EventGenerator) eventExists(ctx context.Context, studentID int64, date time.Time, log *slog.Logger) (bool, error) {
	var pgDate pgtype.Date
	if err := pgDate.Scan(date); err != nil {
		log.Debug("date conversion error: ", sl.Err(err))
		return false, err
	}
	events, err := g.queries.GetEventsByDate(ctx, pgDate)
	if err != nil {
		log.Debug("database request error: ", sl.Err(err))
		return false, err
	}

	//Check if students has existing event for this date
	for _, event := range events {
		if event.StudentID == studentID {
			return true, nil
		}
	}
	return false, nil
}

// createEvent insert a new event into the database
func (g *EventGenerator) createEvent(ctx context.Context, student postgres.Student, date time.Time, log *slog.Logger) error {
	var pgDate pgtype.Date
	if err := pgDate.Scan(date); err != nil {
		log.Debug("date conversion error: ", sl.Err(err))
		return err
	}

	params := postgres.AddEventsForDayParams{
		StudentID: student.ID,
		Column2:   pgDate,
		OrderTime: student.OrderTime,
		OrderCost: student.OrderCost,
	}
	if err := g.queries.AddEventsForDay(ctx, params); err != nil {
		log.Debug("failed to create event: ", sl.Err(err))
		return err
	}
	return nil
}
