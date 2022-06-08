package ctrlutil

import (
	"fmt"
	"github.com/ShatteredRealms/GoUtils/pkg/model"
	"github.com/gin-gonic/gin"
)

func ValidatePresent(c *gin.Context, fieldName string, field string) bool {
	if field == "" {
		resp := model.NewBadRequestResponse(c, fmt.Sprintf("missing %s", fieldName))
		c.JSON(resp.StatusCode, resp)
		return false
	}

	return true
}
