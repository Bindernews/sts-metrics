package web

import (
	"context"
	"os"

	"github.com/bindernews/sts-msr/orm"
	"github.com/bindernews/sts-msr/tools"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/jackc/pgx/v4"
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
	poolCfg, err := pgxpool.ParseConfig(os.Getenv(tools.EnvPostgresConn))
	if err != nil {
		return err
	}
	poolCfg.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		if err := (orm.CardSpec{}).RegisterType(ctx, c); err != nil {
			return err
		}
		return nil
	}
	s.Pool, err = pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return err
	}

	// Setup session
	s.SeStore = cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))

	// Initialize config
	s.Config = NewConfig()
	return nil
}
