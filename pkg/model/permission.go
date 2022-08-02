package model

// TODO(wil) Refactor Permission/Other inside UserPermission and RolePermission and refactor functions that use these
//           structs.

// TODO(wil) Refactor `Other` to `Global` for readability and conveying meaning?

// UserPermission Database model for customized user and permissions join table
type UserPermission struct {
	// User The User to grant the permission to
	UserID uint `gorm:"primaryKey" json:"user_id"`

	// Permission The permission that is assigned to the user
	Permission string `gorm:"primaryKey" json:"permission"`

	// Whether the permission applies to users besides itself. If true, then the permission applies even if
	// the target of the method is not itself
	Other bool `gorm:"not null" json:"other"`
}

// RolePermission Database model for customized role and permissions join table
type RolePermission struct {
	// Role The Role to grant the permission to
	RoleID uint `gorm:"primaryKey" json:"role_id"`

	// Permission The permission that is assigned to the user
	Permission string `gorm:"primaryKey" json:"permission"`

	// Whether the permission applies to users besides itself. If true, then the permission applies even if
	// the target of the method is not itself
	Other bool `gorm:"not null" json:"other"`
}
