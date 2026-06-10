package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/jftuga/ellipsis"
)

func PrintStatistics(statistics *Statistics, timeFormat string) {
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	if statistics.HasWhen {
		table.Header.Cells = append(table.Header.Cells,
			&simpletable.Cell{Align: simpletable.AlignCenter, Text: "WHEN"})
	}
	if statistics.HasID {
		table.Header.Cells = append(table.Header.Cells,
			&simpletable.Cell{Align: simpletable.AlignCenter, Text: "ID"})
	}
	if statistics.HasWho {
		table.Header.Cells = append(table.Header.Cells,
			&simpletable.Cell{Align: simpletable.AlignCenter, Text: "WHO"})
	}
	if statistics.HasEmail {
		table.Header.Cells = append(table.Header.Cells,
			&simpletable.Cell{Align: simpletable.AlignCenter, Text: "EMAIL"})
	}
	table.Header.Cells = append(table.Header.Cells,
		&simpletable.Cell{Align: simpletable.AlignCenter, Text: "COMMITS"},
		&simpletable.Cell{Align: simpletable.AlignCenter, Text: "FILES"},
		&simpletable.Cell{Align: simpletable.AlignCenter, Text: "ADDITIONS"},
		&simpletable.Cell{Align: simpletable.AlignCenter, Text: "DELETIONS"},
	)
	if statistics.HasMessage {
		table.Header.Cells = append(table.Header.Cells,
			&simpletable.Cell{Align: simpletable.AlignCenter, Text: "MESSAGE"})
	}

	for _, s := range statistics.Data {
		line := []*simpletable.Cell{}

		if statistics.HasWhen {
			line = append(line, &simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  s.When.Local().Format(timeFormat),
			})
		}
		if statistics.HasID {
			line = append(line, &simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  s.ID,
			})
		}
		if statistics.HasWho {
			line = append(line, &simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  ellipsis.Shorten(s.Who, 32),
			})
		}
		if statistics.HasEmail {
			line = append(line, &simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  ellipsis.Shorten(s.Email, 32),
			})
		}
		line = append(line,
			&simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  strconv.FormatInt(s.Commits, 10),
			},
			&simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  strconv.FormatInt(s.ModifiedFiles, 10),
			},
			&simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  strconv.FormatInt(s.Additions, 10),
			},
			&simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  strconv.FormatInt(s.Deletions, 10),
			},
		)
		if statistics.HasMessage {
			line = append(line, &simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  s.Message,
			})
		}

		table.Body.Cells = append(table.Body.Cells, line)
	}

	fmt.Fprintln(os.Stdout, table.String())
}

func PrintHeatMapDayHour(
	heatMap map[time.Weekday]map[int]*Statistic,
	display func(stat *Statistic) string,
) {
	disp := display
	if disp == nil {
		disp = func(stat *Statistic) string {
			return strconv.FormatInt(stat.Commits, 10)
		}
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	table.Header.Cells = append(table.Header.Cells,
		&simpletable.Cell{Align: simpletable.AlignCenter, Text: "Day\\Hour"})

	for j := range 24 {
		table.Header.Cells = append(table.Header.Cells,
			&simpletable.Cell{Align: simpletable.AlignCenter, Text: strconv.Itoa(j)})
	}
	for i := time.Sunday; i <= time.Saturday; i++ {
		line := make([]*simpletable.Cell, 0, 25)

		line = append(line, &simpletable.Cell{
			Align: simpletable.AlignLeft,
			Text:  i.String(),
		})
		for j := range 24 {
			line = append(line, &simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  disp(heatMap[i][j]),
			})
		}
		table.Body.Cells = append(table.Body.Cells, line)
	}

	fmt.Fprintln(os.Stdout, table.String())
}

func PrintHeatMapMonthDay(
	heatMap map[time.Month]map[time.Weekday]*Statistic,
	display func(stat *Statistic) string,
) {
	disp := display
	if disp == nil {
		disp = func(stat *Statistic) string {
			return strconv.FormatInt(stat.Commits, 10)
		}
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	table.Header.Cells = append(table.Header.Cells,
		&simpletable.Cell{Align: simpletable.AlignCenter, Text: "Month\\Day"})

	for j := time.Sunday; j <= time.Saturday; j++ {
		table.Header.Cells = append(table.Header.Cells,
			&simpletable.Cell{Align: simpletable.AlignCenter, Text: j.String()})
	}
	for i := time.January; i <= time.December; i++ {
		line := make([]*simpletable.Cell, 0, 25)

		line = append(line, &simpletable.Cell{
			Align: simpletable.AlignLeft,
			Text:  i.String(),
		})
		for j := time.Sunday; j <= time.Saturday; j++ {
			line = append(line, &simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  disp(heatMap[i][j]),
			})
		}
		table.Body.Cells = append(table.Body.Cells, line)
	}

	fmt.Fprintln(os.Stdout, table.String())
}

func PrintTags(tags Tags, timeFormat string) {
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "HASH"},
			{Align: simpletable.AlignCenter, Text: "AUTHOR"},
			{Align: simpletable.AlignCenter, Text: "WHEN"},
		},
	}

	for _, t := range tags {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: t.Name},
			{Align: simpletable.AlignLeft, Text: t.ID},
			{Align: simpletable.AlignLeft, Text: ellipsis.Shorten(t.Author, 32)},
			{Align: simpletable.AlignLeft, Text: t.When.Local().Format(timeFormat)},
		})
	}
	fmt.Fprintln(os.Stdout, table.String())
}
