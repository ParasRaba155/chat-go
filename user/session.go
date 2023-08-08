package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	basePrefix = "blog-session:"
)

var (
	ErrSessionDoesNotExist = errors.New("session does not exist")
	ErrRetrivingSession    = errors.New("could not retrieve the session")
)

type sessionRepo struct {
	redis *redis.Client
}

type session struct {
	Email     string `json:"email" redis:"email"`
	SessionID string `json:"sessionID" redis:"sessionID"`
}

func NewSessionRepo(rc *redis.Client) sessionRepo {
	return sessionRepo{
		redis: rc,
	}
}

func (s sessionRepo) CreateSession(ctx context.Context, email string, expireMinutes int) (string, error) {
	sess := session{
		Email:     email,
		SessionID: uuid.NewString(),
	}
	sessKeys := s.createKey(sess.SessionID)

	sessBytes, err := json.Marshal(sess)
	if err != nil {
		return "", fmt.Errorf("session.CreateSession json Marshal : %w", err)
	}
	err = s.redis.Set(ctx, sessKeys, sessBytes, time.Minute*time.Duration(expireMinutes)).Err()
	if err != nil {
		return "", fmt.Errorf("session.CreateSession redisClient Set : %w", err)
	}
	return sessKeys, nil
}

func (s sessionRepo) GetSessionByID(ctx context.Context, id string) (session, error) {
	out := s.redis.Get(ctx, id)
	err := out.Err()
	if errors.Is(err, redis.Nil) {
		return session{}, fmt.Errorf("%w: session id: %s", ErrSessionDoesNotExist, id)
	}
	sessBytes, err := out.Bytes()
	if err != nil {
		return session{}, fmt.Errorf("%w: %w: session id: %s", ErrRetrivingSession, err, id)
	}
	var sess session
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return session{}, fmt.Errorf("%w: %w: session id: %s", ErrRetrivingSession, err, id)
	}
	return sess, nil
}

func (s sessionRepo) DeleteSessionByID(ctx context.Context, id string) error {
	if err := s.redis.Del(ctx, id).Err(); err != nil {
		return fmt.Errorf(":%w could not delete: session id: %s", err, id)
	}
	return nil
}

func (s sessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, sessionID)
}
