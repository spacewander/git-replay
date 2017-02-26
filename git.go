package main

import (
	go_git "github.com/src-d/go-git"
	"srcd.works/go-git.v4/plumbing"
)

var (
	repo *go_git.Repository
)

func InitRepo(path string) (err error) {
	repo, err = go_git.PlainOpen(path)
	if err != nil {
		return err
	}
	return nil
}

func SearchCommit(hash string) (string, error) {
	commit, err := repo.Commit(plumbing.NewHash(hash))
	if err != nil {
		return "", err
	}
	return commit.String(), nil
}
