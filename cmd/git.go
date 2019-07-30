package cmd

import (
	"bytes"
	"errors"
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type group interface {
	getGroups() []group
	getProjects() []project
	rootFullPath() string
	rootLocation() string
}

type cloner interface {
	// getPath() string
	// getURL() string
	// getToken() string
	clone() status
}

func (p project) clone() status {

	var auth *git_http.BasicAuth
	if p.Token != "" {
		auth = &git_http.BasicAuth{
			Username: "token",
			Password: p.Token,
		}
	}
	var buf bytes.Buffer

	repo, err := git.PlainOpen(p.Location)
	if err == git.ErrRepositoryNotExists {

		_, err := git.PlainClone(p.Location, false, &git.CloneOptions{
			URL:      p.URL,
			Progress: &buf,
			Auth:     auth,
		})
		if err != nil {
			return status{p.Location, "", buf.String(), fmt.Errorf("unable to clone repo: %v", err)}
		}
		return status{p.Location, "cloned", buf.String(), nil}

	} else if err != nil {
		return status{p.Location, "", buf.String(), fmt.Errorf("unable to open repo: %v", err)}
	}

	ref, err := repo.Head()
	if err != nil {
		return status{p.Location, "", buf.String(), err}
	}

	if ref.Name() != "refs/heads/master" {

		err = repo.Fetch(&git.FetchOptions{
			Progress: &buf,
			Auth:     auth,
		})
		if err == git.NoErrAlreadyUpToDate {
			return status{p.Location, "", buf.String(), errors.New("not on master branch but fetched")}
		}
		return status{p.Location, "", buf.String(), fmt.Errorf("not on master branch and: %v", err)}

	}

	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return status{p.Location, "", buf.String(), fmt.Errorf("unable to get worktree: %v", err)}
	}

	err = worktree.Pull(&git.PullOptions{
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

// func processRepo(p.Location, p.URL string) error {

// 	repoExists, err := repoExist(p.Location)
// 	if err != nil {
// 		return err
// 	}

// 	if repoExists {
// 		return clone(p.Location, repo)
// 	}
// 	if repoOnMasterBranch(repo) {
// 		return pull()
// 	}
// 	fetch()
// }

// func repoExists(p.Location) (bool, status) {
// 	repo, err := git.PlainOpen(p.Location)
// 	if err == git.ErrRepositoryNotExists {
// 		return repo, nil
// 	} else if err != nil {

// 	}
// }
