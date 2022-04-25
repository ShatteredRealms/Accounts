package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleWare is the middleware function that sets the necessary headers to setup CORS configuration. If a OPTION
// request is made, then the handle aborts and returns with HTTP status No Content (204)
func CORSMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Request.Header.Get("Origin")) == 0 {
			c.Next()
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, UPDATE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	}
}
