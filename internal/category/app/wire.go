//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/sergio-id/go-notes/cmd/category/config"
	grpcCategory "github.com/sergio-id/go-notes/internal/category/delivery/grpc"
	infrasGRPC "github.com/sergio-id/go-notes/internal/category/infrastructure/network"
	postgresRepoCategory "github.com/sergio-id/go-notes/internal/category/infrastructure/repository/postgres"
	usecaseCategory "github.com/sergio-id/go-notes/internal/category/usecase"
	"github.com/sergio-id/go-notes/pkg/logger"
	"github.com/sergio-id/go-notes/pkg/postgres"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	grpcServer *grpc.Server,
	log logger.Logger,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,

		grpcCategory.CategoryGRPCServerSet,
		postgresRepoCategory.RepositorySet,
		usecaseCategory.UseCaseSet,
		infrasGRPC.AuthGRPCClientSet,
	))
}

func dbEngineFunc(url postgres.DBConnString, log logger.Logger) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("Connected to postgres: %s", url)

	return db, func() { db.Close() }, nil
}
