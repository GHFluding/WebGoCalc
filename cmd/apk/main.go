package main

import (
	"fmt"
	"log/slog"
	"test/internal/config"
	"test/internal/database/postgres"
	handler "test/internal/server/http/handlers"
	"test/internal/server/http/middleware"

	// setup logger
	sl "test/internal/services/slogger"

	"github.com/gin-gonic/gin"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	//init config: env
	cfg := config.MustLoad()

	//init logger: slog
	log := sl.SetupLogger(cfg.Env)

	//Logger is up
	log.Info("starting rep_cal", slog.String("env", cfg.Env))
	log.Debug("debug massages are enabled")

	// connection to DB
	dbpool, err := postgres.Connect(*cfg)
	if err != nil {
		log.Info("Failed to connect to database: ", sl.Err(err))
	}
	defer dbpool.Close()

	//add migration
	m, err := migrate.New(
		"file:///app/migrations", // migartion path
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Storage.User,
			cfg.Storage.Password,
			"webgocalc_postgres",
			cfg.Storage.Port,
			cfg.Storage.DBName,
		),
	)
	if err != nil {
		log.Error("Failed to connect migration: ", sl.Err(err))
	}
	if m == nil {
		log.Error("Migration object is nil")
		return
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error("Failed to use migration: ", sl.Err(err))
	} else {
		log.Info("Migrations applied successfully!")
	}

	// use Queries
	queries := postgres.New(dbpool)

	// use DB here
	_ = queries
	//init router: gin-gonic
	r := gin.Default()
	r.Use(middleware.RequestIdMiddleware())
	r.Use(sl.LoggingMiddleware(log)) //all loggers in package sl
	r.Use(middleware.RecoveryMiddleware(log))

	//router group for work with table students
	studentGroup := r.Group("/api/students")
	{
		studentGroup.GET("/list", handler.ListStudentsHandler(*queries, log))
		studentGroup.POST("/create", handler.CreateStudentHandler(*queries, log))
		studentGroup.POST("/delete", handler.DeleteStudentByNameHandler(*queries, log))
		studentGroup.POST("/get", handler.GetStudentByNameHandler(*queries, log))
		studentGroup.POST("/update", handler.UpdateStudentByNameHandler(*queries, log))
	}

	//router group for work with table calendar
	calendarGroup := r.Group("/api/calendar")
	{
		calendarGroup.GET("/daylist", handler.DayListHandler(*queries, log))
	}
	r.Run(cfg.HTTPServer.Address)
	//TODO: init controllers

	//TODO: init Services
}
