package cmd

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

type GetItemsFromCfg func() ([]ProviderProcessor, []Project)
type GitSyncer func(Git, string) Status

func getItemsFromCfg() ([]ProviderProcessor, []Project) {

	var groups []ProviderProcessor
	var projects []Project

	if len(cfg.Gitlab.Groups) > 0 || len(cfg.Gitlab.Projects) > 0 {

		// TODO improve the handline of no / bad token
		token := os.Getenv("GITLAB_TOKEN")
		if len(token) == 0 {
			panic("bad token?")
		}
		c := gitlab.NewClient(nil, token)

		for _, group := range cfg.Gitlab.Groups {
			groups = append(groups, &gitlabGroupProvider{c, token, "", group.Location, group.Group})
		}

		for _, project := range cfg.Gitlab.Projects {
			if project.Token == "" {
				project.Token = token
			}
			projects = append(projects, project)
		}
	}

	for _, project := range cfg.Anon.Projects {
		projects = append(projects, project)
	}

	return groups, projects
}

func GitSync(p Git, location string) Status {

	repo, err := p.PlainOpen()

	if err == git.ErrRepositoryNotExists {

		progress, err := p.PlainClone()
		if err != nil {
			return Status{location, "", progress, fmt.Errorf("unable to clone repo: %v", err)}
		}
		return Status{location, "cloned", progress, nil}

	} else if err != nil {
		return Status{location, "", "", fmt.Errorf("unable to open repo: %v", err)}
	}

	ref, err := repo.Head()
	if err != nil {
		return Status{location, "", "", fmt.Errorf("unable to get head: %v", err)}
	}

	if ref.Name() != "refs/heads/master" {
		progress, err := p.Fetch(repo)

		if err == git.NoErrAlreadyUpToDate {
			return Status{location, "", progress, errors.New("not on master branch but fetched")}
		}
		return Status{location, "", progress, fmt.Errorf("not on master branch and: %v", err)}

	}

	// TODO gracefully support bare repos
	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return Status{location, "", "", fmt.Errorf("unable to get worktree: %v", err)}
	}

	progress, err := p.Pull(worktree)
	if err == nil {
		return Status{location, "fetched", progress, nil}
	} else if err == git.NoErrAlreadyUpToDate {
		return Status{location, "uptodate", progress, nil}
	}
	return Status{location, "", progress, fmt.Errorf("unable to pull master: %v", err)}
}
