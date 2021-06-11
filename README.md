# Calendar-Heatmap

Generate svg image of calendar heatmap like GitHub contribution graph.

## Guide

### Installation
```
go get github.com/blurfx/calendar-heatmap
```


### Example

```go
package main

import (
	"fmt"
	"github.com/blurfx/calendar-heatmap"
)

func main() {
	data := make(map[heatmap.Date]int)

	data[heatmap.Date{Year: 2021, Month: 1, Day: 1}] = 1
	data[heatmap.Date{Year: 2021, Month: 2, Day: 21}] = 3
	data[heatmap.Date{Year: 2021, Month: 3, Day: 17}] = 4

	h := heatmap.New(nil)
	buffer := h.Generate(
		heatmap.Date{Year: 2021, Month: 1, Day: 1},
		heatmap.Date{Year: 2021, Month: 3, Day: 31},
		data,
    )

	fmt.Println(buffer.String())
	// <svg ...>...</svg>
}
```

### Heatmap Customization

You can customize heat maps using the CalendarHeatmapConfig structure.

#### Default Configs
```go
&CalendarHeatmapConfig{
	Colors:           []string{"#EBEDF0", "#9BE9A8", "#40C463", "#30A14E", "#216E39"},
	BlockSize:        11,
	BlockMargin:      2,
	BlockRoundness:   2,
	MonthLabels:      []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	MonthLabelHeight: 15,
	WeekdayLabels:    []string{"", "Mon", "", "Wed", "", "Fri", ""},
}
```

#### Using your own configuration 
```go
myConfig := &heatmap.CalendarHeatmapConfig{
    // ...
}

h := heatmap.New(myConfig)
```

## Example Output
![Generated SVG](https://raw.githubusercontent.com/blurfx/calendar-heatmap/static/demo.svg)
