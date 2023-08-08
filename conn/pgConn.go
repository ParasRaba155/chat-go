// conn to various database
package conn

import (
	"context"
	"fmt"

	"app/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func getDBURL(c *config.Database) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		c.Host, c.Port, c.User, c.Password, c.DefaultDB, c.SSLMODE, c.TimeoutSecs,
	)
}

func ConnectToPG(c *config.Database, l *zap.Logger) (*pgxpool.Pool, error) {
	dbURL := getDBURL(c)

	pool, err := pgxpool.New(context.TODO(), dbURL)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error connecting to db : %w", err)
	}
	l.Info("database connected successfully")
	return pool, nil
}
