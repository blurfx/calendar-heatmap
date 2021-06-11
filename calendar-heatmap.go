package heatmap

import (
	"bytes"
	"fmt"
	svg "github.com/ajstarks/svgo/float"
	"math"
	"time"
)


type CalendarHeatmapConfig struct {
	Colors           []string
	BlockSize        float64
	BlockRoundness   float64
	BlockMargin      float64
	MonthLabels      []string
	MonthLabelHeight float64
	WeekdayLabels    []string
	weekLabelWidth   float64
}

type CalendarHeatmap struct {
	Config *CalendarHeatmapConfig
}

type Date struct {
	Year int
	Month time.Month
	Day int
}

type point struct {
	Y float64
	X float64
}

var defaultConfig = &CalendarHeatmapConfig{
	Colors:           []string{"#EBEDF0", "#9BE9A8", "#40C463", "#30A14E", "#216E39"},
	BlockSize:        11,
	BlockRoundness:	  2,
	BlockMargin:      2,
	MonthLabels:      []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	MonthLabelHeight: 15,
	WeekdayLabels:    []string{"", "Mon", "", "Wed", "", "Fri", ""},
	weekLabelWidth:   0,
}

func New(config *CalendarHeatmapConfig) *CalendarHeatmap {
	if config == nil {
		config = defaultConfig
	}

	return &CalendarHeatmap{
		Config: config,
	}
}

func (date Date) Time() time.Time {
	return time.Date(date.Year, date.Month, date.Day, 0, 0, 0, 0, time.Local)
}

func (c *CalendarHeatmap) getPosition(row, column int) *point {
	bounds := c.Config.BlockSize + c.Config.BlockMargin

	return &point{
		Y: c.Config.MonthLabelHeight + bounds * float64(column),
		X: c.Config.weekLabelWidth + bounds * float64(row),
	}
}

func (c *CalendarHeatmap) diffWeeks(from, to Date) int {
	fromTime := from.Time()
	toTimestamp := to.Time().Unix()
	weeks := 0
	for fromTime.Unix() <= toTimestamp {
		fromTime = fromTime.AddDate(0, 0, 7)
		weeks += 1
	}
	return weeks
}

func (c *CalendarHeatmap) Generate(dateFrom, dateTo Date, data map[Date]int) *bytes.Buffer {
	const days = 7
	const monthLabelFontSize = 10
	const weekdayLabelFontSize = 9
	config := c.Config
	weeks := c.diffWeeks(dateFrom, dateTo)
	currentDate := dateFrom.Time()
	prevMonth := -1

	endTimestamp := dateTo.Time().Unix()
	buffer := new(bytes.Buffer)
	canvas := svg.New(buffer)

	// calculate label width
	for _, s := range config.WeekdayLabels {
		config.weekLabelWidth = math.Max(config.weekLabelWidth, float64(len(s) * weekdayLabelFontSize))
	}

	// if month of the first week is different from the month of the second week,
	// don't draw label of month of the first week.
	if currentDate.Month() != currentDate.AddDate(0, 0, 7).Month() {
		prevMonth = int(currentDate.Month())
	}

	// draw svg
	canvasPos := c.getPosition(weeks, days)
	canvas.Start(canvasPos.X, canvasPos.Y)

	for week := 0; week < weeks && currentDate.Unix() <= endTimestamp; week += 1 {
		// draw month label
		currentMonth := int(currentDate.Month())
		if prevMonth != currentMonth {
			pos := c.getPosition(week, 0)
			prevMonth = currentMonth
			canvas.Text(
				pos.X,
				pos.Y + (config.BlockSize / 2) - config.MonthLabelHeight,
				config.MonthLabels[prevMonth-1],
				fmt.Sprintf("font-size: %dpx;alignment-baseline: central; fill: #aaa;", monthLabelFontSize),
			)
		}

		// draw heatmap blocks
		for currentDate.Weekday() <= time.Saturday && currentDate.Unix() <= endTimestamp {
			fillColor := config.Colors[0]
			pos := c.getPosition(week, int(currentDate.Weekday()) - 1)
			year, month, day := currentDate.Date()
			date := Date{year, month, day}

			if value, ok := data[date]; ok {
				fillColor = config.Colors[value]
			}

			canvas.Roundrect(
				pos.X,
				pos.Y + config.MonthLabelHeight,
				config.BlockSize,
				config.BlockSize,
				config.BlockRoundness,
				config.BlockRoundness,
				fmt.Sprintf("fill:%s", fillColor),
			)

			currentDate = currentDate.AddDate(0, 0, 1)
			if currentDate.Weekday() == time.Sunday {
				break
			}
		}
	}

	// draw weekday labels
	for day := 0; day < days; day++ {
		pos := c.getPosition(0, day+1)
		style := fmt.Sprintf("font-size: %dpx; fill:#aaa", weekdayLabelFontSize)

		canvas.Text(
			0,
			pos.Y - (config.BlockSize / 2),
			config.WeekdayLabels[day],
			style,
		)
	}

	canvas.End()

	return buffer
}
