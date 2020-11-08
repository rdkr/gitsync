package concurrency

import (
	"context"

	"github.com/google/go-github/v30/github"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type Config struct {
	Github []githubConfig `yaml:"github"`
	Gitlab []gitlabConfig `yaml:"gitlab"`
	Anon   anonConfig     `yaml:"anon"`
}

// Github
type githubConfig struct {
	Users   []githubUser `yaml:"users"`
	Orgs    []githubOrg  `yaml:"orgs"`
	Token   string       `yaml:"token"`
	BaseURL string       `yaml:"baseurl"`
}

type githubUser struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

type githubOrg struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

// Gitlab
type gitlabConfig struct {
	Groups   []gitlabGroup `yaml:"groups"`
	Projects []Project     `yaml:"projects"`
	Token    string        `yaml:"token"`
	BaseURL  string        `yaml:"baseurl"`
}

type gitlabGroup struct {
	Group    int    `yaml:"group"`
	Location string `yaml:"location"`
}

// Anon
type anonConfig struct {
	Projects []Project `yaml:"projects"`
}

func GetGithubItemsFromCfg(cfg Config) ([]Group, []Project) {

	var groups []Group
	var projects []Project

	for _, gh := range cfg.Github {

		if len(gh.Users) > 0 || len(gh.Orgs) > 0 {

			if gh.Token == "" {
				logrus.Fatal("a token is required to sync GitHub users / orgs")
			}

			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: gh.Token},
			)
			tc := oauth2.NewClient(ctx, ts)

			var c *github.Client
			var err error

			if gh.BaseURL != "" {
				c, err = github.NewEnterpriseClient(gh.BaseURL, "", tc)
				if err != nil {
					logrus.Fatal(err)
				}
			} else {
				c = github.NewClient(tc)
			}

			for _, user := range gh.Users {
				groups = append(groups, &GithubUserGroup{c, user.Name, user.Location, gh.Token})
			}

			for _, org := range gh.Orgs {
				groups = append(groups, &GithubOrgGroup{c, org.Name, org.Location, gh.Token})
			}
		}
	}

	return groups, projects
}

func GetGitlabItemsFromCfg(cfg Config) ([]Group, []Project) {

	var groups []Group
	var projects []Project

	for _, gl := range cfg.Gitlab {

		if len(gl.Groups) > 0 || len(gl.Projects) > 0 {

			baseurlOption := gitlab.ClientOptionFunc(nil)
			if gl.BaseURL != "" {
				baseurlOption = gitlab.ClientOptionFunc(gitlab.WithBaseURL(gl.BaseURL))
			}

			c, err := gitlab.NewClient(gl.Token, baseurlOption)
			if err != nil {
				logrus.Fatalf("GitLab error: %v", err)
			}

			for _, group := range gl.Groups {
				groups = append(groups, &GitlabGroup{c, gl.Token, "", group.Location, group.Group})
			}

			for _, project := range gl.Projects {
				if project.Token == "" {
					project.Token = gl.Token
				}
				projects = append(projects, project)
			}
		}
	}

	// we also get the anon projects in this function...
	projects = append(projects, cfg.Anon.Projects...)

	return groups, projects
}
