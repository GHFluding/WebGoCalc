package main
//
import (
	"log/slog"
	"os"
	"test/internal/config"
	"test/internal/database/postgres"
	"test/internal/services/sl"

	"github.com/gin-gonic/gin"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//init config: env
	cfg := config.MustLoad()

	//init logger: slog
	log := setupLogger(cfg.Env)

	//Logger is up
	log.Info("starting rep_cal", slog.String("env", cfg.Env))
	log.Debug("debug masseges are enabled")

	// connection to DB
	dbpool, err := postgres.Connect(*cfg)
	if err != nil {
		log.Info("Failed to connect to database: ", sl.Err(err))
	}
	defer dbpool.Close()

	// use Queries
	queries := postgres.New(dbpool)

	// use DB here
	_ = queries
	//TODO: init router: gin-gonic
	router := gin.Default()
	router.Use()

	//TODO: init controllers

	//TODO: init Services

}

// setup Logger configuraton
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}
	return log
}
