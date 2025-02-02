package mocks

import (
	"log/slog"
	"os"
)

func SetupLoggerMock() *slog.Logger {
	log := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)
	return log
}
