package srv

import (
	"context"
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/pb"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type authorizationServiceServer struct {
	pb.UnimplementedAuthorizationServiceServer
	userService       service.UserService
	permissionService service.PermissionService
	roleService       service.RoleService
	logger            log.LoggerService
	allPermissions    *pb.UserPermissions
}

func NewAuthorizationServiceServer(
	u service.UserService,
	permissionService service.PermissionService,
	roleService service.RoleService,
	logger log.LoggerService,
) *authorizationServiceServer {
	return &authorizationServiceServer{
		userService:       u,
		permissionService: permissionService,
		roleService:       roleService,
		logger:            logger,
	}
}

func (s *authorizationServiceServer) GetAuthorization(
	ctx context.Context,
	message *pb.IDMessage,
) (*pb.AuthorizationMessage, error) {
	user := s.userService.FindById(uint(message.Id))
	if user == nil || !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	permissions := ConvertUserPermissions(s.permissionService.FindPermissionsForUserID(user.ID))
	roles := ConvertRolesWithoutPermissions(user.Roles)
	for i, role := range roles {
		roles[i].Permissions = ConvertRolePermissions(s.permissionService.FindPermissionsForRoleID(uint(role.Id)))
	}

	resp := &pb.AuthorizationMessage{
		UserId:      message.Id,
		Roles:       roles,
		Permissions: permissions,
	}

	return resp, nil
}

func (s *authorizationServiceServer) AddAuthorization(ctx context.Context, message *pb.AuthorizationMessage) (*emptypb.Empty, error) {
	user := s.userService.FindById(uint(message.UserId))
	if user == nil || !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	for _, v := range message.Permissions {
		_ = s.permissionService.AddPermissionForUser(&model.UserPermission{
			UserID:     user.ID,
			Permission: v.Permission,
			Other:      v.Other,
		})
		//if err != nil {
		//    return nil, status.Error(codes.Internal, err.Error())
		//}
	}

	for _, v := range message.Roles {
		_ = s.userService.AddToRole(
			user,
			&model.Role{
				Model: gorm.Model{
					ID: uint(v.Id),
				},
			})
		//if err != nil {
		//    return nil, status.Error(codes.Internal, err.Error())
		//}
	}

	return &emptypb.Empty{}, nil
}

func (s *authorizationServiceServer) RemoveAuthorization(ctx context.Context, message *pb.AuthorizationMessage) (*emptypb.Empty, error) {
	user := s.userService.FindById(uint(message.UserId))
	if user == nil || !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	for _, v := range message.Permissions {
		err := s.permissionService.RemPermissionForUser(&model.UserPermission{
			UserID:     user.ID,
			Permission: v.Permission,
			Other:      v.Other,
		})

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	for _, v := range message.Roles {
		err := s.userService.RemFromRole(
			user,
			&model.Role{
				Model: gorm.Model{
					ID: uint(v.Id),
				},
			})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *authorizationServiceServer) GetRoles(ctx context.Context, message *emptypb.Empty) (*pb.UserRoles, error) {
	resp := &pb.UserRoles{
		Roles: ConvertRolesNamesOnly(s.roleService.FindAll()),
	}
	for i, v := range resp.Roles {
		resp.Roles[i].Permissions = ConvertRolePermissions(s.permissionService.FindPermissionsForRoleID(uint(v.Id)))
	}

	return resp, nil
}

func (s *authorizationServiceServer) GetRole(ctx context.Context, message *pb.IDMessage) (*pb.UserRole, error) {
	resp := ConvertRoleNameOnly(s.roleService.FindById(uint(message.Id)))
	resp.Permissions = ConvertRolePermissions(s.permissionService.FindPermissionsForRoleID(uint(message.Id)))

	return resp, nil
}

func (s *authorizationServiceServer) CreateRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	_, err := s.roleService.Create(&model.Role{
		Name: message.Name,
	})

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *authorizationServiceServer) EditRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	err := s.roleService.Update(&model.Role{
		Model: gorm.Model{
			ID: uint(message.Id),
		},
		Name: message.Name,
	})

	newPermissions := make([]*model.RolePermission, len(message.Permissions))
	for i, permission := range message.Permissions {
		newPermissions[i] = &model.RolePermission{
			RoleID:     uint(message.Id),
			Permission: permission.Permission,
			Other:      permission.Other,
		}
	}

	err = s.permissionService.ResetPermissionsForRole(uint(message.Id), newPermissions)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *authorizationServiceServer) DeleteRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	err := s.roleService.Delete(&model.Role{
		Model: gorm.Model{
			ID: uint(message.Id),
		},
	})

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *authorizationServiceServer) GetAllPermissions(ctx context.Context, message *emptypb.Empty) (*pb.UserPermissions, error) {
	return s.allPermissions, nil
}

func (s *authorizationServiceServer) SetupAllPermissions(servicesInfo map[string]grpc.ServiceInfo) {
	count := 0
	for _, serviceInfo := range servicesInfo {
		count += len(serviceInfo.Methods)
	}

	methods := make([]*pb.UserPermission, count)
	index := 0
	for serviceName, serviceInfo := range servicesInfo {
		for _, method := range serviceInfo.Methods {
			methods[index] = &pb.UserPermission{Permission: fmt.Sprintf("/%s/%s", serviceName, method.Name)}
			index++
		}
	}

	s.allPermissions = &pb.UserPermissions{Permissions: methods}
}
