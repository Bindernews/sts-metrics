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
	s.AddChart("characters", s.Wrap(dv.ListCharacters))
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
	stats, err := qq.StatsGetOverview(ctx, char_name)
	if err != nil {
		c.Error(err)
		return fmt.Errorf("unknown character: %s", char_name)
	}

	percentiles := []string{"p25", "p50", "p75"}
	data := [][]any{
		{"Character", stats.Name},
		{"Runs", stats.Runs},
		{"Wins", stats.Wins},
		{"Avg Win Rate", stats.AvgWinRate},
	}
	data = append(data, lo.Map(stats.PDeckSize, func(v float32, i int) []any {
		return []any{"Deck Size " + percentiles[i], v}
	})...)
	data = append(data, lo.Map(stats.PFloorReached, func(v float32, i int) []any {
		return []any{"Floor Reached " + percentiles[i], v}
	})...)

	tmpl := StatsTableTempl{
		Headers: []string{"Name", "Value"},
		Rows:    data,
	}
	c.HTML(200, "chart_table.html", tmpl)
	return nil
}

func (*DefaultViews) ListCharacters(c *gin.Context, qq *orm.Queries) (err error) {
	var rows []orm.CharacterList
	ctx := c.Request.Context()
	if rows, err = qq.StatsListCharacters(ctx); err != nil {
		return
	}
	tmpl := StatsTableTempl{
		Headers: []string{"ID", "Name"},
		Rows: lo.Map(rows, func(r orm.CharacterList, _ int) []any {
			return []any{fmt.Sprint(r.ID), r.Name}
		}),
	}
	c.HTML(200, "chart_table.html", tmpl)
	return nil
}

type StatsTableTempl struct {
	Headers []string
	Rows    [][]any
}

func ToAny[T any](v T, _ int) any {
	return v
}
