package repository

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(model.Permission) (model.Permission, error)
	Save(model.Permission) (model.Permission, error)

	All() []model.Permission
	FindById(id uint) model.Permission
	FindByMethod(name string) model.Permission

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

func (r permissionRepository) FindByMethod(name string) model.Permission {
	panic("implement me")
}

func (r permissionRepository) Create(permission model.Permission) (model.Permission, error) {
	panic("implement me")
}

func (r permissionRepository) Save(permission model.Permission) (model.Permission, error) {
	panic("implement me")
}

func (r permissionRepository) All() []model.Permission {
	panic("implement me")
}

func (r permissionRepository) FindById(id uint) model.Permission {
	panic("implement me")
}

func (r permissionRepository) WithTrx(db *gorm.DB) PermissionRepository {
	panic("implement me")
}

func (r permissionRepository) Migrate() error {
	panic("implement me")
}
