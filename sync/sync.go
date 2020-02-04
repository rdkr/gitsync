package sync

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"

	"gopkg.in/src-d/go-git.v4"
)

type ConfigParser func(Config) ([]ProviderProcessor, []Project)
type GitSyncer func(Git, string) Status

func GetItemsFromCfg(cfg Config) ([]ProviderProcessor, []Project) {

	var groups []ProviderProcessor
	var projects []Project

	if len(cfg.Gitlab.Groups) > 0 || len(cfg.Gitlab.Projects) > 0 {

		c := gitlab.NewClient(nil, cfg.Gitlab.Token)

		for _, group := range cfg.Gitlab.Groups {
			groups = append(groups, &gitlabGroupProvider{c, cfg.Gitlab.Token, "", group.Location, group.Group})
		}

		for _, project := range cfg.Gitlab.Projects {
			if project.Token == "" {
				project.Token = cfg.Gitlab.Token
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
			return Status{location, StatusError, progress, fmt.Errorf("unable to clone repo: %v", err)}
		}
		return Status{location, StatusCloned, progress, nil}

	} else if err != nil {
		return Status{location, StatusError, "", fmt.Errorf("unable to open repo: %v", err)}
	}

	ref, err := repo.Head()
	if err != nil {
		return Status{location, StatusError, "", fmt.Errorf("unable to get head: %v", err)}
	}

	if ref.Name() != "refs/heads/master" {
		progress, err := p.Fetch(repo)

		if err == git.NoErrAlreadyUpToDate || err == nil {
			return Status{location, StatusError, progress, errors.New("not on master branch but fetched")}
		}
		return Status{location, StatusError, progress, fmt.Errorf("not on master branch and: %v", err)}

	}

	// TODO gracefully support bare repos
	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return Status{location, StatusError, "", fmt.Errorf("unable to get worktree: %v", err)}
	}

	progress, err := p.Pull(worktree)
	if err == nil {
		return Status{location, StatusFetched, progress, nil}
	} else if err == git.NoErrAlreadyUpToDate {
		return Status{location, StatusUpToDate, progress, nil}
	}
	return Status{location, StatusError, progress, fmt.Errorf("unable to pull master: %v", err)}
}
