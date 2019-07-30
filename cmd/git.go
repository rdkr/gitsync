package cmd

import (
	"bytes"
	"errors"
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func clone(path, url, token string) status {

	var auth *git_http.BasicAuth
	if token != "" {
		auth = &git_http.BasicAuth{
			Username: "token",
			Password: token,
		}
	}
	var buf bytes.Buffer

	repo, err := git.PlainOpen(path)
	if err == git.ErrRepositoryNotExists {

		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      url,
			Progress: &buf,
			Auth:     auth,
		})
		if err != nil {
			return status{path, "", buf.String(), fmt.Errorf("unable to clone repo: %v", err)}
		}
		return status{path, "cloned", buf.String(), nil}

	} else if err != nil {
		return status{path, "", buf.String(), fmt.Errorf("unable to open repo: %v", err)}
	}

	ref, err := repo.Head()
	if err != nil {
		return status{path, "", buf.String(), err}
	}

	if ref.Name() != "refs/heads/master" {

		err = repo.Fetch(&git.FetchOptions{
			Progress: &buf,
			Auth:     auth,
		})
		if err == git.NoErrAlreadyUpToDate {
			return status{path, "", buf.String(), errors.New("not on master branch but fetched")}
		}
		return status{path, "", buf.String(), fmt.Errorf("not on master branch and: %v", err)}

	}

	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return status{path, "", buf.String(), fmt.Errorf("unable to get worktree: %v", err)}
	}

	err = worktree.Pull(&git.PullOptions{
		Progress: &buf,
		Auth:     auth,
	})
	if err == nil {
		return status{path, "fetched", buf.String(), nil}
	} else if err == git.NoErrAlreadyUpToDate {
		return status{path, "uptodate", buf.String(), nil}
	}
	return status{path, "", buf.String(), fmt.Errorf("unable to pull master: %v", err)}

}

// func processRepo(path, url string) error {

// 	repoExists, err := repoExist(path)
// 	if err != nil {
// 		return err
// 	}

// 	if repoExists {
// 		return clone(path, repo)
// 	}
// 	if repoOnMasterBranch(repo) {
// 		return pull()
// 	}
// 	fetch()
// }

// func repoExists(path) (bool, status) {
// 	repo, err := git.PlainOpen(path)
// 	if err == git.ErrRepositoryNotExists {
// 		return repo, nil
// 	} else if err != nil {

// 	}
// }
