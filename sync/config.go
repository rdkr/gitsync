package sync

import (
	"github.com/google/go-github/v30/github"
	"github.com/rdkr/gitsync/concurrency"
	"github.com/xanzy/go-gitlab"
)

type Config struct {
	Github githubConfig `yaml:"github"`
	Gitlab gitlabConfig `yaml:"gitlab"`
	Anon   anonConfig   `yaml:"anon"`
}

// Github
type githubConfig struct {
	// Groups   []gitlabGroup         `yaml:"groups"`
	// Projects []concurrency.Project `yaml:"projects"`
	Users []githubUser `yaml:"users"`
	Token string       `yaml:"token"`
}

type githubUser struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

// Gitlab
type gitlabConfig struct {
	Groups   []gitlabGroup         `yaml:"groups"`
	Projects []concurrency.Project `yaml:"projects"`
	Token    string                `yaml:"token"`
}

type gitlabGroup struct {
	Group    int    `yaml:"group"`
	Location string `yaml:"location"`
}

// Anon
type anonConfig struct {
	Projects []concurrency.Project `yaml:"projects"`
}

type ConfigParser func(Config) ([]concurrency.User, []concurrency.Group, []concurrency.Project)

func GetGithubItemsFromCfg(cfg Config) ([]concurrency.User, []concurrency.Group, []concurrency.Project) {

	var users []concurrency.User
	var groups []concurrency.Group
	var projects []concurrency.Project

	// if len(cfg.Github.Groups) > 0 || len(cfg.Github.Projects) > 0 || len(cfg.Github.Users) > 0 {
	if len(cfg.Github.Users) > 0 {

		c := github.NewClient(nil)

		for _, user := range cfg.Github.Users {
			users = append(users, &concurrency.GithubUser{c, user.Name, user.Location, cfg.Github.Token})
		}
	}

	projects = append(projects, cfg.Anon.Projects...)

	return users, groups, projects
}

func GetGitlabItemsFromCfg(cfg Config) ([]concurrency.User, []concurrency.Group, []concurrency.Project) {

	var users []concurrency.User
	var groups []concurrency.Group
	var projects []concurrency.Project

	if len(cfg.Gitlab.Groups) > 0 || len(cfg.Gitlab.Projects) > 0 {

		c := gitlab.NewClient(nil, cfg.Gitlab.Token)

		for _, group := range cfg.Gitlab.Groups {
			groups = append(groups, &concurrency.GitlabGroup{c, cfg.Gitlab.Token, "", group.Location, group.Group})
		}

		for _, project := range cfg.Gitlab.Projects {
			if project.Token == "" {
				project.Token = cfg.Gitlab.Token
			}
			projects = append(projects, project)
		}
	}

	projects = append(projects, cfg.Anon.Projects...)

	return users, groups, projects
}
