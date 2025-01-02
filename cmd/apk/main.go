package main

import (
	"log/slog"
	"test/internal/config"
	"test/internal/database/postgres"
	handler "test/internal/server/http/handlers"
	"test/internal/server/http/middleware"

	// setup logger
	sl "test/internal/services/slogger"

	"github.com/gin-gonic/gin"
)

func main() {
	//init config: env
	cfg := config.MustLoad()

	//init logger: slog
	log := sl.SetupLogger(cfg.Env)

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
	//init router: gin-gonic
	r := gin.Default()
	r.Use(middleware.RequestIdMiddleware())
	r.Use(sl.LoggingMiddleware(log)) //all loggers in package sl
	r.Use(middleware.RecoveryMiddleware(log))
	r.GET("/studentslist", handler.ListStudentsHandler(*queries, log))
	r.POST("/create-student", handler.CreateStudentHandler(*queries, log))

	r.Run(cfg.HTTPServer.Address)
	//TODO: init controllers

	//TODO: init Services
}
