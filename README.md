# gitstats

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://go.dev/)
[![Test Status](https://github.com/alex-cos/gitstats/actions/workflows/test.yml/badge.svg)](https://github.com/alex-cos/gitstats/actions/workflows/test.yml)
[![Lint Status](https://github.com/alex-cos/gitstats/actions/workflows/lint.yml/badge.svg)](https://github.com/alex-cos/gitstats/actions/workflows/lint.yml)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/alex-cos/gitstats)](https://goreportcard.com/report/github.com/alex-cos/gitstats)

A command-line tool written in Go that produces statistics for a Git repository.

## Features

- **List commits** — Display all commits with author, email, date, additions, deletions, and file count.
- **Aggregate by day** — Group commit statistics per calendar day.
- **Aggregate by author** — Group commit statistics per author.
- **Heatmap** — Generate a day/hour heatmap of commit activity.
- **List tags** — Display all tags with hash, author, and date.

## Usage

```bash
gitstats [command] [options]
```

### Global options

| Flag | Alias | Description |
| ------ | ------- | ------------- |
| `--path` | `-p` | Path to a local Git repository. Defaults to current directory. |
| `--url` | `-u` | URL of a remote Git repository to clone. |
| `--since` | `--after` | Include commits after this date (format: `2006-01-02` or `2006/01/02`). |
| `--until` | `--before` | Include commits before this date (format: `2006-01-02` or `2006/01/02`). |
| `--sort` | | Sort direction: `asc` (default) or `desc`. |
| `--version` | | Show version and build date. |

### Commands

#### `commits`

List all commits of the repository.

```bash
gitstats commits --path /path/to/repo
gitstats commits --since 2024-01-01 --until 2024-06-30
```

#### `day`

Aggregate commit statistics by day.

```bash
gitstats day --path /path/to/repo --sort desc
```

#### `author`

Aggregate commit statistics by author.

```bash
gitstats author --path /path/to/repo
```

#### `heatmap`

Generate a day-of-week / hour heatmap of commit activity.

```bash
gitstats heatmap --path /path/to/repo
```

#### `tags`

List all tags of the repository.

```bash
gitstats tags --path /path/to/repo
```

### Remote repository

You can pass a URL instead of a local path:

```bash
gitstats commits --url https://github.com/user/repo.git
```

## Excluded paths

The following paths are excluded from file counts by default:

- `public/fonts`
- `public/images`
- `node_modules`
- `build`
- `dist`

## Development

### Prerequisites

- Go 1.25+
- [golangci-lint](https://golangci-lint.run/) (for linting)

### Makefile targets

| Target | Description |
| -------- | ------------- |
| `make build` | Build the binary |
| `make test` | Run tests |
| `make lint` | Run golangci-lint |
| `make critic` | Run go-critic |
| `make deadcode` | Run deadcode detection |
| `make doc` | Start godoc server on `:8085` |
| `make clean` | Remove build artifacts and vendor |
| `make install-tools` | Install development tools |
