package stms

import (
	"fmt"

	"github.com/bindernews/sts-msr/orm"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type StatsView struct {
	Srv    *Services
	charts map[string]gin.HandlerFunc
}

func (s *StatsView) AddChart(name string, fn gin.HandlerFunc) {
	if s.charts == nil {
		s.charts = make(map[string]gin.HandlerFunc)
	}
	s.charts[name] = fn
}

func (s *StatsView) DefaultCharts() *StatsView {
	dv := &DefaultViews{}
	s.AddChart("overview", s.Wrap(dv.OverviewView))
	return s
}

func (s *StatsView) Queries() *orm.Queries {
	return orm.New(s.Srv.Pool)
}

func (s *StatsView) Wrap(fn func(c *gin.Context, qq *orm.Queries) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c, s.Queries()); err != nil {
			AbortErr(c, 500, err)
		}
	}
}

func (s *StatsView) Init(group *gin.RouterGroup) error {
	// Register charts
	group.Use(s.Srv.AuthRequireScopes([]string{"stats:view"}))
	for name, fn := range s.charts {
		group.GET(name, fn)
	}
	return nil
}

type DefaultViews struct{}

func (*DefaultViews) OverviewView(c *gin.Context, qq *orm.Queries) error {
	ctx := c.Request.Context()
	char_name := c.Query("character")
	char_id, err := qq.GetStr(ctx, char_name)
	if err == nil {
		return err
	}
	stats, err := qq.StatsGetOverall(ctx, char_id)
	if err != nil {
		c.Error(err)
		return fmt.Errorf("unknown character: %s", char_name)
	}
	tmpl := StatsTableTempl{
		Headers: []string{"Name", "Value"},
		Data:    lo.Map(stats, ToAny[orm.StatsGetOverallRow]),
	}
	c.HTML(200, "chart_table.html", tmpl)
	return nil
}

type StatsTableTempl struct {
	Headers []string
	Data    []any
}

func ToAny[T any](v T, _ int) any {
	return v
}
