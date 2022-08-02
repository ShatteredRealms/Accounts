package srv

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/pb"
)

func ConvertRolePermission(inPermission *model.RolePermission) *pb.UserPermission {
	return &pb.UserPermission{
		Permission: inPermission.Permission,
		Other:      inPermission.Other,
	}
}

func ConvertRolePermissions(inPermissions []*model.RolePermission) []*pb.UserPermission {
	permissions := make([]*pb.UserPermission, len(inPermissions))
	for i, permission := range inPermissions {
		permissions[i] = ConvertRolePermission(permission)
	}

	return permissions
}

func ConvertUserPermission(inPermission *model.UserPermission) *pb.UserPermission {
	return &pb.UserPermission{
		Permission: inPermission.Permission,
		Other:      inPermission.Other,
	}
}

func ConvertUserPermissions(inPermissions []*model.UserPermission) []*pb.UserPermission {
	permissions := make([]*pb.UserPermission, len(inPermissions))
	for i, permission := range inPermissions {
		permissions[i] = ConvertUserPermission(permission)
	}

	return permissions
}

func ConvertRoleWithoutPermissions(inRole *model.Role) *pb.UserRole {
	return &pb.UserRole{
		Id:   uint64(inRole.ID),
		Name: inRole.Name,
	}
}

func ConvertRolesWithoutPermissions(inRoles []*model.Role) []*pb.UserRole {
	roles := make([]*pb.UserRole, len(inRoles))
	for i, role := range inRoles {
		roles[i] = ConvertRoleWithoutPermissions(role)
	}

	return roles
}

func ConvertRoleNameOnly(inRole *model.Role) *pb.UserRole {
	return &pb.UserRole{
		Id:   uint64(inRole.ID),
		Name: inRole.Name,
	}
}

func ConvertRolesNamesOnly(inRoles []*model.Role) []*pb.UserRole {
	roles := make([]*pb.UserRole, len(inRoles))
	for i, role := range inRoles {
		roles[i] = ConvertRoleNameOnly(role)
	}

	return roles
}
