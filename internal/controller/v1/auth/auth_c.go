package auth

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/controller/v1/ctrlutil"
	"github.com/ShatteredRealms/Accounts/internal/log"
	accountModel "github.com/ShatteredRealms/Accounts/pkg/model"
	accountService "github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/ShatteredRealms/GoUtils/pkg/model"
	"github.com/ShatteredRealms/GoUtils/pkg/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authController struct {
	userService accountService.UserService
	jwtService  service.JWTService
	logger      log.LoggerService
}

func NewAuthController(u accountService.UserService, jwt service.JWTService, logger log.LoggerService) AuthController {
	return authController{
		userService: u,
		jwtService:  jwt,
		logger:      logger,
	}
}

// Login godoc
// @Summary Handles when a request to login is being made.
// @Schemes object
// @Description Checks if the given login information is correct, and responds with the
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseModel JWT token for user in data
// @Failure 401 {object} model.ResponseModel Explanation of why
// @Router /api/v1/login [post]
func (a authController) Login(c *gin.Context) {
	login := accountModel.LoginRequest{}
	err := ctrlutil.ParseBody(c, &login)
	if err != nil {
		return
	}

	if !ctrlutil.ValidatePresent(c, "email", login.Email) {
		return
	}

	if !ctrlutil.ValidatePresent(c, "password", login.Password) {
		return
	}

	user := a.userService.FindByEmail(login.Email)
	if !user.Exists() {
		resp := model.NewFailedLoginResponse(c, fmt.Errorf("email does not exist"))
		c.JSON(resp.StatusCode, resp)
		return
	}

	err = user.Login(login.Password)
	if err != nil {
		resp := model.NewFailedLoginResponse(c, err)
		c.JSON(resp.StatusCode, resp)
		return
	}

	t, err := a.tokenForUser(&user)
	if err != nil {
		resp := model.NewInternalServerResponse(c, "unable to create auth token")
		c.JSON(resp.StatusCode, resp)
		return
	}
	data := accountModel.LoginResponse{
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		Token:     t,
	}

	a.logger.LogLoginRequest()
	resp := model.NewSuccessResponse(c, "Success", data)
	c.JSON(resp.StatusCode, resp)
}

// Register godoc
// @Summary Handles when a request to register is being made.
// @Schemes object
// @Description Checks if the given login registration is valid and has no conflicts with existing accounts and creates
// the account if true.
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseModel JWT token for user in data
// @Failure 401 {object} model.ResponseModel Explanation of why
// @Router /api/v1/register [post]
func (a authController) Register(c *gin.Context) {
	register := accountModel.RegisterRequest{}
	err := ctrlutil.ParseBody(c, &register)
	if err != nil {
		return
	}

	user := accountModel.User{
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Email:     register.Email,
		Password:  register.Password,
		Username:  register.Username,
	}

	user, err = a.userService.Create(user)
	if err != nil {
		var errors []model.ErrorModel
		e := model.RegistrationFailedError
		e.Info = err.Error()
		errors = append(errors, e)

		resp := model.NewSuccessResponse(c, "Fail", nil)
		resp.Errors = errors

		c.JSON(resp.StatusCode, resp)
		return
	}

	t, err := a.tokenForUser(&user)
	if err != nil {
		resp := model.NewInternalServerResponse(c, "unable to create auth token")
		c.JSON(resp.StatusCode, resp)
		return
	}

	a.logger.LogRegisterRequest()
	data := accountModel.LoginResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     t,
	}

	resp := model.NewSuccessResponse(c, "Success", data)
	c.JSON(resp.StatusCode, resp)
}

func (a *authController) tokenForUser(u *accountModel.User) (t string, err error) {
	claims := jwt.MapClaims{
		"sub":         u.ID,
		"given_name":  u.FirstName,
		"family_name": u.LastName,
		"email":       u.Email,
	}
	t, err = a.jwtService.Create(time.Hour, "shatteredrealmsonline.com/accounts/v1", claims)

	return t, err
}
