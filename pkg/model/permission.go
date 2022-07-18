package model

import (
	"fmt"
	"gorm.io/gorm"
)

const (
	MinMethodLength = 1
	MaxMethodLength = 1000
)

type Permission struct {
	gorm.Model

	// The gRPC method name the permission belongs to
	Method string `gorm:"not null" json:"name"`

	// Whether the permission applies to users besides itself. If true, then the permission applies even if
	// the target of the method is not itself
	Other bool `gorm:"not null" json:"other"`
}

func (r *Permission) Validate() error {
	err := r.validateMethod()
	if err != nil {
		return err
	}

	return nil
}

func (r *Permission) validateMethod() error {
	if len(r.Method) < MinMethodLength {
		return fmt.Errorf("minimum method length is %d", MinMethodLength)
	}

	if len(r.Method) > MaxMethodLength {
		return fmt.Errorf("maximum method length is %d", MaxMethodLength)
	}

	return nil
}
