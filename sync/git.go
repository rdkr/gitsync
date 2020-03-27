package sync

//go:generate mockgen -destination=../mocks/mock_git.go -package=mocks gitsync/sync Git

import (
	"bytes"

	"github.com/rdkr/gitsync/concurrency"
	"gopkg.in/src-d/go-git.v4"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type GitSyncProject struct {
	concurrency.Project
}

// Git interface for network operations
type Git interface {
	PlainOpen() (*git.Repository, error)
	PlainClone() (string, error)
	Fetch(*git.Repository) (string, error)
	Pull(*git.Worktree) (string, error)
}

func (p GitSyncProject) PlainOpen() (*git.Repository, error) {
	return git.PlainOpen(p.Location)
}

func (p GitSyncProject) PlainClone() (string, error) {

	var buf bytes.Buffer

	_, err := git.PlainClone(p.Location, false, &git.CloneOptions{
		URL:      p.URL,
		Progress: &buf,
		Auth:     p.getAuth(),
	})

	return buf.String(), err
}

func (p GitSyncProject) Fetch(repo *git.Repository) (string, error) {

	var buf bytes.Buffer

	err := repo.Fetch(&git.FetchOptions{
		Progress: &buf,
		Auth:     p.getAuth(),
	})

	return buf.String(), err
}

func (p GitSyncProject) Pull(worktree *git.Worktree) (string, error) {

	var buf bytes.Buffer

	err := worktree.Pull(&git.PullOptions{
		Progress: &buf,
		Auth:     p.getAuth(),
	})

	return buf.String(), err
}

func (p GitSyncProject) getAuth() *githttp.BasicAuth {
	if p.Token != "" {
		return &githttp.BasicAuth{
			Username: "token",
			Password: p.Token,
		}
	}
	return nil
}
