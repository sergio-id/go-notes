package app

import (
	"github.com/sergio-id/go-notes/cmd/auth/config"
	"github.com/sergio-id/go-notes/internal/auth/domain"
	"github.com/sergio-id/go-notes/pkg/logger"
	pkgRedis "github.com/sergio-id/go-notes/pkg/redis"
	"github.com/sergio-id/go-notes/proto/gen"
)

type App struct {
	Cfg            *config.Config
	Redis          pkgRedis.RedisEngine
	UC             domain.SessionUseCase
	UserDomainSvc  domain.UserDomainService
	AuthGRPCServer gen.AuthServiceServer
	Log            logger.Logger
}

func New(
	cfg *config.Config,
	redis pkgRedis.RedisEngine,
	uc domain.SessionUseCase,
	userDomainSvc domain.UserDomainService,
	authGRPCServer gen.AuthServiceServer,
	log logger.Logger,
) *App {
	return &App{
		Cfg:            cfg,
		Redis:          redis,
		UC:             uc,
		UserDomainSvc:  userDomainSvc,
		AuthGRPCServer: authGRPCServer,
		Log:            log,
	}
}
