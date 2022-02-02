package middlewares

import (
	"github.com/ShatteredRealms/Web/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContentTypeMiddleWare ensures content-type is application/json for all non-GET requests. If it is not, the request is
// aborted and an HTTP status Unsupported Media Type (415) is returned with more JSON information regarding the error.
func ContentTypeMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPatch:
		case http.MethodPost:
		case http.MethodPut:
			if c.Request.Header.Get("Content-Type") != "application/json" {
				c.JSON(http.StatusUnsupportedMediaType, model.NewGenericUnsupportedMediaResponse(c))
				c.Abort()
				return
			}
		default:
			c.Next()
		}
	}
}
