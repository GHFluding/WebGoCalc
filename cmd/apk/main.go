package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"test/internal/config"
	"test/internal/database/postgres"
	handler "test/internal/server/http/handlers"
	"test/internal/server/http/middleware"

	// setup logger
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
// @contact.url   --
// @contact.email  --

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

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

	//health check
	r.GET("/health", func(c *gin.Context) {
		// Здесь можно добавить проверку подключения к БД, очередям и т. д.
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	r.Run(cfg.HTTPServer.Address)
	//TODO: init controllers

	//TODO: init Services
}
