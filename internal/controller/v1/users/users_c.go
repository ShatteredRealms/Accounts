package users

import (
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UsersController interface {
	ListAll(c *gin.Context)
	GetUser(c *gin.Context)
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

	c.JSON(200, model.NewSuccessResponse(c, "Success", users))
}

func (u usersController) GetUser(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("user"), 10, 32)
	id := uint(id64)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewBadRequestResponse(c, "Invalid user ID"))
	}

	user := u.userService.FindById(id)

	if (model.User{} == user) {
		c.JSON(http.StatusNotFound, model.NewGenericNotFoundResponse(c))
		return
	}

	parsedUser := model.StrippedUserModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		ID:        user.ID,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(c, "Success", parsedUser))
}

func NewUserController(u service.UserService, logger log.LoggerService) UsersController {
	return usersController{
		userService: u,
		logger:      logger,
	}
}
