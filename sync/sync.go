package sync

import (
	"errors"
	"fmt"

	"github.com/rdkr/gitsync/concurrency"

	"gopkg.in/src-d/go-git.v4"
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

func GitSyncHelper(g concurrency.Project) interface{} {
	return GitSync(GitSyncProject{g})
}

func GitSync(g Git) Status {

	p := g // todo rename this variable

	repo, err := p.PlainOpen()

	if err == git.ErrRepositoryNotExists {

		progress, err := p.PlainClone()
		if err != nil {
			return Status{g.GetLocation(), StatusError, progress, fmt.Errorf("unable to clone repo: %v", err)}
		}
		return Status{g.GetLocation(), StatusCloned, progress, nil}

	} else if err != nil {
		return Status{g.GetLocation(), StatusError, "", fmt.Errorf("unable to open repo: %v", err)}
	}

	// TODO gracefully support bare repos
	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return Status{g.GetLocation(), StatusError, "", fmt.Errorf("unable to get worktree: %v", err)}
	}

	ref, err := repo.Head()
	if err != nil {
		return Status{g.GetLocation(), StatusError, "", fmt.Errorf("unable to get head: %v", err)}
	}

	if ref.Name() != "refs/heads/master" {
		progress, err := p.Fetch(repo)

		if err == git.NoErrAlreadyUpToDate || err == nil {
			return Status{g.GetLocation(), StatusError, progress, errors.New("not on master branch but fetched")}
		}
		return Status{g.GetLocation(), StatusError, progress, fmt.Errorf("not on master branch and: %v", err)}

	}

	progress, err := p.Pull(worktree)
	if err == nil {
		return Status{g.GetLocation(), StatusFetched, progress, nil}
	} else if err == git.NoErrAlreadyUpToDate {
		return Status{g.GetLocation(), StatusUpToDate, progress, nil}
	}
	return Status{g.GetLocation(), StatusError, progress, fmt.Errorf("unable to pull master: %v", err)}
}
