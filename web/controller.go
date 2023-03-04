package web

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	_ "golang.org/x/oauth2"

	"github.com/bindernews/sts-msr/orm"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/patrickmn/go-cache"
)

// Generic unauthorized error
var ErrUnauthorized = errors.New("unauthorized")

var ErrRunAlreadyUploaded = errors.New("run already uploaded")

const (
	// gin Context key for the user's email
	CtxEmail = "Email"
	// gin context key for the string cache
	CtxStrCache = "StrCache"
	// gin context key for the go-cache instance
	CtxCache = "gocache"
	// key for the run data
	ctxRunData = "run-data"
)

type MainController struct {
	Srv      *Services
	strcache StrCache
	gcache   *cache.Cache
}

func (s *MainController) Init(r *gin.Engine) error {
	db := orm.New(s.Srv.Pool)
	s.strcache = NewStrCache(db.StrCacheToId, db.StrCacheAdd)
	s.gcache = cache.New(5*time.Minute, 10*time.Minute)

	// Set the gin run mode
	if s.Srv.Config.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(sessions.Sessions("main", s.Srv.SeStore))
	r.LoadHTMLGlob("templates/*.html")

	// Make directory to store runs in
	os.MkdirAll(s.Srv.Config.RunsDir, fs.FileMode(0755))

	// Register main routes
	r.RouterGroup.
		POST("/upload", s.PostUpload, s.handleUpload).
		POST("/upload-file", s.PostUploadFile, s.handleUpload).
		GET("/ping", s.PingDB).
		GET("/", s.Srv.CtxSetEmail(), s.GetIndex)

	// Serve static files
	r.StaticFS("/static", gin.Dir("static", false))
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

func (s *MainController) PostUploadFile(c *gin.Context) {
	var runData RunSchemaJson
	// Parse file
	mpf, err := c.FormFile("run-file")
	if err != nil {
		AbortErr(c, 400, err)
		return
	}
	fd, err := mpf.Open()
	if err != nil {
		AbortErr(c, 400, err)
		return
	}
	defer fd.Close()
	if err := json.NewDecoder(fd).Decode(&runData); err != nil {
		AbortErr(c, 400, err)
		return
	}
	c.Set(ctxRunData, runData)
	c.Next()
}

func (s *MainController) PostUpload(c *gin.Context) {
	var runData RunSchemaJson
	// Parse data
	if err := c.BindJSON(&runData); err != nil {
		AbortErr(c, 400, err)
		return
	}
	c.Set(ctxRunData, runData)
	c.Next()
}

func (s *MainController) handleUpload(c *gin.Context) {
	runData := c.MustGet(ctxRunData).(RunSchemaJson)
	// Check it's not a duplicate
	// TODO
	// Store in file for later
	go s.saveRun(c, &runData)
	// Store in DB
	ctx := c.Request.Context()
	if err := s.Srv.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := runData.AddToDb(ctx, s.strcache, orm.New(tx))
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

func (s *MainController) saveRun(c *gin.Context, data *RunSchemaJson) {
	wr, err := os.Create(path.Join(s.Srv.Config.RunsDir, data.PlayId.String()+".run.gz"))
	if err != nil {
		c.Error(err)
		return
	}
	defer wr.Close()
	gzwr := gzip.NewWriter(wr)
	defer gzwr.Close()
	if err := json.NewEncoder(gzwr).Encode(data); err != nil {
		c.Error(err)
		return
	}
}
