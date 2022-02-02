package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// DateFormat The date format used for formatting the time in the Response
const DateFormat = "2022-01-01 13:13:13"

// ResponseModel Data that will be returned to the user when a verbose message needs to be sent back to the request
type ResponseModel struct {
	Time       time.Time    `json:"time"`
	StatusCode int          `json:"statuscode"`
	Message    string       `json:"message"`
	Data       interface{}  `json:"data"`
	Errors     []ErrorModel `json:"errors"`
	Endpoint   string       `json:"endpoint"`
	Method     string       `json:"method"`
}

func NewSuccessResponse(c *gin.Context, message string, data interface{}) ResponseModel {
	response := newBaseResponse(c, http.StatusOK)
	response.Message = message
	response.Data = data
	return response
}

// NewGenericUnsupportedMediaResponse Creates a generic response when an HTTP Unsupported Media (415) response is needed.
func NewGenericUnsupportedMediaResponse(c *gin.Context) ResponseModel {
	return newGenericError(c, http.StatusUnsupportedMediaType, UnsupportedMediaError)
}

func NewGenericNotFoundResponse(c *gin.Context) ResponseModel {
	return newGenericError(c, http.StatusNotFound, NotFoundError)
}

func NewInternalServerResponse(c *gin.Context, info string) ResponseModel {
	return newError(c, http.StatusInternalServerError, InternalServError, info)
}

func NewBadRequestResponse(c *gin.Context, info string) ResponseModel {
	return newError(c, http.StatusBadRequest, BadRequestError, info)
}

func NewFailedLoginResponse(c *gin.Context, err error) ResponseModel {
	return newError(c, http.StatusUnauthorized, UnauthorizedError, err.Error())
}

func newGenericError(c *gin.Context, status int, e ErrorModel) ResponseModel {
	response := newBaseResponse(c, status)

	var errors []ErrorModel
	errors = append(errors, e)
	response.Errors = errors
	response.Message = "Fail"

	return response
}

func newError(c *gin.Context, status int, e ErrorModel, info string) ResponseModel {
	e.Info = info
	return newGenericError(c, status, e)
}

func newBaseResponse(c *gin.Context, status int) ResponseModel {
	return ResponseModel{
		Time:       time.Now(),
		StatusCode: status,
		Endpoint:   c.Request.URL.String(),
		Method:     c.Request.Method,
	}
}
