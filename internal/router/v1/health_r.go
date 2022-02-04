package v1

import (
	"github.com/ShatteredRealms/Accounts/internal/controller/v1/health"
	"github.com/gin-gonic/gin"
)

// SetHealthRoutes initializes all health routes
func SetHealthRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", health.Health)
}
