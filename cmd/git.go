package cmd

//go:generate mockgen -destination=../mocks/mock_git.go -package=mocks gitsync/cmd Git

import (
	"bytes"

	"gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Git interface for network operations
type Git interface {
	PlainOpen() (*git.Repository, error)
	PlainClone() (string, error)
	Fetch(*git.Repository) (string, error)
	Pull(*git.Worktree) (string, error)
}

func (p project) PlainOpen() (*git.Repository, error) {
	return git.PlainOpen(p.Location)
}

func (p project) PlainClone() (string, error) {

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

	return buf.String(), err
}

func (p project) Fetch(repo *git.Repository) (string, error) {

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

	return buf.String(), err
}

func (p project) Pull(worktree *git.Worktree) (string, error) {

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

	return buf.String(), err
}
