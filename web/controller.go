package web

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"

	_ "golang.org/x/oauth2"

	"github.com/bindernews/sts-msr/orm"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/samber/lo"
)

var (
	// Generic unauthorized error
	ErrUnauthorized = errors.New("unauthorized")
	// Error when a run with the same play_id already exists
	ErrRunAlreadyUploaded = errors.New("run already uploaded")
	// Basically import alias for lo.Ternary
	tern = lo.Ternary[gin.HandlerFunc]
)

const (
	// gin Context key for the user's email
	CtxEmail = "Email"
	// gin context key for the string cache
	CtxStrCache = "StrCache"
	// context key for the database pool
	CtxDbPool = "DbPool"
	// gin context key for the go-cache instance
	CtxCache = "gocache"
	// key for the run data
	ctxRunData = "run-data"
	// play_id of the run for easy access and in case validation fails
	ctxPlayId = "play-id"
	// body data as in-memory byte array
	ctxBodyBytes = "body-bytes"
)

type StrCache = DbCache[string]

type MainController struct {
	Srv    *Services
	ormCtx *OrmContext
}

func (s *MainController) Init(r *gin.Engine) error {
	cfg := s.Srv.Config
	db := orm.New(s.Srv.Pool)

	s.ormCtx = &OrmContext{
		Sc: NewDbCache(db.StrCacheToId, db.StrCacheAdd),
		Cc: NewDbCache(db.CardSpecToId, db.CardSpecAdd),
	}

	// Set the gin run mode
	if cfg.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(sessions.Sessions("main", s.Srv.SeStore))
	r.LoadHTMLGlob("templates/*.html")

	// Make directory to store runs in
	if cfg.Upload.SaveRawToDisk {
		os.MkdirAll(cfg.Upload.RunsDir, fs.FileMode(0755))
	}

	// Register main routes
	g := r.Group(strings.TrimSuffix(cfg.BasePath, "/"))
	g.Use(s.CtxInject)

	// Create proxy for stats server
	if cfg.Stats.Upstream != "" {
		targetUrl, err := url.Parse(cfg.Stats.Upstream)
		if err != nil {
			return err
		}
		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		statsPrefix := g.BasePath() + cfg.Stats.Route
		g.Any(cfg.Stats.Route, HandlerChain(
			tern(cfg.Stats.Auth, s.authScopes([]string{"stats:view"}), nil),
			StripRequestPrefix(statsPrefix),
			gin.WrapH(proxy),
		)...)
	}

	// Create upload handler
	g.POST(cfg.Upload.Route, HandlerChain(
		s.postUploadParse,
		tern(cfg.Upload.SaveRawToDb || cfg.Upload.SaveRawToDisk, s.archiveRawData, nil),
		tern(cfg.Upload.StoreToDb, s.storeToDb, nil),
		func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Thank you!"})
		},
	)...)

	// Add getrun route
	if cfg.GetRun.Route != "" {
		g.GET(cfg.GetRun.Route, HandlerChain(
			tern(cfg.GetRun.Auth, s.authScopes([]string{"getrun"}), nil),
			s.GetRunJson,
		)...)
	}

	// Health check route
	if cfg.HealthRoute != "" {
		g.GET(cfg.HealthRoute, s.healthCheck)
	}

	if cfg.DebugMode {
		g.GET("/", s.GetIndex)
	}
	return nil
}

func (s *MainController) GetIndex(c *gin.Context) {
	email := c.GetString(CtxEmail)
	c.HTML(200, "index.html", gin.H{
		"Email": email,
	})
}

