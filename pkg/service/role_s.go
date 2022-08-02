package service

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/repository"
	"gorm.io/gorm"
)

type RoleService interface {
	Create(*model.Role) (*model.Role, error)
	Save(*model.Role) (*model.Role, error)
	Delete(*model.Role) error
	Update(*model.Role) error

	All() []*model.Role
	FindById(id uint) *model.Role
	FindByName(name string) *model.Role

	WithTrx(*gorm.DB) RoleService

	FindAll() []*model.Role
}

type roleService struct {
	roleRepository repository.RoleRepository
}

func NewRoleService(r repository.RoleRepository) RoleService {
	return roleService{
		roleRepository: r,
	}
}

func (s roleService) Create(role *model.Role) (*model.Role, error) {
	return s.roleRepository.Create(role)
}

func (s roleService) Save(role *model.Role) (*model.Role, error) {
	return s.roleRepository.Save(role)
}

func (s roleService) Delete(role *model.Role) error {
	return s.roleRepository.Delete(role)
}

func (s roleService) Update(role *model.Role) error {
	if len(role.Name) == 0 {
		return nil
	}

	return s.roleRepository.Update(role)
}

func (s roleService) All() []*model.Role {
	return s.roleRepository.All()
}

func (s roleService) FindById(id uint) *model.Role {
	return s.roleRepository.FindById(id)
}

func (s roleService) FindByName(name string) *model.Role {
	return s.roleRepository.FindByName(name)
}

func (s roleService) WithTrx(db *gorm.DB) RoleService {
	s.roleRepository = s.roleRepository.WithTrx(db)
	return s
}

func (s roleService) FindAll() []*model.Role {
	return s.roleRepository.All()
}
