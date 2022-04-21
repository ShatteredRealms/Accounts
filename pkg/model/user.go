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
)

// User Database model for a User
type User struct {
	gorm.Model
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Username  string `gorm:"not null"`
	Email     string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
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
