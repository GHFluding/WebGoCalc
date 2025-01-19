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

	_ "test/docs"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	// Group for working with student-related routes
	studentGroup := r.Group("/api/students")
	{
		// ListStudentsHandler возвращает список студентов.
		// @Summary      Получить список студентов
		// @Description  Возвращает список всех студентов.
		// @Tags         students
		// @Accept       json
		// @Produce      json
		// @Success      200  {array}  postgres.Student
		// @Failure      500  {object}  gin.H
		// @Router       /api/students [get]

		studentGroup.GET("/", handler.ListStudentsHandler(*queries, log))

		// CreateStudentHandler создает нового студента.
		// @Summary      Создать студента
		// @Description  Добавляет нового студента в базу данных.
		// @Tags         students
		// @Accept       json
		// @Produce      json
		// @Param        student  body  postgres.CreateStudentParams  true  "Данные студента"
		// @Success      201  {object}  postgres.Student
		// @Failure      400  {object}  gin.H
		// @Failure      500  {object}  gin.H
		// @Router       /api/students [post]

		studentGroup.POST("/", handler.CreateStudentHandler(*queries, log))
		// DeleteStudentByIdHandler удаляет студента по ID.
		// @Summary      Удалить студента
		// @Description  Удаляет студента из базы данных по его ID.
		// @Tags         students
		// @Accept       json
		// @Produce      json
		// @Param        id   path      int  true  "ID студента"
		// @Success      204
		// @Failure      400  {object}  gin.H
		// @Failure      404  {object}  gin.H
		// @Router       /api/students/{id} [delete]

		studentGroup.DELETE("/:id", handler.DeleteStudentByIdHandler(*queries, log))

		// GetStudentByIdHandler возвращает информацию о студенте по ID.
		// @Summary      Получить информацию о студенте
		// @Description  Возвращает полную информацию о студенте по его ID.
		// @Tags         students
		// @Accept       json
		// @Produce      json
		// @Param        id   path      int  true  "ID студента"
		// @Success      200  {object}  postgres.Student
		// @Failure      400  {object}  gin.H
		// @Failure      404  {object}  gin.H
		// @Router       /api/students/{id} [get]

		studentGroup.GET("/:id", handler.GetStudentByIdHandler(*queries, log))
		// UpdateStudentByIdHandler обновляет информацию о студенте по ID.
		// @Summary      Обновить информацию о студенте
		// @Description  Позволяет обновить определенные данные студента по его ID.
		// @Tags         students
		// @Accept       json
		// @Produce      json
		// @Param        id      path      int                          true  "ID студента"
		// @Param        student body      postgres.UpdateStudentParams true  "Данные для обновления"
		// @Success      200     {object}  postgres.Student
		// @Failure      400     {object}  gin.H
		// @Failure      404     {object}  gin.H
		// @Router       /api/students/{id} [patch]

		studentGroup.PATCH("/:id", handler.UpdateStudentByIdHandler(*queries, log))
	}

	// Group for working with calendar-related routes
	calendarGroup := r.Group("/api/calendar")
	{
		// DayListHandler возвращает список событий на день.
		// @Summary      Получить список событий
		// @Description  Возвращает события календаря на указанный день.
		// @Tags         calendar
		// @Accept       json
		// @Produce      json
		// @Param        date  query  string  true  "Дата в формате YYYY-MM-DD"
		// @Success      200  {array}  postgres.Event
		// @Failure      400  {object}  gin.H
		// @Failure      500  {object}  gin.H
		// @Router       /api/calendar [get]

		calendarGroup.GET("/", handler.DayListHandler(*queries, log))
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(cfg.HTTPServer.Address)
	//TODO: init controllers

	//TODO: init Services
}
