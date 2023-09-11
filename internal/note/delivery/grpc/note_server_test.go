package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/sergio-id/go-notes/internal/note/domain"
	"github.com/sergio-id/go-notes/proto/gen"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestListNotes(t *testing.T) {
	server := &noteGRPCServer{
		cfg:           nil,
		uc:            &mockNoteUseCase{},
		authDomainSvc: &mockAuthDomainService{},
		log:           nil,
	}

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{"authorization": "bearer 97f1950d-7120-451c-8545-615bb845b78c"}),
	)

	limit := int32(10)
	offset := int32(0)

	// Test case: Valid list notes request
	req := &gen.ListNotesRequest{
		Limit:  &limit,
		Offset: &offset,
	}
	resp, err := server.ListNotes(ctx, req)
	if err != nil {
		t.Errorf("ListNotes failed with error: %v", err)
	}
	if resp == nil || len(resp.Notes) != 2 {
		t.Errorf("ListNotes response is incorrect")
	}
}

func TestGetNote(t *testing.T) {
	server := &noteGRPCServer{
		cfg:           nil,
		uc:            &mockNoteUseCase{},
		authDomainSvc: &mockAuthDomainService{},
		log:           nil,
	}

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{"authorization": "bearer 97f1950d-7120-451c-8545-615bb845b78c"}),
	)

	noteID := uuid.New().String()
	// Test case: Valid get note request
	req := &gen.GetNoteRequest{
		Id: noteID,
	}
	resp, err := server.GetNote(ctx, req)
	if err != nil {
		t.Errorf("GetNote failed with error: %v", err)
	}
	if resp == nil || resp.Id != noteID {
		t.Errorf("GetNote response is incorrect")
	}
}

func TestCreateNote(t *testing.T) {
	server := &noteGRPCServer{
		cfg:           nil,
		uc:            &mockNoteUseCase{},
		authDomainSvc: &mockAuthDomainService{},
		log:           nil,
	}

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{"authorization": "bearer 97f1950d-7120-451c-8545-615bb845b78c"}),
	)

	categoryId := uuid.New().String()
	// Test case: Valid create note request
	req := &gen.CreateNoteRequest{
		CategoryId: &categoryId,
		Title:      "New Note",
		Content:    "This is a new note.",
		Pinned:     true,
		Priority:   2,
	}
	resp, err := server.CreateNote(ctx, req)
	if err != nil {
		t.Errorf("CreateNote failed with error: %v", err)
	}
	if resp == nil || resp.Title != "New Note" {
		t.Errorf("CreateNote response is incorrect")
	}
}

func TestUpdateNote(t *testing.T) {
	server := &noteGRPCServer{
		cfg:           nil,
		uc:            &mockNoteUseCase{},
		authDomainSvc: &mockAuthDomainService{},
		log:           nil,
	}

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{"authorization": "bearer 97f1950d-7120-451c-8545-615bb845b78c"}),
	)

	// Test case: Valid update note request
	req := &gen.UpdateNoteRequest{
		Id:       uuid.New().String(),
		Title:    "Updated Note",
		Content:  "This is an updated note.",
		Pinned:   false,
		Priority: 3,
	}
	resp, err := server.UpdateNote(ctx, req)
	if err != nil {
		t.Errorf("UpdateNote failed with error: %v", err)
	}
	if resp == nil || resp.Title != "Updated Note" {
		t.Errorf("UpdateNote response is incorrect")
	}
}

func TestDeleteNote(t *testing.T) {
	server := &noteGRPCServer{
		cfg:           nil,
		uc:            &mockNoteUseCase{},
		authDomainSvc: &mockAuthDomainService{},
		log:           nil,
	}

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{"authorization": "bearer 97f1950d-7120-451c-8545-615bb845b78c"}),
	)

	// Test case: Valid delete note request
	req := &gen.DeleteNoteRequest{
		Id: uuid.New().String(),
	}
	_, err := server.DeleteNote(ctx, req)
	if err != nil {
		t.Errorf("DeleteNote failed with error: %v", err)
	}
}

// --------------------------------------------Mocks for dependencies--------------------------------------------

type mockNoteUseCase struct{}

func (m *mockNoteUseCase) ListNotes(ctx context.Context, params *domain.ListNotesParams) ([]*domain.Note, error) {
	return []*domain.Note{
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			Title:    "Test Note 1",
			Content:  "This is a test note.",
			Pinned:   true,
			Priority: 1,
		},
		{
			ID:       uuid.New(),
			UserID:   uuid.New(),
			Title:    "Test Note 2",
			Content:  "Another test note.",
			Pinned:   false,
			Priority: 2,
		},
	}, nil
}

func (m *mockNoteUseCase) GetNote(ctx context.Context, userID uuid.UUID, noteID uuid.UUID) (*domain.Note, error) {
	return &domain.Note{
		ID:       noteID,
		UserID:   userID,
		Title:    "Test Note",
		Content:  "This is a test note.",
		Pinned:   false,
		Priority: 1,
	}, nil
}

func (m *mockNoteUseCase) CreateNote(ctx context.Context, params *domain.CreateNoteParams) (*domain.Note, error) {
	return &domain.Note{
		ID:       uuid.New(),
		UserID:   params.UserID,
		Title:    params.Title,
		Content:  params.Content,
		Pinned:   params.Pinned,
		Priority: params.Priority,
	}, nil
}

func (m *mockNoteUseCase) UpdateNote(ctx context.Context, params *domain.UpdateNoteParams) (*domain.Note, error) {
	return &domain.Note{
		ID:       params.ID,
		UserID:   params.UserID,
		Title:    params.Title,
		Content:  params.Content,
		Pinned:   params.Pinned,
		Priority: params.Priority,
	}, nil
}

func (m *mockNoteUseCase) DeleteNote(ctx context.Context, userID uuid.UUID, noteID uuid.UUID) error {
	return nil
}

type mockAuthDomainService struct{}

func (m *mockAuthDomainService) GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error) {
	return uuid.New(), nil
}
