package srv

import (
	"context"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/internal/option"
	"github.com/ShatteredRealms/Accounts/pkg/pb"
	accountService "github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/ShatteredRealms/GoUtils/pkg/interceptor"
	"github.com/ShatteredRealms/GoUtils/pkg/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(
	u accountService.UserService,
	p accountService.PermissionService,
	r accountService.RoleService,
	jwt service.JWTService,
	logger log.LoggerService,
	config option.Config,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	publicRPCs := make(map[string]struct{})
	publicRPCs["/sro.accounts.HealthService/Health"] = struct{}{}
	publicRPCs["/sro.accounts.AuthenticationService/Login"] = struct{}{}
	publicRPCs["/sro.accounts.AuthenticationService/Register"] = struct{}{}

	authInterceptor := interceptor.NewAuthInterceptor(jwt, publicRPCs, getPermissions(p, u))

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	authenticationServiceServer := NewAuthenticationServiceServer(u, jwt, p, logger)
	pb.RegisterAuthenticationServiceServer(grpcServer, authenticationServiceServer)
	err := pb.RegisterAuthenticationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	userServiceServer := NewUserServiceServer(u, logger)
	pb.RegisterUserServiceServer(grpcServer, userServiceServer)
	err = pb.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	healthServiceServer := NewHealthServiceServer()
	pb.RegisterHealthServiceServer(grpcServer, healthServiceServer)
	err = pb.RegisterHealthServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	authorizationServiceServer := NewAuthorizationServiceServer(u, p, r, logger)
	pb.RegisterAuthorizationServiceServer(grpcServer, authorizationServiceServer)
	err = pb.RegisterAuthorizationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	// Compute the AllPermissions method once and save in memory
	authorizationServiceServer.SetupAllPermissions(grpcServer.GetServiceInfo())

	return grpcServer, gwmux, nil
}

func getPermissions(
	permissionService accountService.PermissionService,
	userService accountService.UserService,
) func(userID uint) map[string]bool {
	return func(userID uint) map[string]bool {
		user := userService.FindById(userID)
		if user == nil || !user.Exists() {
			return map[string]bool{}
		}

		resp := make(map[string]bool)

		for _, role := range user.Roles {
			for _, rolePermission := range permissionService.FindPermissionsForRoleID(role.ID) {
				resp[rolePermission.Permission] = resp[rolePermission.Permission] || rolePermission.Other
			}
		}

		for _, userPermission := range permissionService.FindPermissionsForUserID(userID) {
			resp[userPermission.Permission] = resp[userPermission.Permission] || userPermission.Other
		}

		return resp
	}
}
