package model_test

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/productivestudy/auth/pkg/model"
	"github.com/productivestudy/auth/tests/helpers"
)

var _ = Describe("ResponseModel", func() {
	var message string
	var response model.ResponseModel
	var ctx *gin.Context

	BeforeEach(func() {
		ctx = &gin.Context{
			Request: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   helpers.RandString(10),
				},
				Method: helpers.RandString(5),
			},
		}
	})

	It("generates success responses correctly", func() {
		message = "message"
		response = model.NewSuccessResponse(ctx, message, nil)
		Expect(response.Data).To(BeNil())
		Expect(response.Errors).To(BeNil())
		Expect(response.Message).To(Equal(message))
		Expect(response.StatusCode).To(Equal(http.StatusOK))
	})

	Context("error responses", func() {
		var expectedError model.ErrorModel

		Context("generic", func() {
			It("generates generic unsupported media type correctly", func() {
				response = model.NewGenericUnsupportedMediaResponse(ctx)
				expectedError = model.UnsupportedMediaError
				Expect(response.StatusCode).To(Equal(http.StatusUnsupportedMediaType))
			})

			It("generates generic not found correctly", func() {
				response = model.NewGenericNotFoundResponse(ctx)
				expectedError = model.NotFoundError
				Expect(response.StatusCode).To(Equal(http.StatusNotFound))

			})

			AfterEach(func() {
				Expect(len(response.Errors)).To(Equal(1))
				Expect(response.Errors[0]).To(Equal(expectedError))
			})
		})

		Context("generic with custom info", func() {
			var info string
			BeforeEach(func() {
				info = helpers.RandString(10)
			})

			It("generates internal server error correctly", func() {
				response = model.NewInternalServerResponse(ctx, info)
				expectedError = model.InternalServError
				Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
			})

			It("generates bad request error correctly", func() {
				response = model.NewBadRequestResponse(ctx, info)
				expectedError = model.BadRequestError
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("generates failed login error correctly", func() {
				response = model.NewFailedLoginResponse(ctx, fmt.Errorf(info))
				expectedError = model.UnauthorizedError
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})

			AfterEach(func() {
				Expect(len(response.Errors)).To(Equal(1))
				Expect(response.Errors[0].ErrorCode).To(Equal(expectedError.ErrorCode))
				Expect(response.Errors[0].Hints).To(Equal(expectedError.Hints))
				Expect(response.Errors[0].Text).To(Equal(expectedError.Text))
				Expect(response.Errors[0].Info).To(Equal(info))
			})
		})

		AfterEach(func() {
			Expect(response.Data).To(BeNil())
			Expect(response.Message).To(Equal("Fail"))
		})
	})

	AfterEach(func() {
		Expect(response.Time).NotTo(BeNil())
		Expect(response.Method).To(Equal(ctx.Request.Method))
		Expect(response.Endpoint).To(Equal(ctx.Request.URL.String()))
	})
})
