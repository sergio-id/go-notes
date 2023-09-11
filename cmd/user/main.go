package main

import (
	"context"
	"fmt"
	"github.com/sergio-id/go-notes/pkg/postgres"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sergio-id/go-notes/cmd/user/config"
	"github.com/sergio-id/go-notes/internal/user/app"
	"github.com/sergio-id/go-notes/pkg/logger"
	"go.uber.org/automaxprocs/maxprocs"
	"google.golang.org/grpc"
)

func main() {
	log.Println("🚀 Starting user service...")

	// set GOMAXPROCS
	_, err := maxprocs.Set()
	if err != nil {
		log.Fatalf("Failed set max procs: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	appLog := logger.NewAppLogger(cfg.Logger)
	appLog.InitLogger()
	appLog.Named(cfg.App.Name)
	appLog.Infof("CFG APP: %#v", cfg)

	server := grpc.NewServer()

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	_, cleanup, err := app.InitApp(cfg, postgres.DBConnString(cfg.PG.DsnURL), server, appLog) //wire
	if err != nil {
		appLog.Errorf("Failed init app: %v", err)
		cancel()
	}

	// gRPC Server.
	address := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		appLog.Errorf("Failed to listen to address, network: %s, address: %s, err: %v", network, address, err)
		cancel()
	}

	appLog.Infof("🌏 Start server...", "address", address)

	defer func() {
		if err1 := l.Close(); err1 != nil {
			appLog.Errorf("Failed to close a listener, network: %s, address: %s, err: %v", network, address, err1)
		}
	}()

	err = server.Serve(l)
	if err != nil {
		appLog.Errorf("Failed start gRPC server, network: %s, address: %s, err: %v", network, address, err)
		cancel()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		cleanup()
		appLog.Infof("Signal.Notify: %v", v)
	case done := <-ctx.Done():
		cleanup()
		appLog.Infof("ctx.Done: %v", done)
	}
}
