package auth

import (
	"encoding/json"
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authController struct {
	userService service.UserService
	jwtService  service.JWTService
	logger      log.LoggerService
}

func NewAuthController(u service.UserService, jwt service.JWTService, logger log.LoggerService) AuthController {
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
	reqBody := c.Request.Body
	if reqBody == nil {
		resp := model.NewBadRequestResponse(c, "Payload missing")
		c.JSON(resp.StatusCode, resp)
		return
	}

	body, err := ioutil.ReadAll(reqBody)
	if err != nil {
		resp := model.NewInternalServerResponse(c, "unable to process payload")
		c.JSON(resp.StatusCode, resp)
		return
	}

	login := model.LoginRequest{}
	err = json.Unmarshal(body, &login)
	if err != nil {
		resp := model.NewBadRequestResponse(c, "Expected JSON body")
		c.JSON(resp.StatusCode, resp)
		return
	}

	if login.Email == "" {
		resp := model.NewBadRequestResponse(c, "Missing email")
		c.JSON(resp.StatusCode, resp)
		return
	}

	if login.Password == "" {
		resp := model.NewBadRequestResponse(c, "Missing password")
		c.JSON(resp.StatusCode, resp)
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

	data := model.LoginResponse{
		Token: t,
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
	reqBody := c.Request.Body
	if reqBody == nil {
		resp := model.NewBadRequestResponse(c, "Payload missing")
		c.JSON(resp.StatusCode, resp)
		return
	}

	body, err := ioutil.ReadAll(reqBody)
	if err != nil {
		resp := model.NewInternalServerResponse(c, "unable to process payload")
		c.JSON(resp.StatusCode, resp)
		return
	}

	register := model.RegisterRequest{}
	err = json.Unmarshal(body, &register)
	if err != nil {
		resp := model.NewBadRequestResponse(c, "Expected JSON body")
		c.JSON(resp.StatusCode, resp)
		return
	}

	user := model.User{
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Email:     register.Email,
		Password:  register.Password,
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
	data := model.LoginResponse{
		Token: t,
	}
	resp := model.NewSuccessResponse(c, "Success", data)
	c.JSON(resp.StatusCode, resp)
}

func (a *authController) tokenForUser(u *model.User) (t string, err error) {
	claims := jwt.MapClaims{
		"sub":         u.ID,
		"given_name":  u.FirstName,
		"family_name": u.LastName,
		"email":       u.Email,
	}
	t, err = a.jwtService.Create(time.Hour, claims)

	return t, err
}
