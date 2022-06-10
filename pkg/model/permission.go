package model

import "gorm.io/gorm"

type Permission struct {
	gorm.Model

	// The gRPC method name the permission belongs to
	Method string `gorm:"not null" json:"name"`

	// Whether the permission applies to users besides itself. If true, then the permission applies even if
	// the target of the rpc is not itself
	Other bool `gorm:"not null" json:"other"`
}
