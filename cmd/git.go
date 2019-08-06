package cmd

import (
	"bytes"
	"errors"
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type Cloner interface {
	Clone() status
	PlainOpen() (*git.Repository, error)
	PlainClone() status
	Fetch(repo *git.Repository) status
	Pull(worktree *git.Worktree) status
}

func (p project) PlainOpen() (*git.Repository, error) {
	return git.PlainOpen(p.Location)
}

func (p project) PlainClone() status {

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
		return status{p.Location, "", buf.String(), fmt.Errorf("unable to clone repo: %v", err)}
	}
	return status{p.Location, "cloned", buf.String(), nil}

}

func (p project) Fetch(repo *git.Repository) status {

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
		return status{p.Location, "", buf.String(), errors.New("not on master branch but fetched")}
	}
	return status{p.Location, "", buf.String(), fmt.Errorf("not on master branch and: %v", err)}

}

func (p project) Pull(worktree *git.Worktree) status {

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
		return status{p.Location, "fetched", buf.String(), nil}
	} else if err == git.NoErrAlreadyUpToDate {
		return status{p.Location, "uptodate", buf.String(), nil}
	}
	return status{p.Location, "", buf.String(), fmt.Errorf("unable to pull master: %v", err)}

}

func (p project) Clone() status {

	repo, err := p.PlainOpen()

	if err == git.ErrRepositoryNotExists {
		return p.PlainClone()
	} else if err != nil {
		return status{p.Location, "", "", fmt.Errorf("unable to open repo: %v", err)}
	}

	ref, err := repo.Head()
	if err != nil {
		return status{p.Location, "", "", err}
	}

	if ref.Name() != "refs/heads/master" {
		return p.Fetch(repo)
	}

	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return status{p.Location, "", "", fmt.Errorf("unable to get worktree: %v", err)}
	}

	return p.Pull(worktree)

}
