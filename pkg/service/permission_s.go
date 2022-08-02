package service

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/repository"
	"gorm.io/gorm"
)

type PermissionService interface {
	AddPermissionForUser(permission *model.UserPermission) error
	AddPermissionForRole(permission *model.RolePermission) error

	RemPermissionForUser(permission *model.UserPermission) error
	RemPermissionForRole(permission *model.RolePermission) error

	FindPermissionsForUserID(id uint) []*model.UserPermission
	FindPermissionsForRoleID(id uint) []*model.RolePermission

	ClearPermissionsForRole(id uint) error
	ClearPermissionsForUser(id uint) error

	ResetPermissionsForRole(id uint, permissions []*model.RolePermission) error
	ResetPermissionsForUser(id uint, permissions []*model.UserPermission) error

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

func (s permissionService) AddPermissionForUser(permission *model.UserPermission) error {
	return s.permissionRepository.AddPermissionForUser(permission)
}

func (s permissionService) AddPermissionForRole(permission *model.RolePermission) error {
	return s.permissionRepository.AddPermissionForRole(permission)
}

func (s permissionService) RemPermissionForUser(permission *model.UserPermission) error {
	return s.permissionRepository.RemPermissionForUser(permission)
}

func (s permissionService) RemPermissionForRole(permission *model.RolePermission) error {
	return s.permissionRepository.RemPermissionForRole(permission)
}

func (s permissionService) WithTrx(db *gorm.DB) PermissionService {
	s.permissionRepository = s.permissionRepository.WithTrx(db)
	return s
}

func (s permissionService) FindPermissionsForUserID(id uint) []*model.UserPermission {
	return s.permissionRepository.FindPermissionsForUserID(id)
}

func (s permissionService) FindPermissionsForRoleID(id uint) []*model.RolePermission {
	return s.permissionRepository.FindPermissionsForRoleID(id)
}

func (s permissionService) ClearPermissionsForRole(id uint) error {
	return s.permissionRepository.ClearPermissionsForRole(id)
}

func (s permissionService) ClearPermissionsForUser(id uint) error {
	return s.permissionRepository.ClearPermissionsForUser(id)
}

func (s permissionService) ResetPermissionsForRole(id uint, permissions []*model.RolePermission) error {
	err := s.ClearPermissionsForRole(id)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		if permission.RoleID == id {
			fmt.Println("new permission")
			err = s.AddPermissionForRole(permission)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s permissionService) ResetPermissionsForUser(id uint, permissions []*model.UserPermission) error {
	err := s.ClearPermissionsForUser(id)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		if permission.UserID == 0 || permission.UserID == id {
			err = s.AddPermissionForUser(permission)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s permissionService) Migrate() error {
	return s.permissionRepository.Migrate()
}
