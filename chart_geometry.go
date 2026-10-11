package popui

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/miguel-raspberry-pi/popui.go/props"
)

// The Chart plot is drawn in a chartVB x chartVB SVG viewBox stretched to the
// plot box (preserveAspectRatio="none"); strokes use non-scaling-stroke so
// lines keep their width. HTML labels and points use the same geometry as
// percentages of the box.
const (
	chartVB              = 1000.0
	chartDefaultHeight   = 200
	chartSparklineHeight = 32
	chartMaxYTicks       = 5
	chartMaxXTicks       = 6
	// Series with more points than this only get hover targets when
	// ShowPoints is set, to keep the markup small.
	chartMaxHoverPoints = 200
)

// chartPalette is the order colors are assigned to series without one.
var chartPalette = []string{
	props.ChartColorAccent,
	props.ChartColorAttention,
	props.ChartColorPositive,
	props.ChartColorWarning,
	props.ChartColorCritical,
	props.ChartColorInfo,
}

// chartColorClass maps a semantic color name to the text color class the
// series is drawn with (strokes and fills use currentColor).
func chartColorClass(color string) string {
	switch color {
	case props.ChartColorPositive, "success":
		return "text-foreground-success"
	case props.ChartColorCritical:
		return "text-foreground-critical"
	case props.ChartColorWarning:
		return "text-foreground-warning"
	case props.ChartColorInfo:
		return "text-foreground-info"
	case props.ChartColorAttention:
		return "text-foreground-attention"
	case props.ChartColorNeutral:
		return "text-foreground-default-secondary"
	default:
		return "text-foreground-accent"
	}
}

type chartTick struct {
	Pos   float64 // percent along the axis (left for x, top for y)
	Label string
	Align string // x ticks: "start", "center" or "end"
}

type chartPointView struct {
	X, Y  float64 // percent of the plot box from the left and top
	Title string
}

type chartSeriesView struct {
	Label      string
	ColorClass string
	Line       string // SVG path in viewBox units
	Area       string
	Points     []chartPointView
	ShowPoints bool // dots always visible
	Hover      bool // invisible hover targets with tooltips
}

type chartView struct {
	Empty     bool
	Height    int
	YTicks    []chartTick
	XTicks    []chartTick
	YWidest   string // longest y label, sizes the label column
	Grid      string // SVG path of the horizontal grid lines
	Series    []chartSeriesView
	Legend    bool
	AriaLabel string
	EmptyText string
}

type chartXY struct{ X, Y float64 }

