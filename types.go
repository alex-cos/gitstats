package main

import (
	"time"
)

type FileStat struct {
	Name     string
	Addition int
	Deletion int
}

type FileStats []FileStat

type Commit struct {
	When      time.Time
	Who       string
	Email     string
	ID        string
	Message   string
	FileStats FileStats
}

type Commits []*Commit

type Tag struct {
	Name   string
	ID     string
	Author string
	When   time.Time
}

type Tags []*Tag

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
}

type Statistics struct {
	HasWhen    bool
	HasWho     bool
	HasEmail   bool
	HasID      bool
	HasMessage bool
	Data       []*Statistic
}
