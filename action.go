package main

import (
	"context"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

func cmdListCommits(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:      "commits",
		Usage:     "list all commits",
		UsageText: "gitstats list <options>",
		Action: func(c context.Context, cmd *cli.Command) error {
			repo, err := Repository(cmd.String("path"), cmd.String("url"))
			if err != nil {
				return err
			}
			author := strings.ToLower(cmd.String("author"))
			s, u := parseTimes(cmd)
			commits, err := RetrieveCommits(repo, author, s, u)
			if err != nil {
				return err
			}
			statistics := ProduceStats(commits)
			SortStats(statistics, "when", cmd.String("sort"))
			PrintStatistics(statistics, "2006-01-02 15:04:05")
			return nil
		},
		Flags: flags,
	}
}

func cmdByDay(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:      "day",
		Usage:     "aggregate git statistics by day",
		UsageText: "gitstats day <options>",
		Action: func(c context.Context, cmd *cli.Command) error {
			repo, err := Repository(cmd.String("path"), cmd.String("url"))
			if err != nil {
				return err
			}
			author := strings.ToLower(cmd.String("author"))
			s, u := parseTimes(cmd)
			commits, err := RetrieveCommits(repo, author, s, u)
			if err != nil {
				return err
			}
			statistics := ProduceStats(commits)
			statistics = AggregByDay(statistics)
			SortStats(statistics, "when", cmd.String("sort"))
			PrintStatistics(statistics, "2006-01-02")
			return nil
		},
		Flags: flags,
	}
}

func cmdByAuthor(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:      "author",
		Usage:     "aggregate git statistics by author",
		UsageText: "gitstats author <options>",
		Action: func(c context.Context, cmd *cli.Command) error {
			repo, err := Repository(cmd.String("path"), cmd.String("url"))
			if err != nil {
				return err
			}
			author := strings.ToLower(cmd.String("author"))
			s, u := parseTimes(cmd)
			commits, err := RetrieveCommits(repo, author, s, u)
			if err != nil {
				return err
			}
			statistics := ProduceStats(commits)
			statistics = AggregByAuthor(statistics)
			SortStats(statistics, "author", cmd.String("sort"))
			PrintStatistics(statistics, "2006-01-02")
			return nil
		},
		Flags: flags,
	}
}

func cmdHeatMap(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:      "heatmap",
		Usage:     "aggregate git statistics by heatmap",
		UsageText: "gitstats heatmap <options>",
		Action: func(c context.Context, cmd *cli.Command) error {
			repo, err := Repository(cmd.String("path"), cmd.String("url"))
			if err != nil {
				return err
			}
			author := strings.ToLower(cmd.String("author"))
			s, u := parseTimes(cmd)
			commits, err := RetrieveCommits(repo, author, s, u)
			if err != nil {
				return err
			}
			statistics := ProduceStats(commits)
			PrintHeatMapDayHour(HeatMapDayHour(statistics), nil)
			return nil
		},
		Flags: flags,
	}
}

func cmdListTags(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:      "tags",
		Usage:     "list all tags",
		UsageText: "gitstats tags <options>",
		Action: func(c context.Context, cmd *cli.Command) error {
			repo, err := Repository(cmd.String("path"), cmd.String("url"))
			if err != nil {
				return err
			}
			tags, err := retrieveTags(repo, cmd.String("sort"))
			if err != nil {
				return err
			}
			PrintTags(tags, "2006-01-02 15:04:05")
			return nil
		},
		Flags: flags,
	}
}

func parseTimes(cmd *cli.Command) (*time.Time, *time.Time) {
	since := cmd.Timestamp("since")
	until := cmd.Timestamp("until")
	s := &since
	if since.Equal(time.Time{}) {
		s = nil
	}
	u := &until
	if until.Equal(time.Time{}) {
		u = nil
	}
	return s, u
}
