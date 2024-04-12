package main

import (
	"sort"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Commits []*object.Commit

func localRepo(path string) (*git.Repository, error) {
	// We instantiate a new repository targeting the given path (the .git folder)
	fs := osfs.New(path)
	if _, err := fs.Stat(git.GitDirName); err == nil {
		fs, err = fs.Chroot(git.GitDirName)
		if err != nil {
			return nil, err
		}
	}

	storage := filesystem.NewStorageWithOptions(fs, cache.NewObjectLRUDefault(), filesystem.Options{KeepDescriptors: true})
	defer storage.Close()
	repo, err := git.Open(storage, fs)
	if err != nil {
		return nil, err
	}

	return repo, nil
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

func retrieveCommits(repo *git.Repository) (Commits, error) {
	// ... retrieves the branch pointed by HEAD
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash(), Since: nil, Until: nil})
	if err != nil {
		return nil, err
	}

	commits := Commits{}
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(commits, func(i, j int) bool {
		return (commits[i].Author.When.UnixNano() < commits[j].Author.When.UnixNano())
	})

	return commits, nil
}
