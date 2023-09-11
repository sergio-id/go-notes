package app

import (
	"github.com/sergio-id/go-notes/cmd/user/config"
	"github.com/sergio-id/go-notes/internal/user/domain"
	"github.com/sergio-id/go-notes/pkg/logger"
	"github.com/sergio-id/go-notes/pkg/postgres"
	"github.com/sergio-id/go-notes/proto/gen"
)

type App struct {
	Cfg            *config.Config
	PG             postgres.DBEngine
	UC             domain.UserUseCase
	UserGRPCServer gen.UserServiceServer
	AuthDomainSvc  domain.AuthDomainService
	Log            logger.Logger
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc domain.UserUseCase,
	userGRPCServer gen.UserServiceServer,
	authDomainSvc domain.AuthDomainService,
	log logger.Logger,
) *App {
	return &App{
		Cfg:            cfg,
		PG:             pg,
		UC:             uc,
		UserGRPCServer: userGRPCServer,
		AuthDomainSvc:  authDomainSvc,
		Log:            log,
	}
}
