package conn

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"app/config"
)

// NewRedisClient will create a redis client and ping the redis and return any error while pinging
func NewRedisClient(cfg *config.Redis, l *zap.Logger) (redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:        cfg.Address,
		Username:    cfg.Username,
		Password:    cfg.Password,
		DB:          cfg.DB,
		MaxRetries:  cfg.MaxRetries,
		PoolSize:    cfg.PoolSize,
		PoolTimeout: cfg.PoolTimeout,
	})
	err := rc.Ping(context.TODO()).Err()
	if err != nil {
		l.Fatal("error in connecting redis", zap.Error(err))
		return *rc, err
	}
	l.Info("successfully connected to redis")
	return *rc, err
}
