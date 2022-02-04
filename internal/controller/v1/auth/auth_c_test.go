package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/controller/v1/auth"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/helpers"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var _ = Describe("Auth controller ", func() {
	var userService testUserService
	var authController auth.AuthController
	var w *httptest.ResponseRecorder
	var c *gin.Context
	l := log.NewLogger(log.Error, "")

	BeforeEach(func() {
		userService = testUserService{
			CreateReturn:      nil,
			SaveReturn:        nil,
			FindByEmailReturn: model.User{},
			FindByIdReturn:    model.User{},
		}

		authController = auth.NewAuthController(userService, testJWT(true), l)

		w, c, _ = helpers.SetupTestEnvironment("POST")
	})

	Context("Making a login request", func() {
		var body model.LoginRequest

		BeforeEach(func() {
			body = model.LoginRequest{
				Email:    helpers.RandString(10),
				Password: helpers.RandString(10),
			}
		})

		Context("with an invalid body", func() {
			It("should error when the body is empty", func() {
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("Payload missing"))
			})

			It("should error with an invalid body", func() {
				c.Request, _ = http.NewRequest(http.MethodPost, "/", errReader(1))
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("unable to process payload"))
			})

			It("should fail with an empty body", func() {
				c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte{}))
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("Expected JSON body"))
			})
		})

		Context("with a valid body", func() {
			It("should require email field", func() {
				body.Email = ""
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("Missing email"))
			})

			It("should require password field", func() {
				body.Password = ""
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("Missing password"))
			})

			It("should fail if the email does not exist", func() {
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController = auth.NewAuthController(userService, testJWT(true), l)
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusUnauthorized))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("email does not exist"))
			})

			It("should fail if the password is saved incorrectly", func() {
				userService.FindByEmailReturn.ID = 1
				authController = auth.NewAuthController(userService, testJWT(true), l)
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusUnauthorized))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("crypto/bcrypt: hashedSecret too short to be a bcrypted password"))
			})

			It("should fail if the password is incorrect", func() {
				userService.FindByEmailReturn.ID = 1
				password, err := bcrypt.GenerateFromPassword([]byte(body.Password+"a"), 0)
				userService.FindByEmailReturn.Password = string(password)
				authController = auth.NewAuthController(userService, testJWT(true), l)

				Expect(err).ShouldNot(HaveOccurred())
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusUnauthorized))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("invalid password"))
			})

			It("should fail if the password is correct but JWT service is bad", func() {
				userService.FindByEmailReturn.ID = 1
				password, err := bcrypt.GenerateFromPassword([]byte(body.Password), 0)
				userService.FindByEmailReturn.Password = string(password)
				authController = auth.NewAuthController(userService, testJWT(false), l)

				Expect(err).ShouldNot(HaveOccurred())
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("unable to create auth token"))
			})

			It("should succeed if the password is correct", func() {
				userService.FindByEmailReturn.ID = 1
				password, err := bcrypt.GenerateFromPassword([]byte(body.Password), 0)
				userService.FindByEmailReturn.Password = string(password)
				authController = auth.NewAuthController(userService, testJWT(true), l)

				Expect(err).ShouldNot(HaveOccurred())
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController.Login(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(resp.Data.(map[string]interface{})["token"]).To(Equal("ok"))
			})
		})
	})

	Context("Making a register request", func() {
		Context("with an invalid body", func() {
			It("should error when the body is empty", func() {
				authController.Register(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("Payload missing"))
			})

			It("should error with an invalid body", func() {
				c.Request, _ = http.NewRequest(http.MethodPost, "/", errReader(1))
				authController.Register(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("unable to process payload"))
			})

			It("should fail with an empty body", func() {
				c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte{}))
				authController.Register(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("Expected JSON body"))
			})
		})

		Context("with a valid body", func() {
			var body model.User
			BeforeEach(func() {
				body = model.User{
					FirstName: helpers.RandString(10),
					LastName:  helpers.RandString(10),
					Email:     helpers.RandString(10) + "@example.com",
					Password:  helpers.RandString(10),
				}
			})

			It("should fail if create fails", func() {
				userService.CreateReturn = fmt.Errorf(helpers.RandString(10))
				authController = auth.NewAuthController(userService, testJWT(true), l)
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController.Register(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal(userService.CreateReturn.Error()))
			})

			It("should fail if jwt service fails", func() {
				authController = auth.NewAuthController(userService, testJWT(false), l)
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController.Register(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(resp.Errors[0]).ToNot(BeNil())
				Expect(resp.Errors[0].Info).To(Equal("unable to create auth token"))
			})

			It("should suceed with a valid request", func() {
				authController = auth.NewAuthController(userService, testJWT(true), l)
				var buf bytes.Buffer
				Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
				c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)
				authController.Register(c)

				resp := model.ResponseModel{}
				Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(resp.Data.(map[string]interface{})["token"]).To(Equal("ok"))
			})
		})
	})
})

type errReader int

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("test error")
}

// testJWT If true, returns no errors with string ok, otherwise returns error.
type testJWT bool

func (t testJWT) Create(time.Duration, jwt.MapClaims) (string, error) {
	if t {
		return "ok", nil
	} else {
		return "", fmt.Errorf("error")
	}
}

func (t testJWT) Validate(token string) (interface{}, error) {
	if t {
		return "ok", nil
	} else {
		return "", fmt.Errorf("error")
	}
}

type testUserService struct {
	CreateReturn      error
	SaveReturn        error
	FindByIdReturn    model.User
	FindByEmailReturn model.User
}

func (t testUserService) Create(u model.User) (model.User, error) {
	return u, t.CreateReturn
}
func (t testUserService) Save(u model.User) (model.User, error) {
	return u, t.SaveReturn
}
func (t testUserService) WithTrx(*gorm.DB) service.UserService {
	return t
}

func (t testUserService) FindById(id uint) model.User {
	return t.FindByIdReturn
}
func (t testUserService) FindByEmail(email string) model.User {
	return t.FindByEmailReturn
}
