package users_test

import (
	"encoding/json"
	"github.com/ShatteredRealms/Accounts/internal/controller/v1/users"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/helpers"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/test/mocks"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http/httptest"
)

var _ = Describe("Users controller", func() {
	var userService mocks.UserService
	var usersController users.UsersController
	var logger = log.NewLogger(log.Info, "")
	var w *httptest.ResponseRecorder
	var c *gin.Context

	BeforeEach(func() {
		userService = mocks.UserService{
			CreateReturn:      nil,
			SaveReturn:        nil,
			FindByEmailReturn: model.User{},
			FindByIdReturn:    model.User{},
		}

		usersController = users.NewUserController(userService, logger)
		w, c, _ = helpers.SetupTestEnvironment("POST")
	})

	Context("Requesting all users", func() {
		It("should return all users", func() {
			usersController.ListAll(c)
			var resp gin.H
			Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
			Expect(resp).To(Equal(gin.H{"message": "service pending"}))
		})
	})
})
