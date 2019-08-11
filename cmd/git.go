package cmd

//go:generate mockgen -destination=../mocks/mock_git.go -package=mocks gitsync/cmd Git

import (
	"bytes"
	"errors"
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type Git interface {
	PlainOpen() (*git.Repository, error)
	PlainClone() Status
	Fetch(*git.Repository) Status
	Pull(*git.Worktree) Status
}

func (p project) PlainOpen() (*git.Repository, error) {
	return git.PlainOpen(p.Location)
}

func (p project) PlainClone() Status {

	var auth *git_http.BasicAuth
	if p.Token != "" {
		auth = &git_http.BasicAuth{
			Username: "token",
			Password: p.Token,
		}
	}

	var buf bytes.Buffer

	_, err := git.PlainClone(p.Location, false, &git.CloneOptions{
		URL:      p.URL,
		Progress: &buf,
		Auth:     auth,
	})
	if err != nil {
		return Status{p.Location, "", buf.String(), fmt.Errorf("unable to clone repo: %v", err)}
	}
	return Status{p.Location, "cloned", buf.String(), nil}

}

func (p project) Fetch(repo *git.Repository) Status {

	var auth *git_http.BasicAuth
	if p.Token != "" {
		auth = &git_http.BasicAuth{
			Username: "token",
			Password: p.Token,
		}
	}

	var buf bytes.Buffer

	err := repo.Fetch(&git.FetchOptions{
		Progress: &buf,
		Auth:     auth,
	})
	if err == git.NoErrAlreadyUpToDate {
		return Status{p.Location, "", buf.String(), errors.New("not on master branch but fetched")}
	}
	return Status{p.Location, "", buf.String(), fmt.Errorf("not on master branch and: %v", err)}

}

func (p project) Pull(worktree *git.Worktree) Status {

	var auth *git_http.BasicAuth
	if p.Token != "" {
		auth = &git_http.BasicAuth{
			Username: "token",
			Password: p.Token,
		}
	}

	var buf bytes.Buffer

	err := worktree.Pull(&git.PullOptions{
		Progress: &buf,
		Auth:     auth,
	})
	if err == nil {
		return Status{p.Location, "fetched", buf.String(), nil}
	} else if err == git.NoErrAlreadyUpToDate {
		return Status{p.Location, "uptodate", buf.String(), nil}
	}
	return Status{p.Location, "", buf.String(), fmt.Errorf("unable to pull master: %v", err)}

}

func Sync(p Git, location string) Status {

	repo, err := p.PlainOpen()

	if err == git.ErrRepositoryNotExists {
		return p.PlainClone()
	} else if err != nil {
		return Status{location, "", "", fmt.Errorf("unable to open repo: %v", err)}
	}

	ref, err := repo.Head()
	if err != nil {
		return Status{location, "", "", fmt.Errorf("unable to get head: %v", err)}
	}

	if ref.Name() != "refs/heads/master" {
		return p.Fetch(repo)
	}

	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return Status{location, "", "", fmt.Errorf("unable to get worktree: %v", err)}
	}

	return p.Pull(worktree)

}
