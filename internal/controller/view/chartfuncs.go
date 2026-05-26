package view

import (
	"fmt"
	"strings"
)

// svgX returns the SVG x-coordinate for day index i (0–6),
// centered within each equal-width flex slot of the chart-days label row.
func svgX(i int) float64 {
	return float64(2*i+1) * 300.0 / 14.0
}

// svgY maps a metric value (1–10) to an SVG y-coordinate within viewBox 0 0 300 100.
func svgY(val int) float64 {
	return 95 - float64(val-1)*85/9
}

func metricPtr(day ChartDay, metric string) *int {
	switch metric {
	case "energy":
		return day.Energy
	case "sleep":
		return day.Sleep
	case "depression":
		return day.Depression
	case "happiness":
		return day.Happiness
	case "pain":
		return day.Pain
	}
	return nil
}

func ChartLinePath(days []ChartDay, metric string) string {
	var sb strings.Builder
	started := false
	for i, day := range days {
		val := metricPtr(day, metric)
		if val == nil {
			started = false
			continue
		}
		x, y := svgX(i), svgY(*val)
		if !started {
			fmt.Fprintf(&sb, "M%.1f,%.1f", x, y)
			started = true
		} else {
			fmt.Fprintf(&sb, " L%.1f,%.1f", x, y)
		}
	}
	return sb.String()
}

func ChartAreaPath(days []ChartDay, metric string) string {
	var sb strings.Builder
	type seg struct{ from, to int }
	var segs []seg
	start := -1
	for i, day := range days {
		if metricPtr(day, metric) != nil {
			if start < 0 {
				start = i
			}
		} else if start >= 0 {
			segs = append(segs, seg{start, i - 1})
			start = -1
		}
	}
	if start >= 0 {
		segs = append(segs, seg{start, len(days) - 1})
	}
	for _, s := range segs {
		startVal := metricPtr(days[s.from], metric)
		fmt.Fprintf(&sb, "M%.1f,%.1f", svgX(s.from), svgY(*startVal))
		for i := s.from + 1; i <= s.to; i++ {
			val := metricPtr(days[i], metric)
			fmt.Fprintf(&sb, " L%.1f,%.1f", svgX(i), svgY(*val))
		}
		fmt.Fprintf(&sb, " L%.1f,95 L%.1f,95 Z", svgX(s.to), svgX(s.from))
	}
	return sb.String()
}

func ChartHasDot(days []ChartDay, metric string) bool {
	for _, day := range days {
		if metricPtr(day, metric) != nil {
			return true
		}
	}
	return false
}

func ChartDotX(days []ChartDay, metric string) float64 {
	for i := len(days) - 1; i >= 0; i-- {
		if metricPtr(days[i], metric) != nil {
			return svgX(i)
		}
	}
	return 0
}

func ChartDotY(days []ChartDay, metric string) float64 {
	for i := len(days) - 1; i >= 0; i-- {
		if val := metricPtr(days[i], metric); val != nil {
			return svgY(*val)
		}
	}
	return 0
}
