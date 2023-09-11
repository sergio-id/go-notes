package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sergio-id/go-notes/cmd/user/config"
	"github.com/sergio-id/go-notes/internal/user/domain"
	"github.com/sergio-id/go-notes/pkg/logger"
	gen "github.com/sergio-id/go-notes/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userGRPCServer struct {
	gen.UnimplementedUserServiceServer
	cfg           *config.Config
	uc            domain.UserUseCase
	authDomainSvc domain.AuthDomainService
	log           logger.Logger
}

var _ gen.UserServiceServer = (*userGRPCServer)(nil)

var UserGRPCServerSet = wire.NewSet(NewGRPCUserServer)

func NewGRPCUserServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc domain.UserUseCase,
	authDomainSvc domain.AuthDomainService,
	log logger.Logger,
) gen.UserServiceServer {
	svc := userGRPCServer{
		cfg:           cfg,
		uc:            uc,
		authDomainSvc: authDomainSvc,
		log:           log,
	}

	gen.RegisterUserServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (s *userGRPCServer) GetMe(ctx context.Context, request *gen.GetMeRequest) (*gen.User, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	user, err := s.uc.GetUser(ctx, userID)
	if err != nil {
		s.log.Errorf("[userGRPCServer] GetMe-uc.GetUser: %v", err)
		return nil, status.Error(codes.Internal, "Get user failed")
	}

	user.Password = ""

	return s.convertUserToProto(user), nil
}

func (s *userGRPCServer) CreateUser(ctx context.Context, request *gen.CreateUserRequest) (*gen.User, error) {
	arg := &domain.CreateUserParams{
		Email:    request.Email,
		Password: request.Password,
	}

	createdUser, err := s.uc.CreateUser(ctx, arg)
	if err != nil {
		s.log.Errorf("[userGRPCServer] CreateUser-uc.CreateUser: %v", err)
		return nil, status.Error(codes.Internal, "Create user failed")
	}

	createdUser.Password = ""

	return s.convertUserToProto(createdUser), nil
}

func (s *userGRPCServer) UpdateUser(ctx context.Context, request *gen.UpdateUserRequest) (*gen.User, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	arg := &domain.UpdateUserParams{
		ID:        userID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}

	updatedUser, err := s.uc.UpdateUser(ctx, arg)
	if err != nil {
		s.log.Errorf("[userGRPCServer] UpdateUser-uc.UpdateUser: %v", err)
		return nil, status.Error(codes.Internal, "Update user failed")
	}

	updatedUser.Password = ""

	return s.convertUserToProto(updatedUser), nil
}

func (s *userGRPCServer) DeleteUser(ctx context.Context, request *gen.DeleteUserRequest) (*emptypb.Empty, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	err = s.uc.DeleteUser(ctx, userID)
	if err != nil {
		s.log.Errorf("[userGRPCServer] DeleteUser-uc.DeleteUser: %v", err)
		return nil, status.Error(codes.Internal, "Delete user failed")
	}

	return new(emptypb.Empty), nil
}

func (s *userGRPCServer) GetUserByEmail(ctx context.Context, request *gen.GetUserByEmailRequest) (*gen.User, error) {
	user, err := s.uc.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetUserByEmail")
	}

	return s.convertUserToProto(user), nil
}

//--------------------------------------------------internal methods--------------------------------------------------//

func (s *userGRPCServer) convertUserToProto(user *domain.User) *gen.User {
	return &gen.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		CreatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func (s *userGRPCServer) getUserIDByTokenFromCtx(ctx context.Context) (uuid.UUID, error) {
	token, err := s.getTokenFromCtx(ctx)
	if err != nil {
		return uuid.UUID{}, status.Errorf(codes.Unauthenticated, "sessUC.getSessionIDFromCtx: %v", err)
	}

	userID, err := s.authDomainSvc.GetUserIDByToken(ctx, token)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return uuid.UUID{}, status.Errorf(codes.NotFound, "sessUC.GetSessionByID")
		}
		return uuid.UUID{}, status.Errorf(codes.Unknown, "sessUC.GetSessionByID: %v", err)
	}

	return userID, nil
}

func (s *userGRPCServer) getTokenFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata.FromIncomingContext")
	}

	//parse bearer token
	token := md.Get("authorization")
	if len(token) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "md.Get")
	}

	//extract token
	bearerToken := token[0][7:]

	return bearerToken, nil
}
