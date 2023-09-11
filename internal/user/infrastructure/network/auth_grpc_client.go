package network

import (
	"context"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sergio-id/go-notes/cmd/user/config"
	"github.com/sergio-id/go-notes/internal/user/domain"
	gen "github.com/sergio-id/go-notes/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.AuthDomainService = (*authGRPCClient)(nil)

var AuthGRPCClientSet = wire.NewSet(NewGRPCAuthClient)

func NewGRPCAuthClient(cfg *config.Config) (domain.AuthDomainService, error) {
	conn, err := grpc.Dial(cfg.AuthClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &authGRPCClient{
		conn: conn,
	}, nil
}

func (p *authGRPCClient) GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error) {
	c := gen.NewAuthServiceClient(p.conn)

	res, err := c.GetSession(ctx, &gen.GetSessionRequest{Token: token})
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, "authGRPCClient-c.GetSession")
	}

	return uuid.Parse(res.UserId)
}
