package service

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/repository"
	"gorm.io/gorm"
)

type PermissionService interface {
	Create(*model.Permission) (*model.Permission, error)
	Save(*model.Permission) (*model.Permission, error)

	FindByMethodAndOtherOrCreate(method string, other bool) *model.Permission

	All() []*model.Permission
	FindById(id uint) *model.Permission
	FindByMethodAndOther(method string, other bool) *model.Permission
	FindWithMethod(name string) []*model.Permission

	WithTrx(*gorm.DB) PermissionService
	Migrate() error
}

type permissionService struct {
	permissionRepository repository.PermissionRepository
}

func NewPermissionService(r repository.PermissionRepository) PermissionService {
	return permissionService{
		permissionRepository: r,
	}
}

func (s permissionService) Create(role *model.Permission) (*model.Permission, error) {
	return s.permissionRepository.Create(role)
}

func (s permissionService) Save(role *model.Permission) (*model.Permission, error) {
	return s.permissionRepository.Save(role)
}

func (s permissionService) All() []*model.Permission {
	return s.permissionRepository.All()
}

func (s permissionService) FindById(id uint) *model.Permission {
	return s.permissionRepository.FindById(id)
}

func (s permissionService) FindWithMethod(name string) []*model.Permission {
	return s.permissionRepository.FindWithMethod(name)
}

func (s permissionService) WithTrx(db *gorm.DB) PermissionService {
	s.permissionRepository = s.permissionRepository.WithTrx(db)
	return s
}

func (s permissionService) FindAll() []*model.Permission {
	return s.permissionRepository.All()
}

func (s permissionService) FindByMethodAndOtherOrCreate(method string, other bool) *model.Permission {
	return s.permissionRepository.FindByMethodAndOtherOrCreate(method, other)
}

func (s permissionService) FindByMethodAndOther(method string, other bool) *model.Permission {
	return s.permissionRepository.FindByMethodAndOther(method, other)
}

func (s permissionService) Migrate() error {
	return s.permissionRepository.Migrate()
}
