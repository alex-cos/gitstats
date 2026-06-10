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
			&cli.StringFlag{
				Name:     "author",
				Usage:    "filter stats for a given author",
				Aliases:  []string{"a"},
				Required: false,
			},
			&cli.TimestampFlag{
				Name:     "since",
				Usage:    "commits to retrieve since the given date",
				Aliases:  []string{"after"},
				Required: false,
				Config: cli.TimestampConfig{
					Layouts: []string{"2006-01-02", "2006/01/02"},
				},
			},
			&cli.TimestampFlag{
				Name:     "until",
				Usage:    "commits to retrieve until the given date",
				Aliases:  []string{"before"},
				Required: false,
				Config: cli.TimestampConfig{
					Layouts: []string{"2006-01-02", "2006/01/02"},
				},
			},
			&cli.StringFlag{
				Name:        "sort",
				Usage:       "sort direction ['asc' for ascending, 'desc' for descending]",
				Aliases:     []string{"s"},
				Required:    false,
				DefaultText: "asc",
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
		cmdHeatMapDayHour(appcmd.Flags),
		cmdHeatMapMonthDay(appcmd.Flags),
		cmdListTags(appcmd.Flags),
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
