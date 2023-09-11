package main

import (
	"context"
	"errors"
	"fmt"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sergio-id/go-notes/cmd/proxy/config"
	"github.com/sergio-id/go-notes/pkg/logger"
	"github.com/sergio-id/go-notes/proto/gen"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strings"
	"time"
)

const _readHeaderTimeout = 1

func main() {
	log.Println("🚀starting notes service proxy...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	appLog := logger.NewAppLogger(cfg.Logger)
	appLog.InitLogger()
	appLog.Named(cfg.App.Name)
	appLog.Infof("CFG APP -: %#v", cfg)

	mux := http.NewServeMux()

	gw, err := newGateway(ctx, cfg, nil)
	if err != nil {
		appLog.Errorf("failed to create a new gateway: %s", err)
	}

	mux.Handle("/", gw)

	s := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           allowCORS(withLogger(mux)),
		ReadHeaderTimeout: time.Duration(_readHeaderTimeout) * time.Second,
	}

	go func() {
		<-ctx.Done()
		appLog.Infof("shutting down the http server")

		if err := s.Shutdown(context.Background()); err != nil {
			appLog.Errorf("failed to shutdown http server: %s", err)
		}
	}()

	appLog.Infof("start listening... Address: %s", fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port))

	if err := s.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		appLog.Errorf("failed to listen and serve: %s", err)
	}
}

// newGateway returns a new gateway server which translates HTTP into gRPC.
func newGateway(
	ctx context.Context,
	cfg *config.Config,
	opts []gwruntime.ServeMuxOption,
) (http.Handler, error) {
	authEndpoint := cfg.GRPC.AuthURL
	noteEndpoint := cfg.GRPC.NoteURL
	categoryEndpoint := cfg.GRPC.CategoryURL
	userEndpoint := cfg.GRPC.UserURL

	mux := gwruntime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := gen.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, authEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}

	err = gen.RegisterNoteServiceHandlerFromEndpoint(ctx, mux, noteEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}

	err = gen.RegisterCategoryServiceHandlerFromEndpoint(ctx, mux, categoryEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}

	err = gen.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}

	return mux, nil
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	headers := []string{"*"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	slog.Info("preflight request", "http_path", r.URL.Path)
}

func withLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Run request", "http_method", r.Method, "http_url", r.URL)

		h.ServeHTTP(w, r)
	})
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
