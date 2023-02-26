package chart

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnknownPredef = errors.New("unknown query predef - ")

type ChartType string

const (
	ChartTable ChartType = "table"
	ChartBar   ChartType = "bar"
	ChartLine  ChartType = "line"
)

// Context for converting config data to parsed the in-memory representation.
type ConvertContext struct {
	// Map of predef names to query parameters
	Predefs map[string]IQueryParam
	// Map of chart names to converted charts
	Charts map[string]*ChartSpec
}

func (cx *ConvertContext) Init() {
	cx.Predefs = make(map[string]IQueryParam)
	cx.Charts = make(map[string]*ChartSpec)
}

func (cx *ConvertContext) Add(data ChartToml) error {
	spec, err := data.ToChartSpec(cx)
	if err != nil {
		return err
	}
	cx.Charts[spec.Path] = spec
	return nil
}

type ChartToml struct {
	// Chart name
	Name string `toml:"name" json:"name"`
	// SQL query used to obtain chart data
	Sql string `toml:"sql" json:"sql"`
	// Chart type
	Type ChartType `toml:"type" json:"type"`
	// Chart parameters
	Params []ParamToml `toml:"params" json:"params"`
	// Name of JS function which will transform the data.
	// This may also be an anonymous inline JS function.
	Transformer string `toml:"transformer" json:"transformer"`
	// type=table: List of columns
	Columns []TabularColumn `toml:"columns" json:"columns"`
}

func (c ChartToml) ToChartSpec(cx *ConvertContext) (*ChartSpec, error) {
	params := make([]IQueryParam, len(c.Params))
	for i, p := range c.Params {
		qp, err := p.ToQueryParam(cx)
		if err != nil {
			return nil, err
		}
		params[i] = qp
	}
	path := strings.ToLower(c.Name)
	s := ChartSpec{
		Name:          c.Name,
		Path:          path,
		Params:        params,
		Sql:           c.Sql,
		Type:          string(c.Type),
		TransformCode: c.Transformer,
	}
	switch c.Type {
	case ChartTable:
		s.Options = TabularOpts{
			Layout:  "fitColumns",
			Columns: c.Columns,
		}
	}
	return &s, nil
}

type ParamToml struct {
	// Use a pre-defined parameter type
	Def string `toml:"def" json:"def"`
	// Parameter label text
	Label string `toml:"label" json:"label"`
	// URL query-string name
	Name string `toml:"name" json:"name"`
	// Parameter type ('enum', 'text', 'number')
	Type string `toml:"type" json:"type"`
	// type=enum: SQL to execute to get list of options
	OptionsSql string `toml:"options_sql" json:"options_sql"`
	// type=text: Validation regex
	TextRegex string `toml:"text_regex" json:"text_regex"`
}

func (p ParamToml) ToQueryParam(cx *ConvertContext) (IQueryParam, error) {
	// Check for predef
	if p.Def != "" {
		if def, ok := cx.Predefs[p.Def]; ok {
			return def, nil
		} else {
			return nil, fmt.Errorf("%w %s", ErrUnknownPredef, p.Def)
		}
	}
	// TODO
	return nil, nil
}
