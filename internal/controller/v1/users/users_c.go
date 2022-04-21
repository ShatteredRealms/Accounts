package users

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/controller/v1/ctrlutil"
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
	EditUser(c *gin.Context)
}

type usersController struct {
	userService service.UserService
	logger      log.LoggerService
}

func NewUserController(u service.UserService, logger log.LoggerService) UsersController {
	return usersController{
		userService: u,
		logger:      logger,
	}
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
	user, err := u.getUserFromParam(c, "user")
	if err != nil {
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

func (u usersController) EditUser(c *gin.Context) {
	userInfo := model.User{}
	err := ctrlutil.ParseBody(c, &userInfo)
	if err != nil {
		return
	}

	user, err := u.getUserFromParam(c, "user")
	if err != nil {
		return
	}

	err = user.UpdateInfo(userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewBadRequestResponse(c, err.Error()))
		return
	}

	user, err = u.userService.Save(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewInternalServerResponse(c, "unable to update database"))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(c, "Success", nil))
}

func (u usersController) getUserFromParam(c *gin.Context, param string) (model.User, error) {
	id64, err := strconv.ParseUint(c.Param(param), 10, 32)
	id := uint(id64)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewBadRequestResponse(c, "Invalid user ID"))
		return model.User{}, err
	}

	user := u.userService.FindById(id)

	if !user.Exists() {
		c.JSON(http.StatusNotFound, model.NewGenericNotFoundResponse(c))
		return model.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}