// buildChart computes everything the Chart template draws.
func buildChart(p props.Chart) chartView {
	v := chartView{Height: p.Height, EmptyText: chartEmptyText(p.Empty)}
	if v.Height <= 0 {
		v.Height = chartDefaultHeight
		if p.Sparkline {
			v.Height = chartSparklineHeight
		}
	}

	series := make([][]props.ChartPoint, len(p.Series))
	var xmin, xmax time.Time
	ymin, ymax := math.Inf(1), math.Inf(-1)
	n := 0
	for i, s := range p.Series {
		pts := make([]props.ChartPoint, 0, len(s.Points))
		for _, pt := range s.Points {
			if math.IsNaN(pt.Y) || math.IsInf(pt.Y, 0) {
				continue
			}
			pts = append(pts, pt)
		}
		sort.SliceStable(pts, func(a, b int) bool { return pts[a].X.Before(pts[b].X) })
		series[i] = pts
		for _, pt := range pts {
			if n == 0 || pt.X.Before(xmin) {
				xmin = pt.X
			}
			if n == 0 || pt.X.After(xmax) {
				xmax = pt.X
			}
			ymin = math.Min(ymin, pt.Y)
			ymax = math.Max(ymax, pt.Y)
			n++
		}
	}

	yfmt := p.YFormat
	v.AriaLabel = chartAriaLabel(p, series, yfmt)
	if n == 0 {
		v.Empty = true
		return v
	}

	// Y domain and its mapping to percent from the top.
	var ylo, yhi float64
	ypad := 0.0
	if p.Sparkline {
		ylo, yhi = ymin, ymax
		if ylo == yhi {
			ylo, yhi = ylo-1, yhi+1
		}
		ypad = 10 // keep the stroke off the edges of a tiny box
	} else {
		var step float64
		var ticks []float64
		ylo, yhi, step, ticks = niceTicks(ymin, ymax, chartMaxYTicks)
		axisfmt := yfmt
		if axisfmt == nil {
			// Same decimals on every tick: 0.5, 1.0, 1.5.
			dec := chartDecimals(step)
			axisfmt = func(f float64) string { return strconv.FormatFloat(f, 'f', dec, 64) }
		}
		var grid strings.Builder
		for _, t := range ticks {
			pos := chartScaleY(t, ylo, yhi, 0)
			label := axisfmt(t)
			v.YTicks = append(v.YTicks, chartTick{Pos: pos, Label: label})
			if len([]rune(label)) > len([]rune(v.YWidest)) {
				v.YWidest = label
			}
			if pos > 0.001 && pos < 99.999 {
				y := chartNum(pos * chartVB / 100)
				grid.WriteString("M0 " + y + "H" + chartNum(chartVB))
			}
		}
		v.Grid = grid.String()
	}
	if yfmt == nil {
		yfmt = chartDefaultYFormat
	}

	if !p.Sparkline {
		v.XTicks = chartXTicks(xmin, xmax, p.XFormat)
	}
	v.Legend = !p.Sparkline && len(p.Series) > 1

	multi := len(p.Series) > 1
	for i, s := range p.Series {
		pts := series[i]
		color := s.Color
		if color == "" {
			color = chartPalette[i%len(chartPalette)]
		}
		sv := chartSeriesView{
			Label:      s.Label,
			ColorClass: chartColorClass(color),
			ShowPoints: p.ShowPoints || len(pts) == 1,
		}
		sv.Hover = sv.ShowPoints || len(pts) <= chartMaxHoverPoints
		xy := make([]chartXY, len(pts))
		for j, pt := range pts {
			xy[j] = chartXY{
				X: chartScaleX(pt.X, xmin, xmax),
				Y: chartScaleY(pt.Y, ylo, yhi, ypad),
			}
			if sv.Hover {
				sv.Points = append(sv.Points, chartPointView{
					X:     xy[j].X,
					Y:     xy[j].Y,
					Title: chartPointTitle(pt, s.Label, multi, yfmt, p.XFormat),
				})
			}
		}
		sv.Line = chartLinePath(xy, s.Step)
		if p.Area && len(xy) > 1 {
			sv.Area = chartAreaPath(xy, s.Step)
		}
		v.Series = append(v.Series, sv)
	}
	return v
}

// niceTicks extends [min, max] to round bounds and returns evenly spaced
// ticks at a "nice" step (1, 2 or 5 times a power of ten), with at most
// about maxTicks ticks. A zero-width range is widened around its value.
func niceTicks(min, max float64, maxTicks int) (lo, hi, step float64, ticks []float64) {
	if min > max {
		min, max = max, min
	}
	if min == max {
		d := math.Abs(min) * 0.1
		if d == 0 {
			d = 1
		}
		min, max = min-d, max+d
	}
	if maxTicks < 2 {
		maxTicks = 2
	}
	rng := niceNum(max-min, false)
	step = niceNum(rng/float64(maxTicks-1), true)
	// Snap away float noise such as 0.6000000000000001.
	scale := math.Pow(10, float64(chartDecimals(step)+1))
	snap := func(f float64) float64 {
		f = math.Round(f*scale) / scale
		if f == 0 {
			f = 0 // no negative zero
		}
		return f
	}
	lo = snap(math.Floor(min/step) * step)
	hi = snap(math.Ceil(max/step) * step)
	count := int(math.Round((hi-lo)/step)) + 1
	for i := 0; i < count; i++ {
		ticks = append(ticks, snap(lo+float64(i)*step))
	}
	return lo, hi, step, ticks
}

// niceNum returns a "nice" number close to x: 1, 2, 5 or 10 times a power
// of ten, rounded when round is set and otherwise the smallest one >= x.
func niceNum(x float64, round bool) float64 {
	if x <= 0 {
		return 1
	}
	exp := math.Floor(math.Log10(x))
	f := x / math.Pow(10, exp)
	var nf float64
	if round {
		switch {
		case f < 1.5:
			nf = 1
		case f < 3:
			nf = 2
		case f < 7:
			nf = 5
		default:
			nf = 10
		}
	} else {
		switch {
		case f <= 1:
			nf = 1
		case f <= 2:
			nf = 2
		case f <= 5:
			nf = 5
		default:
			nf = 10
		}
	}
	return nf * math.Pow(10, exp)
}

// chartDecimals is the number of decimals needed to tell ticks step apart.
func chartDecimals(step float64) int {
	if step <= 0 || step >= 1 {
		return 0
	}
	d := int(math.Ceil(-math.Log10(step) - 1e-9))
	if d > 6 {
		d = 6
	}
	return d
}

