package srv

import (
	"context"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/accountspb"
	"github.com/ShatteredRealms/Accounts/pkg/service"
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
	return nil, nil
}

func (s *authorizationServiceServer) SetAuthorization(
	ctx context.Context,
	message *accountspb.AuthorizationMessage,
) (*emptypb.Empty, error) {
	return nil, nil
}
