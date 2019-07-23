package cmd

import (
	"bytes"
	"os"

	git "gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func clone(path, url string) status {

	var buf bytes.Buffer

	repo, err := git.PlainOpen(path)
	if err == git.ErrRepositoryNotExists {

		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      url,
			Progress: &buf,
			Auth: &git_http.BasicAuth{
				Username: "token",
				Password: os.Getenv("GITLAB_TOKEN"),
			},
		})
		if err != nil {
			return status{path, "cloned", buf.String(), err}
		}
		return status{path, "cloned", buf.String(), nil}

	} else if err == nil {

		err = repo.Fetch(&git.FetchOptions{
			Progress: &buf,
			Auth: &git_http.BasicAuth{
				Username: "token",
				Password: os.Getenv("GITLAB_TOKEN"),
			},
		})
		if err == git.NoErrAlreadyUpToDate {
			return status{path, "uptodate", buf.String(), nil}
		}
		if err == nil {
			return status{path, "fetched", buf.String(), nil}
		}
		return status{path, "fetched", buf.String(), err}

	}
	return status{path, "", buf.String(), err}
}
