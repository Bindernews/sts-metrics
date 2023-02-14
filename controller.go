package stms

import "github.com/gin-gonic/gin"

type StsMs struct {
}

func NewStsMs() *StsMs {
	return &StsMs{}
}

func (s *StsMs) Register(group *gin.RouterGroup) {
	group.POST("/upload", s.PostUpload)
}

func (s *StsMs) PostUpload(c *gin.Context) {
	var runData RunSchemaJson
	if err := c.BindJSON(&runData); err != nil {
		abortBadReq(c, err)
		return
	}
	// TODO convert to SQL request
}

func abortBadReq(c *gin.Context, err error) {
	c.AbortWithStatusJSON(400, c.Error(err).JSON())
}
