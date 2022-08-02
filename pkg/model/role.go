package model

import (
	"fmt"
	"gorm.io/gorm"
)

const (
	MinRoleNameLength = 3
	MaxRoleNameLength = 255
)

type Role struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
}

func (r *Role) Validate() error {
	err := r.validateName()
	if err != nil {
		return err
	}

	return nil
}

func (r *Role) validateName() error {
	if len(r.Name) < MinRoleNameLength {
		return fmt.Errorf("minimum name length is %d", MinRoleNameLength)
	}

	if len(r.Name) > MaxRoleNameLength {
		return fmt.Errorf("maximum name length is %d", MaxRoleNameLength)
	}

	return nil
}
