package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sergio-id/go-notes/cmd/category/config"
	"github.com/sergio-id/go-notes/internal/category/domain"
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

type categoryGRPCServer struct {
	gen.UnimplementedCategoryServiceServer
	cfg           *config.Config
	uc            domain.CategoryUseCase
	authDomainSvc domain.AuthDomainService
	log           logger.Logger
}

var _ gen.CategoryServiceServer = (*categoryGRPCServer)(nil)

var CategoryGRPCServerSet = wire.NewSet(NewGRPCCategoryServer)

func NewGRPCCategoryServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc domain.CategoryUseCase,
	authDomainSvc domain.AuthDomainService,
	log logger.Logger,
) gen.CategoryServiceServer {
	svc := categoryGRPCServer{
		cfg:           cfg,
		uc:            uc,
		authDomainSvc: authDomainSvc,
		log:           log,
	}

	gen.RegisterCategoryServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (s *categoryGRPCServer) ListCategories(ctx context.Context, request *gen.ListCategoriesRequest) (*gen.ListCategoriesResponse, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	arg := &domain.ListCategoriesParams{
		UserID: userID,
		Limit:  request.Limit,
		Offset: request.Offset,
	}

	categories, err := s.uc.ListCategories(ctx, arg)
	if err != nil {
		s.log.Errorf("[categoryGRPCServer] ListCategories-uc.ListCategories: %v", err)
		return nil, status.Error(codes.Internal, "Get categories failed")
	}

	var resCategories []*gen.Category
	for _, category := range categories {
		resCategories = append(resCategories, s.convertCategoryToProto(category))
	}

	res := gen.ListCategoriesResponse{
		Categories: resCategories,
	}

	return &res, nil
}

func (s *categoryGRPCServer) GetCategory(ctx context.Context, request *gen.GetCategoryRequest) (*gen.Category, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		s.log.Errorf("[categoryGRPCServer] GetCategory-uuid.Parse: %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid category id")
	}

	category, err := s.uc.GetCategory(ctx, userID, id)
	if err != nil {
		s.log.Errorf("[categoryGRPCServer] GetCategory-uc.GetCategory: %v", err)
		return nil, status.Error(codes.Internal, "Get category failed")
	}

	return s.convertCategoryToProto(category), nil
}

func (s *categoryGRPCServer) CreateCategory(ctx context.Context, request *gen.CreateCategoryRequest) (*gen.Category, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	arg := &domain.CreateCategoryParams{
		UserID:   userID,
		Title:    request.Title,
		Pinned:   request.Pinned,
		Priority: request.Priority,
	}

	createdCategory, err := s.uc.CreateCategory(ctx, arg)
	if err != nil {
		s.log.Errorf("[categoryGRPCServer] CreateCategory-uc.CreateCategory: %v", err)
		return nil, status.Error(codes.Internal, "Create category failed")
	}

	return s.convertCategoryToProto(createdCategory), nil
}

func (s *categoryGRPCServer) UpdateCategory(ctx context.Context, request *gen.UpdateCategoryRequest) (*gen.Category, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		s.log.Errorf("[categoryGRPCServer] UpdateCategory-uuid.Parse: %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid category id")
	}

	arg := &domain.UpdateCategoryParams{
		ID:       id,
		UserID:   userID,
		Title:    request.Title,
		Pinned:   request.Pinned,
		Priority: request.Priority,
	}

	updatedCategory, err := s.uc.UpdateCategory(ctx, arg)
	if err != nil {
		s.log.Errorf("[categoryGRPCServer] UpdateCategory-uc.UpdateCategory: %v", err)
		return nil, status.Error(codes.Internal, "Update category failed")
	}

	return s.convertCategoryToProto(updatedCategory), nil
}

func (s *categoryGRPCServer) DeleteCategory(ctx context.Context, request *gen.DeleteCategoryRequest) (*emptypb.Empty, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		s.log.Errorf("[categoryGRPCServer] DeleteCategory-uuid.Parse: %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid category id")
	}

	err = s.uc.DeleteCategory(ctx, userID, id)
	if err != nil {
		s.log.Errorf("[categoryGRPCServer] DeleteCategory-uc.DeleteCategory: %v", err)
		return nil, status.Error(codes.Internal, "Delete category failed")
	}

	return new(emptypb.Empty), nil
}

//--------------------------------------------------internal methods--------------------------------------------------//

func (s *categoryGRPCServer) convertCategoryToProto(category *domain.Category) *gen.Category {
	return &gen.Category{
		Id:        category.ID.String(),
		UserId:    category.UserID.String(),
		Title:     category.Title,
		Pinned:    category.Pinned,
		Priority:  category.Priority,
		UpdatedAt: timestamppb.New(category.UpdatedAt),
		CreatedAt: timestamppb.New(category.UpdatedAt),
	}
}

func (s *categoryGRPCServer) getUserIDByTokenFromCtx(ctx context.Context) (uuid.UUID, error) {
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

func (s *categoryGRPCServer) getTokenFromCtx(ctx context.Context) (string, error) {
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
