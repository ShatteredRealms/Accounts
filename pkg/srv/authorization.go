package srv

import (
	"context"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/pb"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authorizationServiceServer struct {
	pb.UnimplementedAuthorizationServiceServer
	userService service.UserService
	logger      log.LoggerService
}

func NewAuthorizationServiceServer(u service.UserService, logger log.LoggerService) *authorizationServiceServer {
	return &authorizationServiceServer{
		userService: u,
		logger:      logger,
	}
}

func (s *authorizationServiceServer) GetAuthorization(
	ctx context.Context,
	message *pb.GetAuthorizationRequest,
) (*pb.AuthorizationMessage, error) {
	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	roles := make([]*pb.UserRole, len(user.Roles))

	for i, v := range user.Roles {
		roles[i] = &pb.UserRole{
			Id:   uint64(v.ID),
			Name: v.Name,
		}
	}

	permissions := make([]*pb.UserPermission, len(user.Permissions))
	for i, v := range user.Permissions {
		permissions[i] = &pb.UserPermission{
			Method: v.Method,
			Other:  v.Other,
		}
	}

	resp := &pb.AuthorizationMessage{
		UserId:      message.UserId,
		Roles:       roles,
		Permissions: permissions,
	}

	return resp, nil
}

func (s *authorizationServiceServer) SetAuthorization(
	ctx context.Context,
	message *pb.AuthorizationMessage,
) (*emptypb.Empty, error) {
	return nil, nil
}
