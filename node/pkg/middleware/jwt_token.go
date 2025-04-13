package middleware

import (
	"github.com/gin-gonic/gin"
	"go-job/node/pkg/auth"
)

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if err := auth.RefreshToken(); err == nil {
			authHeader := c.Writer.Header().Get("Authorization")
			if authHeader == "" {
				c.Header("Authorization", auth.GetJwtToken())
			}
		}
	}
}
