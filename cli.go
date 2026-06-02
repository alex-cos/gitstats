// nolint:forbidigo,exhaustivestruct
package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v3"
)

var (
	version   = "unknown"
	buildDate = "unknown"
)

// parseCLI parses command lines arguments.
func parseCLI() error {
	appcmd := &cli.Command{
		Name:        "gitstats",
		Usage:       "Produces statistics for git repository",
		UsageText:   "gitstats <command> [options]",
		Description: "Build: " + buildDate,
		Version:     version,
		CommandNotFound: func(c context.Context, cmd *cli.Command, name string) {
			fmt.Fprintf(os.Stderr, "Error. Unknown command: '%s'\n\n", name)
			cli.ShowAppHelpAndExit(cmd, 1)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "path",
				Usage:     "Path to git repository containing '.git' folder",
				Aliases:   []string{"p"},
				Required:  false,
				TakesFile: true,
			},
			&cli.StringFlag{
				Name:     "url",
				Usage:    "URL to git repository containing",
				Aliases:  []string{"u"},
				Required: false,
			},
		},
	}

	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintln(os.Stdout, "Version:\t", cmd.Version)
		fmt.Fprintln(os.Stdout, "Build Date:\t", buildDate)
	}

	appcmd.Commands = []*cli.Command{
		cmdListCommits(appcmd.Flags),
		cmdByDay(appcmd.Flags),
		cmdByAuthor(appcmd.Flags),
	}

	sort.Sort(cli.FlagsByName(appcmd.Flags))
	sort.Slice(appcmd.Commands, func(i, j int) bool {
		return appcmd.Commands[i].Name < appcmd.Commands[j].Name
	})

	if err := appcmd.Run(context.Background(), os.Args); err != nil {
		return fmt.Errorf("failed to parse command line arguments: %w", err)
	}

	return nil
}

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
			commits, err := retrieveCommits(repo)
			if err != nil {
				return err
			}
			statistics, err := ProduceStats(commits)
			if err != nil {
				return err
			}
			PrintStatitics(statistics, "2006-01-02T15:04:05")
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
			commits, err := retrieveCommits(repo)
			if err != nil {
				return err
			}
			statistics, err := ProduceStats(commits)
			if err != nil {
				return err
			}
			PrintStatitics(AggregByDay(statistics), "2006-01-02")
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
			commits, err := retrieveCommits(repo)
			if err != nil {
				return err
			}
			statistics, err := ProduceStats(commits)
			if err != nil {
				return err
			}
			PrintStatitics(AggregByAuthor(statistics), "2006-01-02")
			return nil
		},
		Flags: flags,
	}
}
