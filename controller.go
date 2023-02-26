package stms

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	_ "golang.org/x/oauth2"

	"github.com/Masterminds/sprig/v3"
	"github.com/bindernews/sts-msr/chart"
	"github.com/bindernews/sts-msr/orm"
	"github.com/bindernews/sts-msr/util"
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
	// gin Context key for the user's email
	CtxEmail = "Email"
	// gin context key for the string cache
	CtxStrCache = "StrCache"
)

type MainController struct {
	Srv      *Services
	strcache StrCache
}

func (s *MainController) Init(r *gin.Engine) error {
	db := orm.New(s.Srv.Pool)
	s.strcache = NewStrCache(db.StrCacheToId, db.StrCacheAdd)

	r.Use(sessions.Sessions("main", s.Srv.SeStore))
	r.SetFuncMap(sprig.FuncMap())
	r.LoadHTMLGlob("templates/*.html")

	// Make directory to store runs in
	os.MkdirAll(s.Srv.Config.RunsDir, fs.FileMode(0755))

	// Register main routes
	r.RouterGroup.
		POST("/upload", s.PostUpload).
		GET("/pingdb", s.PingDB).
		GET("/test1", s.Test1).
		GET("/", s.Srv.CtxSetEmail(), s.GetIndex)

	// Serve static files
	r.StaticFS("/static", gin.Dir("static", false))

	// Load and register charts
	chartConv := new(chart.ConvertContext)
	chartConv.Init()
	chartConv.Predefs["character"] = chart.NewEnumParam(
		"Character",
		"character",
		"SELECT name FROM character_list",
		true,
	)
	if err := util.TryEach(s.Srv.Config.Charts, chartConv.Add); err != nil {
		return err
	}
	stats2 := r.Group("/stats2")
	stats2.Use(func(c *gin.Context) {
		c.Set(chart.CtxDb, s.Srv.Pool)
		c.Next()
	})
	for k, c := range chartConv.Charts {
		stats2.GET(k, c.Handle)
	}

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

func (s *MainController) Test1(c *gin.Context) {
	rows, err := s.Srv.Pool.Query(
		context.Background(),
		`SELECT row_to_json(stats_overview.*) FROM stats_overview WHERE "name" = $1`,
		"GRACKLE",
	)
	if err != nil {
		AbortErr(c, 500, err)
		return
	}
	data := make([]any, 0)
	for rows.Next() {
		if r, err := rows.Values(); err != nil {
			AbortErr(c, 500, err)
			return
		} else {
			data = append(data, r)
		}
	}
	c.JSON(200, data)
}

func (s *MainController) PostUpload(c *gin.Context) {
	var runData RunSchemaJson
	// Parse data
	if err := c.BindJSON(&runData); err != nil {
		AbortErr(c, 400, err)
		return
	}
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
	if err == nil {
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

func AbortErr(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, c.Error(err).JSON())
}
