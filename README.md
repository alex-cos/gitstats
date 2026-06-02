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
- **Heatmap** — Generate a day/hour heatmap of commit activity (function available, not exposed via CLI).

## Installation

```bash
make build
```

The binary is output to `bin/gitstats`.

## Usage

```bash
gitstats [command] [options]
```

### Options

| Flag | Alias | Description |
| ------ | ------- | ------------- |
| `--path` | `-p` | Path to a local Git repository (must contain a `.git` folder). Defaults to current directory. |
| `--url` | `-u` | URL of a remote Git repository to clone. |
| `--version` | | Show version and build date. |

### Commands

#### `commits`

List all commits of the repository.

```bash
gitstats commits --path /path/to/repo
```

#### `day`

Aggregate commit statistics by day.

```bash
gitstats day --path /path/to/repo
```

#### `author`

Aggregate commit statistics by author.

```bash
gitstats author --path /path/to/repo
```

### Remote repository

You can pass a URL instead of a local path:

```bash
gitstats commits --url https://github.com/user/repo.git
```

## Output format

Output is printed to stdout in a pipe-delimited format:

```txt
date|id|author|email|commits|files|additions|deletions|total_lines|message
```

Fields that are empty or zero may be omitted depending on the command.

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
| `make lint` | Run linter |
| `make clean` | Remove build artifacts |
| `make install-tools` | Install development tools |
