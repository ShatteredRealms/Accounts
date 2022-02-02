package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/productivestudy/auth/internal/requests"
	"github.com/productivestudy/auth/pkg/model"
)

// ContentTypeMiddleWare ensures content-type is application/json for all non-GET requests. If it is not, the request is
// aborted and a HTTP status Unsupported Media Type (415) is returned with more JSON information regarding the error.
func ContentTypeMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPatch:
		case http.MethodPost:
		case http.MethodPut:
			if c.Request.Header.Get(requests.ContentType) != requests.JSONContent {
				c.JSON(http.StatusUnsupportedMediaType, model.NewGenericUnsupportedMediaResponse(c))
				c.Abort()
			}
			return
		default:
			c.Next()
		}
	}
}
