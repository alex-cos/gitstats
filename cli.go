// nolint:forbidigo,exhaustivestruct
package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

var (
	version   = "unknown"
	buildDate = "unknown"
)

// parseCLI parses command lines arguments.
func parseCLI() error {
	cliapp := cli.NewApp()
	cliapp.Name = "gitstats"
	cliapp.Usage = "Produces statistics for git repository"
	cliapp.UsageText = "gitstats <command> [options]"
	cliapp.Description = "Build: %s" + buildDate
	cliapp.Version = version

	cliapp.Flags = []cli.Flag{
		&cli.PathFlag{
			Name:     "path",
			Usage:    "Path to git repository containing '.git' folder",
			Aliases:  []string{"p"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "url",
			Usage:    "URL to git repository containing",
			Aliases:  []string{"u"},
			Required: false,
		},
	}
	cliapp.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Printf("Error. Unknown command: '%s'\n\n", command)
		cli.ShowAppHelpAndExit(c, 1)
	}
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("Version:\t", c.App.Version)
		fmt.Println("Build Date:\t", buildDate)
	}
	cliapp.Commands = []*cli.Command{
		cmdListCommits(cliapp.Flags),
		cmdByDay(cliapp.Flags),
		cmdByAuthor(cliapp.Flags),
	}
	sort.Sort(cli.FlagsByName(cliapp.Flags))
	sort.Sort(cli.CommandsByName(cliapp.Commands))

	if err := cliapp.Run(os.Args); err != nil {
		return fmt.Errorf("failed to parse command line arguments: %w", err)
	}
	return nil
}

func cmdListCommits(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:      "commits",
		Usage:     "list all commits",
		UsageText: "gitstats list <options>",
		Action: func(c *cli.Context) error {
			repo, err := Repository(c.Path("path"), c.String("url"))
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
		Action: func(c *cli.Context) error {
			repo, err := Repository(c.Path("path"), c.String("url"))
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
		Action: func(c *cli.Context) error {
			repo, err := Repository(c.Path("path"), c.String("url"))
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
