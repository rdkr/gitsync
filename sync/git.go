package sync

//go:generate mockgen -destination=../mocks/mock_git.go -package=mocks gitsync/sync Git

import (
	"bytes"

	"gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	StatusError = iota
	StatusCloned
	StatusFetched
	StatusUpToDate
)

type Status struct {
	Path   string
	Status int
	Output string
	Err    error
}

// Git interface for network operations
type Git interface {
	PlainOpen() (*git.Repository, error)
	PlainClone() (string, error)
	Fetch(*git.Repository) (string, error)
	Pull(*git.Worktree) (string, error)
}

func (p Project) PlainOpen() (*git.Repository, error) {
	return git.PlainOpen(p.Location)
}

func (p Project) PlainClone() (string, error) {

	var buf bytes.Buffer

	_, err := git.PlainClone(p.Location, false, &git.CloneOptions{
		URL:      p.URL,
		Progress: &buf,
		Auth:     p.getAuth(),
	})

	return buf.String(), err
}

func (p Project) Fetch(repo *git.Repository) (string, error) {

	var buf bytes.Buffer

	err := repo.Fetch(&git.FetchOptions{
		Progress: &buf,
		Auth:     p.getAuth(),
	})

	return buf.String(), err
}

func (p Project) Pull(worktree *git.Worktree) (string, error) {

	var buf bytes.Buffer

	err := worktree.Pull(&git.PullOptions{
		Progress: &buf,
		Auth:     p.getAuth(),
	})

	return buf.String(), err
}

func (p Project) getAuth() *git_http.BasicAuth {
	if p.Token != "" {
		return &git_http.BasicAuth{
			Username: "token",
			Password: p.Token,
		}
	}
	return nil
}
