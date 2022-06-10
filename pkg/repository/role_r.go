package repository

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(model.Role) (model.Role, error)
	Save(model.Role) (model.Role, error)

	All() []model.Role
	FindById(id uint) model.Role
	FindByName(name string) model.Role

	WithTrx(*gorm.DB) RoleRepository
	Migrate() error
}

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return roleRepository{
		DB: db,
	}
}

func (r roleRepository) Create(role model.Role) (model.Role, error) {
	panic("implement me")
}

func (r roleRepository) Save(role model.Role) (model.Role, error) {
	panic("implement me")
}

func (r roleRepository) All() []model.Role {
	panic("implement me")
}

func (r roleRepository) FindById(id uint) model.Role {
	panic("implement me")
}

func (r roleRepository) FindByName(name string) model.Role {
	panic("implement me")
}

func (r roleRepository) WithTrx(db *gorm.DB) RoleRepository {
	panic("implement me")
}

func (r roleRepository) Migrate() error {
	panic("implement me")
}
