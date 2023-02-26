package stsweb

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/bindernews/sts-msr/chart"
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
	var char_list []string
	if err := listChars(qq, &char_list); err != nil {
		return err
	}
	// Default character if none specified
	if char_name == "" && len(char_list) > 0 {
		char_name = char_list[0]
	}
	stats, err := qq.StatsGetOverview(ctx, char_name)
	if err != nil {
		c.Error(err)
		return fmt.Errorf("unknown character: %s", char_name)
	}

	const percentiles = "(p25, p50, p75)"
	data := []struct {
		Name  string
		Value any
	}{
		{"Character", stats.Name},
		{"Runs", stats.Runs},
		{"Wins", stats.Wins},
		{"Avg Win Rate", stats.AvgWinRate},
		{"Deck Size Quartiles " + percentiles, ListJoin(stats.PDeckSize, ", ")},
		{"Floor Reached Quartiles " + percentiles, ListJoin(stats.PFloorReached, ", ")},
	}

	tmpl := chart.ChartTemplate{
		Title: "Stats Overview",
		Options: chart.TabularOpts{
			Layout: "fitColumns",
			Columns: []chart.TabularColumn{
				{Title: "Name", Field: "Name"},
				{Title: "Value", Field: "Value"},
			},
		},
		Filters: []chart.Filter{
			{
				Label:   "Character",
				Name:    "character",
				Type:    "select",
				Value:   char_name,
				Options: char_list,
			},
		},
		ChartCode: template.JS(`
		function showChart() { opts.data = rawdata; new Tabulator('#table', opts); }
		`),
		RawData: lo.ToAnySlice(data),
		IsTable: true,
	}
	c.HTML(200, "chartview.html", tmpl)
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

func listChars(qq *orm.Queries, outList *[]string) error {
	list, err := qq.StatsListCharacters(context.Background())
	if err != nil {
		return err
	}
	*outList = lo.Map(list, func(v orm.CharacterList, _ int) string {
		return v.Name
	})
	return nil
}

type StatsTableTempl struct {
	Headers []string
	Rows    [][]any
}

func ToStr[T any](v T, _ int) string {
	return fmt.Sprint(v)
}

func ListJoin[T any](co []T, sep string) string {
	return strings.Join(lo.Map(co, ToStr[T]), sep)
}
