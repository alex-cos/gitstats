package main

import (
	"slices"
	"sort"
	"strings"
	"time"
)

var (
	excludePaths = [...]string{
		"bin",
		"node_modules",
		"build",
		"dist",
		"tmp",
		"vendor",
		"__pycache__",
	}
)

func IsInExcludePaths(path string) bool {
	for _, excludePath := range excludePaths {
		list := strings.Split(path, "/")
		if slices.Contains(list, excludePath) {
			return true
		}
	}
	return false
}

func ProduceStats(commits Commits) *Statistics {
	statistics := Statistics{
		HasWhen:    true,
		HasWho:     true,
		HasEmail:   true,
		HasID:      true,
		HasMessage: true,
		Data:       []*Statistic{},
	}

	for _, c := range commits {
		var nbModifiedFiles, additions, deletions int64
		for _, stat := range c.FileStats {
			name := stat.Name
			filename, _, _ := strings.Cut(name, " => ")
			if !IsInExcludePaths(filename) {
				nbModifiedFiles++
				additions += int64(stat.Addition)
				deletions += int64(stat.Deletion)
			}
		}

		statistics.Data = append(statistics.Data, &Statistic{
			When:          c.When,
			Who:           c.Who,
			Email:         c.Email,
			ID:            c.ID,
			Message:       c.Message,
			Commits:       1,
			ModifiedFiles: nbModifiedFiles,
			Additions:     additions,
			Deletions:     deletions,
		})
	}

	sort.Slice(statistics.Data, func(i, j int) bool {
		return (statistics.Data[i].When.UnixNano() < statistics.Data[j].When.UnixNano())
	})

	return &statistics
}

func AggregByDay(statistics *Statistics) *Statistics {
	aggreg := map[time.Time]*Statistic{}

	for _, s := range statistics.Data {
		day := truncateToDay(s.When)
		_, ok := aggreg[day]
		if ok {
			aggreg[day].Commits += s.Commits
			aggreg[day].ModifiedFiles += s.ModifiedFiles
			aggreg[day].Additions += s.Additions
			aggreg[day].Deletions += s.Deletions
		} else {
			aggreg[day] = &Statistic{
				When:          day,
				Commits:       s.Commits,
				ModifiedFiles: s.ModifiedFiles,
				Additions:     s.Additions,
				Deletions:     s.Deletions,
			}
		}
	}
	minimum := truncateToDay(statistics.Data[0].When)
	maximum := truncateToDay(statistics.Data[len(statistics.Data)-1].When)
	day := minimum
	for day.UnixNano() < maximum.UnixNano() {
		_, ok := aggreg[day]
		if !ok {
			aggreg[day] = &Statistic{
				When:          day,
				Commits:       0,
				ModifiedFiles: 0,
				Additions:     0,
				Deletions:     0,
			}
		}
		day = day.AddDate(0, 0, 1)
	}

	results := Statistics{
		HasWhen:    true,
		HasWho:     false,
		HasEmail:   false,
		HasID:      false,
		HasMessage: false,
		Data:       []*Statistic{},
	}
	for _, s := range aggreg {
		results.Data = append(results.Data, s)
	}
	sort.Slice(results.Data, func(i, j int) bool {
		return (results.Data[i].When.UnixNano() < results.Data[j].When.UnixNano())
	})
	return &results
}

func AggregByAuthor(statistics *Statistics) *Statistics {
	aggreg := map[string]*Statistic{}

	for _, s := range statistics.Data {
		author := Capitalize(s.Who)
		_, ok := aggreg[author]
		if ok {
			aggreg[author].Commits += s.Commits
			aggreg[author].ModifiedFiles += s.ModifiedFiles
			aggreg[author].Additions += s.Additions
			aggreg[author].Deletions += s.Deletions
		} else {
			aggreg[author] = &Statistic{
				When:          time.Unix(0, 0),
				Who:           s.Who,
				Email:         s.Email,
				Commits:       s.Commits,
				ModifiedFiles: s.ModifiedFiles,
				Additions:     s.Additions,
				Deletions:     s.Deletions,
			}
		}
	}

	results := Statistics{
		HasWhen:    false,
		HasWho:     true,
		HasEmail:   true,
		HasID:      false,
		HasMessage: false,
		Data:       []*Statistic{},
	}
	for _, s := range aggreg {
		results.Data = append(results.Data, s)
	}
	sort.Slice(results.Data, func(i, j int) bool {
		return strings.Compare(results.Data[i].Who, results.Data[j].Who) < 0
	})
	return &results
}

