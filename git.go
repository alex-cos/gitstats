package main

import (
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func localRepo(path string) (*git.Repository, error) {
	return git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
}

func remoteRepo(url string) (*git.Repository, error) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func Repository(path, url string) (*git.Repository, error) {
	if path != "" {
		return localRepo(path)
	}
	if url != "" {
		return remoteRepo(url)
	}
	return localRepo(".")
}

func RetrieveCommits(repo *git.Repository, author string, since, until *time.Time) (Commits, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash(), Since: since, Until: until})
	if err != nil {
		return nil, err
	}

	commits := Commits{}
	err = cIter.ForEach(func(c *object.Commit) error {
		fs, err := c.Stats()
		if err != nil {
			return err
		}
		authorName := strings.ToLower(c.Author.Name)
		if author != "" && author != authorName {
			return nil
		}
		stats := make(FileStats, 0, len(fs))
		for _, s := range fs {
			stats = append(stats, FileStat{
				Name:     s.Name,
				Addition: s.Addition,
				Deletion: s.Deletion,
			})
		}
		commits = append(commits, &Commit{
			When:      c.Author.When,
			Who:       authorName,
			Email:     strings.ToLower(c.Author.Email),
			ID:        c.ID().String(),
			Message:   c.Message,
			FileStats: stats,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(commits, func(i, j int) bool {
		return (commits[i].When.UnixNano() < commits[j].When.UnixNano())
	})

	return commits, nil
}

func retrieveTags(repo *git.Repository, direction string) (Tags, error) {
	tagrefs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	tags := Tags{}
	err = tagrefs.ForEach(func(ref *plumbing.Reference) error {
		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}
		tags = append(tags, &Tag{
			Name:   strings.ReplaceAll(ref.Name().String(), "refs/tags/", ""),
			ID:     ref.Hash().String(),
			Author: strings.ToLower(commit.Author.Name),
			When:   commit.Author.When,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	dir := strings.ToLower(strings.TrimSpace(direction))
	sort.Slice(tags, func(i, j int) bool {
		if dir == "desc" {
			return (tags[i].When.UnixNano() > tags[j].When.UnixNano())
		}
		return (tags[i].When.UnixNano() < tags[j].When.UnixNano())
	})

	return tags, nil
}
