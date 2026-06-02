package main

import (
	"fmt"
	"os"
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

		fmt.Fprintln(os.Stdout, strings.Join(fields, separator))
	}
}

func PrintHeatMapDayHour(heatMap map[time.Weekday]map[int]*Statistic) {
	fmt.Fprintf(os.Stdout, "Hour%s", separator)
	for j := range 24 {
		fmt.Fprintf(os.Stdout, "%02d|", j)
	}
	fmt.Fprintln(os.Stdout)
	for i := time.Sunday; i <= time.Saturday; i++ {
		fmt.Fprintf(os.Stdout, "%v%s", i, separator)
		for j := range 24 {
			fmt.Fprintf(os.Stdout, "%v%s", heatMap[i][j].Commits, separator)
		}
		fmt.Fprintln(os.Stdout)
	}
}
