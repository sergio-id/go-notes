package grpc

import (
	"context"
	"github.com/google/wire"
	"github.com/sergio-id/go-notes/cmd/auth/config"
	"github.com/sergio-id/go-notes/internal/auth/domain"
	"github.com/sergio-id/go-notes/pkg/logger"
	"github.com/sergio-id/go-notes/pkg/security"
	"github.com/sergio-id/go-notes/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type authGRPCServer struct {
	gen.UnimplementedAuthServiceServer
	cfg           *config.Config
	uc            domain.SessionUseCase
	userDomainSvc domain.UserDomainService
	log           logger.Logger
}

var _ gen.AuthServiceServer = (*authGRPCServer)(nil)

var AuthGRPCServerSet = wire.NewSet(NewGRPCAuthServer)

func NewGRPCAuthServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc domain.SessionUseCase,
	userDomainSvc domain.UserDomainService,
	log logger.Logger,
) gen.AuthServiceServer {
	svc := authGRPCServer{
		cfg:           cfg,
		uc:            uc,
		userDomainSvc: userDomainSvc,
		log:           log,
	}

	gen.RegisterAuthServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (s *authGRPCServer) SignUp(ctx context.Context, req *gen.SignUpRequest) (*gen.SignUpResponse, error) {
	arg := &domain.CreateUserParams{
		Email:    req.Email,
		Password: req.Password,
	}

	//check if user already exists
	_, err := s.userDomainSvc.GetUserByEmail(ctx, arg.Email)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "User already exists")
	}

	//hash password
	hashedPassword, err := security.HashPassword(s.cfg.Security, arg.Password)
	if err != nil {
		s.log.Errorf("[authGRPCServer] SignUp-hashPassword: %v", err)
		return nil, status.Error(codes.Internal, "Hashing password failed")
	}
	arg.Password = hashedPassword

	//create user
	userDTO, err := s.userDomainSvc.CreateUser(ctx, arg)
	if err != nil {
		s.log.Errorf("[authGRPCServer] SignUp-CreateUser: %v", err)
		return nil, status.Error(codes.Internal, "Creating user failed")
	}

	arg1 := &domain.CreateSessionParams{
		UserID:   userDTO.ID,
		Duration: time.Second * time.Duration(s.cfg.SessionDuration),
	}

	//create session with duration
	session, err := s.uc.CreateSession(ctx, arg1)
	if err != nil {
		s.log.Errorf("[authGRPCServer] SignUp-CreateSession: %v", err)
		return nil, status.Error(codes.Internal, "Creating session failed")
	}

	return &gen.SignUpResponse{
		Token: session.SessionID,
	}, status.Error(codes.OK, "User created")
}

func (s *authGRPCServer) SignIn(ctx context.Context, req *gen.SignInRequest) (*gen.SignInResponse, error) {
	//get user by email
	userDTO, err := s.userDomainSvc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	//check password
	if !security.IsEqual(userDTO.Password, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "Invalid password")
	}

	//create session with duration
	arg := &domain.CreateSessionParams{
		UserID:   userDTO.ID,
		Duration: time.Second * time.Duration(s.cfg.SessionDuration),
	}

	session, err := s.uc.CreateSession(ctx, arg)
	if err != nil {
		s.log.Errorf("[authGRPCServer] SignIn-CreateSession: %v", err)
		return nil, status.Error(codes.Internal, "Creating session failed")
	}

	return &gen.SignInResponse{
		Token: session.SessionID,
	}, nil
}

func (s *authGRPCServer) SignOut(ctx context.Context, req *gen.SignOutRequest) (*emptypb.Empty, error) {
	token, err := s.getTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	err = s.uc.DeleteByID(ctx, token)
	if err != nil {
		s.log.Errorf("[authGRPCServer] SignOut-DeleteByID: %v", err)
		return nil, status.Error(codes.Internal, "Deleting session failed")
	}

	return new(emptypb.Empty), nil
}

func (s *authGRPCServer) GetSession(ctx context.Context, req *gen.GetSessionRequest) (*gen.Session, error) {
	session, err := s.uc.GetSessionByID(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	return &gen.Session{
		Token:  session.SessionID,
		UserId: session.UserID,
	}, nil
}

func (s *authGRPCServer) getTokenFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "metadata.FromIncomingContext")
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
