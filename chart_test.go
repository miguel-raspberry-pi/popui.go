package popui

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/miguel-raspberry-pi/popui.go/props"
)

func TestNiceTicks(t *testing.T) {
	tests := []struct {
		min, max     float64
		lo, hi, step float64
		ticks        []float64
	}{
		{0, 10, 0, 10, 2, []float64{0, 2, 4, 6, 8, 10}},
		{1180, 1360, 1150, 1400, 50, []float64{1150, 1200, 1250, 1300, 1350, 1400}},
		{0.1, 0.9, 0, 1, 0.2, []float64{0, 0.2, 0.4, 0.6, 0.8, 1}},
		{-3, 7, -4, 8, 2, []float64{-4, -2, 0, 2, 4, 6, 8}},
		// Flat series widen by 10% around the value.
		{100, 100, 90, 110, 5, []float64{90, 95, 100, 105, 110}},
		// Flat zero widens to ±1.
		{0, 0, -1, 1, 0.5, []float64{-1, -0.5, 0, 0.5, 1}},
		// Reversed input.
		{10, 0, 0, 10, 2, []float64{0, 2, 4, 6, 8, 10}},
	}
	for _, tt := range tests {
		lo, hi, step, ticks := niceTicks(tt.min, tt.max, 5)
		if lo != tt.lo || hi != tt.hi || step != tt.step || !reflect.DeepEqual(ticks, tt.ticks) {
			t.Errorf("niceTicks(%v, %v) = %v, %v, %v, %v; want %v, %v, %v, %v",
				tt.min, tt.max, lo, hi, step, ticks, tt.lo, tt.hi, tt.step, tt.ticks)
		}
		if lo > tt.min && lo > tt.max || hi < tt.max && hi < tt.min {
			t.Errorf("niceTicks(%v, %v) does not cover the range", tt.min, tt.max)
		}
	}
}

func TestNiceNum(t *testing.T) {
	for _, tt := range []struct {
		x     float64
		round bool
		want  float64
	}{
		{0.7, false, 1}, {1.3, false, 2}, {3, false, 5}, {7, false, 10}, {120, false, 200},
		{1.2, true, 1}, {2.4, true, 2}, {4, true, 5}, {8, true, 10}, {0, true, 1},
	} {
		if got := niceNum(tt.x, tt.round); got != tt.want {
			t.Errorf("niceNum(%v, %v) = %v, want %v", tt.x, tt.round, got, tt.want)
		}
	}
}

func TestChartDecimals(t *testing.T) {
	for step, want := range map[float64]int{50: 0, 1: 0, 0.5: 1, 0.2: 1, 0.1: 1, 0.05: 2, 0.001: 3} {
		if got := chartDecimals(step); got != want {
			t.Errorf("chartDecimals(%v) = %d, want %d", step, got, want)
		}
	}
}

func TestChartScale(t *testing.T) {
	if got := chartScaleY(10, 0, 10, 0); got != 0 {
		t.Errorf("top = %v", got)
	}
	if got := chartScaleY(0, 0, 10, 0); got != 100 {
		t.Errorf("bottom = %v", got)
	}
	if got := chartScaleY(5, 0, 10, 10); got != 50 {
		t.Errorf("padded middle = %v", got)
	}
	if got := chartScaleY(0, 0, 10, 10); got != 90 {
		t.Errorf("padded bottom = %v", got)
	}
	if got := chartScaleY(3, 3, 3, 0); got != 50 {
		t.Errorf("flat = %v", got)
	}
	a := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	b := a.AddDate(0, 0, 10)
	if got := chartScaleX(a.AddDate(0, 0, 4), a, b); got != 40 {
		t.Errorf("x = %v", got)
	}
	if got := chartScaleX(a, a, a); got != 50 {
		t.Errorf("single moment x = %v", got)
	}
}

func TestChartPaths(t *testing.T) {
	pts := []chartXY{{0, 100}, {50, 25.5}, {100, 0}}
	if got, want := chartLinePath(pts, false), "M0 1000L500 255L1000 0"; got != want {
		t.Errorf("line = %s, want %s", got, want)
	}
	if got, want := chartLinePath(pts, true), "M0 1000H500V255H1000V0"; got != want {
		t.Errorf("step = %s, want %s", got, want)
	}
	if got, want := chartAreaPath(pts, false), "M0 1000L500 255L1000 0V1000H0Z"; got != want {
		t.Errorf("area = %s, want %s", got, want)
	}
	if got := chartAreaPath(pts[:1], false); got != "" {
		t.Errorf("single point area = %s", got)
	}
	if got := chartLinePath(nil, false); got != "" {
		t.Errorf("empty line = %s", got)
	}
	if got := chartNum(1.23456); got != "1.23" {
		t.Errorf("chartNum = %s", got)
	}
}

