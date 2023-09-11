// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/sergio-id/go-notes/cmd/note/config"
	grpc2 "github.com/sergio-id/go-notes/internal/note/delivery/grpc"
	"github.com/sergio-id/go-notes/internal/note/infrastructure/network"
	postgres2 "github.com/sergio-id/go-notes/internal/note/infrastructure/repository/postgres"
	"github.com/sergio-id/go-notes/internal/note/usecase"
	"github.com/sergio-id/go-notes/pkg/logger"
	"github.com/sergio-id/go-notes/pkg/postgres"
	"google.golang.org/grpc"
)

// Injectors from wire.go:

func InitApp(cfg *config.Config, dbConnStr postgres.DBConnString, grpcServer *grpc.Server, log logger.Logger) (*App, func(), error) {
	dbEngine, cleanup, err := dbEngineFunc(dbConnStr, log)
	if err != nil {
		return nil, nil, err
	}
	noteRepository := postgres2.NewNoteRepo(dbEngine)
	noteUseCase := usecase.NewNoteUseCase(noteRepository)
	authDomainService, err := network.NewGRPCAuthClient(cfg)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	noteServiceServer := grpc2.NewGRPCNoteServer(grpcServer, cfg, noteUseCase, authDomainService, log)
	app := New(cfg, dbEngine, noteUseCase, noteServiceServer, log)
	return app, func() {
		cleanup()
	}, nil
}

// wire.go:

func dbEngineFunc(url postgres.DBConnString, log logger.Logger) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("Connected to postgres: %s", url)

	return db, func() { db.Close() }, nil
}
