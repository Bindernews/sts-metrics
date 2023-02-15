package stms

import (
	"context"
	"os"

	"github.com/bindernews/sts-msr/orm"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type StsMs struct {
	Pool *pgxpool.Pool
}

func NewStsMs() (*StsMs, error) {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, os.Getenv("POSTGRES_CONN"))
	if err != nil {
		return nil, err
	}
	s := &StsMs{
		Pool: pool,
	}
	return s, nil
}

func (s *StsMs) Register(group *gin.RouterGroup) {
	group.POST("/upload", s.PostUpload)
	group.GET("/pingdb", s.PingDB)
}

func (s *StsMs) PingDB(c *gin.Context) {
	ctx := c.Request.Context()
	if err := s.Pool.Ping(ctx); err != nil {
		c.AbortWithStatusJSON(500, c.Error(err).JSON())
		return
	}
	c.JSON(200, gin.H{"message": "Success"})
}

func (s *StsMs) PostUpload(c *gin.Context) {
	var runData RunSchemaJson
	if err := c.BindJSON(&runData); err != nil {
		abortBadReq(c, err)
		return
	}

	ctx := c.Request.Context()
	if err := s.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := runData.AddToDb(ctx, orm.New(tx))
		return err
	}); err != nil {
		abortBadReq(c, err)
		return
	}
	c.JSON(200, gin.H{"message": "Thank you!"})
}

func abortBadReq(c *gin.Context, err error) {
	c.AbortWithStatusJSON(400, c.Error(err).JSON())
}
