package v1

import (
	"github.com/ShatteredRealms/Accounts/internal/controller/v1/auth"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	utilService "github.com/ShatteredRealms/GoUtils/pkg/service"
	"github.com/gin-gonic/gin"
)

// SetAuthRoutes initializes all auth routes
func SetAuthRoutes(rg *gin.RouterGroup, s service.UserService, jwt utilService.JWTService, logger log.LoggerService) {
	authController := auth.NewAuthController(s, jwt, logger)
	rg.POST("/login", authController.Login)
	rg.POST("/register", authController.Register)
}
