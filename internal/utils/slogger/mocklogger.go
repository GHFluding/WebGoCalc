package sl

import (
	"log/slog"
	"os"
)

// MockLogger
func SetupMockLogger() *slog.Logger {

	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)

	return log
}
