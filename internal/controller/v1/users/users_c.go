package users

import (
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/model"
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
	users := u.userService.FindAll()
	parsedUsers := make([]model.StrippedUserModel, len(users))
	for i, u := range users {
		parsedUsers[i] = model.StrippedUserModel{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
			ID:        u.ID,
		}
	}

	c.JSON(200, parsedUsers)
}

func NewUserController(u service.UserService, logger log.LoggerService) UsersController {
	return usersController{
		userService: u,
		logger:      logger,
	}
}
