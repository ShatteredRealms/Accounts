package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/controller/v1/users"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/test/factory"
	"github.com/ShatteredRealms/Accounts/test/mocks"
	"github.com/ShatteredRealms/GoUtils/pkg/helpers"
	utilModel "github.com/ShatteredRealms/GoUtils/pkg/model"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var _ = Describe("Users controller", func() {
	var userService mocks.UserService
	var usersController users.UsersController
	var logger = log.NewLogger(log.Info, "")
	var w *httptest.ResponseRecorder
	var c *gin.Context
	f := factory.NewFactory()

	BeforeEach(func() {
		userService = mocks.UserService{
			CreateReturn:      nil,
			SaveReturn:        nil,
			FindByEmailReturn: model.User{},
			FindByIdReturn:    model.User{},
			FindAllReturn: []model.User{
				f.UserFactory().User(),
				f.UserFactory().User(),
				f.UserFactory().User(),
			},
		}

		usersController = users.NewUserController(userService, logger)
		w, c, _ = helpers.SetupTestEnvironment("POST")
	})

	Context("Requesting all users", func() {
		It("should return all users", func() {
			usersController.ListAll(c)

			var resp utilModel.ResponseModel
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(resp.Data).NotTo(BeNil())
		})
	})

	Context("Requesting a user", func() {
		It("should return nil if no user given", func() {
			usersController.GetUser(c)

			var resp utilModel.ResponseModel
			data := resp.Data
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(data).To(BeNil())
		})

		It("should return the use user if found", func() {
			c.Params = gin.Params{
				gin.Param{
					Key:   "user",
					Value: strconv.FormatUint(uint64(userService.FindAllReturn[0].ID), 10),
				},
			}
			user := f.UserFactory().User()
			userService.FindByIdReturn = user
			usersController = users.NewUserController(userService, logger)
			usersController.GetUser(c)

			var resp utilModel.ResponseModel
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(resp.Data).NotTo(BeNil())
		})
	})

	Context("Editing a user", func() {
		It("should fail on invalid request body", func() {
			usersController.EditUser(c)

			var resp utilModel.ResponseModel
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(resp.Errors).NotTo(BeNil())
		})

		It("should fail on invalid user param", func() {
			body := f.UserFactory().User()
			var buf bytes.Buffer
			Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
			c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)

			usersController = users.NewUserController(userService, logger)
			usersController.EditUser(c)

			var resp utilModel.ResponseModel
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(resp.Errors).NotTo(BeNil())
		})

		It("should fail if update information has issues", func() {
			body := f.UserFactory().User()
			body.Email = "asdf"
			var buf bytes.Buffer
			Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
			c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)

			user := f.UserFactory().User()
			c.Params = gin.Params{
				gin.Param{
					Key:   "user",
					Value: strconv.FormatUint(uint64(user.ID), 10),
				},
			}
			userService.FindByIdReturn = user
			usersController = users.NewUserController(userService, logger)
			usersController.EditUser(c)

			var resp utilModel.ResponseModel
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(resp.Errors).NotTo(BeNil())
		})

		It("should fail if the user in param is not found", func() {
			body := f.UserFactory().User()
			var buf bytes.Buffer
			Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
			c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)

			user := f.UserFactory().User()
			c.Params = gin.Params{
				gin.Param{
					Key:   "user",
					Value: strconv.FormatUint(uint64(user.ID), 10),
				},
			}
			usersController = users.NewUserController(userService, logger)
			usersController.EditUser(c)

			var resp utilModel.ResponseModel
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(resp.Errors).NotTo(BeNil())
		})

		It("should fail if there were database issues", func() {
			body := f.UserFactory().User()
			var buf bytes.Buffer
			Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
			c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)

			user := f.UserFactory().User()
			c.Params = gin.Params{
				gin.Param{
					Key:   "user",
					Value: strconv.FormatUint(uint64(user.ID), 10),
				},
			}
			userService.FindByIdReturn = user
			userService.SaveReturn = fmt.Errorf("error")
			usersController = users.NewUserController(userService, logger)
			usersController.EditUser(c)

			var resp utilModel.ResponseModel
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(resp.Errors).NotTo(BeNil())
		})

		It("should succeed for a valid request", func() {
			body := f.UserFactory().User()
			var buf bytes.Buffer
			Expect(json.NewEncoder(&buf).Encode(body)).ShouldNot(HaveOccurred())
			c.Request, _ = http.NewRequest(http.MethodPost, "/", &buf)

			user := f.UserFactory().User()
			c.Params = gin.Params{
				gin.Param{
					Key:   "user",
					Value: strconv.FormatUint(uint64(user.ID), 10),
				},
			}
			userService.FindByIdReturn = user
			usersController = users.NewUserController(userService, logger)
			usersController.EditUser(c)

			var resp utilModel.ResponseModel
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(resp.Errors).To(BeNil())
		})
	})
})
