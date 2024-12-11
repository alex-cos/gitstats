package main

import (
	"sort"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02T15:04:05"

var (
	separator    = "|"
	excludePaths = [...]string{"public/fonts", "public/images", "node_modules", "build", "dist"}
	AllFiles     = map[string]string{}
)

type Statistic struct {
	When          time.Time
	Who           string
	Email         string
	ID            string
	Message       string
	Commits       int64
	ModifiedFiles int64
	Additions     int64
	Deletions     int64
	TotalLines    int64
	ByExt         map[string]struct {
		Additions int64
		Deletions int64
	}
}

type Statistics []*Statistic

func IsInexcludePaths(path string) bool {
	for _, excludePath := range excludePaths {
		if strings.Contains(path, excludePath) {
			return true
		}
	}
	return false
}

func ProduceStats(commits Commits) (Statistics, error) {
	var totalLines int64

	statistics := Statistics{}

	for _, c := range commits {
		var nbModifiedFiles, additions, deletions int64

		stats, err := c.Stats()
		if err != nil {
			return nil, err
		}

		for _, stat := range stats {
			name := stat.Name
			filename := strings.Split(name, " => ")[0]
			if !IsInexcludePaths(filename) {
				AllFiles[filename] = ""
				nbModifiedFiles++
				additions += int64(stat.Addition)
				deletions += int64(stat.Deletion)
			}
		}
		totalLines += additions - deletions

		statistics = append(statistics, &Statistic{
			When:          c.Author.When,
			Who:           strings.ToLower(c.Author.Name),
			Email:         strings.ToLower(c.Author.Email),
			ID:            c.ID().String(),
			Message:       c.Message,
			Commits:       1,
			ModifiedFiles: nbModifiedFiles,
			Additions:     additions,
			Deletions:     deletions,
			TotalLines:    totalLines,
		})
	}

	sort.Slice(statistics, func(i, j int) bool {
		return (statistics[i].When.UnixNano() < statistics[j].When.UnixNano())
	})

	return statistics, nil
}

func AggregByDay(statistics Statistics) Statistics {
	results := Statistics{}
	aggreg := map[time.Time]*Statistic{}

	for _, s := range statistics {
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
	minimum := truncateToDay(statistics[0].When)
	maximum := truncateToDay(statistics[len(statistics)-1].When)
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
	for _, s := range aggreg {
		results = append(results, s)
	}
	sort.Slice(results, func(i, j int) bool {
		return (results[i].When.UnixNano() < results[j].When.UnixNano())
	})
	return results
}

func AggregByAuthor(statistics Statistics) Statistics {
	results := Statistics{}
	aggreg := map[string]*Statistic{}

	for _, s := range statistics {
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
				Who:           strings.ToLower(s.Email),
				Email:         s.Email,
				Commits:       s.Commits,
				ModifiedFiles: s.ModifiedFiles,
				Additions:     s.Additions,
				Deletions:     s.Deletions,
			}
		}
	}
	for _, s := range aggreg {
		results = append(results, s)
	}
	sort.Slice(results, func(i, j int) bool {
		return strings.Compare(results[i].Who, results[j].Who) < 0
	})
	return results
}

func HeatMapDayHour(statistics Statistics) map[time.Weekday]map[int]*Statistic {
	aggreg := map[time.Weekday]map[int]*Statistic{}

	for _, s := range statistics {
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
		for j := 0; j < 24; j++ {
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
