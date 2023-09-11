package app

import (
	"github.com/sergio-id/go-notes/cmd/note/config"
	"github.com/sergio-id/go-notes/internal/note/domain"
	"github.com/sergio-id/go-notes/pkg/logger"
	"github.com/sergio-id/go-notes/pkg/postgres"
	"github.com/sergio-id/go-notes/proto/gen"
)

type App struct {
	Cfg            *config.Config
	PG             postgres.DBEngine
	UC             domain.NoteUseCase
	NoteGRPCServer gen.NoteServiceServer
	Log            logger.Logger
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc domain.NoteUseCase,
	noteGRPCServer gen.NoteServiceServer,
	log logger.Logger,
) *App {
	return &App{
		Cfg:            cfg,
		PG:             pg,
		UC:             uc,
		NoteGRPCServer: noteGRPCServer,
		Log:            log,
	}
}