func (s *MainController) healthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	if err := s.Srv.Pool.Ping(ctx); err != nil {
		AbortMsg(c, 500, err)
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

// Parse body into RunSchemaJson
func (s *MainController) postUploadParse(c *gin.Context) {
	// Read the body since we parse it multiple times
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.Request.Body.Close()
	c.Set(ctxBodyBytes, body)

	var runData RunSchemaJson
	if err := binding.JSON.BindBody(body, &runData); err != nil {
		AbortMsg(c, 400, err)
		return
	}
	c.Set(ctxRunData, runData)
	c.Set(ctxPlayId, runData.PlayId.String())
}

func (s *MainController) archiveRawData(c *gin.Context) {
	body := c.MustGet(ctxBodyBytes).([]byte)
	playId := c.MustGet(ctxPlayId).(string)
	params := orm.ArchiveAddParams{
		Bdata:  pgtype.JSON{Bytes: body, Status: pgtype.Present},
		PlayID: playId,
	}
	cfg := s.Srv.Config
	if cfg.Upload.SaveRawToDb {
		if err := s.saveRawToDb(c, params); err != nil {
			c.Error(err)
		}
	}
	if cfg.Upload.SaveRawToDisk {
		if err := s.saveRawToFile(c, params); err != nil {
			c.Error(err)
		}
	}
}

// Save the raw contents of the json body into the database
func (s *MainController) saveRawToDb(c *gin.Context, req orm.ArchiveAddParams) error {
	ctx := context.Background()
	if err := orm.New(s.Srv.Pool).ArchiveAdd(ctx, req); err != nil {
		return err
	}
	return nil
}

// Save the raw contents of the json body to the disk
func (s *MainController) saveRawToFile(c *gin.Context, req orm.ArchiveAddParams) error {
	fpath := path.Join(s.Srv.Config.Upload.RunsDir, req.PlayID+".run")
	if wr, err := CreateNewFile(fpath); err == nil {
		defer wr.Close()
		if _, err := wr.Write(req.Bdata.Bytes); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (s *MainController) storeToDb(c *gin.Context) {
	ctx := c.Request.Context()
	runData := c.MustGet(ctxRunData).(RunSchemaJson)
	oc := s.ormCtx.Copy()
	// Store in DB
	err := s.Srv.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := runData.AddToDb(ctx, oc, orm.New(tx))
		return err
	})
	if err != nil {
		// Duplicate play id is a bad request
		if strings.Contains(err.Error(), "\"runsdata_play_id_key\"") {
			AbortMsg(c, 400, fmt.Errorf("%w - play_id = %s", ErrRunAlreadyUploaded, runData.PlayId))
		} else {
			c.AbortWithError(500, err)
		}
	}
}

func (s *MainController) GetRunJson(c *gin.Context) {
	ctx := c.Request.Context()
	db := orm.New(s.Srv.Pool)

	var params struct {
		PlayId string `form:"play_id"`
	}
	if err := c.BindQuery(&params); err != nil {
		AbortMsg(c, 400, err)
	}
	data, err := RunToJson(ctx, db, params.PlayId)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	c.JSON(200, data)
}

// Middleware that sets the CtxEmail value for the context, regardless of if
// the user is authenticated or not. If the user is not logged in, sets to the empty string.
func (s *MainController) CtxInject(c *gin.Context) {
	email := c.GetHeader("X-Email")
	c.Set(CtxEmail, email)
	c.Set(CtxDbPool, s.Srv.Pool)
}

// Returns middleware that checks if the user has the required scopes.
// The user must have ALL scopes listed in the array to be allowed access.
// If the user is not authenticated, they will be denied access.
//
// Uses email from CtxEmail
func (s *MainController) authScopes(scopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const query = `SELECT user_has_scopes($1,$2)`
		ctx := c.Request.Context()
		// Attempt to get the user's email address
		email := c.GetString(CtxEmail)
		if email == "" {
			AbortMsg(c, 403, ErrUnauthorized)
			return
		}
		ok := false
		err := s.Srv.Pool.QueryRow(ctx, query, email, scopes).Scan(&ok)
		if err != nil || !ok {
			AbortMsg(c, 403, ErrUnauthorized)
			return
		}
	}
}
