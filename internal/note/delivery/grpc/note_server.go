package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sergio-id/go-notes/cmd/note/config"
	"github.com/sergio-id/go-notes/internal/note/domain"
	"github.com/sergio-id/go-notes/pkg/logger"
	"github.com/sergio-id/go-notes/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type noteGRPCServer struct {
	gen.UnimplementedNoteServiceServer
	cfg           *config.Config
	uc            domain.NoteUseCase
	authDomainSvc domain.AuthDomainService
	log           logger.Logger
}

var _ gen.NoteServiceServer = (*noteGRPCServer)(nil)

var NoteGRPCServerSet = wire.NewSet(NewGRPCNoteServer)

func NewGRPCNoteServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc domain.NoteUseCase,
	authDomainSvc domain.AuthDomainService,
	log logger.Logger,
) gen.NoteServiceServer {
	svc := noteGRPCServer{
		cfg:           cfg,
		uc:            uc,
		authDomainSvc: authDomainSvc,
		log:           log,
	}

	gen.RegisterNoteServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (s *noteGRPCServer) ListNotes(ctx context.Context, request *gen.ListNotesRequest) (*gen.ListNotesResponse, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	arg := &domain.ListNotesParams{
		UserID: userID,
		Limit:  request.Limit,
		Offset: request.Offset,
	}

	notes, err := s.uc.ListNotes(ctx, arg)
	if err != nil {
		s.log.Errorf("[noteGRPCServer] ListNotes-uc.ListNotes: %v", err)
		return nil, status.Error(codes.Internal, "Get notes failed")
	}

	var resNotes []*gen.Note
	for _, note := range notes {
		resNotes = append(resNotes, s.convertNoteToProto(note))
	}

	res := gen.ListNotesResponse{
		Notes: resNotes,
	}

	return &res, nil
}

func (s *noteGRPCServer) GetNote(ctx context.Context, request *gen.GetNoteRequest) (*gen.Note, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		s.log.Errorf("[noteGRPCServer] GetNote-uuid.Parse: %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid note id")
	}

	note, err := s.uc.GetNote(ctx, userID, id)
	if err != nil {
		s.log.Errorf("[noteGRPCServer] GetNote-uc.GetNote: %v", err)
		return nil, status.Error(codes.Internal, "Get note failed")
	}

	return s.convertNoteToProto(note), nil
}

func (s *noteGRPCServer) CreateNote(ctx context.Context, request *gen.CreateNoteRequest) (*gen.Note, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	arg := &domain.CreateNoteParams{
		UserID:     userID,
		CategoryID: s.getCategoryUIIDOrNullUIID(request.CategoryId),
		Title:      request.Title,
		Content:    request.Content,
		Pinned:     request.Pinned,
		Priority:   request.Priority,
	}

	createdNote, err := s.uc.CreateNote(ctx, arg)
	if err != nil {
		s.log.Errorf("[noteGRPCServer] CreateNote-uc.CreateNote: %v", err)
		return nil, status.Error(codes.Internal, "Create note failed")
	}

	return s.convertNoteToProto(createdNote), nil
}

func (s *noteGRPCServer) UpdateNote(ctx context.Context, request *gen.UpdateNoteRequest) (*gen.Note, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		s.log.Errorf("[noteGRPCServer] UpdateNote-uuid.Parse: %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid note id")
	}
	arg := &domain.UpdateNoteParams{
		ID:       id,
		UserID:   userID,
		Title:    request.Title,
		Content:  request.Content,
		Pinned:   request.Pinned,
		Priority: request.Priority,
	}

	updatedNote, err := s.uc.UpdateNote(ctx, arg)
	if err != nil {
		s.log.Errorf("[noteGRPCServer] UpdateNote-uc.UpdateNote: %v", err)
		return nil, status.Error(codes.Internal, "Update note failed")
	}

	return s.convertNoteToProto(updatedNote), nil
}

func (s *noteGRPCServer) DeleteNote(ctx context.Context, request *gen.DeleteNoteRequest) (*emptypb.Empty, error) {
	// check if user is authenticated and get user id
	userID, err := s.getUserIDByTokenFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		s.log.Errorf("[noteGRPCServer] GetCategory-uuid.Parse: %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid category id")
	}

	err = s.uc.DeleteNote(ctx, userID, id)
	if err != nil {
		s.log.Errorf("[noteGRPCServer] DeleteNote-uc.DeleteNote: %v", err)
		return nil, status.Error(codes.Internal, "Delete note failed")
	}

	return new(emptypb.Empty), nil
}

//--------------------------------------------------internal methods--------------------------------------------------//

func (s *noteGRPCServer) convertNoteToProto(note *domain.Note) *gen.Note {
	protoNote := gen.Note{
		Id:         note.ID.String(),
		UserId:     note.UserID.String(),
		CategoryId: nil,
		Title:      note.Title,
		Content:    note.Content,
		Pinned:     note.Pinned,
		Priority:   note.Priority,
		UpdatedAt:  timestamppb.New(note.UpdatedAt),
		CreatedAt:  timestamppb.New(note.CreatedAt),
	}

	if note.CategoryID != nil {
		c := note.CategoryID.String()
		protoNote.CategoryId = &c
	}

	return &protoNote
}

func (s *noteGRPCServer) getUserIDByTokenFromCtx(ctx context.Context) (uuid.UUID, error) {
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

func (s *noteGRPCServer) getTokenFromCtx(ctx context.Context) (string, error) {
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

func (s *noteGRPCServer) getCategoryUIIDOrNullUIID(categoryID *string) uuid.NullUUID {
	if categoryID == nil || *categoryID == "" {
		return uuid.NullUUID{
			UUID:  uuid.UUID{},
			Valid: false,
		}
	}

	categoryUIID, err := uuid.Parse(*categoryID)
	if err != nil {
		s.log.Warnf("[noteGRPCServer] CreateNote-uuid.Parse: %v", err)
		return uuid.NullUUID{}
	}

	return uuid.NullUUID{
		UUID:  categoryUIID,
		Valid: true,
	}
}
