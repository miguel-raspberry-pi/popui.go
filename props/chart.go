package props

import (
	"time"

	"github.com/a-h/templ"
)

// Chart series colors. Each name maps to a semantic foreground token, so
// series follow the theme (including dark mode).
const (
	ChartColorAccent    string = "accent"    // foreground-accent
	ChartColorPositive  string = "positive"  // foreground-success
	ChartColorCritical  string = "critical"  // foreground-critical
	ChartColorWarning   string = "warning"   // foreground-warning
	ChartColorInfo      string = "info"      // foreground-info
	ChartColorAttention string = "attention" // foreground-attention
	ChartColorNeutral   string = "neutral"   // foreground-default-secondary
)

// Chart defines the props of a server-rendered SVG line chart over a time
// axis. No client JavaScript is involved: the plot is an SVG stretched to
// the container's width, and the axis labels and points are HTML positioned
// over it, so text stays legible at any width.
type Chart struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Title is the accessible name of the chart, e.g. "Price history". A
	// summary of each series (range of dates and values) is appended to it
	// in the aria-label.
	Title string
	// Series are the lines to draw. All series share the same axes. A legend
	// is shown when there is more than one.
	Series []ChartSeries
	// Height of the plot area in px. Defaults to 200, or 32 for sparklines.
	Height int
	// YFormat formats values for the y-axis labels and the point tooltips,
	// e.g. a currency formatter. Defaults to a plain number.
	YFormat func(float64) string
	// XFormat formats dates for the x-axis labels and the point tooltips.
	// Defaults to a format that fits the time span ("15:04", "2 Jan",
	// "Jan 2006") for the axis and "2 Jan 2006" for tooltips.
	XFormat func(time.Time) string
	// Area fills the space under each line with a translucent tint.
	Area bool
	// ShowPoints draws a dot on every point. Without it dots appear only
	// while hovered (and always for a series with a single point).
	ShowPoints bool
	// Sparkline renders a compact inline chart: no axes, grid, labels or
	// legend, and a y range fitted to the data.
	Sparkline bool
	// Empty is the text shown when no series has points. Defaults to
	// "No data".
	Empty string
}

// ChartSeries is one line of a Chart.
type ChartSeries struct {
	// Label names the series in the legend, tooltips and aria-label.
	Label string
	// Points of the series, in any order; they are sorted by X.
	Points []ChartPoint
	// Color is a semantic color name, see the ChartColor constants.
	// Defaults to a color picked by series position (accent, attention,
	// positive, warning, critical, info).
	Color string
	// Step draws a step line: each value holds until the next point, which
	// suits prices and other values that change at discrete moments.
	Step bool
}

// ChartPoint is a single value at a moment in time.
type ChartPoint struct {
	X time.Time
	Y float64
	// Label is optional text added to the point's tooltip, e.g. "Price
	// drop".
	Label string
}
