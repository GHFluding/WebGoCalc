package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// func make request id in gin context
func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader("X-Request-Id")
		if requestId == "" {
			requestId = c.GetHeader("RequestId") // generate uniq id
		}
		if requestId == "" {
			requestId = uuid.New().String() // generate uniq id
		}
		c.Set("RequestId", requestId) // saved in gin context
		c.Writer.Header().Set("X-Request-Id", requestId)
		c.Next()
	}
}

// func for get id from gin context
func RequestIdFromContext(c *gin.Context) string {
	if requestId, exists := c.Get("RequestId"); exists {
		return requestId.(string)
	}
	return ""
}
