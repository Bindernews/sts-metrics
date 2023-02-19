package stms

import (
	"errors"
	"fmt"
	"strings"

	_ "golang.org/x/oauth2"

	"github.com/bindernews/sts-msr/orm"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// Generic unauthorized error
var ErrUnauthorized = errors.New("unauthorized")

// Error for unknown chart id
var ErrUnknownChart = errors.New("unknown chart")
var ErrRunAlreadyUploaded = errors.New("run already uploaded")

const (
	// gin Context key for the user's email, set in
	CtxEmail = "Email"
)

type MainController struct {
	Srv *Services
}

func (s *MainController) Init(r *gin.Engine) error {
	r.Use(sessions.Sessions("main", s.Srv.SeStore))
	r.LoadHTMLGlob("templates/*.html")

	// Register main routes
	r.RouterGroup.
		POST("/upload", s.PostUpload).
		GET("/pingdb", s.PingDB).
		GET("/", s.Srv.CtxSetEmail(), s.GetIndex)

	return nil
}

func (s *MainController) GetIndex(c *gin.Context) {
	email := c.GetString(CtxEmail)
	c.HTML(200, "index.html", gin.H{
		"Email": email,
	})
}

func (s *MainController) PingDB(c *gin.Context) {
	ctx := c.Request.Context()
	if err := s.Srv.Pool.Ping(ctx); err != nil {
		AbortErr(c, 500, err)
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func (s *MainController) PostUpload(c *gin.Context) {
	var runData RunSchemaJson
	if err := c.BindJSON(&runData); err != nil {
		AbortErr(c, 400, err)
		return
	}

	ctx := c.Request.Context()
	if err := s.Srv.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := runData.AddToDb(ctx, orm.New(tx))
		return err
	}); err != nil {
		// Duplicate play id is a bad request
		if strings.Contains(err.Error(), "\"runsdata_play_id_key\"") {
			AbortErr(c, 400, fmt.Errorf("%w - play_id = %s", ErrRunAlreadyUploaded, runData.PlayId))
			return
		}
		c.Error(err)
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, gin.H{"message": "Thank you!"})
}

func AbortErr(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, c.Error(err).JSON())
}

// Returns the value of the session key, or an empty string if
// the value doesn't exist, or is not a string.
func sessGetString(s sessions.Session, key string) string {
	val := s.Get(key)
	if val == nil {
		return ""
	} else if val2, ok := val.(string); ok {
		return val2
	} else {
		return ""
	}
}
