package chart

import (
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/bindernews/sts-msr/orm"
	"github.com/bindernews/sts-msr/util"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// Gin context key for an instance of orm.DBTX
const CtxDb = "chartdb"

// Error when a parameter has an invalid value
var ErrInvalidParam = errors.New("parameter value invalid")

const chartJsCode = `
function showChart() { opts.data = %s(rawdata); const ctx = document.getElementById('chartc'); new Chart(ctx, opts); }
`

const tableJsCode = `
function showChart() { opts.data = %s(rawdata); new Tabulator('#table', opts); }
`

// Specification of a chart. These may be specified in the config file
// and are then parsed and converted into 'ChartSpec's.
type ChartSpec struct {
	// Chart name
	Name string
	// Chart URL-safe path
	Path string
	// Query parameters, positional
	Params []IQueryParam
	// SQL statement to execute
	Sql string
	// Chart type
	Type string
	// JS code to transform the SQL output to appropriate format
	// for the chart or table library to use
	TransformCode string
	// Chart-specific options
	Options any
}

func (spec *ChartSpec) Handle(c *gin.Context) {
	ctx := c.Request.Context()
	db := c.MustGet(CtxDb).(orm.DBTX)
	// Parse query params
	args := make([]any, len(spec.Params))
	filters := make([]Filter, len(args))
	for i, param := range spec.Params {
		if a, err := param.Parse(c); err != nil {
			util.AbortErr(c, 400, err)
			return
		} else {
			args[i] = a
			filters[i] = param.ToFilter(a)
		}
	}
	// Perform the query
	rows, err := db.Query(ctx, spec.Sql, args...)
	if err != nil {
		util.AbortErr(c, 500, err)
		return
	}
	// Load rows into JSON data
	data := make([]any, 0)
	for rows.Next() {
		if r, err := rows.Values(); err != nil {
			util.AbortErr(c, 500, err)
			return
		} else {
			data = append(data, r)
		}
	}

	// Render and serve HTML
	tmpl := ChartTemplate{
		Title:   spec.Name,
		Filters: filters,
		RawData: data,
		Options: spec.Options,
		IsTable: spec.Type == "table",
	}

	switch spec.Type {
	case "table":
		tmpl.ChartCode = template.JS(fmt.Sprintf(tableJsCode, spec.TransformCode))
	case "bar":
		tmpl.ChartCode = template.JS(fmt.Sprintf(chartJsCode, spec.TransformCode))
	default:
		util.AbortErr(c, 500, errors.New("invalid chart type"))
		return
	}
	c.HTML(200, "chartview.html", tmpl)
}

type IQueryParam interface {
	Parse(c *gin.Context) (any, error)
	ToFilter(value any) Filter
}

func NewEnumParam(label, name, sql string, isArray bool) IQueryParam {
	return &enumParam{
		Label:      label,
		Name:       name,
		OptionsSql: sql,
		IsArray:    isArray,
		options:    make([]string, 0),
	}
}

type enumParam struct {
	Label      string
	Name       string
	OptionsSql string
	IsArray    bool
	// List of known options, refreshed during parse
	options []string
}

func (p *enumParam) Parse(c *gin.Context) (any, error) {
	ctx := c.Request.Context()
	qvalue := c.Query(p.Name)
	db := c.MustGet(CtxDb).(orm.DBTX)
	// Reset options
	p.options = p.options[:0]
	optRows, _ := db.Query(ctx, p.OptionsSql)
	if err := util.ScanRows(optRows, &p.options); err != nil {
		return nil, err
	}
	// Check if value(s) are in options
	var values []string
	if p.IsArray {
		values = strings.Split(qvalue, ",")
	} else {
		values = []string{qvalue}
	}
	for _, v := range values {
		if !lo.Contains(p.options, v) {
			return nil, fmt.Errorf("%w (%s, %s)", ErrInvalidParam, p.Name, v)
		}
	}
	return values, nil
}

func (p *enumParam) ToFilter(value any) Filter {
	return Filter{
		Label:   p.Label,
		Name:    p.Name,
		Type:    "select",
		Value:   strings.Join(value.([]string), ","),
		Options: p.options,
	}
}
