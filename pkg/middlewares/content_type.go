package middlewares

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContentTypeMiddleWare ensures content-type is application/json for all non-GET requests. If it is not, the request is
// aborted and an HTTP status Unsupported Media Type (415) is returned with more JSON information regarding the error.
func ContentTypeMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if (c.Request.Method == http.MethodPatch ||
			c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut) &&
			c.Request.Header.Get("Content-Type") != "application/json" {

			c.JSON(http.StatusUnsupportedMediaType, model.NewGenericUnsupportedMediaResponse(c))
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
