package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"not null" json:"name"`
	Permissions []Permission `json:"permissions"`
}
