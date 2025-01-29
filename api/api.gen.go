// Package api provides primitives to interact with the openapi HTTP API.
// Modified for gin-gonic, before echo/v4 
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
)

// NocsqlcpgCreateStudentSwagger defines model for nocsqlcpg.CreateStudentSwagger.
type NocsqlcpgCreateStudentSwagger struct {
	Name      *string             `json:"name,omitempty"`
	OrderCost *int                `json:"orderCost,omitempty"`
	OrderDay  *int                `json:"orderDay,omitempty"`
	OrderTime *NocsqlcpgOrderTime `json:"orderTime,omitempty"`
	School    *string             `json:"school,omitempty"`
	Sclass    *string             `json:"sclass,omitempty"`
}

// NocsqlcpgEventDate defines model for nocsqlcpg.EventDate.
type NocsqlcpgEventDate struct {
	InfinityModifier *int    `json:"infinityModifier,omitempty"`
	Time             *string `json:"time,omitempty"`
	Valid            *bool   `json:"valid,omitempty"`
}

// NocsqlcpgOrderCheck defines model for nocsqlcpg.OrderCheck.
type NocsqlcpgOrderCheck struct {
	Bool  *bool `json:"bool,omitempty"`
	Valid *bool `json:"valid,omitempty"`
}

// NocsqlcpgOrderTime defines model for nocsqlcpg.OrderTime.
type NocsqlcpgOrderTime struct {
	Microseconds *int  `json:"microseconds,omitempty"`
	Valid        *bool `json:"valid,omitempty"`
}

// NocsqlcpgStudentEventSwagger defines model for nocsqlcpg.StudentEventSwagger.
type NocsqlcpgStudentEventSwagger struct {
	EventDate  *NocsqlcpgEventDate  `json:"eventDate,omitempty"`
	Id         *int                 `json:"id,omitempty"`
	OrderCheck *NocsqlcpgOrderCheck `json:"orderCheck,omitempty"`
	OrderCost  *int                 `json:"orderCost,omitempty"`
	OrderTime  *NocsqlcpgOrderTime  `json:"orderTime,omitempty"`
	StudentID  *int                 `json:"studentID,omitempty"`
}

// NocsqlcpgStudentSwagger defines model for nocsqlcpg.StudentSwagger.
type NocsqlcpgStudentSwagger struct {
	Id        *int                `json:"id,omitempty"`
	Name      *string             `json:"name,omitempty"`
	OrderCost *int                `json:"orderCost,omitempty"`
	OrderDay  *int                `json:"orderDay,omitempty"`
	OrderTime *NocsqlcpgOrderTime `json:"orderTime,omitempty"`
	School    *string             `json:"school,omitempty"`
	Sclass    *string             `json:"sclass,omitempty"`
}

// GetApiCalendarParams defines parameters for GetApiCalendar.
type GetApiCalendarParams struct {
	// Date Дата в формате YYYY-MM-DD
	Date string `form:"date" json:"date"`
}

// PostApiStudentsJSONRequestBody defines body for PostApiStudents for application/json ContentType.
type PostApiStudentsJSONRequestBody = NocsqlcpgCreateStudentSwagger

// ServerInterface represents all server handlers.
type ServerInterface interface {
    GetApiCalendar(*gin.Context, GetApiCalendarParams) error
    GetApiStudents(*gin.Context) error
    PostApiStudents(*gin.Context) error
}

// GinServerInterfaceWrapper обертка для Gin
type GinServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ConvertHandler преобразует обработчики для Gin
func ConvertHandler(handler func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			c.Error(err) // Обработка ошибок через Gin
		}
	}
}

// Реализация методов для Gin

// GetApiCalendar Gin-обработчик для GET /api/calendar
func (w *GinServerInterfaceWrapper) GetApiCalendar(c *gin.Context) {
	var params GetApiCalendarParams

	// Извлекаем параметр date из query
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'date' is required"})
		return
	}
	params.Date = date

	// Вызываем оригинальный обработчик
	if err := w.Handler.GetApiCalendar(c, params); err != nil {
		c.Error(err)
	}
}

// GetApiStudents Gin-обработчик для GET /api/students
func (w *GinServerInterfaceWrapper) GetApiStudents(c *gin.Context) {
	if err := w.Handler.GetApiStudents(c); err != nil {
		c.Error(err)
	}
}

// PostApiStudents Gin-обработчик для POST /api/students
func (w *GinServerInterfaceWrapper) PostApiStudents(c *gin.Context) {
	var requestBody NocsqlcpgCreateStudentSwagger
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Вызываем оригинальный обработчик
	if err := w.Handler.PostApiStudents(c); err != nil {
		c.Error(err)
	}
}

// RegisterGinHandlers регистрация маршрутов в Gin
func RegisterGinHandlers(router *gin.Engine, si ServerInterface) {
	wrapper := &GinServerInterfaceWrapper{Handler: si}

	router.GET("/api/calendar", wrapper.GetApiCalendar)
	router.GET("/api/students", wrapper.GetApiStudents)
	router.POST("/api/students", wrapper.PostApiStudents)
}