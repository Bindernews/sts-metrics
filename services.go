package stms

import (
	"context"
	"os"

	"github.com/bindernews/sts-msr/tonoauth"
	"github.com/bindernews/sts-msr/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	pgxstdlib "github.com/jackc/pgx/v4/stdlib"
	csrf "github.com/utrack/gin-csrf"
)

type Services struct {
	Pool      *pgxpool.Pool
	SeStore   sessions.Store
	CsrfGuard gin.HandlerFunc
	Config    *Config
}

func (s *Services) LoadDefaults() error {
	// Connect to DB
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, os.Getenv("POSTGRES_CONN"))
	if err != nil {
		return err
	}
	s.Pool = pool

	// Setup sessions
	db := pgxstdlib.OpenDB(*s.Pool.Config().ConnConfig)
	store, err := postgres.NewStore(db, []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		return err
	}
	s.SeStore = store

	s.CsrfGuard = csrf.Middleware(csrf.Options{
		Secret: os.Getenv("CSRF_SECRET"),
	})

	s.Config = NewConfig()

	return nil
}

// Returns middleware that checks if the user has the required scopes before invoking
// the next handler. The user must have ALL scopes listed in the array to be allowed
// access. If the user is not authenticated, they will be denied access.
func (s *Services) AuthRequireScopes(scopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		ctx := c.Request.Context()

		// Attempt to get the user's email address
		email := util.SessGetString(sess, tonoauth.KeyUserEmail)
		if email == "" {
			AbortErr(c, 403, ErrUnauthorized)
			return
		}

		// Get scopes
		var count int
		row := s.Pool.QueryRow(ctx, "SELECT user_has_scopes($1, $2) AS cn", email, scopes)
		if err := row.Scan(&count); err != nil {
			c.Error(err)
			AbortErr(c, 403, ErrUnauthorized)
			return
		}
		if count != len(scopes) {
			AbortErr(c, 403, ErrUnauthorized)
			return
		}
		// Success! Put the email in the context for convenience.
		c.Set(CtxEmail, email)
		c.Next()
	}
}

// Returns middleware that sets the CtxEmail value for the context, regardless of if
// the user is authenticated or not. If the user is not logged in, sets to the empty string.
func (s *Services) CtxSetEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		email := util.SessGetString(sess, tonoauth.KeyUserEmail)
		c.Set(CtxEmail, email)
		c.Next()
	}
}
