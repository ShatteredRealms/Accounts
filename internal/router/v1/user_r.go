package v1

import (
	"github.com/ShatteredRealms/Accounts/internal/controller/v1/users"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/gin-gonic/gin"
)

// SetUsersRoutes initializes all users routes
func SetUsersRoutes(rg *gin.RouterGroup, s service.UserService, logger log.LoggerService) {
	usersController := users.NewUserController(s, logger)
	rg.GET("/users", usersController.ListAll)
	rg.GET("/users/:user", usersController.GetUser)
	rg.PUT("/users/:user", usersController.EditUser)
}
