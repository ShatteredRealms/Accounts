package srv

import (
	"context"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/pb"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
	logger      log.LoggerService
}

func NewUserServiceServer(u service.UserService, logger log.LoggerService) *userServiceServer {
	return &userServiceServer{
		userService: u,
		logger:      logger,
	}
}

func (s *userServiceServer) GetAll(
	ctx context.Context,
	message *emptypb.Empty,
) (*pb.GetAllUsersResponse, error) {
	users := s.userService.FindAll()
	resp := &pb.GetAllUsersResponse{
		Users: []*pb.UserMessage{},
	}
	for _, u := range users {
		resp.Users = append(resp.Users, &pb.UserMessage{
			Id:       uint64(u.ID),
			Email:    u.Email,
			Username: u.Username,
		})
	}

	return resp, nil
}

func (s *userServiceServer) Get(
	ctx context.Context,
	message *pb.GetUserMessage,
) (*pb.GetUserResponse, error) {
	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{
		Id:        uint64(user.ID),
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (s *userServiceServer) Edit(
	ctx context.Context,
	message *pb.UserDetails,
) (*emptypb.Empty, error) {
	user := s.userService.FindById(uint(message.UserId))
	if !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	newUserData := model.User{
		FirstName: message.FirstName,
		LastName:  message.LastName,
		Username:  message.Username,
		Email:     message.Email,
		Password:  message.Password,
	}

	err := user.UpdateInfo(newUserData)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.userService.Save(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
