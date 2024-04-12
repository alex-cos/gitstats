// nolint:forbidigo
package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func PrintStatitics(statistics Statistics, timeFormat string) {
	for _, s := range statistics {
		fields := []string{}
		if s.When.UnixNano() > 0 {
			fields = append(fields, s.When.Local().Format(timeFormat))
		}
		if s.ID != "" {
			fields = append(fields, s.ID)
		}
		if s.Who != "" {
			fields = append(fields, s.Who)
		}
		if s.Email != "" {
			fields = append(fields, s.Email)
		}
		fields = append(fields,
			strconv.FormatInt(s.Commits, 10),
			strconv.FormatInt(s.ModifiedFiles, 10),
			strconv.FormatInt(s.Additions, 10),
			strconv.FormatInt(s.Deletions, 10),
		)
		if s.TotalLines != 0 {
			fields = append(fields, strconv.FormatInt(s.TotalLines, 10))
		}
		if s.Message != "" {
			fields = append(fields, s.Message)
		}

		fmt.Println(strings.Join(fields, separator))
	}
}

func PrintHeatMapDayHour(heatMap map[time.Weekday]map[int]*Statistic) {
	fmt.Printf("Hour%s", separator)
	for j := 0; j < 24; j++ {
		fmt.Printf("%02d|", j)
	}
	fmt.Println()
	for i := time.Sunday; i <= time.Saturday; i++ {
		fmt.Printf("%v%s", i, separator)
		for j := 0; j < 24; j++ {
			fmt.Printf("%v%s", heatMap[i][j].Commits, separator)
		}
		fmt.Println()
	}
}