func TestChartTimeTicks(t *testing.T) {
	d := func(y int, m time.Month, day int) time.Time { return time.Date(y, m, day, 0, 0, 0, 0, time.UTC) }
	for _, tt := range []struct {
		name     string
		min, max time.Time
		want     []time.Time
		layout   string
	}{
		{
			name: "days",
			min:  d(2026, 3, 1), max: d(2026, 3, 5),
			want:   []time.Time{d(2026, 3, 1), d(2026, 3, 2), d(2026, 3, 3), d(2026, 3, 4), d(2026, 3, 5)},
			layout: "2 Jan",
		},
		{
			name: "weeks start on monday",
			min:  d(2026, 3, 4), max: d(2026, 4, 10),
			want:   []time.Time{d(2026, 3, 9), d(2026, 3, 16), d(2026, 3, 23), d(2026, 3, 30), d(2026, 4, 6)},
			layout: "2 Jan",
		},
		{
			name: "quarters",
			min:  d(2025, 2, 10), max: d(2026, 8, 1),
			want:   []time.Time{d(2025, 4, 1), d(2025, 7, 1), d(2025, 10, 1), d(2026, 1, 1), d(2026, 4, 1), d(2026, 7, 1)},
			layout: "Jan 2006",
		},
		{
			name: "hours",
			min:  time.Date(2026, 3, 1, 1, 30, 0, 0, time.UTC), max: time.Date(2026, 3, 1, 5, 0, 0, 0, time.UTC),
			want: []time.Time{
				time.Date(2026, 3, 1, 2, 0, 0, 0, time.UTC), time.Date(2026, 3, 1, 3, 0, 0, 0, time.UTC),
				time.Date(2026, 3, 1, 4, 0, 0, 0, time.UTC), time.Date(2026, 3, 1, 5, 0, 0, 0, time.UTC),
			},
			layout: "15:04",
		},
		{
			name: "single moment",
			min:  d(2026, 3, 1), max: d(2026, 3, 1),
			want:   []time.Time{d(2026, 3, 1)},
			layout: "2 Jan 2006",
		},
	} {
		got, layout := chartTimeTicks(tt.min, tt.max, chartMaxXTicks)
		if !reflect.DeepEqual(got, tt.want) || layout != tt.layout {
			t.Errorf("%s: got %v %q, want %v %q", tt.name, got, layout, tt.want, tt.layout)
		}
	}
}

func TestChartXTicksAlign(t *testing.T) {
	a := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	ticks := chartXTicks(a, a.AddDate(0, 0, 4), nil)
	if len(ticks) != 5 || ticks[0].Align != "start" || ticks[2].Align != "center" || ticks[4].Align != "end" {
		t.Errorf("ticks = %+v", ticks)
	}
	if ticks[0].Label != "1 Mar" {
		t.Errorf("label = %s", ticks[0].Label)
	}
	ticks = chartXTicks(a, a, func(time.Time) string { return "custom" })
	if len(ticks) != 1 || ticks[0].Pos != 50 || ticks[0].Align != "center" || ticks[0].Label != "custom" {
		t.Errorf("single = %+v", ticks)
	}
}

func eur(f float64) string { return fmt.Sprintf("%.0f €", f) }

