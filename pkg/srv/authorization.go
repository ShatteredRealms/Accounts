package srv

import (
	"context"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/accountspb"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authorizationServiceServer struct {
	accountspb.UnimplementedAuthorizationServiceServer
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
	message *accountspb.GetAuthorizationRequest,
) (*accountspb.AuthorizationMessage, error) {
	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	roles := make([]*accountspb.UserRole, len(user.Roles))

	for i, v := range user.Roles {
		roles[i] = &accountspb.UserRole{
			Id:   uint64(v.ID),
			Name: v.Name,
		}
	}

	permissions := make([]*accountspb.UserPermission, len(user.Permissions))
	for i, v := range user.Permissions {
		permissions[i] = &accountspb.UserPermission{
			Method: v.Method,
			Other:  v.Other,
		}
	}

	resp := &accountspb.AuthorizationMessage{
		UserId:      message.UserId,
		Roles:       roles,
		Permissions: permissions,
	}

	return resp, nil
}

func (s *authorizationServiceServer) SetAuthorization(
	ctx context.Context,
	message *accountspb.AuthorizationMessage,
) (*emptypb.Empty, error) {
	return nil, nil
}
