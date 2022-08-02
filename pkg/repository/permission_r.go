package repository

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	AddPermissionForUser(permission *model.UserPermission) error
	AddPermissionForRole(permission *model.RolePermission) error

	RemPermissionForUser(permission *model.UserPermission) error
	RemPermissionForRole(permission *model.RolePermission) error

	FindPermissionsForUserID(id uint) []*model.UserPermission
	FindPermissionsForRoleID(id uint) []*model.RolePermission

	ClearPermissionsForRole(id uint) error
	ClearPermissionsForUser(id uint) error

	WithTrx(*gorm.DB) PermissionRepository
	Migrate() error
}

type permissionRepository struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return permissionRepository{
		DB: db,
	}
}

func (r permissionRepository) AddPermissionForUser(permission *model.UserPermission) error {
	return r.DB.Create(&permission).Error
}

func (r permissionRepository) AddPermissionForRole(permission *model.RolePermission) error {
	return r.DB.Create(&permission).Error
}

func (r permissionRepository) RemPermissionForUser(permission *model.UserPermission) error {
	return r.DB.Delete(&permission).Error
}

func (r permissionRepository) RemPermissionForRole(permission *model.RolePermission) error {
	return r.DB.Delete(&permission).Error
}

func (r permissionRepository) FindPermissionsForUserID(id uint) []*model.UserPermission {
	var permissions []*model.UserPermission
	r.DB.Where("user_id = ?", id).Find(&permissions)
	return permissions
}

func (r permissionRepository) FindPermissionsForRoleID(id uint) []*model.RolePermission {
	var permissions []*model.RolePermission
	r.DB.Where("role_id = ?", id).Find(&permissions)
	return permissions
}

func (r permissionRepository) ClearPermissionsForRole(id uint) error {
	return r.DB.Delete(&model.RolePermission{}, "role_id = ?", id).Error
}

func (r permissionRepository) ClearPermissionsForUser(id uint) error {
	return r.DB.Delete(&model.UserPermission{}, "user_id = ?", id).Error
}

func (r permissionRepository) WithTrx(trx *gorm.DB) PermissionRepository {
	if trx == nil {
		return r
	}

	r.DB = trx
	return r
}

func (r permissionRepository) Migrate() error {
	err := r.DB.AutoMigrate(&model.UserPermission{})
	if err != nil {
		return err
	}

	return r.DB.AutoMigrate(&model.RolePermission{})
}