func HeatMapDayHour(statistics *Statistics) map[time.Weekday]map[int]*Statistic {
	aggreg := map[time.Weekday]map[int]*Statistic{}

	for _, s := range statistics.Data {
		weekday := truncateToDay(s.When).Weekday()
		hour := s.When.Hour()
		_, ok := aggreg[weekday]
		if ok {
			_, ok := aggreg[weekday][hour]
			if ok {
				aggreg[weekday][hour].Commits += s.Commits
				aggreg[weekday][hour].ModifiedFiles += s.ModifiedFiles
				aggreg[weekday][hour].Additions += s.Additions
				aggreg[weekday][hour].Deletions += s.Deletions
			} else {
				aggreg[weekday][hour] = &Statistic{
					When:          time.Unix(0, 0),
					Commits:       s.Commits,
					ModifiedFiles: s.ModifiedFiles,
					Additions:     s.Additions,
					Deletions:     s.Deletions,
				}
			}
		} else {
			aggreg[weekday] = map[int]*Statistic{}
			aggreg[weekday][hour] = &Statistic{
				When:          time.Unix(0, 0),
				Commits:       s.Commits,
				ModifiedFiles: s.ModifiedFiles,
				Additions:     s.Additions,
				Deletions:     s.Deletions,
			}
		}
	}
	for i := time.Sunday; i <= time.Saturday; i++ {
		_, ok := aggreg[i]
		if !ok {
			aggreg[i] = map[int]*Statistic{}
		}
		for j := range 24 {
			_, ok := aggreg[i][j]
			if !ok {
				aggreg[i][j] = &Statistic{
					When:          time.Unix(0, 0),
					Commits:       0,
					ModifiedFiles: 0,
					Additions:     0,
					Deletions:     0,
				}
			}
		}
	}
	return aggreg
}

func HeatMapMonthDay(statistics *Statistics) map[time.Month]map[time.Weekday]*Statistic {
	aggreg := map[time.Month]map[time.Weekday]*Statistic{}

	for _, s := range statistics.Data {
		month := s.When.Month()
		weekday := truncateToDay(s.When).Weekday()
		_, ok := aggreg[month]
		if ok {
			_, ok := aggreg[month][weekday]
			if ok {
				aggreg[month][weekday].Commits += s.Commits
				aggreg[month][weekday].ModifiedFiles += s.ModifiedFiles
				aggreg[month][weekday].Additions += s.Additions
				aggreg[month][weekday].Deletions += s.Deletions
			} else {
				aggreg[month][weekday] = &Statistic{
					When:          time.Unix(0, 0),
					Commits:       s.Commits,
					ModifiedFiles: s.ModifiedFiles,
					Additions:     s.Additions,
					Deletions:     s.Deletions,
				}
			}
		} else {
			aggreg[month] = map[time.Weekday]*Statistic{}
			aggreg[month][weekday] = &Statistic{
				When:          time.Unix(0, 0),
				Commits:       s.Commits,
				ModifiedFiles: s.ModifiedFiles,
				Additions:     s.Additions,
				Deletions:     s.Deletions,
			}
		}
	}
	for i := time.January; i <= time.December; i++ {
		_, ok := aggreg[i]
		if !ok {
			aggreg[i] = map[time.Weekday]*Statistic{}
		}
		for j := time.Sunday; j <= time.Saturday; j++ {
			_, ok := aggreg[i][j]
			if !ok {
				aggreg[i][j] = &Statistic{
					When:          time.Unix(0, 0),
					Commits:       0,
					ModifiedFiles: 0,
					Additions:     0,
					Deletions:     0,
				}
			}
		}
	}
	return aggreg
}

func SortStats(statistics *Statistics, by, direction string) {
	dir := strings.ToLower(strings.TrimSpace(direction))

	switch by {
	case "author":
		sort.Slice(statistics.Data, func(i, j int) bool {
			if dir == "desc" {
				return strings.Compare(statistics.Data[i].Who, statistics.Data[j].Who) > 0
			}
			return strings.Compare(statistics.Data[i].Who, statistics.Data[j].Who) < 0
		})
	default:
		sort.Slice(statistics.Data, func(i, j int) bool {
			if dir == "desc" {
				return (statistics.Data[i].When.UnixNano() > statistics.Data[j].When.UnixNano())
			}
			return (statistics.Data[i].When.UnixNano() < statistics.Data[j].When.UnixNano())
		})
	}
}