func chartDefaultYFormat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// chartScaleY maps v in [lo, hi] to percent from the top, leaving pad
// percent at both ends.
func chartScaleY(v, lo, hi, pad float64) float64 {
	if hi == lo {
		return 50
	}
	return pad + (hi-v)/(hi-lo)*(100-2*pad)
}

// chartScaleX maps t in [min, max] to percent from the left; a single
// moment sits in the middle.
func chartScaleX(t, min, max time.Time) float64 {
	span := max.Sub(min)
	if span <= 0 {
		return 50
	}
	return float64(t.Sub(min)) / float64(span) * 100
}

// chartNum formats a coordinate with at most two decimals.
func chartNum(f float64) string {
	f = math.Round(f*100) / 100
	if f == 0 {
		f = 0
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func chartPt(p chartXY) string {
	return chartNum(p.X*chartVB/100) + " " + chartNum(p.Y*chartVB/100)
}

// chartLinePath builds the SVG path of a series from percent coordinates.
// Step lines hold each value until the next point.
func chartLinePath(pts []chartXY, step bool) string {
	if len(pts) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("M" + chartPt(pts[0]))
	for i := 1; i < len(pts); i++ {
		if step {
			b.WriteString("H" + chartNum(pts[i].X*chartVB/100))
			b.WriteString("V" + chartNum(pts[i].Y*chartVB/100))
		} else {
			b.WriteString("L" + chartPt(pts[i]))
		}
	}
	return b.String()
}

// chartAreaPath closes the series line down to the bottom of the plot.
func chartAreaPath(pts []chartXY, step bool) string {
	if len(pts) < 2 {
		return ""
	}
	base := chartNum(chartVB)
	return chartLinePath(pts, step) +
		"V" + base +
		"H" + chartNum(pts[0].X*chartVB/100) + "Z"
}

// chartTimeStep is one candidate interval between x-axis ticks.
type chartTimeStep struct {
	unit   string // minute, hour, day, week, month, year
	n      int
	approx time.Duration
	layout string // default label layout
}

const chartDay = 24 * time.Hour

var chartTimeSteps = []chartTimeStep{
	{"minute", 1, time.Minute, "15:04"},
	{"minute", 5, 5 * time.Minute, "15:04"},
	{"minute", 15, 15 * time.Minute, "15:04"},
	{"minute", 30, 30 * time.Minute, "15:04"},
	{"hour", 1, time.Hour, "15:04"},
	{"hour", 3, 3 * time.Hour, "15:04"},
	{"hour", 6, 6 * time.Hour, "15:04"},
	{"hour", 12, 12 * time.Hour, "2 Jan 15:04"},
	{"day", 1, chartDay, "2 Jan"},
	{"day", 2, 2 * chartDay, "2 Jan"},
	{"week", 1, 7 * chartDay, "2 Jan"},
	{"week", 2, 14 * chartDay, "2 Jan"},
	{"month", 1, 30 * chartDay, "Jan 2006"},
	{"month", 2, 61 * chartDay, "Jan 2006"},
	{"month", 3, 91 * chartDay, "Jan 2006"},
	{"month", 6, 182 * chartDay, "Jan 2006"},
	{"year", 1, 365 * chartDay, "2006"},
	{"year", 2, 2 * 365 * chartDay, "2006"},
	{"year", 5, 5 * 365 * chartDay, "2006"},
	{"year", 10, 10 * 365 * chartDay, "2006"},
	{"year", 25, 25 * 365 * chartDay, "2006"},
	{"year", 50, 50 * 365 * chartDay, "2006"},
	{"year", 100, 100 * 365 * chartDay, "2006"},
}

// chartTimeTicks returns tick times on calendar boundaries (in min's
// location) between min and max, using the finest step that gives at most
// maxTicks ticks, and the step's default label layout.
func chartTimeTicks(min, max time.Time, maxTicks int) ([]time.Time, string) {
	span := max.Sub(min)
	if span <= 0 {
		return []time.Time{min}, "2 Jan 2006"
	}
	for _, s := range chartTimeSteps {
		if int(span/s.approx) > maxTicks {
			continue
		}
		ticks := s.ticks(min, max)
		if len(ticks) <= maxTicks {
			return ticks, s.layout
		}
	}
	return []time.Time{min, max}, "2006"
}

// start is the first step boundary at or after t.
func (s chartTimeStep) start(t time.Time) time.Time {
	loc := t.Location()
	y, mo, d := t.Date()
	var b time.Time
	switch s.unit {
	case "minute":
		b = time.Date(y, mo, d, t.Hour(), t.Minute()/s.n*s.n, 0, 0, loc)
	case "hour":
		b = time.Date(y, mo, d, t.Hour()/s.n*s.n, 0, 0, 0, loc)
	case "day":
		b = time.Date(y, mo, d, 0, 0, 0, 0, loc)
	case "week":
		b = time.Date(y, mo, d, 0, 0, 0, 0, loc)
		// Back to Monday.
		b = b.AddDate(0, 0, -((int(b.Weekday()) + 6) % 7))
	case "month":
		b = time.Date(y, time.Month((int(mo)-1)/s.n*s.n+1), 1, 0, 0, 0, 0, loc)
	default:
		b = time.Date(y/s.n*s.n, time.January, 1, 0, 0, 0, 0, loc)
	}
	for b.Before(t) {
		b = s.next(b)
	}
	return b
}

func (s chartTimeStep) next(t time.Time) time.Time {
	switch s.unit {
	case "minute":
		return t.Add(time.Duration(s.n) * time.Minute)
	case "hour":
		return t.Add(time.Duration(s.n) * time.Hour)
	case "day":
		return t.AddDate(0, 0, s.n)
	case "week":
		return t.AddDate(0, 0, 7*s.n)
	case "month":
		return t.AddDate(0, s.n, 0)
	default:
		return t.AddDate(s.n, 0, 0)
	}
}

func (s chartTimeStep) ticks(min, max time.Time) []time.Time {
	var out []time.Time
	for t := s.start(min); !t.After(max) && len(out) <= chartMaxXTicks*4; t = s.next(t) {
		out = append(out, t)
	}
	return out
}

// chartXTicks positions and labels the x-axis ticks.
func chartXTicks(min, max time.Time, format func(time.Time) string) []chartTick {
	times, layout := chartTimeTicks(min, max, chartMaxXTicks)
	if len(times) == 0 {
		times = []time.Time{min, max}
		if !max.After(min) {
			times = times[:1]
		}
	}
	out := make([]chartTick, 0, len(times))
	for _, t := range times {
		pos := chartScaleX(t, min, max)
		label := t.Format(layout)
		if format != nil {
			label = format(t)
		}
		align := "center"
		switch {
		case max.After(min) && pos < 8:
			align = "start"
		case max.After(min) && pos > 92:
			align = "end"
		}
		out = append(out, chartTick{Pos: pos, Label: label, Align: align})
	}
	return out
}

// chartDate formats a point's date for tooltips and the aria-label.
func chartDate(t time.Time, format func(time.Time) string) string {
	if format != nil {
		return format(t)
	}
	if t.Hour() != 0 || t.Minute() != 0 {
		return t.Format("2 Jan 2006 15:04")
	}
	return t.Format("2 Jan 2006")
}

func chartPointTitle(pt props.ChartPoint, series string, multi bool, yfmt func(float64) string, xfmt func(time.Time) string) string {
	value := yfmt(pt.Y)
	if multi && series != "" {
		value = series + ": " + value
	}
	parts := []string{chartDate(pt.X, xfmt), value}
	if pt.Label != "" {
		parts = append(parts, pt.Label)
	}
	return strings.Join(parts, " · ")
}

// chartAriaLabel names the chart and summarises each series.
func chartAriaLabel(p props.Chart, series [][]props.ChartPoint, yfmt func(float64) string) string {
	if yfmt == nil {
		yfmt = chartDefaultYFormat
	}
	title := p.Title
	if title == "" {
		title = "Line chart"
	}
	var parts []string
	for i, pts := range series {
		if len(pts) == 0 {
			continue
		}
		label := p.Series[i].Label
		if label == "" {
			label = "Series " + strconv.Itoa(i+1)
		}
		first, last := pts[0], pts[len(pts)-1]
		if len(pts) == 1 {
			parts = append(parts, label+": "+yfmt(first.Y)+" on "+chartDate(first.X, p.XFormat))
			continue
		}
		lo, hi := first.Y, first.Y
		for _, pt := range pts {
			lo = math.Min(lo, pt.Y)
			hi = math.Max(hi, pt.Y)
		}
		parts = append(parts, label+": "+strconv.Itoa(len(pts))+" points from "+
			chartDate(first.X, p.XFormat)+" to "+chartDate(last.X, p.XFormat)+
			", between "+yfmt(lo)+" and "+yfmt(hi)+", latest "+yfmt(last.Y))
	}
	if len(parts) == 0 {
		return title + ": " + chartEmptyText(p.Empty)
	}
	return title + ". " + strings.Join(parts, "; ")
}

func chartEmptyText(s string) string {
	if s == "" {
		return "No data"
	}
	return s
}
