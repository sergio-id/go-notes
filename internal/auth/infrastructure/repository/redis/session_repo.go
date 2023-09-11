package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sergio-id/go-notes/internal/auth/domain"
	pkgRedis "github.com/sergio-id/go-notes/pkg/redis"
)

const (
	basePrefix = "sessions:"
)

type sessionRepo struct {
	redis      pkgRedis.RedisEngine
	basePrefix string
}

var _ domain.SessionRepository = (*sessionRepo)(nil)

var RepositorySet = wire.NewSet(NewSessionRepo)

func NewSessionRepo(redis pkgRedis.RedisEngine) domain.SessionRepository {
	return &sessionRepo{redis: redis, basePrefix: basePrefix}
}

func (r *sessionRepo) CreateSession(ctx context.Context, arg *domain.CreateSessionParams) (*domain.Session, error) {
	session := &domain.Session{
		UserID:    arg.UserID.String(),
		SessionID: uuid.New().String(),
	}

	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		return nil, errors.Wrap(err, "sessionRepo.CreateSession.json.Marshal")
	}

	sessionKey := r.createKey(session.SessionID)
	err = r.redis.GetRedisClient().Set(ctx, sessionKey, sessionBytes, arg.Duration).Err()
	if err != nil {
		return nil, errors.Wrap(err, "sessionRepo.CreateSession.redisClient.Set")
	}

	return session, nil
}

func (r *sessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*domain.Session, error) {
	sessionBytes, err := r.redis.GetRedisClient().Get(ctx, r.createKey(sessionID)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.redisClient.Get")
	}

	session := &domain.Session{}
	if err = json.Unmarshal(sessionBytes, &session); err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.json.Unmarshal")
	}
	return session, nil
}

func (r *sessionRepo) DeleteByID(ctx context.Context, sessionID string) error {
	if err := r.redis.GetRedisClient().Del(ctx, r.createKey(sessionID)).Err(); err != nil {
		return errors.Wrap(err, "sessionRepo.DeleteByID")
	}
	return nil
}

func (r *sessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, sessionID)
}
