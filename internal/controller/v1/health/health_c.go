package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response The data to respond with when a Health request is made
type Response struct {
	// Health The health status message to response with
	Health string `json:"health"`
}

// Health godoc
// @Summary pings health
// @Schemes
// @Description Responds with health ok JSON response if the server is operational
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} Response Responds with Health of "ok"
// @Router /api/v1/health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Health: "ok"})
}
