package app

import (
	"github.com/sergio-id/go-notes/cmd/category/config"
	"github.com/sergio-id/go-notes/internal/category/domain"
	"github.com/sergio-id/go-notes/pkg/logger"
	"github.com/sergio-id/go-notes/pkg/postgres"
	"github.com/sergio-id/go-notes/proto/gen"
)

type App struct {
	Cfg                *config.Config
	PG                 postgres.DBEngine
	UC                 domain.CategoryUseCase
	CategoryGRPCServer gen.CategoryServiceServer
	AuthDomainSvc      domain.AuthDomainService
	Log                logger.Logger
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc domain.CategoryUseCase,
	categoryGRPCServer gen.CategoryServiceServer,
	authDomainSvc domain.AuthDomainService,
	log logger.Logger,
) *App {
	return &App{
		Cfg:                cfg,
		PG:                 pg,
		UC:                 uc,
		CategoryGRPCServer: categoryGRPCServer,
		AuthDomainSvc:      authDomainSvc,
		Log:                log,
	}
}