func TestBuildChart(t *testing.T) {
	d := func(m time.Month, day int) time.Time { return time.Date(2026, m, day, 0, 0, 0, 0, time.UTC) }
	v := buildChart(props.Chart{
		Title:   "Price",
		YFormat: eur,
		Area:    true,
		Series: []props.ChartSeries{{
			Label: "Asking price",
			Step:  true,
			Color: props.ChartColorPositive,
			// Unsorted on purpose.
			Points: []props.ChartPoint{
				{X: d(3, 20), Y: 1300, Label: "Price drop"},
				{X: d(1, 1), Y: 1350},
				{X: d(2, 14), Y: 1350},
			},
		}},
	})
	if v.Empty || v.Height != 200 || v.Legend {
		t.Fatalf("view = %+v", v)
	}
	if v.YWidest == "" || v.YTicks[0].Label != "1300 €" || v.YTicks[len(v.YTicks)-1].Label != "1350 €" {
		t.Errorf("y ticks = %+v", v.YTicks)
	}
	s := v.Series[0]
	if s.ColorClass != "text-foreground-success" {
		t.Errorf("color = %s", s.ColorClass)
	}
	if !strings.HasPrefix(s.Line, "M0 0H") || !strings.Contains(s.Line, "V1000") {
		t.Errorf("step line = %s", s.Line)
	}
	if !strings.HasSuffix(s.Area, "V1000H0Z") {
		t.Errorf("area = %s", s.Area)
	}
	if len(s.Points) != 3 || s.ShowPoints || !s.Hover {
		t.Errorf("points = %+v", s)
	}
	if got := s.Points[2].Title; got != "20 Mar 2026 · 1300 € · Price drop" {
		t.Errorf("title = %s", got)
	}
	if want := "Price. Asking price: 3 points from 1 Jan 2026 to 20 Mar 2026, between 1300 € and 1350 €, latest 1300 €"; v.AriaLabel != want {
		t.Errorf("aria = %s", v.AriaLabel)
	}

	empty := buildChart(props.Chart{Series: []props.ChartSeries{{Label: "x"}}, Empty: "No prices yet"})
	if !empty.Empty || empty.EmptyText != "No prices yet" || empty.AriaLabel != "Line chart: No prices yet" {
		t.Errorf("empty = %+v", empty)
	}

	single := buildChart(props.Chart{Series: []props.ChartSeries{{Points: []props.ChartPoint{{X: d(1, 1), Y: 5}}}}})
	if p := single.Series[0]; !p.ShowPoints || len(p.Points) != 1 || p.Points[0].X != 50 || math.Abs(p.Points[0].Y-50) > 1e-9 {
		t.Errorf("single = %+v", p)
	}

	spark := buildChart(props.Chart{Sparkline: true, Series: []props.ChartSeries{{Points: []props.ChartPoint{{X: d(1, 1), Y: 1}, {X: d(1, 2), Y: 3}}}}})
	if spark.Height != 32 || len(spark.YTicks) != 0 || len(spark.XTicks) != 0 || spark.Series[0].Line != "M0 900L1000 100" {
		t.Errorf("sparkline = %+v", spark)
	}

	multi := buildChart(props.Chart{Series: []props.ChartSeries{
		{Label: "A", Points: []props.ChartPoint{{X: d(1, 1), Y: 1}}},
		{Label: "B", Points: []props.ChartPoint{{X: d(1, 2), Y: 2}}},
	}})
	if !multi.Legend || multi.Series[0].ColorClass == multi.Series[1].ColorClass {
		t.Errorf("multi = %+v", multi)
	}
	if got := multi.Series[1].Points[0].Title; got != "2 Jan 2026 · B: 2" {
		t.Errorf("multi title = %s", got)
	}
}

func TestChart(t *testing.T) {
	d := func(day int) time.Time { return time.Date(2026, 3, day, 0, 0, 0, 0, time.UTC) }
	got := renderString(t, Chart(props.Chart{
		ID:         "price",
		Title:      "Price history",
		ShowPoints: true,
		Area:       true,
		Series: []props.ChartSeries{
			{Label: "Rent", Points: []props.ChartPoint{{X: d(1), Y: 10}, {X: d(3), Y: 20}}},
			{Label: "Median", Color: props.ChartColorNeutral, Points: []props.ChartPoint{{X: d(2), Y: 15}}},
		},
	}))
	for _, want := range []string{
		`id="price"`,
		`role="img"`,
		`aria-label="Price history. Rent: 2 points`,
		`data-chart="line"`,
		`preserveAspectRatio="none"`,
		`vector-effect="non-scaling-stroke"`,
		`class="text-foreground-accent"`,
		`class="text-foreground-default-secondary"`,
		`fill-opacity="0.12"`,
		`title="3 Mar 2026 · Rent: 20"`,
		`data-chart-point`,
		`style="height: 200px;"`,
		`1 Mar`,
		`Median`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	if strings.Contains(got, "opacity-0") {
		t.Error("ShowPoints should keep points visible")
	}

	empty := renderString(t, Chart(props.Chart{Empty: "Nothing yet"}))
	if !strings.Contains(empty, "Nothing yet") || strings.Contains(empty, "<svg") {
		t.Errorf("empty = %s", empty)
	}

	spark := renderString(t, Chart(props.Chart{Sparkline: true, Class: "w-24", Series: []props.ChartSeries{{Points: []props.ChartPoint{{X: d(1), Y: 1}, {X: d(2), Y: 2}}}}}))
	for _, want := range []string{`data-chart="sparkline"`, `style="height: 32px;"`, `w-24`, `opacity-0`} {
		if !strings.Contains(spark, want) {
			t.Errorf("missing %s in %s", want, spark)
		}
	}
	if strings.Contains(spark, "stroke-border") {
		t.Errorf("sparkline has grid: %s", spark)
	}
}
