package chart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/bindernews/sts-msr/orm"
	"github.com/bindernews/sts-msr/util"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
)

const (
	// Gin context key for an instance of orm.DBTX
	CtxDb = "chartdb"
	// gin context key for go-cache instance
	CtxCache = "gocache"
)

// Error when a parameter has an invalid value
var ErrInvalidParam = errors.New("parameter value invalid")

const (
	chartJsCode = `
function showChart() { opts.data = %s(rawdata); const ctx = document.getElementById('chartc'); new Chart(ctx, opts); }`
	tableJsCode = `
function showChart() { opts.data = %s(rawdata); new Tabulator('#table', opts); }`
	rawJsCode = `
function showChart() { const elem = document.getElementById('raw'); elem.innerText = JSON.stringify(JSON.parse(%s(rawdata)), null, 2); }`
)

// Specification of a chart. These may be specified in the config file
// and are then parsed and converted into 'ChartSpec's.
type ChartSpec struct {
	// Chart name
	Name string
	// Chart URL-safe path
	path string
	// Query parameters, positional
	Params []IQueryParam
	// SQL statement to execute
	Sql string
	// Chart type
	Type string
	// JS code to transform the SQL output to appropriate format
	// for the chart or table library to use
	Transform string
	// Chart-specific options
	Options any
}

func (spec *ChartSpec) Path() string {
	if spec.path == "" {
		p := strings.ToLower(spec.Name)
		p = strings.ReplaceAll(p, " ", "_")
		spec.path = p
	}
	return spec.path
}

func (spec *ChartSpec) Handle(gctx *gin.Context) {
	c := newChartContext(gctx)
	// Parse query params
	args := make([]any, len(spec.Params))
	filters := make([]Filter, len(args))
	for i, param := range spec.Params {
		if a, err := param.Parse(c); err != nil {
			util.AbortErr(gctx, 400, err)
			return
		} else {
			c.Set(param.Name(), a)
			args[i] = a
			filters[i] = param.ToFilter(c)
		}
	}
	// Perform the query
	rows, err := c.Db().Query(c.Ctx, spec.Sql, args...)
	if err != nil {
		util.AbortErr(gctx, 500, err)
		return
	}
	// Load rows into JSON data
	data := make([]any, 0)
	for rows.Next() {
		var r json.RawMessage
		if err := rows.Scan(&r); err != nil {
			util.AbortErr(gctx, 500, err)
			return
		} else {
			data = append(data, r)
		}
	}

	// Render and serve HTML
	tmpl := ChartTemplate{
		Title:   spec.Name,
		Type:    spec.Type,
		Filters: filters,
		RawData: data,
		Options: spec.Options,
		IsTable: spec.Type == "table",
	}

	switch spec.Type {
	case "table":
		tmpl.ChartCode = template.JS(fmt.Sprintf(tableJsCode, spec.Transform))
	case "raw":
		tmpl.ChartCode = template.JS(fmt.Sprintf(rawJsCode, spec.Transform))
	case "bar":
		tmpl.ChartCode = template.JS(fmt.Sprintf(chartJsCode, spec.Transform))
	default:
		util.AbortErr(gctx, 500, errors.New("invalid chart type"))
		return
	}
	gctx.HTML(200, "chartview.html", tmpl)
}

type ChartContext struct {
	data  map[string]any
	db    orm.DBTX
	cache *cache.Cache
	gctx  *gin.Context
	Ctx   context.Context
}

func newChartContext(g *gin.Context) *ChartContext {
	return &ChartContext{
		data:  make(map[string]any),
		db:    g.MustGet(CtxDb).(orm.DBTX),
		cache: g.MustGet(CtxCache).(*cache.Cache),
		gctx:  g,
		Ctx:   g.Request.Context(),
	}
}
func (cx *ChartContext) Set(key string, value any) {
	cx.data[key] = value
}

// Retrieve values stored with Set, as well as the results of 'IQueryParam.Parse()'.
// Query param values are stored with 'IQueryParam.Name()' as the key.
func (cx *ChartContext) Get(key string) (any, bool) {
	value, ok := cx.data[key]
	return value, ok
}

// Exactly like Get but without the bool return.
func (cx *ChartContext) MustGet(key string) any {
	return cx.data[key]
}
func (cx *ChartContext) Query(key string) string {
	return cx.Gin().Query(key)
}
func (cx *ChartContext) Db() orm.DBTX        { return cx.db }
func (cx *ChartContext) Cache() *cache.Cache { return cx.cache }
func (cx *ChartContext) Gin() *gin.Context   { return cx.gctx }

type IQueryParam interface {
	// Returns the query param's URL query field name
	Name() string
	Parse(c *ChartContext) (any, error)
	ToFilter(c *ChartContext) Filter
	// Copy the query parameter
	Clone() IQueryParam
}

type enumParam struct {
	Label      string
	name       string
	OptionsSql string
	IsArray    bool
	// If true, allow * to mean all options
	allowStar bool
	// List of options, pulled from cache
	options []string
}

func NewEnumParam(label, name, sql string, isArray, allowStar bool) IQueryParam {
	return &enumParam{
		Label:      label,
		name:       name,
		OptionsSql: sql,
		IsArray:    isArray,
		allowStar:  allowStar,
	}
}

func (p *enumParam) Name() string {
	return p.name
}

func (p *enumParam) Parse(c *ChartContext) (any, error) {
	qvalue := c.Query(p.name)
	// Load options
	options, err := p.getOptions(c)
	if err != nil {
		return nil, err
	}
	// Check if value(s) are in options
	var values []string
	if p.IsArray {
		values = strings.Split(qvalue, ",")
	} else {
		values = []string{qvalue}
	}
	if p.allowStar && values[0] == "*" {
		values = options[1:]
	} else {
		for _, v := range values {
			if !lo.Contains(options, v) {
				return nil, fmt.Errorf("%w (%s, %s)", ErrInvalidParam, p.name, v)
			}
		}
	}
	p.options = options
	return values, nil
}

func (p *enumParam) getOptions(c *ChartContext) ([]string, error) {
	cachekey := p.Label + "-options"
	if opts, ok := c.Cache().Get(cachekey); ok {
		return opts.([]string), nil
	} else {
		optRows, _ := c.Db().Query(c.Ctx, p.OptionsSql)
		opts := make([]string, 0)
		if err := util.ScanRows(optRows, &opts); err != nil {
			return nil, err
		}
		if p.allowStar {
			opts = append([]string{"*"}, opts...)
		}
		c.Cache().Set(cachekey, opts, 0)
		return opts, nil
	}
}

func (p *enumParam) ToFilter(c *ChartContext) Filter {
	value := c.MustGet(p.name).([]string)
	return Filter{
		Label:   p.Label,
		Name:    p.name,
		Type:    "select",
		Value:   strings.Join(value, ","),
		Options: p.options,
	}
}

func (p *enumParam) Clone() IQueryParam {
	return &enumParam{
		Label:      p.Label,
		name:       p.name,
		OptionsSql: p.OptionsSql,
		IsArray:    p.IsArray,
		allowStar:  p.allowStar,
	}
}
