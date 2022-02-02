package model_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"

	"github.com/productivestudy/auth/pkg/model"
	"github.com/productivestudy/auth/tests/helpers"
)

var _ = Describe("User", func() {
	var user model.User

	BeforeEach(func() {
		user = model.User{
			FirstName: helpers.RandString(10),
			LastName:  helpers.RandString(10),
			Email:     helpers.RandString(10) + "@test.com",
			Password:  helpers.RandString(10),
		}
	})

	Context("Login", func() {
		var expectedError error
		var password string
		var passwordBytes []byte

		BeforeEach(func() {
			expectedError = nil
			password = user.Password
			passwordBytes, _ = bcrypt.GenerateFromPassword([]byte(password), 0)
			user.Password = string(passwordBytes)
		})

		It("should work if the password is correct", func() {
		})

		It("should fail if the db password isn't encrypted", func() {
			user.Password = password
			expectedError = fmt.Errorf("crypto/bcrypt: hashedSecret too short to be a bcrypted password")
		})

		It("should fail if the password does not match", func() {
			password = password + "a"
			expectedError = fmt.Errorf("invalid password")
		})

		AfterEach(func() {
			user.Validate()
			if expectedError == nil {
				Expect(user.Login(password)).To(BeNil())
			} else {
				Expect(user.Login(password)).To(Equal(expectedError))
			}
		})
	})

	Context("Validation", func() {
		var expectedError error

		BeforeEach(func() {
			expectedError = nil
		})

		It("should require an email", func() {
			user.Email = ""
			expectedError = fmt.Errorf("cannot create a user without an email")
		})

		It("should require a valid email", func() {
			user.Email = helpers.RandString(10)
			expectedError = fmt.Errorf("email is not valid")
		})

		It("should require a first name", func() {
			user.FirstName = ""
			expectedError = fmt.Errorf("cannot create a user without a first name")
		})

		It("should require a last name", func() {
			user.LastName = ""
			expectedError = fmt.Errorf("cannot create a user without a last name")
		})

		It("should require a password", func() {
			user.Password = ""
			expectedError = fmt.Errorf("cannot create a user without a password")
		})

		It(fmt.Sprintf("should require a password with minimum length of %d", model.MinPasswordLength), func() {
			user.Password = helpers.RandString(model.MinPasswordLength - 1)
			expectedError = fmt.Errorf("less than minimum password length of %d", model.MinPasswordLength)
		})

		It(fmt.Sprintf("should require a password with maximum length of %d", model.MaxPasswordLength), func() {
			user.Password = helpers.RandString(model.MaxPasswordLength + 1)
			expectedError = fmt.Errorf("exeeded maximum password length of %d", model.MaxPasswordLength)
		})

		It(fmt.Sprintf("should allow a password of length of %d", model.MaxPasswordLength), func() {
			user.Password = helpers.RandString(model.MaxPasswordLength)
		})

		It(fmt.Sprintf("should allow a password of length of %d", model.MinPasswordLength), func() {
			user.Password = helpers.RandString(model.MinPasswordLength)
		})

		AfterEach(func() {
			if expectedError == nil {
				Expect(user.Validate()).To(BeNil())
			} else {
				Expect(user.Validate()).To(Equal(expectedError))
			}
		})
	})

	It("should know if it exists", func() {
		Expect(user.Exists()).To(BeFalse())
		user.ID = 1
		Expect(user.Exists()).To(BeTrue())
	})
})
