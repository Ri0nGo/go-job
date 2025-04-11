package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			strings.Join(
				[]string{
					http.MethodPost,
					http.MethodOptions,
					http.MethodDelete,
					http.MethodPut,
					http.MethodPatch,
				},
				",",
			),
		)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Session")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization")

		if c.Request.Method == http.MethodOptions { //跨域预检请求
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
