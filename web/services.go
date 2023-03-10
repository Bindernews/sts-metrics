package web

import (
	"context"
	"os"

	"github.com/bindernews/sts-msr/tools"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Services struct {
	Pool    *pgxpool.Pool
	SeStore sessions.Store
	Config  *Config
}

func (s *Services) LoadDefaults() error {
	// Connect to DB
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, os.Getenv(tools.EnvPostgresConn))
	if err != nil {
		return err
	}
	s.Pool = pool
	// Setup session
	s.SeStore = cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	// Initialize config
	s.Config = NewConfig()
	return nil
}
