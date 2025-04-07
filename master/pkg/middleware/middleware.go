package middleware

import "github.com/gin-gonic/gin"

func NewGinMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors(),
	}
}
