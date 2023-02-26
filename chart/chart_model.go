package chart

import "html/template"

type ChartMeta struct {
	Type    string         `json:"type"`
	Options map[string]any `json:"options"`
}

type BarChart struct {
	ChartMeta
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

func (c *BarChart) Defaults() {
	c.ChartMeta.Type = "bar"
}

type Dataset struct {
	Label string    `json:"label"`
	Data  []float32 `json:"data"`
}

type TabularOpts struct {
	Layout  string          `json:"layout"`
	Data    []any           `json:"data"`
	Columns []TabularColumn `json:"columns"`
}

type TabularColumn struct {
	Title string `toml:"title" json:"title"`
	Field string `toml:"field" json:"field"`
}

// Type passed when rendering either 'chartview.html' or
// 'tableview.html' templates.
type ChartTemplate struct {
	// Page title
	Title string
	// Raw HTML to include in the <head>
	Head template.HTML
	// JSON chart options
	Options any
	// Chart setup code
	ChartCode template.JS
	// List of filters the user can change to select more information
	Filters []Filter
	// Raw chart data, likely directly from SQL, optional
	RawData any
	// If true use the table library, otherwise use chartjs
	IsTable bool
}

type Filter struct {
	// Filter description
	Label string
	// HTML name of the filter
	Name string
	// Filter type ("select", "text", "number")
	Type string
	// Initial input value
	Value string
	// If type is "select", these are the options
	Options []string
}
