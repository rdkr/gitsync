package cmd

import (
	"os"

	git "gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func clone(path, url string) status {

	repo, err := git.PlainOpen(path)
	if err == git.ErrRepositoryNotExists {

		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      url,
			Progress: nil,
			Auth: &git_http.BasicAuth{
				Username: "token",
				Password: os.Getenv("GITLAB_TOKEN"),
			},
		})
		if err != nil {
			return status{path, "clone", err}
		}
		return status{path, "clone", nil}

	} else if err == nil {

		err = repo.Fetch(&git.FetchOptions{
			Auth: &git_http.BasicAuth{
				Username: "token",
				Password: os.Getenv("GITLAB_TOKEN"),
			},
		})
		if err == nil || err == git.NoErrAlreadyUpToDate {
			return status{path, "fetch", nil}
		}
		return status{path, "fetch", err}

	}
	return status{path, "open", err}
}
