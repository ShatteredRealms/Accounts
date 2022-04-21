package model

import (
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	MinPasswordLength  = 6
	MaxPasswordLength  = 64
	MaxFirstNameLength = 50
	MaxLastNameLength  = MaxFirstNameLength
	MinUsernameLength  = 3
	MaxUsernameLength  = 25
)

// User Database model for a User
type User struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Username  string `gorm:"not null" json:"username"`
	Email     string `gorm:"not null;unique" json:"email"`
	Password  string `gorm:"not null" json:"password"`
}

// Validate Checks if all user data fields are valid.
func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("cannot create a user without an email")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("email is not valid")
	}

	if err := u.validateFirstName(); err != nil {
		return err
	}

	if err := u.validateLastName(); err != nil {
		return err
	}

	if err := u.validatePassword(); err != nil {
		return err
	}

	if err := u.validateUsername(); err != nil {
		return err
	}

	return nil
}

// Login Checks if the given password belongs to the user
func (u *User) Login(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			err = fmt.Errorf("invalid password")
		}
		return err
	}

	return nil
}

func (u *User) Exists() bool {
	return u.ID != 0
}

// UpdateInfo Updates the info of the user if the fields are present and valid. If any field is present but not valid
// then the user will not be updated and and error is returned. If there are no errors, then the non-empty fields for
// the FirstName, LastName, Email, and Username will be updated.
func (u *User) UpdateInfo(newUser User) error {
	updateFirstName := false
	updateLastName := false
	updateEmail := false
	updateUsername := false

	if newUser.FirstName != "" {
		if err := newUser.validateFirstName(); err != nil {
			return err
		}

		updateFirstName = true
	}

	if newUser.LastName != "" {
		if err := newUser.validateLastName(); err != nil {
			return err
		}

		updateLastName = true
	}

	if newUser.Username != "" {
		if err := newUser.validateUsername(); err != nil {
			return err
		}

		updateUsername = true
	}

	if newUser.Email != "" {
		if _, err := mail.ParseAddress(newUser.Email); err != nil {
			return err
		}

		updateEmail = true
	}

	if updateFirstName {
		u.FirstName = newUser.FirstName
	}
	if updateLastName {
		u.LastName = newUser.LastName
	}
	if updateUsername {
		u.Username = newUser.Username
	}
	if updateEmail {
		u.Email = newUser.Email
	}

	return nil
}

func (u *User) validateFirstName() error {
	if u.FirstName == "" {
		return fmt.Errorf("first name cannot be empty")
	}

	if len(u.FirstName) > MaxFirstNameLength {
		return fmt.Errorf("first name cannot be longer than 50 characters")
	}

	return nil
}

func (u *User) validateLastName() error {
	if u.LastName == "" {
		return fmt.Errorf("last name cannot be empty")
	}

	if len(u.LastName) > MaxLastNameLength {
		return fmt.Errorf("last name cannot be longer than 50 characters")
	}

	return nil
}

func (u *User) validatePassword() error {
	if u.Password == "" {
		return fmt.Errorf("cannot create a user without a password")
	}

	if len(u.Password) < MinPasswordLength {
		return fmt.Errorf("password less than minimum password length of %d", MinPasswordLength)
	}

	if len(u.Password) > MaxPasswordLength {
		return fmt.Errorf("password exeeded maximum password length of %d", MaxPasswordLength)
	}

	return nil
}

func (u *User) validateUsername() error {
	if u.Username == "" {
		return fmt.Errorf("cannot create a user without a username")
	}

	if len(u.Username) < MinUsernameLength {
		return fmt.Errorf("username less than minimum password length of %d", MinUsernameLength)
	}

	if len(u.Username) > MaxUsernameLength {
		return fmt.Errorf("username exeeded maximum password length of %d", MaxUsernameLength)
	}

	return nil
}
