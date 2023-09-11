package network

import (
	"context"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sergio-id/go-notes/cmd/auth/config"
	"github.com/sergio-id/go-notes/internal/auth/domain"
	gen "github.com/sergio-id/go-notes/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.UserDomainService = (*userGRPCClient)(nil)

var UserGRPCClientSet = wire.NewSet(NewGRPCUserClient)

func NewGRPCUserClient(cfg *config.Config) (domain.UserDomainService, error) {
	conn, err := grpc.Dial(cfg.UserClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &userGRPCClient{
		conn: conn,
	}, nil
}

func (p *userGRPCClient) GetUserByEmail(ctx context.Context, email string) (*domain.UserDTO, error) {
	c := gen.NewUserServiceClient(p.conn)

	res, err := c.GetUserByEmail(ctx, &gen.GetUserByEmailRequest{Email: email})
	if err != nil {
		return nil, errors.Wrap(err, "(userGRPCClient) c.GetUserByEmail ")
	}

	id, _ := uuid.Parse(res.Id)
	return &domain.UserDTO{
		ID:       id,
		Email:    res.Email,
		Password: res.Password,
	}, nil
}

func (p *userGRPCClient) CreateUser(ctx context.Context, arg *domain.CreateUserParams) (*domain.UserDTO, error) {
	c := gen.NewUserServiceClient(p.conn)

	res, err := c.CreateUser(ctx, &gen.CreateUserRequest{Email: arg.Email, Password: arg.Password})
	if err != nil {
		return nil, errors.Wrap(err, "(userGRPCClient) c.CreateUser ")
	}

	id, _ := uuid.Parse(res.Id)
	return &domain.UserDTO{
		ID:       id,
		Email:    res.Email,
		Password: res.Password,
	}, nil
}
