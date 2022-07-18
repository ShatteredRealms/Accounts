package repository

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(*model.Permission) (*model.Permission, error)
	Save(*model.Permission) (*model.Permission, error)

	FindByMethodAndOtherOrCreate(method string, other bool) *model.Permission

	All() []*model.Permission
	FindById(id uint) *model.Permission
	FindByMethodAndOther(method string, other bool) *model.Permission
	FindWithMethod(name string) []*model.Permission

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

func (r permissionRepository) FindByMethodAndOtherOrCreate(method string, other bool) *model.Permission {
	permission := r.FindByMethodAndOther(method, other)
	if permission != nil {
		return permission
	}

	permission, err := r.Create(&model.Permission{
		Method: method,
		Other:  other,
	})

	if err != nil {
		return nil
	}

	return permission
}

func (r permissionRepository) FindByMethodAndOther(method string, other bool) *model.Permission {
	var permission *model.Permission
	r.DB.Where("method = ? AND other = ?", method, other).Find(&permission)
	return permission
}

func (r permissionRepository) FindWithMethod(name string) []*model.Permission {
	var permissions []*model.Permission
	r.DB.Where("method = ?", name).Find(&permissions)
	return permissions
}

func (r permissionRepository) Create(permission *model.Permission) (*model.Permission, error) {
	err := permission.Validate()
	if err != nil {
		return nil, err
	}

	conflictPermission := r.FindByMethodAndOther(permission.Method, permission.Other)
	if conflictPermission != nil {
		return nil, fmt.Errorf("permission with method and other state exists")
	}

	return permission, r.DB.Save(&permission).Error
}

func (r permissionRepository) Save(permission *model.Permission) (*model.Permission, error) {
	conflictPermission := r.FindByMethodAndOther(permission.Method, permission.Other)
	if conflictPermission != nil {
		return nil, fmt.Errorf("permission with method and other state exists")
	}

	return permission, r.DB.Save(&permission).Error
}

func (r permissionRepository) All() []*model.Permission {
	var permissions []*model.Permission
	r.DB.Find(&permissions)
	return permissions
}

func (r permissionRepository) FindById(id uint) *model.Permission {
	var permission *model.Permission
	r.DB.Where("id = ?", id).Find(&permission)
	return permission
}

func (r permissionRepository) WithTrx(trx *gorm.DB) PermissionRepository {
	if trx == nil {
		return r
	}

	r.DB = trx
	return r
}

func (r permissionRepository) Migrate() error {
	return r.DB.AutoMigrate(&model.Permission{})
}
