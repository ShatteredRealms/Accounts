package repository

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(*model.Role) (*model.Role, error)
	Save(*model.Role) (*model.Role, error)

	All() []*model.Role
	FindById(id uint) *model.Role
	FindByName(name string) *model.Role

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

func (r roleRepository) Create(role *model.Role) (*model.Role, error) {
	err := role.Validate()
	if err != nil {
		return nil, err
	}

	existingRoleWithName := r.FindByName(role.Name)
	if existingRoleWithName != nil {
		return nil, fmt.Errorf("name already exists")
	}

	return role, r.DB.Save(&role).Error
}

func (r roleRepository) Save(role *model.Role) (*model.Role, error) {
	existingRoleWithName := r.FindByName(role.Name)
	if existingRoleWithName != nil {
		return nil, fmt.Errorf("name already exists")
	}

	return role, r.DB.Save(&role).Error
}

func (r roleRepository) All() []*model.Role {
	var roles []*model.Role
	r.DB.Find(&roles)
	return roles
}

func (r roleRepository) FindById(id uint) *model.Role {
	var role *model.Role
	r.DB.Where("id = ?", id).Find(&role)
	return role
}

func (r roleRepository) FindByName(name string) *model.Role {
	var role *model.Role
	r.DB.Where("name = ?", name).Find(&role)
	return role
}

func (r roleRepository) WithTrx(trx *gorm.DB) RoleRepository {
	if trx == nil {
		return r
	}

	r.DB = trx
	return r
}

func (r roleRepository) Migrate() error {
	return r.DB.AutoMigrate(&model.Role{})
}
