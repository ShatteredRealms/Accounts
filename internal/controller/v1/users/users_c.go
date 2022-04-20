package users

import (
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/gin-gonic/gin"
)

type UsersController interface {
	ListAll(c *gin.Context)
}

type usersController struct {
	userService service.UserService
	logger      log.LoggerService
}

func (u usersController) ListAll(c *gin.Context) {
	c.JSON(200, gin.H{"message": "service pending"})
}

func NewUserController(u service.UserService, logger log.LoggerService) UsersController {
	return usersController{
		userService: u,
		logger:      logger,
	}
}
