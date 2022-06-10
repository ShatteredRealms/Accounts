package srv

import (
	"context"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/internal/option"
	"github.com/ShatteredRealms/Accounts/pkg/accountspb"
	accountService "github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/ShatteredRealms/GoUtils/pkg/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(
	u accountService.UserService,
	jwt service.JWTService,
	logger log.LoggerService,
	config option.Config,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	grpcServer := grpc.NewServer()
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	authenticationServiceServer := NewAuthenticationServiceServer(u, jwt, logger)
	accountspb.RegisterAuthenticationServiceServer(grpcServer, authenticationServiceServer)
	err := accountspb.RegisterAuthenticationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	authorizationServiceServer := NewAuthorizationServiceServer(u, logger)
	accountspb.RegisterAuthorizationServiceServer(grpcServer, authorizationServiceServer)
	err = accountspb.RegisterAuthorizationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	userServiceServer := NewUserServiceServer(u, logger)
	accountspb.RegisterUserServiceServer(grpcServer, userServiceServer)
	err = accountspb.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		config.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	return grpcServer, gwmux, nil
}
