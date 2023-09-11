//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/sergio-id/go-notes/cmd/auth/config"
	grpcAuth "github.com/sergio-id/go-notes/internal/auth/delivery/grpc"
	infrasGRPC "github.com/sergio-id/go-notes/internal/auth/infrastructure/network"
	redisRepoSession "github.com/sergio-id/go-notes/internal/auth/infrastructure/repository/redis"
	usecaseSession "github.com/sergio-id/go-notes/internal/auth/usecase"
	"github.com/sergio-id/go-notes/pkg/logger"
	pkgRedis "github.com/sergio-id/go-notes/pkg/redis"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	cfgRedis pkgRedis.Config,
	grpcServer *grpc.Server,
	log logger.Logger,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		redisEngineFunc,

		grpcAuth.AuthGRPCServerSet,
		redisRepoSession.RepositorySet,
		usecaseSession.UseCaseSet,
		infrasGRPC.UserGRPCClientSet,
	))
}

func redisEngineFunc(cfgRedis pkgRedis.Config, log logger.Logger) (pkgRedis.RedisEngine, func(), error) {
	r, err := pkgRedis.NewRedisClient(cfgRedis)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("Redis connected on %s", cfgRedis.Addr)

	return r, func() { r.Close() }, nil
}
