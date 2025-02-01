package main

import (
	"context"
	"fmt"
	"log/slog"
	"test/internal/config"
	"test/internal/database/postgres"
	"test/internal/models/pgmodels"
	handler "test/internal/transport/rest/handlers"
	"test/internal/transport/rest/middleware"

	// setup logger
	"test/internal/services/event_generator"
	sl "test/internal/services/slogger"

	_ "test/docs"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Student and Calendar API
// @version         1.0
// @description     API для управления студентами и событиями в календаре.
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.url    http://http://81.177.220.96/
// @contact.email  lyoshabura@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @schemes         http https
// @accept          json
// @produce         json

func main() {
	//init config: env
	cfg := config.MustLoad()

	//init logger: slog
	log := sl.SetupLogger(cfg.Env)

	//Logger is up
	log.Info("starting rep_cal", slog.String("env", cfg.Env))
	log.Debug("debug massages are enabled")

	// connection to DB
	dbpool, err := pgmodels.Connect(*cfg)
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
	// Initialize and run event generator
	// Generate 12 weeks ahead
	generator := event_generator.New(queries, 12)
	if err := generator.GenerateEvents(context.Background(), log); err != nil {
		log.Error("Event generation failed", sl.Err(err))
	} else {
		log.Info("Successfully generated future events")
	}
	// use DB here
	_ = queries
	//init router: gin-gonic
	r := gin.Default()
	r.Use(middleware.RequestIdMiddleware())
	//all loggers in package sl
	r.Use(sl.LoggingMiddleware(log))
	r.Use(middleware.RecoveryMiddleware(log))

	//router group for work with table students
	// Group for working with student-related routes
	studentGroup := r.Group("/api/students")
	{

		studentGroup.GET("/", handler.ListStudentsHandler(*queries, log))

		studentGroup.POST("/", handler.CreateStudentHandler(*queries, log))

		studentGroup.DELETE("/:id", handler.DeleteStudentByIdHandler(*queries, log))

		studentGroup.GET("/:id", handler.GetStudentByIdHandler(*queries, log))

		studentGroup.PATCH("/:id", handler.UpdateStudentByIdHandler(*queries, log))
	}

	// Group for working with calendar-related routes
	calendarGroup := r.Group("/api/calendar")
	{

		calendarGroup.GET("/", handler.DayListHandler(*queries, log))
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(cfg.HTTPServer.Address)
	//TODO: init controllers

	//TODO: init Services
}
