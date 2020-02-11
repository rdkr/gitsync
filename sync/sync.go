package sync

import (
	"errors"
	"fmt"
	"github.com/rdkr/gitsync/concurrency"
	"github.com/xanzy/go-gitlab"

	"gopkg.in/src-d/go-git.v4"
)

type ConfigParser func(Config) ([]concurrency.ProviderProcessor, []concurrency.Project)
type GitSyncer func(Git, string) concurrency.Status

func GetItemsFromCfg(cfg Config) ([]concurrency.ProviderProcessor, []concurrency.Project) {

	var groups []concurrency.ProviderProcessor
	var projects []concurrency.Project

	if len(cfg.Gitlab.Groups) > 0 || len(cfg.Gitlab.Projects) > 0 {

		c := gitlab.NewClient(nil, cfg.Gitlab.Token)

		for _, group := range cfg.Gitlab.Groups {
			groups = append(groups, &concurrency.GitlabGroupProvider{c, cfg.Gitlab.Token, "", group.Location, group.Group})
		}

		for _, project := range cfg.Gitlab.Projects {
			if project.Token == "" {
				project.Token = cfg.Gitlab.Token
			}
			projects = append(projects, project)
		}
	}

	projects = append(projects, cfg.Anon.Projects...)

	return groups, projects
}

func GitSyncHelper(g concurrency.Project) concurrency.Status {
	return GitSync(g, func(concurrency.Project) Git {
		return GitSyncProject{g}
	})
}

func GitSync(g concurrency.Project, getGitClient func(concurrency.Project) Git) concurrency.Status {

	p := getGitClient(g) // todo rename this variable

	repo, err := p.PlainOpen()

	if err == git.ErrRepositoryNotExists {

		progress, err := p.PlainClone()
		if err != nil {
			return concurrency.Status{g.Location, concurrency.StatusError, progress, fmt.Errorf("unable to clone repo: %v", err)}
		}
		return concurrency.Status{g.Location, concurrency.StatusCloned, progress, nil}

	} else if err != nil {
		return concurrency.Status{g.Location, concurrency.StatusError, "", fmt.Errorf("unable to open repo: %v", err)}
	}

	// TODO gracefully support bare repos
	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return concurrency.Status{g.Location, concurrency.StatusError, "", fmt.Errorf("unable to get worktree: %v", err)}
	}

	ref, err := repo.Head()
	if err != nil {
		return concurrency.Status{g.Location, concurrency.StatusError, "", fmt.Errorf("unable to get head: %v", err)}
	}

	if ref.Name() != "refs/heads/master" {
		progress, err := p.Fetch(repo)

		if err == git.NoErrAlreadyUpToDate || err == nil {
			return concurrency.Status{g.Location, concurrency.StatusError, progress, errors.New("not on master branch but fetched")}
		}
		return concurrency.Status{g.Location, concurrency.StatusError, progress, fmt.Errorf("not on master branch and: %v", err)}

	}

	progress, err := p.Pull(worktree)
	if err == nil {
		return concurrency.Status{g.Location, concurrency.StatusFetched, progress, nil}
	} else if err == git.NoErrAlreadyUpToDate {
		return concurrency.Status{g.Location, concurrency.StatusUpToDate, progress, nil}
	}
	return concurrency.Status{g.Location, concurrency.StatusError, progress, fmt.Errorf("unable to pull master: %v", err)}
}
