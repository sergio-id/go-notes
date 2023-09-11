//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/sergio-id/go-notes/cmd/note/config"
	grpcNote "github.com/sergio-id/go-notes/internal/note/delivery/grpc"
	infrasGRPC "github.com/sergio-id/go-notes/internal/note/infrastructure/network"
	postgresRepoNote "github.com/sergio-id/go-notes/internal/note/infrastructure/repository/postgres"
	usecaseNote "github.com/sergio-id/go-notes/internal/note/usecase"
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

		grpcNote.NoteGRPCServerSet,
		postgresRepoNote.RepositorySet,
		usecaseNote.UseCaseSet,
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
